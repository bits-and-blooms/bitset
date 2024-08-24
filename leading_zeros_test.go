//go:build !go1.9
// +build !go1.9

package bitset

import (
	"testing"
)

func TestLen64(t *testing.T) {
	for i := 0; i < 64; i++ {
		received := len64(uint64(1) << i)
		expected := uint(i + 1)
		if received != expected {
			t.Errorf("len64(%b) is incorrect: received %d, expected %d", uint64(1)<<i, received, expected)
		}
	}
}
