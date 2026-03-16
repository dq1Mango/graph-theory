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

			vecty.Markup(
				vecty.Class("algorithm"),
			),
			iterator.Render(),
			elem.Div(
				vecty.Markup(
					vecty.Class("algorithm-options"),
				),
				&Button{Text: "Step", OnClick: func(e *vecty.Event) {
					err := iterator.Next()
					if err != nil {
						InfoChan <- err.Error()
					}
					vecty.Rerender(a)
				}},
				&Button{Text: "Iterate", OnClick: func(e *vecty.Event) {
					err := iterator.Iterate()
					if err != nil {
						InfoChan <- err.Error()
					}
					vecty.Rerender(a)
				}},
			),
		)
	} else {
		return nil
	}
}

type PruferCodeIterator struct {
	graph *GraphCanvas

	prufer       []uint
	smallestLeaf *uint
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
			uLabel, vLabel := g.Bijection.Forwards(u.Label()), g.Bijection.Forwards(v.Label())
			if uLabel > vLabel {
				return 1
			} else if uLabel == vLabel {
				return 0
			} else {
				return -1
			}
		},
	)

	iterator := &PruferCodeIterator{
		graph: g,
		// super big capacity optimization here
		prufer:    make([]uint, 0, length),
		verticies: verticies,
		length:    int(length),
	}

	err := iterator.Next()

	return iterator, err
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

	if p.smallestLeaf != nil && p.neighbor != nil {
		p.graph.RemoveVertex(p.smallestLeaf)
		p.prufer = append(p.prufer, *p.neighbor)
	}

	// graph := p.graph.Graph

	for index, vertex := range p.verticies {

		// we shall see if this counts as a 'leaf' for a directed graph
		if vertex.InDegree() == 1 {

			// what phenomenal programming
			bruh := uint(vertex.Label())
			p.smallestLeaf = &bruh

			neighbor := vertex.Neighbors()[0].Label()
			p.neighbor = &neighbor

			// p.prufer = append(p.prufer, neighbor)

			// cantOneline := (p.verticies[index].Label())
			// p.graph.RemoveVertex(&cantOneline)
			p.verticies = slices.Delete(p.verticies, index, index+1)

			p.graph.SelectedVertex = &bruh
			p.graph.ShiftSelected = &neighbor

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

	return &IteratorEnd{}
}

func (p *PruferCodeIterator) Render() vecty.ComponentOrHTML {

	leafText := "Smallest Leaf: "
	neighborText := "Leaf's Neigbor: "

	if p.smallestLeaf != nil {
		leafText += fmt.Sprintf("%d", *p.smallestLeaf)
	}

	if p.neighbor != nil {
		neighborText += fmt.Sprintf("%d", *p.neighbor)
	}

	return elem.Div(

		elem.Heading3(
			vecty.Markup(vecty.Class("centered")),
			vecty.Text("Prufer Code Generation"),
		),

		elem.Paragraph(vecty.Text(fmt.Sprintf("Current Code: %v", p.prufer))),

		elem.Paragraph(
			vecty.Markup(vecty.Class("red")),
			vecty.Text(leafText),
		),
		elem.Paragraph(
			vecty.Markup(vecty.Class("blue")),
			vecty.Text(neighborText),
		),
	)
}
