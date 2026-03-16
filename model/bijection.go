package model

// A bijection maps from vertex "Ids" to vetex "Labels"
// Ids are whats used in the internal graph representation
// Labels are what the user will see the verticies *labeled* as

type Bijection interface {
	Forwards(from uint) uint
	Backwards(to uint) uint

	// Add(from, to uint)
	AddNext()
	Remove(from, to uint)
	Clear()

	Len() int
	Labels() []uint
	Labeling() map[uint]uint

	NextId() uint
	NextLabel() uint
}

type ContinousBijection struct {
	// forwards []uint
	forwards map[uint]uint
	nextId   uint
	// backwards map[uint]uint
}

func NewBijection() Bijection {
	return &ContinousBijection{forwards: make(map[uint]uint, 0)}
}

// Labels implements Bijection.
func (p *ContinousBijection) Labels() []uint {
	labels := make([]uint, p.Len())

	for i, v := range labels {
		labels[i] = v
	}

	return labels
}
func (p *ContinousBijection) Labeling() map[uint]uint {
	return p.forwards

}
func (p *ContinousBijection) NextId() uint {
	return p.nextId
}

func (p *ContinousBijection) NextLabel() uint {
	return uint(len(p.forwards)) + 1
}

// Forwards implements Bijection.
func (p *ContinousBijection) Forwards(from uint) uint {
	return p.forwards[from]
}

// Backwards implements Bijection.
func (p *ContinousBijection) Backwards(to uint) uint {

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
func (p *ContinousBijection) Remove(from uint, to uint) {
	// delete(p.backwards, to)
	// p.forwards = append(p.forwards[0:from], p.forwards[from+1:]...)
	delete(p.forwards, from)

	for k := range p.forwards {
		if k > from {
			p.forwards[k] -= 1
		}
	}
}

func (p *ContinousBijection) Clear() {
	for len(p.forwards) > 0 {
		p.Remove(0, p.forwards[0])
	}
}

// Set implements Bijection.
//
//	func (p *ContinousBijection) Add(from, to uint) {
//		p.forwards = append(p.forwards, to)
//	}
func (p *ContinousBijection) AddNext() {
	p.forwards[p.nextId] = p.NextLabel()
	p.nextId++
}

func (p *ContinousBijection) Len() int {
	// length := len(p.backwards)

	return len(p.forwards)
}

type FakeBijection struct {
	nextLabel uint
}

// AddNext implements Bijection.
func (f *FakeBijection) AddNext(to uint) {

	panic("unimplemented")
}

// Backwards implements Bijection.
func (f *FakeBijection) Backwards(to uint) uint {
	panic("unimplemented")
}

// Clear implements Bijection.
func (f *FakeBijection) Clear() {
	panic("unimplemented")
}

// Forwards implements Bijection.
func (f *FakeBijection) Forwards(from uint) uint {
	panic("unimplemented")
}

// Labels implements Bijection.
func (f *FakeBijection) Labels() []uint {
	panic("unimplemented")
}

// Len implements Bijection.
func (f *FakeBijection) Len() int {
	panic("unimplemented")
}

// NextLabel implements Bijection.
func (f *FakeBijection) NextLabel() uint {
	return f.nextLabel
}

// Remove implements Bijection.
func (f *FakeBijection) Remove(from uint, to uint) {
}

// func NewFakeBijection() Bijection {
// 	return &FakeBijection{
// 		nextLabel: 0,
// 	}
// }

// type NaturalBijection struct {
// 	forwards  map[uint]uint
// 	backwards map[uint]uint
// 	nextLabel uint
// }
//
// func NewNaturalBijection() Bijection {
// 	return &NaturalBijection{
// 		forwards:  map[uint]uint{},
// 		backwards: map[uint]uint{},
// 		nextLabel: 1,
// 	}
// }
//
// func (n *NaturalBijection) NextLabel() uint {
//
//
// 	return n.nextLabel
// }
//
// // AddNext implements Bijection.
// func (n *NaturalBijection) AddNext(to uint) {
//
// 	label := n.NextLabel()
//
// 	n.forwards[label] = to
// 	n.backwards[to] = label
//
// 	i := uint(1)
// 	for k := range n.forwards {
// 		if k != i && k != n.nextLabel {
// 			n.nextLabel = i
// 			break
// 		}
// 		i++
// 	}
//
// 	panic("unimplemented")
// }
//
// // Backwards implements Bijection.
// func (n *NaturalBijection) Backwards(to uint) uint {
// 	panic("unimplemented")
// }
//
// // Clear implements Bijection.
// func (n *NaturalBijection) Clear() {
// 	panic("unimplemented")
// }
//
// // Forwards implements Bijection.
// func (n *NaturalBijection) Forwards(from uint) uint {
// 	panic("unimplemented")
// }
//
// // Labels implements Bijection.
// func (n *NaturalBijection) Labels() []uint {
// 	panic("unimplemented")
// }
//
// // Len implements Bijection.
// func (n *NaturalBijection) Len() int {
// 	panic("unimplemented")
// }
//
// // Remove implements Bijection.
// func (n *NaturalBijection) Remove(from uint, to uint) {
// 	panic("unimplemented")
// }
