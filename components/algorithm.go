package components

import (
	// "errors"
	"fmt"
	"slices"

	"github.com/dq1Mango/gograph"
	"github.com/dq1Mango/gograph/util"
	"github.com/dq1Mango/graph-theory/actions"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type IteratorEnd struct{}

func (i *IteratorEnd) Error() string {
	return "Reached the end"
}

type GraphIterator interface {
	Yield() bool
	UnYield() bool

	Next() error
	Previous() error
	Iterate() error

	Render() vecty.ComponentOrHTML
}

type AlgorithmWalk struct {
	vecty.Core
	iterator *GraphIterator
}

func (a *AlgorithmWalk) Render() vecty.ComponentOrHTML {

	if a.iterator != nil {
		iterator := (*a.iterator)
		return elem.Div(
			iterator.Render(),
			&Button{Text: "Step", OnClick: func(e *vecty.Event) {
				err := iterator.Next()
				if err != nil {
					InfoChan <- err.Error()
				}
				vecty.Rerender(a)
			}},
			&Button{Text: "Iterate", OnClick: func(e *vecty.Event) { iterator.Iterate() }},
		)
	} else {
		return nil
	}
}

type PruferCodeIterator struct {
	graph *GraphCanvas

	prufer       []uint
	smallestLeaf *int
	neighbor     *uint

	verticies []*gograph.Vertex[uint]
	length    int
}

func NewPruferCodeIterator(g *GraphCanvas) (GraphIterator, error) {
	graph := g.Graph
	length := graph.Order() - 2

	// ensure our graph is *probably* a tree
	if graph.Size() < 2 {
		return nil, &util.NonTreeError{}
	}

	if graph.Size() != graph.Order()-1 {
		return nil, &util.NonTreeError{}
	}

	verticies := graph.GetAllVertices()

	// sort all of the verticies
	slices.SortFunc(verticies,
		func(u, v *gograph.Vertex[uint]) int {
			if u.Label() > v.Label() {
				return 1
			} else if u.Label() == v.Label() {
				return 0
			} else {
				return -1
			}
		},
	)

	return &PruferCodeIterator{
		graph: g,
		// super big capacity optimization here
		prufer:    make([]uint, 0, length),
		verticies: verticies,
		length:    int(length),
	}, nil
}

func (p *PruferCodeIterator) Yield() bool {
	if len(p.prufer) < p.length {
		return true
	} else {
		return false
	}

}
func (p *PruferCodeIterator) UnYield() bool {
	if len(p.prufer) > 0 {
		return true
	} else {
		return false
	}
}

func (p *PruferCodeIterator) Next() error {
	if !p.Yield() {
		return &IteratorEnd{}
	}

	// graph := p.graph.Graph

	for index, vertex := range p.verticies {

		// we shall see this counts as a 'leaf' for a directed graph
		if vertex.InDegree() == 1 {

			bruh := int(vertex.Label())
			p.smallestLeaf = &bruh

			neighbor := vertex.Neighbors()[0].Label()
			p.neighbor = &neighbor

			p.prufer = append(p.prufer, neighbor)

			cantOneline := (p.verticies[index].Label())
			p.graph.RemoveVertex(&cantOneline)
			p.verticies = slices.Delete(p.verticies, index, index+1)

			// It should be noted i do not like this, but maybe it will be fine anyway
			p.graph.Actions <- &actions.Draw{}

			return nil

		}
	}

	// if we fail to find a leaf our graph is not a tree
	return &util.NonTreeError{}

}
func (p *PruferCodeIterator) Previous() error { return nil }
func (p *PruferCodeIterator) Iterate() error {
	for p.Yield() {
		err := p.Next()

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PruferCodeIterator) Render() vecty.ComponentOrHTML {

	var leafText *vecty.HTML
	var neighborText *vecty.HTML

	if p.smallestLeaf != nil {
		leafText = elem.Paragraph(vecty.Text(fmt.Sprintf("Smallest Leaf: %d", *p.smallestLeaf)))
	}

	if p.neighbor != nil {
		neighborText = elem.Paragraph(vecty.Text(fmt.Sprintf("Leaf's Neighbor: %d", *p.neighbor)))
	}

	return elem.Div(
		vecty.Markup(
			vecty.Class("algorithm"),
		),

		elem.Heading3(vecty.Text("Prufer Code Generation")),

		elem.Paragraph(vecty.Text(fmt.Sprintf("Current Code: %v", p.prufer))),

		leafText,
		neighborText,
	)
}
