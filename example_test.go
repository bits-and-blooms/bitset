// This file contains runnable, self-verifying examples that double as
// documentation on https://pkg.go.dev. Every example asserts its output via
// an "// Output:" comment, so `go test` fails if the behavior ever changes.

package bitset

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Example gives a tour of the most common operations: creating a set, setting
// and testing bits, chaining, and querying the population count and length.
func Example() {
	// The zero value of a BitSet is an empty set that is ready to use.
	var b BitSet

	// Set returns the BitSet, so calls can be chained.
	b.Set(10).Set(11).Set(1000)

	fmt.Println("bit 10 set:", b.Test(10))
	fmt.Println("bit 500 set:", b.Test(500))
	fmt.Println("count:", b.Count()) // number of set bits
	fmt.Println("len:", b.Len())     // one past the highest set bit

	// Output:
	// bit 10 set: true
	// bit 500 set: false
	// count: 3
	// len: 1001
}

// ExampleNew shows creating a BitSet with a size hint. The hint avoids
// re-allocations; the set still grows automatically past the hint if needed.
func ExampleNew() {
	b := New(64) // hint: expect to use ~64 bits
	b.Set(1).Set(2).Set(3)
	fmt.Println(b.String())
	// Output: {1,2,3}
}

// ExampleFrom builds a BitSet directly from packed 64-bit words, which is handy
// when interoperating with other bitmap representations.
func ExampleFrom() {
	// In word 0, bits 0 and 2 are set: 0b101 == 5.
	b := From([]uint64{0b101})
	fmt.Println(b.String())
	// Output: {0,2}
}

// ExampleBitSet_Set demonstrates chaining, which works because Set (and Clear,
// Flip, ...) return the receiver.
func ExampleBitSet_Set() {
	b := New(100)
	b.Set(0).Set(1).Set(50)
	fmt.Println(b.String())
	// Output: {0,1,50}
}

// ExampleBitSet_Clear turns an individual bit off.
func ExampleBitSet_Clear() {
	b := New(10).Set(1).Set(2).Set(3)
	b.Clear(2)
	fmt.Println(b.String())
	// Output: {1,3}
}

// ExampleBitSet_Flip toggles a single bit.
func ExampleBitSet_Flip() {
	var b BitSet
	b.Set(0).Set(1)
	b.Flip(1) // bit 1 was set, so it is cleared
	b.Flip(2) // bit 2 was clear, so it is set
	fmt.Println(b.String())
	// Output: {0,2}
}

// ExampleBitSet_FlipRange toggles every bit in the half-open range [start, end).
func ExampleBitSet_FlipRange() {
	var b BitSet
	b.FlipRange(0, 5) // sets bits 0,1,2,3,4
	fmt.Println(b.String())
	// Output: {0,1,2,3,4}
}

// ExampleBitSet_Count reports the number of set bits (the cardinality).
func ExampleBitSet_Count() {
	b := New(100).Set(3).Set(30).Set(90)
	fmt.Println(b.Count())
	// Output: 3
}

// ExampleBitSet_NextSet is the idiomatic way to iterate over the set bits in
// ascending order on any Go version.
func ExampleBitSet_NextSet() {
	b := New(100).Set(0).Set(2).Set(4).Set(80)

	var set []uint
	for i, ok := b.NextSet(0); ok; i, ok = b.NextSet(i + 1) {
		set = append(set, i)
	}

	fmt.Println(set)
	// Output: [0 2 4 80]
}

// ExampleBitSet_NextSetMany retrieves many set bits at once into a reusable
// buffer, which is useful for batching in hot loops.
func ExampleBitSet_NextSetMany() {
	b := New(100).Set(1).Set(2).Set(3).Set(70).Set(99)

	buf := make([]uint, 4) // fetch at most 4 bits per call
	_, set := b.NextSetMany(0, buf)

	fmt.Println(set)
	// Output: [1 2 3 70]
}

// ExampleBitSet_NextClear finds the first clear (unset) bit at or after i.
func ExampleBitSet_NextClear() {
	b := New(10).Set(0).Set(1).Set(2)

	i, ok := b.NextClear(0)
	fmt.Println(i, ok)
	// Output: 3 true
}

// ExampleBitSet_Union returns a new set containing every bit set in either set.
func ExampleBitSet_Union() {
	a := New(10).Set(1).Set(2).Set(3)
	b := New(10).Set(3).Set(4).Set(5)
	fmt.Println(a.Union(b).String())
	// Output: {1,2,3,4,5}
}

// ExampleBitSet_Intersection returns a new set containing the bits set in both.
func ExampleBitSet_Intersection() {
	a := New(10).Set(1).Set(2).Set(3)
	b := New(10).Set(3).Set(4).Set(5)
	fmt.Println(a.Intersection(b).String())
	// Output: {3}
}

// ExampleBitSet_Difference returns the bits set in a but not in b.
func ExampleBitSet_Difference() {
	a := New(10).Set(1).Set(2).Set(3)
	b := New(10).Set(3).Set(4).Set(5)
	fmt.Println(a.Difference(b).String())
	// Output: {1,2}
}

// ExampleBitSet_SymmetricDifference returns the bits set in exactly one of the
// two sets.
func ExampleBitSet_SymmetricDifference() {
	a := New(10).Set(1).Set(2).Set(3)
	b := New(10).Set(3).Set(4).Set(5)
	fmt.Println(a.SymmetricDifference(b).String())
	// Output: {1,2,4,5}
}

// ExampleBitSet_Complement flips every bit within the set's length.
func ExampleBitSet_Complement() {
	b := New(8).Set(1).Set(3)
	fmt.Println(b.Complement().String())
	// Output: {0,2,4,5,6,7}
}

// ExampleBitSet_InPlaceUnion mutates the receiver instead of allocating a new
// set. The In-Place variants (and the *Cardinality helpers) avoid allocation in
// performance-sensitive code.
func ExampleBitSet_InPlaceUnion() {
	a := New(10).Set(1).Set(2)
	b := New(10).Set(2).Set(3)
	a.InPlaceUnion(b)
	fmt.Println(a.String())
	// Output: {1,2,3}
}

// ExampleBitSet_Any shows the emptiness predicates. All, Any and None answer
// "are all/any/no bits set?".
func ExampleBitSet_Any() {
	var b BitSet
	fmt.Println("any:", b.Any())   // empty set
	fmt.Println("none:", b.None()) // empty set
	b.Set(5)
	fmt.Println("any after Set:", b.Any())
	// Output:
	// any: false
	// none: true
	// any after Set: true
}

// ExampleBitSet_IsSuperSet reports whether every bit of other is also set here.
func ExampleBitSet_IsSuperSet() {
	super := New(10).Set(1).Set(2).Set(3)
	sub := New(10).Set(1).Set(3)
	fmt.Println(super.IsSuperSet(sub))
	// Output: true
}

// ExampleBitSet_String renders the set in mathematical set notation.
func ExampleBitSet_String() {
	b := New(50).Set(1).Set(17).Set(49)
	fmt.Println(b.String())
	// Output: {1,17,49}
}

// ExampleBitSet_Rank counts how many bits are set at or below the given index.
func ExampleBitSet_Rank() {
	b := New(10).Set(1).Set(3).Set(5)
	fmt.Println(b.Rank(0)) // no bit <= 0 is set
	fmt.Println(b.Rank(3)) // bits 1 and 3
	fmt.Println(b.Rank(9)) // bits 1, 3 and 5
	// Output:
	// 0
	// 2
	// 3
}

// ExampleBitSet_Select returns the index of the j-th set bit (0-indexed).
func ExampleBitSet_Select() {
	b := New(10).Set(1).Set(3).Set(5)
	fmt.Println(b.Select(0)) // first set bit
	fmt.Println(b.Select(1)) // second set bit
	fmt.Println(b.Select(2)) // third set bit
	// Output:
	// 1
	// 3
	// 5
}

// ExampleBitSet_WriteTo serializes a BitSet to a stream and reads it back with
// ReadFrom. The encoding is portable across platforms.
func ExampleBitSet_WriteTo() {
	b := New(100).Set(10).Set(20).Set(99)

	var buf bytes.Buffer

	_, err := b.WriteTo(&buf)
	if err != nil {
		panic(err)
	}

	// ReadFrom reuses the receiver's storage where possible.
	restored := New(0)

	_, err = restored.ReadFrom(&buf)
	if err != nil {
		panic(err)
	}

	fmt.Println("equal:", restored.Equal(b))
	fmt.Println("bits:", restored.String())
	// Output:
	// equal: true
	// bits: {10,20,99}
}

// ExampleBitSet_MarshalJSON round-trips a BitSet through JSON. BitSet
// implements json.Marshaler and json.Unmarshaler.
func ExampleBitSet_MarshalJSON() {
	b := New(20).Set(0).Set(19)

	data, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	var restored BitSet

	err = json.Unmarshal(data, &restored)
	if err != nil {
		panic(err)
	}

	fmt.Println(restored.String())
	// Output: {0,19}
}
