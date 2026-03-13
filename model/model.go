package model

import (
	"math"
)

type Point struct {
	X float64
	Y float64
}

func PointFromTheta(theta float64) Point {
	return Point{X: math.Cos(theta), Y: math.Sin(theta)}
}

func (p *Point) Magnitude() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p *Point) Scale(r float64) {
	p.X *= r
	p.Y *= r
}

func SubtractPoint(p1, p2 Point) Point {
	return Point{p1.X - p2.X, p1.Y - p2.Y}
}

func (p *Point) Distance(point Point) float64 {
	distanceVector := SubtractPoint(*p, point)

	return distanceVector.Magnitude()
}

type Transform struct {
	scaleX, scaleY, shiftX, shiftY float64
}

func NewTransform(scaleX, scaleY, shiftX, shiftY float64) Transform {
	return Transform{scaleX: scaleX, scaleY: scaleY, shiftX: shiftX, shiftY: shiftY}
}

func (t *Transform) Forwards(point Point) Point {
	return Point{X: point.X*t.scaleX + t.shiftX, Y: point.Y*t.scaleY + t.shiftY}
}

func (t *Transform) Backwards(point Point) Point {
	return Point{X: (point.X - t.shiftX) / t.scaleX, Y: (point.Y - t.shiftY) / t.scaleY}
}
