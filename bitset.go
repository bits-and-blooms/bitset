// Copyright 2011 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package bitset implements bitsets.
	
	It provides methods for making a BitSet of an arbitrary
	upper limit, setting and testing bit locations, and clearing
	bit locations as well as the entire set.
	
	BitSets are implemented as arrays of uint64s, so it may be best
	to limit the upper size to a multiple of 64. It is not an error
	to set, test, or clear a bit within a 64 bit boundary.
	
	Example use:
	
	b := MakeBitSet(64000)
	b.SetBit(1000)
	if b.Bit(1000) {
		b.ClearBit(1000)
	}
*/
package bitset

// for MaxUint64
import (
	"math"
)

// BitSets are arrays of uint64. 
type BitSet []uint64

// Make a BitSet with an upper limit on size. Note this is the
// number of bits, not the number of uint64s, which is a kind of
// implementation detail.
func MakeBitSet(max_size uint) BitSet {
	return make(BitSet, (max_size+(64-1))/64)
}

/// Test whether bit i is set. 
func (set BitSet) Bit(i uint) bool {  
	return ((set[i/64] & (1<<(i%64)))!=0)
}

// Set bit i to 1
func (set BitSet) SetBit(i uint) {
	set[i/64]|=(1<<(i%64))
}

// Clear bit i to 0
func (set BitSet) ClearBit( i uint) {
	set[i/64]&=(1<<(i%64))^math.MaxUint64
}

// Clear entire BitSet
func (set BitSet) Clear() {
	for i,_ := range set {
		set[i] = 0
	}
}
