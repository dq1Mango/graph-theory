package model

import "slices"

// A bijection maps from vertex "Ids" to vetex "Labels"
// Ids are whats used in the internal graph representation
// Labels are what the user will see the verticies *labeled* as

type Bijection interface {
	Forwards(from uint) uint
	Backwards(to uint) uint

	Add(from ...uint)
	AddNext()
	Remove(from, to uint)
	Clear()

	Len() int
	Labels() []uint
	Labeling() map[uint]uint

	NextId() uint
	NextLabel() uint

	// a list of Ids sorted by their labels
	SortedIds() []uint
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
	labels := make([]uint, 0, p.Len())

	for _, v := range p.forwards {
		labels = append(labels, v)
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
	return uint(len(p.forwards))
}

// Forwards implements Bijection.
func (p *ContinousBijection) Forwards(from uint) uint {
	return p.forwards[from]
}

// Backwards implements Bijection.
func (p *ContinousBijection) Backwards(to uint) uint {

	for index, value := range p.forwards {
		if value == to {
			return uint(index)
		}
	}

	panic("sorry man, couldnt find ur id")
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
func (p *ContinousBijection) add(from uint) {
	p.forwards[from] = p.NextLabel()

	if from >= p.nextId {
		p.nextId = from + 1
	}
}

func (p *ContinousBijection) Add(from ...uint) {
	for _, value := range from {
		p.add(value)
	}
}

func (p *ContinousBijection) AddNext() {
	p.forwards[p.nextId] = p.NextLabel()
	p.nextId++
}

func (p *ContinousBijection) Len() int {
	// length := len(p.backwards)

	return len(p.forwards)
}

func (p *ContinousBijection) SortedIds() []uint {
	sorted := make([]uint, p.Len())

	labels := p.Labels()

	slices.Sort(labels)

	for i, label := range labels {
		sorted[i] = p.Backwards(label)
	}

	return sorted

}

type FakeBijection struct {
	verticies []uint
	nextLabel uint
}

// Labeling implements Bijection.
func (f *FakeBijection) Labeling() map[uint]uint {
	labeling := make(map[uint]uint)
	for _, id := range f.verticies {
		labeling[id] = id
	}
	return labeling
}

func (f *FakeBijection) findNextId() {
	var i uint = 0

	for _, id := range f.verticies {
		if id > i {
			f.nextLabel = i
			return
		}
		i++
	}

	f.nextLabel = i
}

func (f *FakeBijection) insert(newId uint) {
	for index, id := range f.verticies {
		if id > newId {
			f.verticies = slices.Insert(f.verticies, index, newId)
			f.findNextId()
			return
		}
	}
	f.verticies = append(f.verticies, newId)
	f.findNextId()

}

func (f *FakeBijection) add(id uint) {
	f.insert(id)
}

func (f *FakeBijection) Add(ids ...uint) {
	for _, id := range ids {
		f.add(id)
	}
}

// AddNext implements Bijection.
func (f *FakeBijection) AddNext() {
	f.add(f.nextLabel)
}

// Backwards implements Bijection.
func (f *FakeBijection) Backwards(to uint) uint {
	return to
}

// Clear implements Bijection.
func (f *FakeBijection) Clear() {
	f.verticies = make([]uint, 0)
	f.nextLabel = 0
}

// Forwards implements Bijection.
func (f *FakeBijection) Forwards(from uint) uint {
	return from
}

// Labels implements Bijection.
func (f *FakeBijection) Labels() []uint {
	return f.verticies
}

// Len implements Bijection.
func (f *FakeBijection) Len() int {
	return len(f.verticies)
}

// NextLabel implements Bijection.
func (f *FakeBijection) NextLabel() uint {
	return f.nextLabel
}
func (f *FakeBijection) NextId() uint {
	return f.nextLabel
}

// Remove implements Bijection.
func (f *FakeBijection) Remove(from uint, to uint) {

	// erm, did you know binary search runs in O(ln(n)) 🤓
	for index, value := range f.verticies {
		if value == from {
			f.verticies = slices.Delete(f.verticies, index, index+1)
			f.findNextId()
			return
		}
	}
	panic("ruh roh, vertex does not exist in bijection")
}

// SortedIds implements Bijection.
func (f *FakeBijection) SortedIds() []uint {
	result := make([]uint, f.Len())
	copy(result, f.verticies)
	return result
}

func NewFakeBijection() Bijection {
	return &FakeBijection{
		nextLabel: 0,
	}
}

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
