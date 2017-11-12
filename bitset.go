// Copyright 2014 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package bitset implements bitsets, a mapping
between non-negative integers and boolean values. It should be more
efficient than map[uint] bool.

It provides methods for setting, clearing and testing
individual integers.
*/
package bitset

import (
	"bytes"
	"encoding/binary"
	"io"
)

const (
	wordSize    = uint(64)
	logWordSize = uint(6) // lg(wordSize)
	wordMask    = wordSize - 1
)

var (
	endianness = binary.LittleEndian
)

// BitSet is a bitset set membership, it is safe for concurrent reads
// but not for concurrent read/writes.
type BitSet struct {
	values []uint64
}

// NewBitSet returns a new bitset that can represent a set of n elements
// exactly. It is safe for concurrent reads but not for concurrent read/writes.
func NewBitSet(n uint) *BitSet {
	return &BitSet{values: make([]uint64, bitSetIndexOf(n)+1)}
}

// Test returns true if i is within the bit set or false otherwise.
func (b *BitSet) Test(i uint) bool {
	idx := bitSetIndexOf(i)
	if idx >= len(b.values) {
		return false
	}
	return b.values[idx]&(1<<(i&wordMask)) != 0
}

// Set adds i to the membership of the set.
func (b *BitSet) Set(i uint) {
	idx := bitSetIndexOf(i)
	currLen := len(b.values)
	if idx >= currLen {
		newValues := make([]uint64, 2*(idx+1))
		copy(newValues, b.values)
		b.values = newValues
	}
	b.values[idx] |= 1 << (i & wordMask)
}

// Clear removes i from the membership of the set.
func (b *BitSet) Clear(i uint) {
	idx := bitSetIndexOf(i)
	currLen := len(b.values)
	if idx >= currLen {
		return
	}
	b.values[idx] &^= 1 << (i & (wordSize - 1))
}

// ClearAll clears the set.
func (b *BitSet) ClearAll() {
	for i := range b.values {
		b.values[i] = 0
	}
}

// Write writes the bitset values to a stream.
func (b *BitSet) Write(w io.Writer) error {
	return binary.Write(w, endianness, b.values)
}

// ReadOnly returns a read only bit set created by writing
// the bit set to a buffer.
func (b *BitSet) ReadOnly() (*ReadOnlyBitSet, error) {
	buf := bytes.NewBuffer(make([]byte, 8*len(b.values)))
	buf.Reset()
	if err := b.Write(buf); err != nil {
		return nil, err
	}
	return NewReadOnlyBitSet(buf.Bytes()), nil
}

// ReadOnlyBitSet is a read only bitset set membership, it is safe for
// concurrent reads but not for concurrent read/writes.
type ReadOnlyBitSet struct {
	data []byte
}

// NewReadOnlyBitSet returns a new read only bit set backed
// by a byte slice, this means it can be used with a mmap'd bytes ref.
// It is safe for concurrent reads but not for concurrent read/writes.
func NewReadOnlyBitSet(data []byte) *ReadOnlyBitSet {
	return &ReadOnlyBitSet{data: data}
}

// Test returns true if i is within the bit set or false otherwise.
func (b *ReadOnlyBitSet) Test(i uint) bool {
	idx := bitSetIndexOf(i)
	values := len(b.data) / 8
	if idx >= values {
		return false
	}
	value := endianness.Uint64(b.data[idx*8 : (idx*8)+8])
	return value&(1<<(i&wordMask)) != 0
}

// Write writes the bitset values to a stream.
func (b *ReadOnlyBitSet) Write(w io.Writer) error {
	_, err := w.Write(b.data)
	return err
}

func bitSetIndexOf(i uint) int {
	return int(i >> logWordSize)
}
