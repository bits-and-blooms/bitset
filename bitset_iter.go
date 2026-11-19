//go:build go1.23
// +build go1.23

package bitset

import (
	"iter"
	"math/bits"
)

// EachSet returns an iterator over the indices of all set bits in ascending
// order, for use with Go 1.23+ range-over-function loops:
//
//	for i := range b.EachSet() {
//		// i is the index of a set bit
//	}
//
// Breaking out of the loop stops the iteration early.
func (b *BitSet) EachSet() iter.Seq[uint] {
	return func(yield func(uint) bool) {
		for wordIndex, word := range b.set {
			for word != 0 {
				if !yield(uint(wordIndex<<log2WordSize + bits.TrailingZeros64(word))) {
					return
				}
				word &= word - 1
			}
		}
	}
}
