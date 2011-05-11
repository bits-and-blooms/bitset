Package bitset implements bitsets.

It provides methods for making a BitSet of an arbitrary
upper limit, setting and testing bit locations, and clearing
bit locations as well as the entire set.

Example use:
    
    b := bitset.New(64000)
    b.SetBit(1000)
    if b.Bit(1000) {
        b.ClearBit(1000)
    }
    
Discussion at: [golang-nuts Google Group](https://groups.google.com/d/topic/golang-nuts/7n1VkRTlBf4/discussion)

