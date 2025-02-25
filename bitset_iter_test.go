//go:build go1.23
// +build go1.23

package bitset

import (
	"testing"
)

func TestIter(t *testing.T) {
	var b BitSet
	b.Set(0).Set(3).Set(5).Set(6).Set(63).Set(64).Set(65).Set(127).Set(128).Set(1000000)

	// Expected values that should be set
	expected := []uint{0, 3, 5, 6, 63, 64, 65, 127, 128, 1000000}
	got := make([]uint, 0)

	// Collect all values from iterator
	for i := range b.EachSet() {
		got = append(got, i)
	}

	// Test 1: Check length matches expected
	if len(got) != len(expected) {
		t.Errorf("Expected %d elements, got %d", len(expected), len(got))
	}

	// Test 2: Check all expected values are present and in correct order
	for i, want := range expected {
		if i >= len(got) {
			t.Errorf("Missing expected value %d at position %d", want, i)
			continue
		}
		if got[i] != want {
			t.Errorf("At position %d: expected %d, got %d", i, want, got[i])
		}
	}

	// Test 3: Check no extra values
	if len(got) > len(expected) {
		t.Errorf("Got extra values: %v", got[len(expected):])
	}
}
func BenchmarkIter(b *testing.B) {
	b.StopTimer()
	s := New(10000)
	for i := 0; i < 10000; i += 3 {
		s.Set(uint(i))
	}

	b.StartTimer()
	for j := 0; j < b.N; j++ {
		c := uint(0)
		for range s.EachSet() {
			c++
		}
	}
}

func BenchmarkNonInter(b *testing.B) {
	b.StopTimer()
	s := New(10000)
	for i := 0; i < 10000; i += 3 {
		s.Set(uint(i))
	}

	b.StartTimer()
	for j := 0; j < b.N; j++ {
		c := uint(0)
		for i, e := s.NextSet(0); e; i, e = s.NextSet(i + 1) {
			c++
		}
	}
}
