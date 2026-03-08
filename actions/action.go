package actions

import (
	"github.com/dq1Mango/graph-theory/model"
)

type Action any

type AddVertex struct {
	Connected bool
}

type RecomputeVertexPositions struct{}

type Draw struct{}

type MouseMove struct {
	Pos model.Point
}
