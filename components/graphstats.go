package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hmdsefi/gograph"
)

func (g *GraphStats) SimpleStats() []vecty.MarkupOrChild {

	graph := g.Graph

	stats := make([]string, 2)

	stats[0] = fmt.Sprintf("Order: %d", graph.Order())
	stats[1] = fmt.Sprintf("Size: %d", graph.Size())

	paragraphs := make([]vecty.MarkupOrChild, 0, 2)

	for _, stat := range stats {
		paragraphs = append(paragraphs, elem.Paragraph(vecty.Text(stat)))
	}

	return paragraphs
}

type GraphStats struct {
	vecty.Core

	Graph    gograph.Graph[uint]
	Expanded bool
}

func (s *GraphStats) Render() vecty.ComponentOrHTML {

	stats := s.SimpleStats()

	return elem.Div(
		elem.Paragraph(vecty.Text("Graph Stats")),

		elem.Div(
			stats...,
		),
	)
}
