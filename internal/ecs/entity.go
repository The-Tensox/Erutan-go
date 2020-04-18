package ecs

import (
	"sync/atomic"
)

var (
	idInc uint64
)

// Identifier is an interface for anything that implements the basic ID() uint64,
// as the BasicEntity does.  It is useful as more specific interface for an
// entity registry than just the interface{} interface
type Identifier interface {
	ID() uint64
}

// IdentifierSlice implements the sort.Interface, so you can use the
// store entites in slices, and use the P=n*log n lookup for them
type IdentifierSlice []Identifier

func NewId() uint64 {
	return atomic.AddUint64(&idInc, 1)
}

// Len returns the length of the underlying slice
// part of the sort.Interface
func (is IdentifierSlice) Len() int {
	return len(is)
}

// Less will return true if the ID of element at i is less than j;
// part of the sort.Interface
func (is IdentifierSlice) Less(i, j int) bool {
	return is[i].ID() < is[j].ID()
}

// Swap the elements at positions i and j
// part of the sort.Interface
func (is IdentifierSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}
