package main

import (
	"bytes"
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	// "github.com/hexops/vecty/event"
	// "github.com/hexops/vecty/prop"
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
	return elem.Body(
		// Display a textarea on the right-hand side of the page.
		// elem.Div(
		// vecty.Markup(
		// 	vecty.Style("float", "right"),
		// ),
		// elem.TextArea(
		// 	vecty.Markup(
		// 		vecty.Style("font-family", "monospace"),
		// 		vecty.Property("rows", 14),
		// 		vecty.Property("cols", 70),
		//
		// 		// When input is typed into the textarea, update the local
		// 		// component state and rerender.
		// 		event.Input(func(e *vecty.Event) {
		// 			p.Input = e.Target.Get("value").String()
		// 			vecty.Rerender(p)
		// 		}),
		// 	),
		// 	vecty.Text(p.Input), // initial textarea text.
		// ),
		// ),

		// Render the markdown.
		&Markdown{Input: p.Input},

		elem.Div(
			vecty.Markup(
				vecty.Class("content-box"),
			),
			elem.Div(
				&Button{text: "testing"},
			),
			&GraphCanvas{Size: style.Size("100%")},
		),
	)
}

type GraphCanvas struct {
	vecty.Core
	Size style.Size
}

func (c *GraphCanvas) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("canvas-wrapper")),
		elem.Canvas(
			vecty.Markup(
				vecty.Class("canvas"),
				// style.Height(c.Size),
				style.Width(c.Size),
				vecty.Attribute("width", "500"),
				vecty.Attribute("height", "500"),
				// style.Color("blue"),
			),
			// vecty.Property("height", 20),
		),
	)
}

type Button struct {
	vecty.Core
	text     string
	callback func(*vecty.Event)
}

func (b *Button) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Button(
			vecty.Text(b.text),
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
