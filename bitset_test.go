// Copyright 2011 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file tests bit sets 

package bitset

import (
	"testing"
)

func TestEmptySet(t *testing.T) {
	defer func() {
	        if r := recover(); r != nil {
	            t.Error("A zero-length bitset should be fine")
	        }
	    }()
	b := New(0)
	if b.Cap() != 0 {
		t.Errorf("Empty set should have capacity 0, not %d", b.Cap())
	}
}

func TestBitSetNew(t *testing.T) {
	v := New(16)
	if v.Bit(0) != false {
		t.Errorf("Unable to make a bit set and read its 0th value.")
	}
}

func TestBitSetIsClear(t *testing.T) {
	v := New(1000)
	for i := uint(0); i < 1000; i++ {
		if v.Bit(i) != false {
			t.Errorf("Bit %d is set, and it shouldn't be.", i)
		}
	}
}


func TestBitSetAndGet(t *testing.T) {
	v := New(1000)
	v.SetBit(100)
	if v.Bit(100) != true {
		t.Errorf("Bit %d is clear, and it shouldn't be.", 100)
	}
}

func TestLotsOfSetsAndGets(t *testing.T) {
	tot := uint(100000)
	v := New(tot)
	for i := uint(0); i < tot; i+=2 {
		v.SetBit(i)
	}
	for i := uint(0); i < tot; i++ {
		if i % 2 == 0 {
			if v.Bit(i) != true {
				t.Errorf("Bit %d is clear, and it shouldn't be.", i)
			}
		} else {
			if v.Bit(i) != false {
				t.Errorf("Bit %d is set, and it shouldn't be.", i)
			}
		}
	}
}

func TestClear(t *testing.T) {
	tot := uint(1000)
	v := New(tot)
	for i := uint(0); i < tot; i++ {
		v.SetBit(i)
	}
	v.Clear()
	for i := uint(0); i < tot; i++ {
		if v.Bit(i) != false {
			t.Errorf("Bit %d is set, and it shouldn't be.", i)
			break
		}
	}
}

func TestOutOfBoundsBad(t *testing.T) {
	v := New(64)
	defer func() {
	        if r := recover(); r == nil {
	            t.Error("Long distance out of index error should have caused a panic")
	        }
	    }()
	v.SetBit(1000)
}

func TestOutOfBoundsOK(t *testing.T) {
	v := New(65)
	defer func() {
	        if r := recover(); r == nil {
	            t.Error("Local out of index error should have caused a panic")
	        }
	    }()
	v.SetBit(66)   
}

func TestCap(t *testing.T) {
	v := New(1000)
	if v.Cap() != 1000 {
		t.Errorf("Cap should be 1000, but is %d.", v.Cap())
	}
}

func TestLen(t *testing.T) {
	v := New(1000)
	if v.Len() != 1000 {
		t.Errorf("Len should be 1000, but is %d.", v.Cap())
	}
}

func TestCount(t *testing.T) {
	tot := uint(64*4+11) // just some multi unit64 number
	v := New(tot)
	checkLast := true
	for i := uint(0); i < tot; i++ {
		sz := v.Count()
		if sz != i {
			t.Errorf("Count reported as %d, but it should be %d", sz, i)
			checkLast = false
			break
		} 
		v.SetBit(i)
	}
	if checkLast {
		sz := v.Count()
		if sz != tot {
			t.Errorf("After all bits set, size reported as %d, but it should be %d", sz, tot)   
		}
	}
}

// test setting every 3rd bit, just in case something odd is happening
func TestCount2(t *testing.T) {
	tot := uint(64*4+11) // just some multi unit64 number
	v := New(tot)
	for i := uint(0); i < tot; i+=3 {
		sz := v.Count()
		if sz != i/3 {
			t.Errorf("Count reported as %d, but it should be %d", sz, i)
			break
		} 
		v.SetBit(i)
	}
}

// nil tests

func TestNullBit(t *testing.T) {
	var v *BitSet = nil
	defer func() {
	        if r := recover(); r == nil {
	            t.Error("Checking bit of null reference should have caused a panic")
	        }
	    }()
	v.Bit(66)   
}

func TestNullSetBit(t *testing.T) {
	var v *BitSet = nil
	defer func() {
	        if r := recover(); r == nil {
	            t.Error("Setting bit of null reference should have caused a panic")
	        }
	    }()
	v.SetBit(66)   
} 

func TestNullClearBit(t *testing.T) {
	var v *BitSet = nil
	defer func() {
	        if r := recover(); r == nil {
	            t.Error("Clearning bit of null reference should have caused a panic")
	        }
	    }()
	v.ClearBit(66)   
}

func TestNullClear(t *testing.T) {
	var v *BitSet = nil
	defer func() {
	        if r := recover(); r != nil {
	            t.Error("Clearing null reference should not have caused a panic")
	        }
	    }()
	v.Clear()   
} 

func TestNullCount(t *testing.T) {
	var v *BitSet = nil
	defer func() {
	        if r := recover(); r != nil {
	            t.Error("Counting null reference should not have caused a panic")
	        }
	    }() 
	cnt := v.Count()
	if cnt != 0 {
	   t.Errorf("Count reported as %d, but it should be 0", cnt)
	}  
}

func TestMap(t *testing.T) {
	attended := map[string] bool {
	    "Ann": true,
	    "Joe": true,
	}
	if there := attended["Ann"]; !there {
		t.Errorf("John didn't come: %d", there)
	}
}

func TestEqual(t *testing.T) {
	a := New(100)
	b := New(99)
	c := New(100)
	if a.Equal(b) {
        t.Error("Sets of different sizes should be not be equal")
    }
	if !a.Equal(c){
        t.Error("Two empty sets of the same size should be equal")
    }
 	a.SetBit(99)
 	c.SetBit(0)
    if a.Equal(c){
        t.Error("Two sets with differences should not be equal")
    }
	c.SetBit(99)
 	a.SetBit(0)
    if !a.Equal(c){
        t.Error("Two sets with the same bits set should be equal")
    }
}


// NOTE: These following tests rely on understanding the
// the internal structure of bitsets. It's a little easier
// to test that way.


func TestDifferenceSimple (t *testing.T) {
	a := New(64)
	b := New(64)
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0x0f0f0f0f0f0f0f0f
	c,e := a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.set[0] != 0 {
		t.Errorf("Set difference of exact sets should be 0, but was %x", c.set[0])
	}
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0
	c,e = a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.set[0] != 0x0f0f0f0f0f0f0f0f {
		t.Error("Set difference of set with bits set less one all clear should be original set")
	}
	a.set[0] = 0
	b.set[0] = 0x0f0f0f0f0f0f0f0f
	c,e = a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.set[0] != 0 {
		t.Error("Set difference of set with no bits set less one with bits should be zero")
	}
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0xffffffffffffffff
	c,e = a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.set[0] != 0 {
		t.Error("Set diff of set to all bits set should be zero")
	}
	a.set[0] = 0xffffffffffffffff
	b.set[0] = 0x0f0f0f0f0f0f0f0f
	c,e = a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.set[0] != 0xf0f0f0f0f0f0f0f0 {
		t.Error("Set diff from all bits should remove bits in 2nd set")
	}	
}

func TestSetDifferenceDifferentSizes (t *testing.T) {
	a := New(642)
	b := New(1011)
	c,e := a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.Count() != 0 {
		t.Errorf("Set difference of empty sets should be empty.")
	}
	c,e = b.Difference(a)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.Count() != 0 {
		t.Errorf("Set difference of empty sets should be empty.")
	}
	a.SetBit(450)
	b.SetBit(999)
	c,e = a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.Count() != 1 || !c.Bit(450) {
		t.Errorf("Set difference a-b should be a, basically")
	}
	a.Clear()
	b.Clear()
	a.SetBit(450)
	b.SetBit(450)
	b.SetBit(999)
	c,e = a.Difference(b)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.Count() != 0 || c.Bit(450) {
		t.Errorf("Set difference a-b should now be empty, but count is %d", c.Count())
	}
	c,e = b.Difference(a)
	if e!= nil {
		t.Errorf("Set difference returned error: %d", e)
	}
	if c.Count() != 1 || !c.Bit(999) {
		t.Errorf("Set difference b-a should be b, basically")
	}
}

func TestUnionSimple (t *testing.T) {
	a := New(64)
	b := New(64)
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0xf0f0f0f0f0f0f0f0
	c,e := a.Union(b)
	if e!= nil {
		t.Errorf("Set union returned error: %d", e)
	}
	if c.Count() != 64 {
		t.Error("Union should be of size 64")
	}
}

func TestUnionDifferent (t *testing.T) {
	a := New(64)
	b := New(64*2+5)
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[1] = 0xf0f0f0f0f0f0f0f0
	b.SetBit(64*2)
	c,e := a.Union(b)
	if e!= nil {
		t.Errorf("Set union returned error: %d", e)
	}
	if c.Count() != 65 {
		t.Error("Union should be of size 65")
	}
	if c.set[0] != 0x0f0f0f0f0f0f0f0f {
		t.Error("First word is wrong in union")
	}
	if c.set[1] != 0xf0f0f0f0f0f0f0f0 {
		t.Error("Second word is wrong in union")
	}
	if c.set[2] != 1 {
		t.Errorf("c.set[0]: %x, c.set[1]: %x, c.set[2]: %x", c.set[0],c.set[1],c.set[2])
	}
}

func TestIntersectionSimple (t *testing.T) {
	a := New(64)
	b := New(64)
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0xf0f0f0f0f0f0f0f0
	c,e := a.Intersection(b)
	if e!= nil {
		t.Errorf("Set intersection returned error: %d", e)
	}
	if c.Count() != 0 {
		t.Error("Intersection should be empty")
	}
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0x0f0f0f0f0f0f0f0f
	c,e = a.Intersection(b)
	if e!= nil {
		t.Errorf("Set intersection returned error: %d", e)
	}
	if c.set[0] != 0x0f0f0f0f0f0f0f0f {
		t.Error("Intersection should be the same")
	}
}

func TestIntersectionDifferent (t *testing.T) {
	a := New(64)
	b := New(64*2+5)
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[1] = 0xf0f0f0f0f0f0f0f0
	b.SetBit(64*2+2)
	c,e := a.Intersection(b)
	if e!= nil {
		t.Errorf("Set union returned error: %d", e)
	}
	if c.Count() != 0 {
		t.Error("Union should be of size 0")
	}
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[1] = 0x0f0f0f0f0f0f0f0f
	b.SetBit(64*2+2)
	c,e = a.Intersection(b)
	if e!= nil {
		t.Errorf("Set union returned error: %d", e)
	}
	if c.set[0] != 0x0f0f0f0f0f0f0f0f {
		t.Error("Intersection should be of similar word")
	}	
}

func TestSymmetricDifferenceSimple (t *testing.T) {
	a := New(64)
	b := New(64)
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0xf0f0f0f0f0f0f0f0
	c,e := a.SymmetricDifference(b)
	if e!= nil {
		t.Errorf("Set symmetric difference returned error: %d", e)
	}
	if c.Count() != 64 {
		t.Error("Symmetric difference should be of size 64")
	}
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0x0f0f0f0f0f0f0f0f
	c,e = a.SymmetricDifference(b)
	if e!= nil {
		t.Errorf("Set symmetric difference returned error: %d", e)
	}
	if c.Count() != 0 {
		t.Error("Symmetric difference should be of size 0")
	}
}

func TestSymmetricDifferenceDifferent (t *testing.T) {
	a := New(64)
	b := New(64*2+5)
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[1] = 0xf0f0f0f0f0f0f0f0
	b.SetBit(64*2)
	c,e := a.SymmetricDifference(b)
	if e!= nil {
		t.Errorf("Set SymmetricDifference returned error: %d", e)
	}
	if c.Count() != 65 {
		t.Error("SymmetricDifference should be of size 65")
	}
	if c.set[0] != 0x0f0f0f0f0f0f0f0f {
		t.Error("First word is wrong in SymmetricDifference")
	}
	if c.set[1] != 0xf0f0f0f0f0f0f0f0 {
		t.Error("Second word is wrong in SymmetricDifference")
	}
	if c.set[2] != 1 {
		t.Errorf("c.set[0]: %x, c.set[1]: %x, c.set[2]: %x", c.set[0],c.set[1],c.set[2])
	}
	a.Clear()
	b.Clear()
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	b.set[0] = 0x0f0f0f0f0f0f0f0f
	b.SetBit(64*2)
	c,e = a.SymmetricDifference(b)
	if e!= nil {
		t.Errorf("Set SymmetricDifference returned error: %d", e)
	}
	if c.Count() != 1 {
		t.Errorf("SymmetricDifference should be of size 1, but is %d", c.Count())
	}
	if c.set[0] != 0 {
		t.Error("First word is wrong in SymmetricDifference")
	}
	if c.set[1] != 0 {
		t.Error("Second word is wrong in SymmetricDifference")
	}
	if c.set[2] != 1 {
		t.Errorf("c.set[0]: %x, c.set[1]: %x, c.set[2]: %x", c.set[0],c.set[1],c.set[2])
	}
}

func TestSubset(t *testing.T) {
	a := New(64*3)
	b, err := a.Subset(0,64*3+1)
	if err == nil {
		t.Error("End index should be < length")
	}
	b, err = a.Subset(0,1)
	if err != nil {
		t.Error("Simple subset should create no error")
	}
	if b.Cap() != 1 {
		t.Errorf("capacity should be 1, was %d", b.Cap())
	}
	c, err := a.Subset(0,64*3)
	if err != nil {
		t.Error("Complete subset should create no error")
	}
	if c.Cap() != 64*3 {
		t.Errorf("capacity should be %d, was %d", 64*3, b.Cap())
	}
	a.set[0] = 0x0f0f0f0f0f0f0f0f
	a.set[1] = 0x0f0f0f0f0f0f0f0f
	a.set[2] = 0x0f0f0f0f0f0f0f0f
	d, err := a.Subset(64/2, 64+64/2)
	if err != nil {
		t.Error("Overlap subset")
	}
	if d.Cap() != 64 {
		t.Errorf("capacity should be %d, was %d", 64, d.Cap())
	}
	if d.set[0] != 0x0f0f0f0f0f0f0f0f {
		t.Errorf("bitpattern incorrect, was %f", d.set[0])
	}   
}

func TestFlipAll(t *testing.T) {
	a := New(64*2)
	a.FlipAll()
	if a.Count() != 64*2 {
		t.Errorf("After flipping all, count should be %d, but is %d", 64*2, a.Count())
	}
	a = New(64*11+10)
	a.FlipAll()
	if a.Count() != 64*11+10 {
		t.Errorf("After flipping all, count should be %d, but is %d", 64*11+10, a.Count())
	}
}

func TestFlipBit(t *testing.T) {
	a := New(64*2)
	a.FlipBit(100)
	if a.Count() != 1 {
		t.Errorf("After flipping all, count should be %d, but is %d", 1, a.Count())
	}
	if !a.Bit(100) {
		t.Errorf("After flipping a clear bit, the bit should be set, but is not")
	}
	a.FlipBit(100)
	if a.Count() != 0 {
		t.Errorf("After flipping all, count should be %d, but is %d", 1, a.Count())
	}
	if a.Bit(100) {
		t.Errorf("After flipping a set bit, the bit should be clear, but is not")
	}
}

// TODO: Tests for None, Any, ALL 