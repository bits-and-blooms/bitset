package bitset

// From Wikipedia: http://en.wikipedia.org/wiki/Hamming_weight
const m1 uint64 = 0x5555555555555555  //binary: 0101...
const m2 uint64 = 0x3333333333333333  //binary: 00110011..
const m4 uint64 = 0x0f0f0f0f0f0f0f0f  //binary:  4 zeros,  4 ones ...
const m8 uint64 = 0x00ff00ff00ff00ff  //binary:  8 zeros,  8 ones ...
const m16 uint64 = 0x0000ffff0000ffff //binary: 16 zeros, 16 ones ...
const m32 uint64 = 0x00000000ffffffff //binary: 32 zeros, 32 ones
const hff uint64 = 0xffffffffffffffff //binary: all ones
const h01 uint64 = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...

// From Wikipedia: count number of set bits.
// This is algorithm popcount_2 in the article retrieved May 9, 2011

func popcount_2(x uint64) uint64 {
	x -= (x >> 1) & m1             //put count of each 2 bits into those 2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each 4 bits into those 4 bits
	x = (x + (x >> 4)) & m4        //put count of each 8 bits into those 8 bits
	x += x >> 8                    //put count of each 16 bits into their lowest 8 bits
	x += x >> 16                   //put count of each 32 bits into their lowest 8 bits
	x += x >> 32                   //put count of each 64 bits into their lowest 8 bits
	return x & 0x7f
}

func popcntSliceGo(s []uint64) uint64 {
	cnt := uint64(0)
	for _, x := range s {
		cnt += popcount_2(x)
	}
	return cnt
}

func popcntMaskSliceGo(s, m []uint64) uint64 {
	cnt := uint64(0)
	for i := range s {
		cnt += popcount_2(s[i] &^ m[i])
	}
	return cnt
}

func popcntAndSliceGo(s, m []uint64) uint64 {
	cnt := uint64(0)
	for i := range s {
		cnt += popcount_2(s[i] & m[i])
	}
	return cnt
}

func popcntOrSliceGo(s, m []uint64) uint64 {
	cnt := uint64(0)
	for i := range s {
		cnt += popcount_2(s[i] | m[i])
	}
	return cnt
}

func popcntXorSliceGo(s, m []uint64) uint64 {
	cnt := uint64(0)
	for i := range s {
		cnt += popcount_2(s[i] ^ m[i])
	}
	return cnt
}
