package bitset

import "math/bits"

func popcntSlice(s []uint64) (cnt uint64) {
	for _, x := range s {
		cnt += uint64(bits.OnesCount64(x))
	}
	return
}

func popcntMaskSlice(s, m []uint64) (cnt uint64) {
	_ = m[len(s)-1] // BCE
	for i := range s {
		cnt += uint64(bits.OnesCount64(s[i] &^ m[i]))
	}
	return
}

func popcntAndSlice(s, m []uint64) (cnt uint64) {
	_ = m[len(s)-1] // BCE
	for i := range s {
		cnt += uint64(bits.OnesCount64(s[i] & m[i]))
	}
	return
}

func popcntOrSlice(s, m []uint64) (cnt uint64) {
	_ = m[len(s)-1] // BCE
	for i := range s {
		cnt += uint64(bits.OnesCount64(s[i] | m[i]))
	}
	return
}

func popcntXorSlice(s, m []uint64) (cnt uint64) {
	_ = m[len(s)-1] // BCE
	for i := range s {
		cnt += uint64(bits.OnesCount64(s[i] ^ m[i]))
	}
	return
}
