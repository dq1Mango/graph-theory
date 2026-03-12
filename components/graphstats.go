package components

import (
	"fmt"

	"github.com/dq1Mango/gograph"
	path "github.com/dq1Mango/gograph/path"
	"github.com/dq1Mango/graph-theory/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/style"
)

type GraphStats struct {
	vecty.Core

	Graph    *GraphCanvas
	Expanded bool
}

func NewGraphStats(graph *GraphCanvas) *GraphStats {
	stats := GraphStats{Graph: graph, Expanded: false}

	graph.Stats = &stats

	return &stats
}

func (s *GraphStats) Render() vecty.ComponentOrHTML {

	stats := s.SimpleStats()

	return elem.Div(
		vecty.Markup(
			vecty.Class("stats"),
		),

		elem.Paragraph(
			vecty.Markup(
				vecty.Class("stats-title"),
			),
			vecty.Text("Graph Stats"),
		),

		elem.Div(
			stats...,
		),

		s.PrimaryVertexStats(),
		s.SecondaryVertexStats(),
	)
}

func (g *GraphStats) SimpleStats() []vecty.MarkupOrChild {

	// mhhh great names
	graph := g.Graph.Graph

	stats := make([]string, 2)

	stats[0] = fmt.Sprintf("Order: %d", graph.Order())
	stats[1] = fmt.Sprintf("Size: %d", graph.Size())

	paragraphs := make([]vecty.MarkupOrChild, 0, 2)

	for _, stat := range stats {
		paragraphs = append(paragraphs, elem.Paragraph(vecty.Text(stat)))
	}

	return paragraphs
}

func (g *GraphStats) PrimaryVertexStats() *vecty.HTML {
	graph := g.Graph

	// sure thing guys lets just call everything a 'graph'
	if graph.SelectedVertex != nil &&
		graph.Graph.ContainsVertex(gograph.NewVertex(*graph.SelectedVertex)) {
		return elem.Div(
			elem.Break(),

			elem.Paragraph(
				vecty.Markup(
					vecty.Class("stats-title"),
					style.Color(model.CurrentPalette.Red),
				),
				vecty.Text("Primary Vertex"),
			),

			elem.Div(
				g.VertexStats(*graph.SelectedVertex),
			),
		)
	} else {
		return nil
	}

}
func (g *GraphStats) SecondaryVertexStats() *vecty.HTML {
	graph := g.Graph

	if graph.ShiftSelected != nil &&
		graph.Graph.ContainsVertex(gograph.NewVertex(*graph.ShiftSelected)) {
		return elem.Div(
			elem.Break(),

			elem.Paragraph(
				vecty.Markup(
					vecty.Class("stats-title"),
					style.Color(model.CurrentPalette.Blue),
				),
				vecty.Text("Secondary Vertex"),
			),

			elem.Div(
				g.VertexStats(*graph.ShiftSelected),
			),
		)
	} else {
		return nil
	}
}

func (g *GraphStats) VertexStats(id uint) *vecty.HTML {
	graph := g.Graph.Graph
	vertex := graph.GetVertexByID(id)

	stats := []string{
		fmt.Sprintf("Degree: %d", vertex.OutDegree()),
		fmt.Sprintf("Ecentricity: %v", path.Eccentricity(graph, vertex)),
	}

	paragraphs := make([]vecty.MarkupOrChild, 0, 2)

	for _, stat := range stats {
		paragraphs = append(paragraphs, elem.Paragraph(vecty.Text(stat)))
	}

	return elem.Div(
		paragraphs...,
	)
}
