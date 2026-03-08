package main

import (
	"bytes"
	"fmt"
	"math"
	// "os"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"

	// "github.com/hexops/vecty/event"
	// "github.com/hexops/vecty/prop"
	"syscall/js"

	// "github.com/hexops/vecty/style"
	"github.com/hmdsefi/gograph"
	"github.com/yuin/goldmark"
)

const (
	VERTEX_RADIUS = 5
)

type Point struct {
	X float64
	Y float64
}

func pointFromTheta(theta float64) Point {
	return Point{X: math.Cos(theta), Y: math.Sin(theta)}
}

func (p *Point) Magnitude() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p *Point) scale(r float64) {
	p.X *= r
	p.Y *= r
}

func subtractPoint(p1, p2 Point) Point {
	return Point{p1.X - p2.X, p1.Y - p2.Y}
}

func (p *Point) Distance(point Point) float64 {
	distanceVector := subtractPoint(*p, point)

	return distanceVector.Magnitude()
}

type Action any

type AddVertex struct {
	connected bool
}

type RecomputeVertexPositions struct{}

type Draw struct{}

type MouseMove struct {
	pos Point
}

// type AddEdge struct {
// 	U
// }

func testGraphing() {
	graph := gograph.New[int](gograph.Acyclic())

	graph.AddEdge(gograph.NewVertex(1), gograph.NewVertex(2))
	graph.AddEdge(gograph.NewVertex(2), gograph.NewVertex(3))
	_, err := graph.AddEdge(gograph.NewVertex(3), gograph.NewVertex(1))

	if err != nil {
		fmt.Println("could not add edge")
	} else {
		fmt.Println("could add the edge")
	}
}

func main() {
	fmt.Println("Hello World!")
	// testGraphing()

	vecty.SetTitle("Markdown Demo")
	vecty.AddStylesheet("style.css")
	vecty.RenderBody(&PageView{})
}

// PageView is our main page component.
type PageView struct {
	vecty.Core
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {

	graph := initGraphCanvas(100, "main-canvas")

	return elem.Body(

		elem.Div(
			elem.Header(
				elem.Heading1(vecty.Text("Graph Theory Visualizations")),
				elem.Heading4(vecty.Text("This is a subtitle I'm sure I wont forget about")),
			),
		),

		elem.Div(
			vecty.Markup(
				vecty.Class("content-box"),
			),

			elem.Div(
				vecty.Markup(
					vecty.Class("button-box"),
				),

				&Button{text: "testing", onClick: func(e *vecty.Event) { fmt.Println("Clicked!") }},
				&Button{text: "add node", onClick: func(e *vecty.Event) { graph.Actions <- &AddVertex{connected: true} }},
			),
			&graph,
		),
	)
}

type Bijection interface {
	Forwards(from uint) uint
	Backwards(to uint) uint

	Add(from, to uint)
	Remove(from, to uint)

	Len() int
	Labels() []uint
}

type PosIntBijection struct {
	forwards []uint
	// backwards map[uint]uint
}

// Labels implements Bijection.
func (p *PosIntBijection) Labels() []uint {
	labels := make([]uint, p.Len())

	for i := range labels {
		labels[i] = uint(i)
	}

	return labels
}

// Forwards implements Bijection.
func (p *PosIntBijection) Forwards(from uint) uint {

	// gograph.Graph
	return p.forwards[from]
}

// Backwards implements Bijection.
func (p *PosIntBijection) Backwards(to uint) uint {

	for index, value := range p.forwards {
		if value == to {
			return uint(index)
		}
	}

	return 0
	// return p.backwards[to]
}

// Remove implements Bijection.
func (p *PosIntBijection) Remove(from uint, to uint) {
	// delete(p.backwards, to)
	p.forwards = append(p.forwards[0:from], p.forwards[from+1:]...)

}

// Set implements Bijection.
func (p *PosIntBijection) Add(from, to uint) {
	p.forwards = append(p.forwards, to)

	// panic("unimplemented")
}

func (p *PosIntBijection) Len() int {
	// length := len(p.backwards)

	return len(p.forwards)
}

func NewBijection() Bijection {
	return &PosIntBijection{forwards: make([]uint, 0)}

}

type Transform struct {
	scaleX, scaleY, shiftX, shiftY float64
}

func NewTransform(scaleX, scaleY, shiftX, shiftY float64) Transform {
	return Transform{scaleX: scaleX, scaleY: scaleY, shiftX: shiftX, shiftY: shiftY}
}

func (t *Transform) Forwards(point Point) Point {
	return Point{X: point.X*t.scaleX + t.shiftX, Y: point.Y*t.scaleY + t.shiftY}
}

func (t *Transform) Backwards(point Point) Point {
	return Point{X: (point.X - t.shiftX) / t.scaleX, Y: (point.Y - t.shiftY) / t.scaleY}
}

type GraphCanvas struct {
	vecty.Core
	ctx     js.Value
	Id      string
	Size    uint
	Actions chan any
	Graph   gograph.Graph[uint]

	bijection       Bijection
	nextLabel       uint
	vertexPositions map[uint]Point
	highlitedVertex uint

	transform Transform
}

func initGraphCanvas(size uint, id string) GraphCanvas {

	graph := GraphCanvas{
		Size:            size,
		Id:              id,
		Graph:           gograph.New[uint](),
		Actions:         make(chan any, 10),
		bijection:       NewBijection(),
		nextLabel:       0,
		vertexPositions: make(map[uint]Point)}

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
						point := Point{X: e.Get("offsetX").Float(), Y: e.Get("offsetY").Float()}
						point.scale(2)
						point = g.transform.Backwards(point)

						select {
						case g.Actions <- &MouseMove{
							pos: point,
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

	g.transform = NewTransform(scale, -scale, shift, shift)
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

func (g *GraphCanvas) DrawNodeOutline(point Point) {

	g.ctx.Call("beginPath")
	g.ctx.Call("arc", point.X, point.Y, VERTEX_RADIUS, 0, 2*math.Pi)
	g.ctx.Call("stroke")
}

func (g *GraphCanvas) DrawNode(point Point) {

	// g.ctx.Set("strokeStyle", "blue")

	g.ctx.Call("beginPath")
	g.ctx.Call("arc", point.X, point.Y, VERTEX_RADIUS, 0, 2*math.Pi)

	g.ctx.Call("fill")
	g.ctx.Call("stroke")
}

func (g *GraphCanvas) DrawEdge(from Point, to Point) {

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

func (g *GraphCanvas) HighlightActiveVertex(mousePos Point) *Action {

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

func (g *GraphCanvas) addNextVertex(connected bool) Action {
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

	return &RecomputeVertexPositions{}
}

func (g *GraphCanvas) RecomputeVertexPositions() Action {
	order := g.Graph.Order()
	// verticies := g.Graph.GetAllVertices()
	labels := g.bijection.Labels()
	fmt.Println(labels)
	radius := 30.0

	g.vertexPositions = make(map[uint]Point)

	if order > 1 {
		// vertexPositions := make(map[*gograph.Vertex[uint]]Point, 0)

		deltaTheta := 2 * math.Pi / float64(order)

		for i, label := range labels {
			point := pointFromTheta(deltaTheta * float64(i))
			point.scale(radius)

			v := g.bijection.Forwards(label)

			g.vertexPositions[v] = point

		}

	} else if order == 1 {
		g.vertexPositions[labels[0]] = Point{X: 0, Y: 0}
	}

	fmt.Println("recomputed vertex positions")

	return &Draw{}

}

func (g *GraphCanvas) handleActions() {
	for {
		a := <-g.Actions
		for a != nil {

			switch action := a.(type) {

			case *AddVertex:
				a = g.addNextVertex(action.connected)
				continue

			case *RecomputeVertexPositions:
				a = g.RecomputeVertexPositions()
				continue

			case *Draw:
				g.Draw()

			case *MouseMove:
				// fmt.Println(action.pos)
				g.HighlightActiveVertex(action.pos)

			default:
				// fmt.Fprintln(os.Stderr, "Unhandled action of type: ", action)
				fmt.Println("Unhandled action of type: ", action)
			}

			a = nil
		}
	}
}

type Graphs struct {
	// list of graphs
	graphs   vecty.List
	selected uint
}

func (g *Graphs) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("graph-list"),
		),
		g.graphs,
	)
}

func (_ Graphs) New(amount int) Graphs {
	graphs := make(vecty.List, 0, amount)
	for range amount {
		graphs = append(graphs, &GraphCanvas{})
	}

	return Graphs{graphs: graphs, selected: 0}
}

type Button struct {
	vecty.Core
	text    string
	onClick func(*vecty.Event)
}

func (b *Button) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Button(
			vecty.Text(b.text),
			vecty.Markup(
				event.Click(b.onClick),
			),
		),
	)
}

// Markdown is a simple component which renders the Input markdown as sanitized
// HTML into a div.
type Markdown struct {
	vecty.Core
	Input string `vecty:"prop"`
}

// Render implements the vecty.Component interface.
func (m *Markdown) Render() vecty.ComponentOrHTML {
	// Render the markdown input into HTML using Goldmark.
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(m.Input), &buf); err != nil {
		panic(err)
	}
	// The goldmark README says:
	// "By default, goldmark does not render raw HTML or potentially dangerous links. "
	// So, it should be ok without sanitizing.

	// Return the HTML.
	return elem.Div(
		vecty.Markup(
			vecty.UnsafeHTML(buf.String()),
		),
	)
}
