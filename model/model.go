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

type Bijection interface {
	Forwards(from uint) uint
	Backwards(to uint) uint

	Add(from, to uint)
	Remove(from, to uint)

	Len() int
	Labels() []uint
}

type PosIntBijection struct {
	forwards []uint
	// backwards map[uint]uint
}

// Labels implements Bijection.
func (p *PosIntBijection) Labels() []uint {
	labels := make([]uint, p.Len())

	for i := range labels {
		labels[i] = uint(i)
	}

	return labels
}

// Forwards implements Bijection.
func (p *PosIntBijection) Forwards(from uint) uint {

	// gograph.Graph
	return p.forwards[from]
}

// Backwards implements Bijection.
func (p *PosIntBijection) Backwards(to uint) uint {

	// erm, did you know binary search run in O(ln(n)) 🤓
	for index, value := range p.forwards {
		if value == to {
			return uint(index)
		}
	}

	return 0
	// return p.backwards[to]
}

// Remove implements Bijection.
func (p *PosIntBijection) Remove(from uint, to uint) {
	// delete(p.backwards, to)
	p.forwards = append(p.forwards[0:from], p.forwards[from+1:]...)

}

// Set implements Bijection.
func (p *PosIntBijection) Add(from, to uint) {
	p.forwards = append(p.forwards, to)

	// panic("unimplemented")
}

func (p *PosIntBijection) Len() int {
	// length := len(p.backwards)

	return len(p.forwards)
}

func NewBijection() Bijection {
	return &PosIntBijection{forwards: make([]uint, 0)}

}
