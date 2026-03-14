package model

type Bijection interface {
	Forwards(from uint) uint
	Backwards(to uint) uint

	Add(from, to uint)
	Remove(from, to uint)
	Clear()

	Len() int
	Labels() []uint
}

type PosIntBijection struct {
	forwards []uint
	// backwards map[uint]uint
}

func NewBijection() Bijection {
	return &PosIntBijection{forwards: make([]uint, 0)}

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

	return p.forwards[from]
}

// Backwards implements Bijection.
func (p *PosIntBijection) Backwards(to uint) uint {

	// erm, did you know binary search runs in O(ln(n)) 🤓
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

func (p *PosIntBijection) Clear() {
	for len(p.forwards) > 0 {
		p.Remove(0, p.forwards[0])
	}
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
