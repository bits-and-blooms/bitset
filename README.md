Package bitset implements bitsets.

It provides methods for making a BitSet of an arbitrary
upper limit, setting and testing bit locations, and clearing
bit locations as well as the entire set, counting the number
of bits.

It also supports set intersection, union, difference and 
symmetric difference, cloning, equality testing, and subsetting.

Example use:

        import "bitset"
        b := bitset.New(64000)
        b.SetBit(1000)
        b.SetBit(999)
        if b.Bit(1000) {
        	b.ClearBit(1000)
        }
        b.Clear()
    
Discussion at: [golang-nuts Google Group](https://groups.google.com/d/topic/golang-nuts/7n1VkRTlBf4/discussion)

