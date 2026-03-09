package actions

import (
	"github.com/dq1Mango/graph-theory/model"
)

type Action any

type AddVertex struct {
	Connected bool
}

type RemoveVertex struct {
	Id *uint
}

type AddEdge struct {
	Vertex1, Vertex2 *uint
}

type RemoveEdge struct {
	Vertex1, Vertex2 *uint
}

type RecomputeVertexPositions struct{}

type Draw struct{}

type MouseMove struct {
	Pos model.Point
}

type MouseDown struct {
	Pos   model.Point
	Shift bool
}

type PopupMessage struct {
	Message string
}
