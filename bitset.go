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
	b.SetBit(999)
	if b.Bit(1000) {
		b.ClearBit(1000)
	}
	b.Clear()
	
*/
package bitset

import (
	"fmt"
)

// BitSet internal details 
type BitSet struct {
	length uint
	set      []uint64
}

// Make a BitSet with an upper limit on capacity.
func New(length uint) *BitSet {
	return &BitSet{length, make([]uint64, (length+(64-1))>>6)}
}


// Query maximum size of a bit set
func (b *BitSet) Cap() uint {
	return b.length
}

// Query length of a bit set (which is the same as its capacity,
// by analogy to len and cap functions on Arrays.
func (b *BitSet) Len() uint {
	return b.length
} 

/// Test whether bit i is set. 
func (b *BitSet) Bit(i uint) bool {
	if i >= b.length {
		panic(fmt.Sprintf("index out of range: %v", i))
	}
	return ((b.set[i>>6] & (1 << (i & (64-1)))) != 0)
}

// Set bit i to 1
func (b *BitSet) SetBit(i uint) {
	if i >= b.length {
		panic(fmt.Sprintf("index out of range: %v", i))
	}
	b.set[i>>6] |= (1 << (i & (64-1)))
}

// Clear bit i to 0
func (b *BitSet) ClearBit(i uint) {
	if i >= b.length {
		panic(fmt.Sprintf("index out of range: %v", i))
	}	
	b.set[i>>6] &^= 1 << (i & (64-1))
}

// Clear entire BitSet
func (b *BitSet) Clear() {
	if b != nil {
		for i := range b.set {
			b.set[i] = 0
		}
	}
}

// From Wikipedia: http://en.wikipedia.org/wiki/Hamming_weight                                     
const m1  uint64 = 0x5555555555555555 //binary: 0101...
const m2  uint64 = 0x3333333333333333 //binary: 00110011..
const m4  uint64 = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...

// From Wikipedia: count number of set bits. 
// This is algorithm popcount_2 in the article retrieved May 9, 2011
func popCountUint64(x uint64) uint64 {
    x -= (x >> 1) & m1;             //put count of each 2 bits into those 2 bits
    x = (x & m2) + ((x >> 2) & m2); //put count of each 4 bits into those 4 bits 
    x = (x + (x >> 4)) & m4;        //put count of each 8 bits into those 8 bits 
    x += x >>  8;  //put count of each 16 bits into their lowest 8 bits
    x += x >> 16;  //put count of each 32 bits into their lowest 8 bits
    x += x >> 32;  //put count of each 64 bits into their lowest 8 bits
    return x & 0x7f;
}

// Count (number of set bits)
func (b *BitSet) Count() uint {
   	if b != nil {
		cnt := uint64(0)
		for _, word := range b.set {
			cnt += popCountUint64(word)
		}
		return uint(cnt)
	}
	return 0
}

// Test two sets equivalence
func (b *BitSet) Equ(c *BitSet) bool {
	if c == nil {
		return false
	}
	if b.length != c.length {
		return false
	}
	for p, v := range b.set {
		if c.set[p] != v {
			return false
		}
	}
	return true
}



