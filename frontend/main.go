package main

import (
	"bytes"
	"fmt"

	// "os"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	// "github.com/hexops/vecty/event"
	// "github.com/hexops/vecty/prop"

	// "github.com/hexops/vecty/style"
	"github.com/yuin/goldmark"

	"github.com/dq1Mango/graph-theory/actions"
	"github.com/dq1Mango/graph-theory/components"
	"github.com/dq1Mango/graph-theory/model"
)

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

	theme model.Theme
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {

	graph := components.InitGraphCanvas(100, "main-canvas")

	popup := components.NewPopup()

	themeDropDown := components.NewDropdown("theme", "dark", "catppuccin")
	// fmt.Println(popup)

	return elem.Body(

		elem.Div(
			elem.Header(
				elem.Heading1(vecty.Text("Graph Theory Visualizations")),
				elem.Heading4(vecty.Text("This is a subtitle I'm sure I wont forget about")),
			),
		),

		themeDropDown,

		&popup,

		elem.Div(
			vecty.Markup(
				vecty.Class("content-box"),
			),

			elem.Div(
				vecty.Markup(
					vecty.Class("button-box"),
				),

				&components.Button{
					Text:    "testing",
					OnClick: func(e *vecty.Event) { fmt.Println("Clicked!") },
				},
				&components.Button{
					Text:    "add node",
					OnClick: func(e *vecty.Event) { graph.Actions <- &actions.AddVertex{Connected: true} },
				},
				&components.Button{
					Text:    "remove vertex",
					OnClick: func(e *vecty.Event) { graph.Actions <- &actions.RemoveVertex{Id: graph.SelectedVertex} },
				},
				&components.Button{
					Text: "add edge",
					OnClick: func(e *vecty.Event) {
						graph.Actions <- &actions.AddEdge{Vertex1: graph.SelectedVertex, Vertex2: graph.ShiftSelected}
					}},
				&components.Button{Text: "remove edge", OnClick: func(*vecty.Event) {
					graph.Actions <- &actions.RemoveEdge{Vertex1: graph.SelectedVertex, Vertex2: graph.ShiftSelected}
				}},
			),
			&graph,
		),
	)
}

func (p *PageView) Mount() {
	p.theme.SetTheme("catppuccin")

	// model.GreenFlag = true

	go func() {
		for {
			select {
			case newTheme := <-model.ThemChan:
				p.theme.SetTheme(newTheme)
			}
		}
	}()
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
