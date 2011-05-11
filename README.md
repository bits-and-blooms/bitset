Package bitset implements bitsets.

It provides methods for making a BitSet of an arbitrary
upper limit, setting and testing bit locations, and clearing
bit locations as well as the entire set.

BitSets are implemented as arrays of uint64s, so it may be best
to limit the upper size to a multiple of 64. It is not an error
to set, test, or clear a bit within a 64 bit boundary.

Example use:
    
    b := MakeBitSet(64000)
    b.SetBit(1000)
    if b.Bit(1000) {
        b.ClearBit(1000)
    }
