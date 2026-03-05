package main

import (
	"bytes"
	"fmt"
	"math"

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

type Point struct {
	X float64
	Y float64
}

func pointFromTheta(theta float64) Point {
	return Point{X: math.Cos(theta), Y: math.Sin(theta)}
}

func (p *Point) scale(r float64) {
	p.X *= r
	p.Y *= r
}

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
	testGraphing()

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
				&Button{text: "add node", onClick: func(e *vecty.Event) { graph.addNextVertex() }},
			),
			&graph,
		),
	)
}

type GraphCanvas struct {
	vecty.Core
	ctx   js.Value
	Id    string
	Size  uint
	Graph gograph.Graph[uint]
}

func initGraphCanvas(size uint, id string) GraphCanvas {

	graph := GraphCanvas{Size: size, Id: id, Graph: gograph.New[uint]()}

	return graph
}

func (c *GraphCanvas) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("canvas-wrapper")),
		elem.Canvas(
			vecty.Markup(
				vecty.Class("canvas"),
				prop.ID(c.Id),
				// this should be in the style sheet but for some reason
				// style.Width(style.Size("100%")),
				// style.Color("blue"),
			),
			// vecty.Property("height", 20),
		),
	)
}

func (g *GraphCanvas) SetCanvasTransform() {

	// dpr := js.Global().Get("devicePixelRatio").Float()
	canvas := g.ctx.Get("canvas")

	width := canvas.Get("clientWidth").Float()
	height := canvas.Get("clientHeight").Float()

	canvas.Set("width", width)
	canvas.Set("height", height)

	fmt.Printf("width: %f, height: %f\n", width, height)

	if width != height {
		fmt.Printf("width and height of canvas: %s not equal\n", g.Id)
	}

	scale := math.Min(width, height) / (float64(g.Size))

	g.ctx.Call("setTransform", scale, 0, 0, -scale, width/2, height/2)

}

func (c *GraphCanvas) Mount() {
	canvas := js.Global().Get("document").Call("getElementById", c.Id)
	c.ctx = canvas.Call("getContext", "2d")

	c.SetCanvasTransform()

	// safe to draw here, DOM is ready
	c.ctx.Set("fillStyle", "red")
	c.ctx.Call("fillRect", 0, 0, 50, 25)
}

func (g *GraphCanvas) Clear() {
	s := int(g.Size)
	// fmt.Println(-g.Size / 2)
	// g.ctx.Call("clearRect", -g.Size/2, g.Size/2, g.Size, -g.Size)
	g.ctx.Call("clearRect", -s/2, s/2, s, -s)
}

func (g *GraphCanvas) DrawNode(point Point) {

	fmt.Println("drawing a node")
	radius := 5

	// g.ctx.Set("strokeStyle", "blue")
	g.ctx.Set("lineWidth", 1)
	g.ctx.Call("beginPath")
	g.ctx.Call("arc", point.X, point.Y, radius, 0, 2*math.Pi)
	g.ctx.Call("stroke")
}

func (g *GraphCanvas) Draw() {
	g.Clear()

	order := g.Graph.Order()
	radius := 30.0

	if order > 1 {
		deltaTheta := 2 * math.Pi / float64(order)

		for i := range order {
			point := pointFromTheta(deltaTheta * float64(i))
			point.scale(radius)

			g.DrawNode(point)
		}

	} else if order == 1 {
		g.DrawNode(Point{X: 0, Y: 0})
	}
}

func (c *GraphCanvas) addNextVertex() error {
	order := uint(c.Graph.Order())

	c.Graph.AddVertex(gograph.NewVertex(order))

	c.Draw()

	return nil
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
