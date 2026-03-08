package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

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
