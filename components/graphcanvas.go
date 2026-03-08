package components

import (
	"fmt"
	"math"
	"syscall/js"

	"github.com/dq1Mango/graph-theory/actions"
	"github.com/dq1Mango/graph-theory/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
	"github.com/hmdsefi/gograph"
)

const (
	VERTEX_RADIUS = 5
	CYCLE_RADIUS  = 30.0
)

type GraphCanvas struct {
	vecty.Core
	ctx     js.Value
	Id      string
	Size    uint
	Actions chan any
	Graph   gograph.Graph[uint]

	bijection       model.Bijection
	nextLabel       uint
	vertexPositions map[uint]model.Point
	highlitedVertex uint

	transform model.Transform
}

func InitGraphCanvas(size uint, id string) GraphCanvas {

	graph := GraphCanvas{
		Size:            size,
		Id:              id,
		Graph:           gograph.New[uint](),
		Actions:         make(chan any, 10),
		bijection:       model.NewBijection(),
		nextLabel:       0,
		vertexPositions: make(map[uint]model.Point)}

	return graph
}

func (g *GraphCanvas) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("canvas-wrapper")),
		elem.Canvas(
			vecty.Markup(
				vecty.Class("canvas"),
				prop.ID(g.Id),

				event.MouseMove(func(e *vecty.Event) {
					// mouse position updates aren't super important,
					// and with WASM we only have one os thread,
					// so if we are busy in any way just dont add more fuel to the fire
					if len(g.Actions) < 2 {
						point := model.Point{X: e.Get("offsetX").Float(), Y: e.Get("offsetY").Float()}
						point.Scale(2)
						point = g.transform.Backwards(point)

						// non blocking send just in case
						select {
						case g.Actions <- &actions.MouseMove{
							Pos: point,
						}:
						default:
							fmt.Println("action buffer busy...")
						}
					}
				}),
			),
			// vecty.Property("height", 20),
		),
	)
}

func (g *GraphCanvas) SetCanvasTransform() {

	dpr := js.Global().Get("devicePixelRatio").Float()
	canvas := g.ctx.Get("canvas")
	rect := canvas.Call("getBoundingClientRect")

	width := rect.Get("width").Float()
	height := rect.Get("width").Float()

	if width != height {
		fmt.Printf("width and height of canvas: %s not equal\n", g.Id)
	}

	// fmt.Printf("width: %f, height: %f\n", width, height)

	canvas.Set("width", width*dpr)
	canvas.Set("height", height*dpr)

	canvas.Get("style").Set("width", fmt.Sprintf("%fpx", width))
	canvas.Get("style").Set("height", fmt.Sprintf("%fpx", height))

	scale := dpr * math.Min(width, height) / (float64(g.Size))

	shift := dpr * width / 2

	g.ctx.Call("setTransform", scale, 0, 0, -scale, shift, shift)

	g.transform = model.NewTransform(scale, -scale, shift, shift)
}

func (c *GraphCanvas) Mount() {
	canvas := js.Global().Get("document").Call("getElementById", c.Id)
	c.ctx = canvas.Call("getContext", "2d")

	c.SetCanvasTransform()

	// safe to draw here, DOM is ready
	c.ctx.Set("fillStyle", "red")
	c.ctx.Call("fillRect", 0, 0, 50, 25)

	go func() {
		c.handleActions()
	}()

}

func (g *GraphCanvas) Clear() {
	s := int(g.Size)
	// fmt.Println(-g.Size / 2)
	// g.ctx.Call("clearRect", -g.Size/2, g.Size/2, g.Size, -g.Size)
	g.ctx.Call("clearRect", -s/2, s/2, s, -s)
}

func (g *GraphCanvas) DrawNodeOutline(point model.Point) {

	g.ctx.Call("beginPath")
	g.ctx.Call("arc", point.X, point.Y, VERTEX_RADIUS, 0, 2*math.Pi)
	g.ctx.Call("stroke")
}

func (g *GraphCanvas) DrawNode(point model.Point) {

	// g.ctx.Set("strokeStyle", "blue")

	g.ctx.Call("beginPath")
	g.ctx.Call("arc", point.X, point.Y, VERTEX_RADIUS, 0, 2*math.Pi)

	g.ctx.Call("fill")
	g.ctx.Call("stroke")
}

func (g *GraphCanvas) DrawEdge(from model.Point, to model.Point) {

	g.ctx.Call("beginPath")
	g.ctx.Call("moveTo", from.X, from.Y)
	g.ctx.Call("lineTo", to.X, to.Y)

	g.ctx.Call("stroke")
}

func (g *GraphCanvas) Draw() {
	g.Clear()

	// order := g.Graph.Order()
	// radius := 30.0
	g.ctx.Set("lineWidth", 1)

	g.ctx.Set("strokeStyle", "black")

	for _, edge := range g.Graph.AllEdges() {
		g.DrawEdge(g.vertexPositions[edge.Source().Label()], g.vertexPositions[edge.Destination().Label()])
	}

	g.ctx.Set("fillStyle", "grey")

	// draw the verticies after the edges to draw over them
	for _, point := range g.vertexPositions {
		g.DrawNode(point)

	}

	g.ctx.Call("save")
	g.ctx.Set("fillStyle", "white")
	g.ctx.Set("textBaseline", "middle")
	g.ctx.Set("textAlign", "center")

	// ALERT: floating magic number over here
	fontSize := float64(g.Size) * 0.067
	g.ctx.Set("font", fmt.Sprintf("%.2fpx Arial", fontSize))

	// flip y back to normal for this draw call
	g.ctx.Call("transform", 1, 0, 0, -1, 0, 0)
	for _, label := range g.bijection.Labels() {
		point := g.vertexPositions[g.bijection.Forwards(label)]

		// y coordinate needs to be negated since we flipped
		g.ctx.Call("fillText", label, point.X, -point.Y)
	}

	g.ctx.Call("restore")

	fmt.Println("drew da graph")
}

func (g *GraphCanvas) HighlightActiveVertex(mousePos model.Point) actions.Action {

	g.ctx.Set("lineWidth", 1)
	g.ctx.Set("strokeStyle", "black")

	for id, pos := range g.vertexPositions {
		if mousePos.Distance(pos) <= VERTEX_RADIUS+1 {

			g.highlitedVertex = id
			fmt.Println("highlited vertex id:", id)
			g.ctx.Set("strokeStyle", "red")
		}

		g.DrawNodeOutline(pos)
		g.ctx.Set("strokeStyle", "black")
	}

	return nil
}

func (g *GraphCanvas) addNextVertex(connected bool) actions.Action {
	// order := uint(c.Graph.Order())

	vertex := gograph.NewVertex(g.nextLabel)
	// g.bijection = append(g.bijection, g.nextLabel)
	g.bijection.Add(0, g.nextLabel)
	g.nextLabel++

	if connected {
		for _, v := range g.Graph.GetAllVertices() {
			g.Graph.AddEdge(vertex, v)
		}
	}
	g.Graph.AddVertex(vertex)

	// c.Actions <- Draw{}

	fmt.Println("added vertex")

	return &actions.RecomputeVertexPositions{}
}

func (g *GraphCanvas) RecomputeVertexPositions() actions.Action {
	order := g.Graph.Order()
	// verticies := g.Graph.GetAllVertices()
	labels := g.bijection.Labels()
	fmt.Println(labels)
	radius := CYCLE_RADIUS

	g.vertexPositions = make(map[uint]model.Point)

	if order > 1 {
		// vertexPositions := make(map[*gograph.Vertex[uint]]model.Point, 0)

		deltaTheta := 2 * math.Pi / float64(order)

		for i, label := range labels {
			point := model.PointFromTheta(deltaTheta * float64(i))
			point.Scale(radius)

			v := g.bijection.Forwards(label)

			g.vertexPositions[v] = point

		}

	} else if order == 1 {
		g.vertexPositions[labels[0]] = model.Point{X: 0, Y: 0}
	}

	fmt.Println("recomputed vertex positions")

	return &actions.Draw{}

}

func (g *GraphCanvas) handleActions() {
	for {
		a := <-g.Actions
		for a != nil {

			switch action := a.(type) {

			case *actions.AddVertex:
				a = g.addNextVertex(action.Connected)
				continue

			case *actions.RecomputeVertexPositions:
				a = g.RecomputeVertexPositions()
				continue

			case *actions.Draw:
				g.Draw()

			case *actions.MouseMove:
				// fmt.Println(action.pos)
				g.HighlightActiveVertex(action.Pos)

			default:
				// fmt.Fprintln(os.Stderr, "Unhandled action of type: ", action)
				fmt.Println("Unhandled action of type: ", action)
			}

			a = nil
		}
	}
}
