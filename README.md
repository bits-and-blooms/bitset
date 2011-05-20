Package bitset implements bitsets, a mapping
between non-negative integers and boolean values. It should be more
efficient than map[uint] bool.

It provides methods for setting, clearing, flipping, and testing
individual integers.

But it also provides set intersection, union, difference,
complement, and symmetric operations, as well as tests to 
check whether any, all, or no bits are set, and querying a 
bitset's current length and number of postive bits.

BitSets are expanded to the size of the largest set bit; the
memory allocation is approximately Max bits, where Max is 
the largest set bit. BitSets are never shrunk. On creation,
a hint can be given for the number of bits that will be used.

Many of the methods, including Set,Clear, and Flip, return 
a BitSet pointer, which allows for chaining.

Example use:

    import "bitset"
    b := bitset.New().b.Set(10).Set(11)
    if b.Test(1000) {
    	b.Clear(1000)
    }
    if B.Intersection(bitset.New().Set(10)).Count() > 1 {
    	fmt.Println("Intersection works.")
    }

As an alternative to BitSets, one should check out the 'big' package,
which provides a (less set-theoretical) view of bitsets.
    
Discussion at: [golang-nuts Google Group](https://groups.google.com/d/topic/golang-nuts/7n1VkRTlBf4/discussion)

