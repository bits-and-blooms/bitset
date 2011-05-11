// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package bitset implements bitsets.

	It provides methods for making a BitSet of an arbitrary
	upper limit, setting and testing bit locations, and clearing
	bit locations as well as the entire set.

	Example use:

	b := bitset.New(64000)
	b.SetBit(1000)
	if b.Bit(1000) {
		b.ClearBit(1000)
	}
*/
package bitset

// BitSet internal details 
type BitSet struct {
	max_size uint
	set      []uint64
}

// Make a BitSet with an upper limit on size.
func New(max_size uint) *BitSet {
	return &BitSet{max_size, make([]uint64, (max_size+(64-1))/64)}
}

// Query maximum size of a bit set
func (b *BitSet) MaxSize() uint {
	return b.max_size
}

/// Test whether bit i is set. 
func (b *BitSet) Bit(i uint) bool {
	if b != nil && i < b.max_size {
		return ((b.set[i/64] & (1 << (i % 64))) != 0)
	}
	return false
}

// Set bit i to 1
func (b *BitSet) SetBit(i uint) {
	if b != nil && i < b.max_size {
		b.set[i/64] |= (1 << (i % 64))
	}
}

// Clear bit i to 0
func (b *BitSet) ClearBit(i uint) {
	if b != nil && i < b.max_size {
		b.set[i/64] &^= 1 << (i % 64)
	}
}

// Clear entire BitSet
func (b *BitSet) Clear() {
	if b != nil {
		for i, _ := range b.set {
			b.set[i] = 0
		}
	}
}
