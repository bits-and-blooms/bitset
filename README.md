Package bitset implements bitsets.

It provides methods for making a BitSet of an arbitrary
upper limit, setting and testing bit locations, and clearing
bit locations as well as the entire set.

Example use:
    
    b := MakeBitSet(64000)
    b.SetBit(1000)
    if b.Bit(1000) {
        b.ClearBit(1000)
    }
