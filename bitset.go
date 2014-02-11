/// Copyright 2014 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

	Package bitset implements bitsets, a mapping
	between non-negative integers and boolean values. It should be more
	efficient than map[uint] bool.

	It provides methods for setting, clearing, flipping, and testing
	individual integers.

	But it also provides set intersection, union, difference,
	complement, and symmetric operations, as well as tests to
	check whether any, all, or no bits are set, and querying a
	bitset's current length and number of postive bits.

	BitSets are expanded to the size of the largest set bit; the
	memory allocation is approximately Max bits, where Max is
	the largest set bit. BitSets are never shrunk. On creation,
	a hint can be given for the number of bits that will be used.

    Many of the methods, including Set,Clear, and Flip, return
	a BitSet pointer, which allows for chaining.

	Example use:

	import "bitset"
	var b BitSet
	b.Set(10).Set(11)
	if b.Test(1000) {
		b.Clear(1000)
	}
	if B.Intersection(bitset.New(100).Set(10)).Count() > 1 {
		fmt.Println("Intersection works.")
	}

	As an alternative to BitSets, one should check out the 'big' package,
	which provides a (less set-theoretical) view of bitsets.

*/
package bitset

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
)

// Word size of a bit set
const wordSize = uint(64)

// Mask for cleaning last word
const allBits uint64 = 0xffffffffffffffff

// for laster arith.
const log2WordSize = uint(6)

// The zero value of a BitSet is an empty set of length 0.
type BitSet struct {
	length uint
	set    []uint64
}

type BitSetError string

// fixup b.set to be non-nil and return the field value
func (b *BitSet) safeSet() []uint64 {
	if b.set == nil {
		b.set = make([]uint64, wordsNeeded(0))
	}
	return b.set
}

func wordsNeeded(i uint) uint {
	if i > (math.MaxUint64 - wordSize + 1 ) { // safer?
	// if i == math.MaxUint64 {
		return math.MaxUint64 >> log2WordSize
	} else if i == 0 {
		return 1
	}
	return (i + (wordSize - 1)) >> log2WordSize
}

func New(length uint) *BitSet {
	return &BitSet{length, make([]uint64, wordsNeeded(length))}
}

func (b *BitSet) Cap() uint {
	return uint(math.MaxUint64)
}

func (b *BitSet) Len() uint {
	return b.length
}

//
func (b *BitSet) extendSetMaybe(i uint) {
	if i >= b.length { // if we need more bits, make 'em
		nsize := wordsNeeded(i + 1)
		if b.set == nil {
			b.set = make([]uint64, nsize)
		} else if uint(len(b.set)) < nsize {
			newset := make([]uint64, nsize)
			copy(newset, b.set)
			b.set = newset
		}
		b.length = i + 1
	}
}

/// Test whether bit i is set.
func (b *BitSet) Test(i uint) bool {
	if i >= b.length {
		return false
	}
	return b.set[i>>log2WordSize]&(1<<(i&(wordSize-1))) != 0
}

// Set bit i to 1
func (b *BitSet) Set(i uint) *BitSet {
	b.extendSetMaybe(i)
	//fmt.Printf("length in bits: %d, real size of sets: %d, bits: %d, index: %d\n", b.length, len(b.set), i, i>>log2WordSize)
	b.set[i>>log2WordSize] |= 1 << (i & (wordSize - 1))
	return b
}

// Clear bit i to 0
func (b *BitSet) Clear(i uint) *BitSet {
	if i >= b.length {
		return b
	}
	b.set[i>>log2WordSize] &^= 1 << (i & (wordSize - 1))
	return b
}

// Set bit i to value
func (b *BitSet) SetTo(i uint, value bool) *BitSet {
	if value {
		return b.Set(i)
	}
	return b.Clear(i)
}

// Flip bit at i
func (b *BitSet) Flip(i uint) *BitSet {
	if i >= b.length {
		return b.Set(i)
	}
	b.set[i>>log2WordSize] ^= 1 << (i & (wordSize - 1))
	return b
}

// return the next bit set from the specified index, including possibly the current index
// returns -1 if none is found
// inspired by the Java API: for i:=int64(0); i>=0; i = NextSet(i) {...}
func (b *BitSet) NextSet(i uint) (uint,bool) {
	x := i >> log2WordSize
	if x >= b.length {
		return 0, false
	}
	w := b.set[x]
	w = w >> (i & (wordSize - 1))
	if w != 0 {
		return i + trailingZeroes64(w),true
	}
	x = x + 1
	for x < wordsNeeded(b.length) {
		if b.set[x] != 0 {
			return x * wordSize + trailingZeroes64(b.set[x]),true
		}
		x = x + 1

	}
	return 0, false
}

// Clear entire BitSet
func (b *BitSet) ClearAll() *BitSet {
	if b != nil && b.set != nil {
		for i := range b.set {
			b.set[i] = 0
		}
	}
	return b
}

// Query words used in a bit set
func (b *BitSet) wordCount() uint {
	return wordsNeeded(b.length)
}

// Clone this BitSet
func (b *BitSet) Clone() *BitSet {
	c := New(b.length)
	copy(c.set, b.safeSet())
	return c
}

// Copy this BitSet into a destination BitSet
// Returning the size of the destination BitSet
// like array copy
func (b *BitSet) Copy(c *BitSet) (count uint) {
	if c == nil {
		return
	}
	copy(c.set, b.safeSet())
	count = c.length
	if b.length < c.length {
		count = b.length
	}
	return
}

// From Wikipedia: http://en.wikipedia.org/wiki/Hamming_weight
const m1 uint64 = 0x5555555555555555  //binary: 0101...
const m2 uint64 = 0x3333333333333333  //binary: 00110011..
const m4 uint64 = 0x0f0f0f0f0f0f0f0f  //binary:  4 zeros,  4 ones ...
const m8 uint64 = 0x00ff00ff00ff00ff  //binary:  8 zeros,  8 ones ...
const m16 uint64 = 0x0000ffff0000ffff //binary: 16 zeros, 16 ones ...
const m32 uint64 = 0x00000000ffffffff //binary: 32 zeros, 32 ones
const hff uint64 = 0xffffffffffffffff //binary: all ones
const h01 uint64 = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...

// From Wikipedia: count number of set bits.
// This is algorithm popcount_2 in the article retrieved May 9, 2011

func popcount_2(x uint64) uint64 {
	x -= (x >> 1) & m1             //put count of each 2 bits into those 2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each 4 bits into those 4 bits
	x = (x + (x >> 4)) & m4        //put count of each 8 bits into those 8 bits
	x += x >> 8                    //put count of each 16 bits into their lowest 8 bits
	x += x >> 16                   //put count of each 32 bits into their lowest 8 bits
	x += x >> 32                   //put count of each 64 bits into their lowest 8 bits
	return x & 0x7f
}

// Count (number of set bits)
func (b *BitSet) Count() uint {
	if b != nil && b.set != nil {
		cnt := uint64(0)
		for _, word := range b.set {
			cnt += popcount_2(word)
		}
		return uint(cnt)
	}
	return 0
}

// computes the number of trailing zeroes on the assumption that v is non-zero
func trailingZeroes64(v uint64) uint {
	// NOTE: if 0 == v, then c = 63.
	if v&0x1 != 0 {
		return 0
	}
	c := uint(1)
	if (v & 0xffffffff) == 0 {
		v >>= 32
		c += 32
	}
	if (v & 0xffff) == 0 {
		v >>= 16
		c += 16
	}
	if (v & 0xff) == 0 {
		v >>= 8
		c += 8
	}
	if (v & 0xf) == 0 {
		v >>= 4
		c += 4
	}
	if (v & 0x3) == 0 {
		v >>= 2
		c += 2
	}
	c -= uint(v & 0x1)
	return c
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
	for p, v := range b.safeSet() {
		if c.set[p] != v {
			return false
		}
	}
	return true
}

func panicIfNull(b *BitSet) {
	if b == nil {
		panic(BitSetError("BitSet must not be null"))
	}
}

// Difference of base set and other set
// This is the BitSet equivalent of &^ (and not)
func (b *BitSet) Difference(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	result = b.Clone() // clone b (in case b is bigger than compare)
	szl := compare.wordCount()
	for i, word := range b.safeSet() {
		if uint(i) >= szl {
			break
		}
		result.set[i] = word &^ compare.set[i]
	}
	return
}

// Convenience function: return two bitsets ordered by
// increasing length. Note: neither can be nil
func sortByLength(a *BitSet, b *BitSet) (ap *BitSet, bp *BitSet) {
	if a.length <= b.length {
		ap, bp = a, b
	} else {
		ap, bp = b, a
	}
	return
}

// Intersection of base set and other set
// This is the BitSet equivalent of & (and)
func (b *BitSet) Intersection(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	result = New(b.length)
	for i, word := range b.safeSet() {
		result.set[i] = word & compare.set[i]
	}
	return
}

// Union of base set and other set
// This is the BitSet equivalent of | (or)
func (b *BitSet) Union(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	result = compare.Clone()
	szl := compare.wordCount()
	for i, word := range b.safeSet() {
		if uint(i) >= szl {
			break
		}
		result.set[i] = word | compare.set[i]
	}
	return
}

// SymmetricDifference of base set and other set
// This is the BitSet equivalent of ^ (xor)
func (b *BitSet) SymmetricDifference(compare *BitSet) (result *BitSet) {
	panicIfNull(b)
	panicIfNull(compare)
	b, compare = sortByLength(b, compare)
	// compare is bigger, so clone it
	result = compare.Clone()
	szl := b.wordCount()
	for i, word := range b.safeSet() {
		if uint(i) >= szl {
			break
		}
		result.set[i] = word ^ compare.set[i]
	}
	return
}

// Is the length an exact multiple of word sizes?
func (b *BitSet) isEven() bool {
	return b.length%wordSize == 0
}

// Clean last word by setting unused bits to 0
func (b *BitSet) cleanLastWord() {
	if !b.isEven() {
		b.set[wordsNeeded(b.length)-1] &= allBits >> (wordSize - b.length%wordSize)
	}
}

// Return the (local) Complement of a biset (up to length bits)
func (b *BitSet) Complement() (result *BitSet) {
	panicIfNull(b)
	result = New(b.length)
	for i, word := range b.safeSet() {
		result.set[i] = ^word
	}
	result.cleanLastWord()
	return
}

// Returns true if all bits are set, false otherwise
func (b *BitSet) All() bool {
	panicIfNull(b)
	return b.Count() == b.length
}

// Return true if no bit is set, false otherwise
func (b *BitSet) None() bool {
	panicIfNull(b)
	if b != nil && b.set != nil {
		for _, word := range b.set {
			if word > 0 {
				return false
			}
		}
		return true
	}
	return true
}

// Return true if any bit is set, false otherwise
func (b *BitSet) Any() bool {
	panicIfNull(b)
	return !b.None()
}

// Dump as bits
func (b *BitSet) DumpAsBits() string {
	buffer := bytes.NewBufferString("")
	b.safeSet()
	i := int(wordsNeeded(b.length) - 1)
	for ; i >= 0; i-- {
		fmt.Fprintf(buffer, "%064b.", b.set[i])
	}
	return string(buffer.Bytes())
}

func (b *BitSet) MarshalJSON() ([]byte, error) {
	// Put the bitset length in front of the string
	length := uint64(b.length)
	dataCap := binary.Size(length) + binary.Size(b.set)
	buffer := bytes.NewBuffer(make([]byte, 0, dataCap))

	// Write length
	err := binary.Write(buffer, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}

	// Write set
	err = binary.Write(buffer, binary.BigEndian, b.set)
	if err != nil {
		return nil, err
	}

	// URLEncode all bytes
	return json.Marshal(base64.URLEncoding.EncodeToString(buffer.Bytes()))
}

func (b *BitSet) UnmarshalJSON(data []byte) error {
	// Unmarshal as string
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	// URLDecode string
	buf, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buf)
	newset := New(0)
	var length uint64

	// Read length first
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	newset.length = uint(length)

	// Read remaining bytes as set
	err = binary.Read(reader, binary.BigEndian, newset.set)
	if err != nil {
		return err
	}

	*b = *newset
	return nil
}
