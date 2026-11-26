package bitset

import "testing"

// FromWithLength documents that the backing slice may be longer than the
// number of words needed for the length. The set operations below must only
// touch the words covered by the bitset length and must not panic.

func TestComplementWithPaddedWords(t *testing.T) {
	b := FromWithLength(10, []uint64{0b0000000101, 0xffffffffffffffff})

	got := b.Complement()

	if got.Len() != 10 {
		t.Fatalf("Complement length = %d, want 10", got.Len())
	}
	for i := uint(0); i < 10; i++ {
		want := i != 0 && i != 2
		if got.Test(i) != want {
			t.Fatalf("Complement bit %d = %v, want %v", i, got.Test(i), want)
		}
	}
	if got.Test(10) || got.Test(63) {
		t.Fatalf("Complement set a bit at or beyond the length")
	}
}

func TestIntersectionWithPaddedWords(t *testing.T) {
	b := FromWithLength(10, []uint64{0b0000000101, 0xffffffffffffffff})
	var compare BitSet
	compare.Set(0).Set(1).Set(70)

	got := b.Intersection(&compare)

	if got.Len() != 10 {
		t.Fatalf("Intersection length = %d, want 10", got.Len())
	}
	if got.Count() != 1 || !got.Test(0) {
		t.Fatalf("Intersection = %s, want only bit 0", got)
	}
}

func TestUnionWithPaddedWords(t *testing.T) {
	b := FromWithLength(10, []uint64{0b0000000101, 0xffffffffffffffff})
	var compare BitSet
	compare.Set(15)

	got := b.Union(&compare)

	if got.Len() != 16 {
		t.Fatalf("Union length = %d, want 16", got.Len())
	}
	if got.Count() != 3 || !got.Test(0) || !got.Test(2) || !got.Test(15) {
		t.Fatalf("Union = %s, want bits 0, 2 and 15", got)
	}
}

func TestSymmetricDifferenceWithPaddedWords(t *testing.T) {
	b := FromWithLength(10, []uint64{0b0000000101, 0xffffffffffffffff})
	var compare BitSet
	compare.Set(0).Set(15)

	got := b.SymmetricDifference(&compare)

	if got.Len() != 16 {
		t.Fatalf("SymmetricDifference length = %d, want 16", got.Len())
	}
	if got.Count() != 2 || !got.Test(2) || !got.Test(15) || got.Test(0) {
		t.Fatalf("SymmetricDifference = %s, want bits 2 and 15 only", got)
	}
}
