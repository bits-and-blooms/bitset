package bitset

// bit population count, take from
// https://code.google.com/p/go/issues/detail?id=4988#c11
// credit: https://code.google.com/u/arnehormann/
func popcount(x uint64) (n uint64) {
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return x >> 56
}

func popcntSlice(s []uint64) uint64 {
	if useAsm {
		return popcntSliceAsm(s)
	}
	cnt := uint64(0)
	for _, x := range s {
		cnt += popcount(x)
	}
	return cnt
}

func popcntMaskSlice(s, m []uint64) uint64 {
	if useAsm {
		return popcntMaskSliceAsm(s, m)
	}
	cnt := uint64(0)
	for i := range s {
		cnt += popcount(s[i] &^ m[i])
	}
	return cnt
}

func popcntAndSlice(s, m []uint64) uint64 {
	if useAsm {
		return popcntAndSliceAsm(s, m)
	}
	cnt := uint64(0)
	for i := range s {
		cnt += popcount(s[i] & m[i])
	}
	return cnt
}

func popcntOrSlice(s, m []uint64) uint64 {
	if useAsm {
		return popcntOrSliceAsm(s, m)
	}
	cnt := uint64(0)
	for i := range s {
		cnt += popcount(s[i] | m[i])
	}
	return cnt
}

func popcntXorSlice(s, m []uint64) uint64 {
	if useAsm {
		return popcntXorSliceAsm(s, m)
	}
	cnt := uint64(0)
	for i := range s {
		cnt += popcount(s[i] ^ m[i])
	}
	return cnt
}
