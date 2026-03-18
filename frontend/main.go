package main

import (
	"fmt"

	// "os"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	// "github.com/hexops/vecty/event"
	// "github.com/hexops/vecty/prop"

	// "github.com/hexops/vecty/style"

	"github.com/dq1Mango/graph-theory/actions"
	"github.com/dq1Mango/graph-theory/components"
	"github.com/dq1Mango/graph-theory/model"
)

func main() {
	fmt.Println("Hello World!")
	// testGraphing()

	vecty.SetTitle("Graph Theory")
	vecty.AddStylesheet("style.css")
	vecty.AddStylesheet("colors.css")
	vecty.RenderBody(NewPageView())
}

func NewPageView() *PageView {
	graph := components.InitGraphCanvas(100, "main-canvas")

	return &PageView{
		graph: &graph,
	}
}

// PageView is our main page component.
type PageView struct {
	vecty.Core

	theme model.Theme
	graph *components.GraphCanvas
	// algorithm *components.AlgorithmWalk
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {

	graph := p.graph

	graphStats := components.NewGraphStats(graph)

	popup := components.NewPopup()

	// fmt.Println(popup)

	return elem.Body(

		elem.Div(
			elem.Header(
				elem.Heading1(vecty.Text("Graph Theory Visualizations")),
				elem.Heading4(vecty.Text("This is a subtitle I'm sure I wont forget about")),
			),
		),

		components.NewThemeChanger(),

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
					Text: "add vertex",
					OnClick: func(e *vecty.Event) {
						graph.Actions <- &actions.AddVertex{Connected: false}
					},
				},
				&components.Button{
					Text: "remove vertex",
					OnClick: func(e *vecty.Event) {
						graph.Actions <- &actions.RemoveVertex{Id: graph.SelectedVertex}
					},
				},
				&components.Button{
					Text: "add edge",
					OnClick: func(e *vecty.Event) {
						graph.Actions <- &actions.AddEdge{Vertex1: graph.SelectedVertex, Vertex2: graph.ShiftSelected}
					},
				},
				&components.Button{
					Text: "remove edge",
					OnClick: func(*vecty.Event) {
						graph.Actions <- &actions.RemoveEdge{Vertex1: graph.SelectedVertex, Vertex2: graph.ShiftSelected}
					},
				},
				&components.Button{
					Text: "generate prufer code",
					OnClick: func(*vecty.Event) {
						graph.Actions <- &actions.SetIterator{Iterator: "pruferCode"}
					},
				},

				components.NewGraphFromPrufer(
					graph.Actions,
				),
				components.NewDropdown("Bijection", func(bi string) {
					graph.Actions <- &actions.SetLabeling{Labeling: bi}
				}, "continous", "static"),
			),
			graph,
			graphStats,
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
				p.graph.Actions <- &actions.Draw{}
			}
		}
	}()
}
