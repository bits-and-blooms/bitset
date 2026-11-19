// This file tests the popcnt functions

package bitset

import (
	"testing"
)

func TestPopcntSlice(t *testing.T) {
	s := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	res := popcntSlice(s)
	const l uint64 = 27
	if res != l {
		t.Errorf("Wrong popcount %d != %d", res, l)
	}
}

func TestPopcntMaskSlice(t *testing.T) {
	s := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	m := []uint64{31, 37, 41, 43, 47, 53, 59, 61, 67, 71}
	res := popcntMaskSlice(s, m)
	const l uint64 = 9
	if res != l {
		t.Errorf("Wrong mask %d !=  %d", res, l)
	}
}

func TestPopcntAndSlice(t *testing.T) {
	s := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	m := []uint64{31, 37, 41, 43, 47, 53, 59, 61, 67, 71}
	res := popcntAndSlice(s, m)
	const l uint64 = 18
	if res != l {
		t.Errorf("Wrong And %d !=  %d", res, l)
	}
}

func TestPopcntOrSlice(t *testing.T) {
	s := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	m := []uint64{31, 37, 41, 43, 47, 53, 59, 61, 67, 71}
	res := popcntOrSlice(s, m)
	const l uint64 = 50
	if res != l {
		t.Errorf("Wrong OR %d !=  %d", res, l)
	}
}

func TestPopcntXorSlice(t *testing.T) {
	s := []uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	m := []uint64{31, 37, 41, 43, 47, 53, 59, 61, 67, 71}
	res := popcntXorSlice(s, m)
	const l uint64 = 32
	if res != l {
		t.Errorf("Wrong OR %d !=  %d", res, l)
	}
}

func TestPopcntEmptySlices(t *testing.T) {
	var empty []uint64
	if res := popcntSlice(empty); res != 0 {
		t.Errorf("expected 0, got %d", res)
	}
	if res := popcntMaskSlice(empty, empty); res != 0 {
		t.Errorf("expected 0, got %d", res)
	}
	if res := popcntAndSlice(empty, empty); res != 0 {
		t.Errorf("expected 0, got %d", res)
	}
	if res := popcntOrSlice(empty, empty); res != 0 {
		t.Errorf("expected 0, got %d", res)
	}
	if res := popcntXorSlice(empty, empty); res != 0 {
		t.Errorf("expected 0, got %d", res)
	}
}

func TestPopcntUnrolledLengths(t *testing.T) {
	for n := 0; n <= 32; n++ {
		s := make([]uint64, n)
		m := make([]uint64, n)
		for i := 0; i < n; i++ {
			s[i] = uint64(i*31 + 7)
			m[i] = uint64(i*17 + 13)
		}

		var wantAnd, wantOr, wantXor, wantMask uint64
		for i := 0; i < n; i++ {
			wantAnd += uint64(popcntSlice([]uint64{s[i] & m[i]}))
			wantOr += uint64(popcntSlice([]uint64{s[i] | m[i]}))
			wantXor += uint64(popcntSlice([]uint64{s[i] ^ m[i]}))
			wantMask += uint64(popcntSlice([]uint64{s[i] &^ m[i]}))
		}

		if got := popcntAndSlice(s, m); got != wantAnd {
			t.Errorf("len %d: popcntAndSlice got %d, want %d", n, got, wantAnd)
		}
		if got := popcntOrSlice(s, m); got != wantOr {
			t.Errorf("len %d: popcntOrSlice got %d, want %d", n, got, wantOr)
		}
		if got := popcntXorSlice(s, m); got != wantXor {
			t.Errorf("len %d: popcntXorSlice got %d, want %d", n, got, wantXor)
		}
		if got := popcntMaskSlice(s, m); got != wantMask {
			t.Errorf("len %d: popcntMaskSlice got %d, want %d", n, got, wantMask)
		}
	}
}
