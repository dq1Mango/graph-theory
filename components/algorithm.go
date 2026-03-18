package components

import (
	// "errors"
	"fmt"
	// "slices"

	"github.com/dq1Mango/graph-theory/actions"
	"github.com/dq1Mango/graph-theory/model"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hmdsefi/gograph"
	"github.com/hmdsefi/gograph/util"
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
				&Button{Text: "UnStep", OnClick: func(e *vecty.Event) {
					err := iterator.Previous()
					if err != nil {
						InfoChan <- err.Error()
					}
					vecty.Rerender(a)
				}},
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

type VertexWithPoint struct {
	vertex *gograph.Vertex[uint]
	point  model.Point
}

type PruferCodeIterator struct {
	graph *GraphCanvas

	prufer       []uint
	smallestLeaf *uint
	neighbor     *uint

	verticies        []*gograph.Vertex[uint]
	removedVerticies []VertexWithPoint
	length           int
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
	SortVerticies(verticies)

	iterator := &PruferCodeIterator{
		graph: g,
		// super big capacity optimization here
		prufer:           make([]uint, 0, length),
		verticies:        verticies,
		removedVerticies: make([]VertexWithPoint, 0, length),
		length:           int(length),
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
		p.removedVerticies = append(
			p.removedVerticies,
			VertexWithPoint{
				vertex: p.verticies[*p.smallestLeaf],
				point:  p.graph.VertexPositions[*p.smallestLeaf],
			},
		)
		p.graph.RemoveVertex(p.smallestLeaf)

		// p.verticies = slices.Delete(p.verticies, index, index+1)
		p.verticies[*p.smallestLeaf] = nil

		p.prufer = append(p.prufer, *p.neighbor)
	}

	// graph := p.graph.Graph

	for _, vertex := range p.verticies {

		if vertex == nil {
			continue
		}

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
func (p *PruferCodeIterator) Previous() error {
	if !p.UnYield() {
		return &IteratorEnd{}
	}

	end := len(p.removedVerticies) - 1
	lastRemoved := p.removedVerticies[end]
	removedLabel := lastRemoved.vertex.Label()
	// lastId := lastRemoved.Label()

	p.graph.SelectedVertex = &removedLabel
	p.smallestLeaf = &removedLabel

	p.graph.ShiftSelected = &p.prufer[end]
	p.neighbor = &p.prufer[end]

	// pos := p.graph.VertexPositions[lastRemoved.vertex.Label()]
	p.graph.AddVertex(removedLabel, &lastRemoved.point, false)
	p.graph.AddEdge(&removedLabel, &p.prufer[end])

	p.verticies[removedLabel] = p.graph.Graph.GetVertexByID(removedLabel)
	fmt.Println(p.verticies[removedLabel])

	p.prufer = p.prufer[:end]
	p.removedVerticies = p.removedVerticies[:end]

	p.graph.Actions <- &actions.Draw{}

	return nil

	// p.graph.Bijection.Set(lastRemoved.Label(), i)

	// panic("should't be possible")
}
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
