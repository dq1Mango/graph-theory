package main

import (
	"bytes"
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"github.com/hexops/vecty/prop"

	// "github.com/hexops/vecty/event"
	// "github.com/hexops/vecty/prop"
	"syscall/js"

	"github.com/hexops/vecty/style"
	"github.com/hmdsefi/gograph"
	"github.com/yuin/goldmark"
)

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
	vecty.RenderBody(&PageView{
		Input: `# Markdown Example

This is a live editor, try editing the Markdown on the right of the page.
`,
	})
}

// PageView is our main page component.
type PageView struct {
	vecty.Core
	Input string
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {

	graph := GraphCanvas{Size: 100, Id: "main-canvas"}

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

func (c *GraphCanvas) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("canvas-wrapper")),
		elem.Canvas(
			vecty.Markup(
				vecty.Class("canvas"),
				prop.ID(c.Id),
				// this should be in the style sheet but for some reason
				// style.Width(style.Size("100%")),
				// vecty.Attribute("width", "500"),
				// vecty.Attribute("height", "500"),
				// style.Color("blue"),
			),
			// vecty.Property("height", 20),
		),
	)
}

func (c *GraphCanvas) SetCanvasTransform() {

	canvas := c.ctx.Get("canvas")

	width := canvas.Get("width").Float()
	height := canvas.Get("height").Float()

	scaleX, scaleY := width, height

	c.ctx.Call("setTransform", scaleX, 0, 0, -scaleY, width/2, height/2)

}

func (c *GraphCanvas) Mount() {
	canvas := js.Global().Get("document").Call("getElementById", c.Id)
	c.ctx = canvas.Call("getContext", "2d")

	c.SetCanvasTransform()

	// safe to draw here, DOM is ready
	c.ctx.Set("fillStyle", "red")
	c.ctx.Call("fillRect", 0, 0, 100, 100)
}

func (g *GraphCanvas) Clear() {
	g.ctx.Call("clearRect", 0, 0, 500, 500)
}

func (g *GraphCanvas) DrawNode() {}

func (c *GraphCanvas) Draw() {
	c.Clear()

	order := c.Graph.Order()

	if order == 0 {
		return
	} else if order == 1 {

	}

}

func (c *GraphCanvas) addNextVertex() error {
	order := uint(c.Graph.Order())

	c.Graph.AddVertex(gograph.NewVertex(order))

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
