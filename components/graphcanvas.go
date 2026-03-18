package components

import (
	"fmt"
	"math"
	"slices"
	"syscall/js"

	"github.com/dq1Mango/graph-theory/actions"
	"github.com/dq1Mango/graph-theory/model"
	"github.com/hmdsefi/gograph"
	"github.com/hmdsefi/gograph/util"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"
)

const (
	VERTEX_RADIUS = 5
	CYCLE_RADIUS  = 30.0
	LINE_WIDTH    = 1
	FONT_SIZE     = 6.7
)

// Display modes
const (
	FreeForm uint = iota
	Ring
)

var WhereNumbersStart uint = 0

func SortVerticies(verticies []*gograph.Vertex[uint]) {
	slices.SortFunc(verticies,
		func(u, v *gograph.Vertex[uint]) int {
			// uLabel, vLabel := g.Bijection.Forwards(u.Label()), g.Bijection.Forwards(v.Label())
			if u.Label() > v.Label() {
				return 1
			} else if u.Label() == v.Label() {
				return 0
			} else {
				return -1
			}
		},
	)
}

type GraphCanvas struct {
	vecty.Core
	ctx js.Value

	Id      string
	Size    uint
	Actions chan any
	Graph   gograph.Graph[uint]
	Stats   *GraphStats
	Mode    *ChoiceBar

	Labeling string

	// Bijection       model.Bijection
	VertexPositions map[uint]model.Point
	SelectedVertex  *uint
	ShiftSelected   *uint

	Transform model.Transform
	Iterator  *AlgorithmWalk
	Locked    bool
}

func InitGraphCanvas(size uint, id string) GraphCanvas {

	actionChan := make(chan any, 10)

	// little function to convert the mode selector updates to draw updates
	modeUpdates := make(chan string)
	go func() {
		for {
			// we dont rly care what the mode changes to or from
			<-modeUpdates
			// fmt.Println("sent draw upadte")
			actionChan <- &actions.RecomputeVertexPositions{}
		}
	}()

	graph := GraphCanvas{
		Size:     size,
		Id:       id,
		Graph:    gograph.New[uint](),
		Actions:  actionChan,
		Mode:     NewChoiceBar(modeUpdates, "freeform", "ring"),
		Labeling: "continous",

		// Bijection:       model.NewBijection(),
		VertexPositions: make(map[uint]model.Point),
		Iterator:        &AlgorithmWalk{},
	}

	// Start with the trivial graph
	graph.Actions <- &actions.MouseDown{Pos: model.Point{X: 0, Y: 0}}
	// graph.addNextVertex(&model.Point{X: 0, Y: 0}, false)

	return graph
}

func (g *GraphCanvas) Render() vecty.ComponentOrHTML {
	return elem.Div(

		g.Mode,
		elem.Div(

			vecty.Markup(vecty.Class("canvas-wrapper"), prop.ID("canvas-wrapper")),
			elem.Canvas(
				vecty.Markup(
					vecty.Class("canvas"),
					prop.ID(g.Id),

					event.MouseMove(func(e *vecty.Event) {

						// mouse position updates aren't super important,
						// and with WASM we only have one os thread,
						// so if we are busy in any way just dont add more fuel to the fire
						if len(g.Actions) < 2 {
							point := g.PointFromMouseEvent(e)
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

					event.MouseDown(func(e *vecty.Event) {

						shift := e.Get("shiftKey").Bool()
						point := g.PointFromMouseEvent(e)

						g.Actions <- &actions.MouseDown{Pos: point, Shift: shift}
					}),

					// these closures are rly nice i have to say
				),
			),
		),
		g.Iterator,
	)
}

func (g *GraphCanvas) PointFromMouseEvent(e *vecty.Event) model.Point {
	point := model.Point{
		X: e.Get("offsetX").Float(),
		Y: e.Get("offsetY").Float(),
	}
	point.Scale(2)
	point = g.Transform.Backwards(point)
	return point

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

	wrapper := js.Global().Get("document").Call("getElementById", "canvas-wrapper")

	wrapper_height := wrapper.Get("clientHeight").Float()
	wrapper_width := wrapper.Get("clientWidth").Float()

	limiter := math.Min(wrapper_height, wrapper_width)

	// wrapper.Set("clientWidth", limiter)
	// wrapper.Set("clientHeight", limiter)
	wrapper.Get("style").Set("width", fmt.Sprintf("%fpx", limiter))
	wrapper.Get("style").Set("height", fmt.Sprintf("%fpx", limiter))

	// canvas.Get("style").Set("width", fmt.Sprintf("%fpx", width))
	// canvas.Get("style").Set("height", fmt.Sprintf("%fpx", height))

	scale := dpr * math.Min(width, height) / (float64(g.Size))

	shift := dpr * width / 2

	g.ctx.Call("setTransform", scale, 0, 0, -scale, shift, shift)

	g.Transform = model.NewTransform(scale, -scale, shift, shift)

	g.Actions <- &actions.Draw{}
}

func (c *GraphCanvas) Mount() {
	canvas := js.Global().Get("document").Call("getElementById", c.Id)
	wrapper := js.Global().Get("document").Call("getElementById", "canvas-wrapper")
	c.ctx = canvas.Call("getContext", "2d")

	resizeFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		c.SetCanvasTransform()
		return nil
	})

	observer := js.Global().Get("ResizeObserver").New(resizeFunc)
	observer.Call("observe", wrapper)
	// c.observer = observer

	c.SetCanvasTransform()

	// safe to draw here, DOM is ready
	// c.ctx.Set("fillStyle", model.CurrentPalette.Red)
	// c.ctx.Call("fillRect", 0, 0, 50, 25)

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
	g.ctx.Set("lineWidth", LINE_WIDTH)

	g.ctx.Set("strokeStyle", model.CurrentPalette.Base)

	for _, edge := range g.Graph.AllEdges() {
		g.DrawEdge(
			g.VertexPositions[edge.Source().Label()],
			g.VertexPositions[edge.Destination().Label()],
		)
	}

	g.ctx.Set("fillStyle", model.CurrentPalette.Surface1)

	// draw the verticies after the edges to draw over them
	for _, point := range g.VertexPositions {
		g.DrawNode(point)

	}

	// highlight the selected one
	if g.SelectedVertex != nil {
		g.ctx.Set("strokeStyle", model.CurrentPalette.Red)
		g.DrawNodeOutline(g.VertexPositions[*g.SelectedVertex])
	}

	if g.ShiftSelected != nil {
		g.ctx.Set("strokeStyle", model.CurrentPalette.Blue)
		g.DrawNodeOutline(g.VertexPositions[*g.ShiftSelected])
	}

	g.ctx.Call("save")
	g.ctx.Set("fillStyle", model.CurrentPalette.Text)
	g.ctx.Set("textBaseline", "middle")
	g.ctx.Set("textAlign", "center")

	// ALERT: floating magic number over here
	// fontSize := float64(g.Size) * 0.067
	g.ctx.Set("font", fmt.Sprintf("%.2fpx Arial", FONT_SIZE))

	// flip y back to normal for this draw call
	g.ctx.Call("transform", 1, 0, 0, -1, 0, 0)
	for label, point := range g.VertexPositions {
		// point := g.VertexPositions[id]

		// y coordinate needs to be negated since we flipped
		g.ctx.Call("fillText", label+WhereNumbersStart, point.X, -point.Y)
	}

	g.ctx.Call("restore")

	fmt.Println("drew da graph")
}

func (g *GraphCanvas) HighlightActiveVertex(mousePos model.Point) actions.Action {

	g.ctx.Set("lineWidth", LINE_WIDTH)
	g.ctx.Set("strokeStyle", model.CurrentPalette.Base)

	for _, pos := range g.VertexPositions {

		if mousePos.Distance(pos) <= VERTEX_RADIUS+1 {

			g.ctx.Set("strokeStyle", model.CurrentPalette.Mauve)
		}

		g.DrawNodeOutline(pos)
		g.ctx.Set("strokeStyle", model.CurrentPalette.Base)
	}

	g.HighlightSelectedVerticies()

	return nil
}

func (g *GraphCanvas) HighlightSelectedVerticies() {

	g.ctx.Set("lineWidth", LINE_WIDTH)

	if id := g.SelectedVertex; id != nil {
		g.ctx.Set("strokeStyle", model.CurrentPalette.Red)
		g.DrawNodeOutline(g.VertexPositions[*id])
	}

	if id := g.ShiftSelected; id != nil {
		g.ctx.Set("strokeStyle", model.CurrentPalette.Blue)
		g.DrawNodeOutline(g.VertexPositions[*id])
	}
}

func (g *GraphCanvas) HandleMouseClick(mousePos model.Point, shift bool) actions.Action {
	fmt.Println("got a click...")

	for id, pos := range g.VertexPositions {
		if mousePos.Distance(pos) <= VERTEX_RADIUS+1 {

			if !shift {

				if g.SelectedVertex != nil && *g.SelectedVertex == id {
					g.SelectedVertex = nil
				} else if g.ShiftSelected == nil || id != *g.ShiftSelected {
					g.SelectedVertex = &id
				}

			} else {

				if g.ShiftSelected != nil && *g.ShiftSelected == id {
					g.ShiftSelected = nil
				} else if g.SelectedVertex == nil || id != *g.SelectedVertex {
					g.ShiftSelected = &id

					// this second condition isnt nessecary bc of the prior control flow,
					// but it is here anyway for the spirit of go's readability
					if g.SelectedVertex != nil && *g.SelectedVertex != id {
						return &actions.AddEdge{Vertex1: g.SelectedVertex, Vertex2: g.ShiftSelected}
					}
				}

			}

			return &actions.Draw{}
		}
	}

	nextLabel := g.NextLabel()
	result := g.addNextVertex(&mousePos, false)
	g.SelectedVertex = &nextLabel

	return result
}

func (g *GraphCanvas) NextLabel() uint {

	switch g.Labeling {

	case "static":
		for expected := range g.Graph.Order() {
			if g.Graph.GetVertexByID(uint(expected)) == nil {
				return uint(expected)
			}
		}
		return uint(g.Graph.Order())

	case "continous":
		return uint(g.Graph.Order())

	default:
		panic("how did this happen")

	}
}

func (g *GraphCanvas) SetLabeling(labeling string) actions.Action {

	switch labeling {
	case "continous":

		g.ensureSequential()

	case "static":

	default:
		fmt.Println("Uknown labeling choice: ", labeling)
		return nil
	}

	fmt.Println("set labelinng to: ", labeling)
	g.Labeling = labeling
	return &actions.Draw{}

}

func (g *GraphCanvas) AddVertex(id uint, point *model.Point, connected bool) actions.Action {
	// order := uint(c.Graph.Order())

	vertex := gograph.NewVertex(id)

	if g.Mode.Selected == FreeForm {
		if point != nil {
			g.VertexPositions[vertex.Label()] = *point
		} else {
			fmt.Println("Error: freeform add without pos")
			return nil
		}
	}

	if connected {
		for _, v := range g.Graph.GetAllVertices() {
			g.Graph.AddEdge(vertex, v)
		}
	} else if g.SelectedVertex != nil {
		// fmt.Println(vertex.Label(), *g.SelectedVertex)
		g.Graph.AddEdge(vertex, g.Graph.GetVertexByID(*g.SelectedVertex))
	}

	g.Graph.AddVertex(vertex)

	fmt.Println("added vertex")

	return &actions.RecomputeVertexPositions{}
}

func (g *GraphCanvas) addNextVertex(point *model.Point, connected bool) actions.Action {
	// order := uint(c.Graph.Order())

	return g.AddVertex(g.NextLabel(), point, connected)
}

func (g *GraphCanvas) ensureSequential() {
	// ids := make([])
	verticies := g.Graph.GetAllVertices()

	// slices.SortFunc(verticies)
	SortVerticies(verticies)

	var expectedLabel uint = 0
	for _, vertex := range verticies {
		if vertex.Label() != expectedLabel {
			// for i := index; i < len(verticies); i++ {
			// 	label := verticies[i].Label()
			// 	g.Graph.ChangeLabel(label, expectedLabel+uint(i))
			// }
			// fmt.Println("expected: ", expectedLabel, "actual:", vertex.Label())

			g.VertexPositions[expectedLabel] = g.VertexPositions[vertex.Label()]
			delete(g.VertexPositions, vertex.Label())
			g.Graph.ChangeLabel(vertex.Label(), expectedLabel)
			// if v := g.Graph.GetVertexByID(expectedLabel); v == nil {
			// if err
			//
			// }

			// fmt.Println("expected: ", expectedLabel, "changed:", vertex.Label())
			// fmt.Println("i dont get it;", g.Graph.GetVertexByID(expectedLabel))

		}
		expectedLabel++
	}
}

func (g *GraphCanvas) RemoveVertex(id *uint) actions.Action {

	if id != nil {

		if g.Graph.Order() == 1 {
			return &actions.PopupMessage{Message: "Vertex set cannot be empty"}
		}

		g.Graph.RemoveVertices(gograph.NewVertex(*id))

		delete(g.VertexPositions, *id)

		g.SelectedVertex = nil
		g.ShiftSelected = nil

		if g.Labeling == "continous" {
			g.ensureSequential()
		}

		return &actions.RecomputeVertexPositions{}
	} else {
		return &actions.PopupMessage{Message: "No Selected Vertex"}
	}
}

func (g *GraphCanvas) AddEdge(vertex1, vertex2 *uint) actions.Action {

	if vertex1 != nil {
		if vertex2 != nil {

			g.Graph.AddEdge(gograph.NewVertex(*vertex1), gograph.NewVertex(*vertex2))
			return &actions.Draw{}

		} else {
			return &actions.PopupMessage{Message: "Secondary Vertex Not Selected"}
		}
	} else {
		return &actions.PopupMessage{Message: "Primary Vertex Not Selected"}
	}

}

func (g *GraphCanvas) RemoveEdge(vertex1, vertex2 *uint) actions.Action {

	if vertex1 != nil {
		if vertex2 != nil {
			g.Graph.RemoveEdges(
				gograph.NewEdge(gograph.NewVertex(*vertex1), gograph.NewVertex(*vertex2)),
			)
			return &actions.Draw{}
		} else {
			return &actions.PopupMessage{Message: "Secondary Vertex Not Selected"}
		}
	} else {
		return &actions.PopupMessage{Message: "Primary Vertex Not Selected"}
	}

}

func (g *GraphCanvas) RecomputeVertexPositions() actions.Action {

	if g.Mode.Selected == FreeForm {
		return &actions.Draw{}
	}

	order := g.Graph.Order()
	// verticies := g.Graph.GetAllVertices()
	// fmt.Println(labels)
	radius := CYCLE_RADIUS

	g.VertexPositions = make(map[uint]model.Point)

	if order > 1 {
		// vertexPositions := make(map[*gograph.Vertex[uint]]model.Point, 0)

		deltaTheta := 2 * math.Pi / float64(order)

		for i, vertex := range g.Graph.GetAllVertices() {
			point := model.PointFromTheta(deltaTheta * float64(i))
			point.Scale(radius)

			// v := g.Bijection.Backwards(label)

			g.VertexPositions[vertex.Label()] = point

		}

	} else if order == 1 {
		g.VertexPositions[0] = model.Point{X: 0, Y: 0}
	}

	fmt.Println("recomputed vertex positions")

	return &actions.Draw{}

}

func (g *GraphCanvas) ClearSelections() {
	g.SelectedVertex = nil
	g.ShiftSelected = nil
}

func (g *GraphCanvas) PruferFromGraph() actions.Action {
	g.ClearSelections()

	code, err := util.PruferCodeFromGraph(g.Graph)

	if err != nil {
		InfoChan <- err.Error()
		return &actions.RecomputeVertexPositions{}
	}

	fmt.Println(code)

	return &actions.RecomputeVertexPositions{}
}

func (g *GraphCanvas) GraphFromPrufer(pruferCode []uint) actions.Action {

	for i := range pruferCode {
		pruferCode[i] -= WhereNumbersStart
	}

	newGraph, err := util.GraphFromPruferCode(pruferCode...)

	if err != nil {
		InfoChan <- err.Error()
		return nil
	}

	g.ClearSelections()

	g.Mode.SetMode(Ring)

	g.Graph = newGraph

	return &actions.RecomputeVertexPositions{}

}

func (g *GraphCanvas) UpdateStats() {
	vecty.Rerender(g.Stats)
}

// func (g *GraphCanvas) SetBijection(bi model.Bijection) actions.Action {
// 	ids := g.Bijection.SortedIds()
//
// 	fmt.Println("got these ids", ids)
//
// 	g.Bijection = bi
//
// 	g.Bijection.Add(ids...)
//
// 	return &actions.Draw{}
// }

func (g *GraphCanvas) SetIterator(name string) actions.Action {
	var iterator GraphIterator
	var err error

	switch name {
	case "pruferCode":
		iterator, err = NewPruferCodeIterator(g)
		g.SetLabeling("static")
		g.Mode.SetMode(FreeForm)

	default:
		fmt.Println("Uknown iterator: ", name)
		return nil
	}

	if err != nil {
		InfoChan <- err.Error()
		return nil
	}

	// g.SetBijection(model.NewBijection())
	// g.SetBijection(model.NewFakeBijection())

	g.Locked = true

	if g.Iterator.iterator != nil {
		//g.Iterator.UnIterate()
	}

	g.Iterator.iterator = &iterator

	vecty.Rerender(g.Iterator)

	return &actions.Draw{}
}

func (g *GraphCanvas) handleActions() {
	for {
		a := <-g.Actions
		for a != nil {

			switch action := a.(type) {

			case *actions.AddVertex:
				a = g.addNextVertex(nil, action.Connected)

			case *actions.RemoveVertex:
				a = g.RemoveVertex(action.Id)

			case *actions.AddEdge:
				a = g.AddEdge(action.Vertex1, action.Vertex2)

			case *actions.RemoveEdge:
				a = g.RemoveEdge(action.Vertex1, action.Vertex2)

			case *actions.RecomputeVertexPositions:
				a = g.RecomputeVertexPositions()

			case *actions.Draw:
				g.Draw()
				g.UpdateStats()
				a = nil

			case *actions.MouseMove:
				// fmt.Println(action.pos)
				g.HighlightActiveVertex(action.Pos)
				a = nil

			case *actions.MouseDown:
				a = g.HandleMouseClick(action.Pos, action.Shift)

			case *actions.PruferFromGraph:
				a = g.PruferFromGraph()
				// fmt.Printf("heres the type: %T\n", a)

			case *actions.GraphFromPrufer:
				a = g.GraphFromPrufer(action.Prufer)

			case *actions.SetIterator:
				a = g.SetIterator(action.Iterator)

			case *actions.SetLabeling:
				a = g.SetLabeling(action.Labeling)

			case *actions.PopupMessage:
				fmt.Println("Popup: ", action.Message)
				InfoChan <- action.Message
				a = nil

			default:
				// fmt.Fprintln(os.Stderr, "Unhandled action of type: ", action)
				fmt.Printf("Unhandled action of type: %T\n", action)
				a = nil
			}

		}
	}
}
