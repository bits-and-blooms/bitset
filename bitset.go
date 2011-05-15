// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package bitset implements bitsets.

	It provides methods for making a BitSet of an arbitrary
	upper limit, setting and testing bit locations, and clearing
	bit locations as well as the entire set, counting the number
	of bits.
	
	It also supports set intersection, union, difference and 
	symmetric difference, cloning, equality testing, and subsetting.

	Example use:
    
	import "bitset"
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

type BitSetError string

// Make a BitSet with an upper limit on capacity.
func New(length uint) *BitSet {
	return &BitSet{length, make([]uint64, (length+(64-1))>>6)}
}


// Query maximum size of a bit set
func (b *BitSet) Cap() uint {
	return b.length
}

// Word size of a bit set
const WordSize = uint(64)

// Query words used in a bit set
func (b *BitSet) WordCount() uint {
	return b.length / WordSize
}

// Is the length an exact multiple of word sizes?
func (b *BitSet) isEven() bool {
   return (b.length % WordSize) == 0
}

// Clone this BitSet
func (b *BitSet) Clone() *BitSet {
	c := New(b.length)
	copy(c.set, b.set)
	return c
}
// Copy this BitSet into a destination BitSet
// Returning the size of the destination BitSet
// like array copy
func (b *BitSet) Copy(c *BitSet) (count uint) {
	if c == nil {
		return
	}
	copy(c.set, b.set)
	count = c.length
	if b.length < c.length {
		count = b.length
	}
	return
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

// Test the equvalence of two BitSets. 
// False if they are of different sizes, otherwise true
// only if all the same bits are set
func (b *BitSet) Equal(c *BitSet) bool {
	if c == nil {
		return false
	}
	if b.length != c.length {
		return false
	}
	for p, v := range b.set { 
		// Note: this assumes that the extra bits will be
		// the same. This is not _quite_ guaranteed except
		// by convention. Is it worth special casing the
		// last word?
		if c.set[p] != v {
			return false
		}
	}
	return true
}

func checkBitSetsForNull(b1 *BitSet, b2 *BitSet) (e *BitSetError) {
	if b1==nil || b2==nil { 
		err := BitSetError("Null pointer")
		return &err
	}
	return nil  
}                 
// call this only with non-nil parameters. 
// ie, call checkBitSetsForNull first
func swapToSmaller(b1 *BitSet, b2 *BitSet) (b3 *BitSet, b4 *BitSet) {
	b3,b4 = b1,b2
	if b1.length > b2.length {
		b3,b4 = b2,b1
	}
	return  
}

// Difference of base set and other set
// This is the BitSet equivalent of &^ (and not)
// Sets can be of different capacities, but neither can be nil
func (b1 *BitSet) Difference(b2 *BitSet) (b3 *BitSet, err *BitSetError) {
	err = checkBitSetsForNull(b1, b2)
	if err != nil {
		return nil, err
	}
	b3 = b1.Clone()
	szl := b2.WordCount()
	for i ,word := range b1.set {
		if uint(i) >= szl {
			break
		}
		b3.set[i] = word &^ b2.set[i]
	}
	return
}

// Union of base set and other set
// This is the BitSet equivalent of | (or)
// Sets can be of different capacities, but neither can be nil
func (b1 *BitSet) Union(b2 *BitSet) (b3 *BitSet, err *BitSetError) {
	err = checkBitSetsForNull(b1, b2)
	if err != nil {
		return nil, err
	}
	b1,b2 = swapToSmaller(b1,b2)
	// b2 is bigger, so clone it
	b3 = b2.Clone()
	szl := b1.WordCount()
	for i ,word := range b1.set {
		if uint(i) >= szl {
			break
		}
		b3.set[i] = word | b2.set[i]
	}
	return
}

// Intersection of base set and other set
// This is the BitSet equivalent of & (and)
// Sets can be of different capacities, but neither can be nil
func (b1 *BitSet) Intersection(b2 *BitSet) (b3 *BitSet, err *BitSetError) {
	err = checkBitSetsForNull(b1, b2)
	if err != nil {
		return nil, err
	}
	b1,b2 = swapToSmaller(b1,b2)
	// b1 is smaller; use its size: larger bits will be clear
	b3 = New(b1.Cap())
	for i ,word := range b1.set {
		b3.set[i] = word & b2.set[i]
	}
	return
}

// SymmetricDifference of base set and other set
// This is the BitSet equivalent of ^ (xor)
// Sets can be of different capacities, but neither can be nil
func (b1 *BitSet) SymmetricDifference(b2 *BitSet) (b3 *BitSet, err *BitSetError) {
	err = checkBitSetsForNull(b1, b2)
	if err != nil {
		return nil, err
	}
	b1,b2 = swapToSmaller(b1,b2)
	// b2 is bigger, so clone it
	b3 = b2.Clone()
	szl := b1.WordCount()
	for i ,word := range b1.set {
		if uint(i) >= szl {
			break
		}
		b3.set[i] = word ^ b2.set[i]
	}
	return
}

// Copy, from start to end, a subset of bits from a set
// returns the copy and a possible error
func (b *BitSet) Subset(start, end uint) (c *BitSet, err *BitSetError) {
	if (end - start) < 0 {
		c = nil
		e := BitSetError(fmt.Sprintf("Resulting BitSet would be negative in length: %d", end-start-1))
		return c, &e
	}
	if end > b.length {
		c = nil
		e := BitSetError(fmt.Sprintf("End index %d exceeds length %d", c, b.length))
		return c, &e
	}
	c = New(end - start)
	if start&(64-1) == 0 {
		copy(c.set, b.set[start>>6:(end+63)>>6])
		return c, nil
	}
	ipos := start & (64 - 1)
	ifirst := start >> 6
	ilen := (end - start + 64 - 1) >> 6
	var i uint
	for ; i < ilen; i++ {
		c.set[i] = b.set[i+ifirst] >> ipos
		c.set[i] |= b.set[i+ifirst+1] << (64 - ipos)
	}
	return c, nil
}