//go:build go1.23
// +build go1.23

package bitset

import "fmt"

// ExampleBitSet_EachSet iterates over the set bits using a range-over-function
// loop. This is the most concise way to visit every set bit on Go 1.23+.
func ExampleBitSet_EachSet() {
	var b BitSet
	b.Set(1).Set(2).Set(4).Set(8)

	for i := range b.EachSet() {
		fmt.Println(i)
	}
	// Output:
	// 1
	// 2
	// 4
	// 8
}
