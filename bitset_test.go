// Copyright 2011 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file tests bit sets
package bitset

import (
	"testing"
)

func TestMakeBitSet(t *testing.T) {
	v := MakeBitSet(10000)
	if v.Bit(0) != false {
		t.Errorf("Unable to make a bit set and read its 0th value.")
	}
}

func TestBitSetIsClear(t *testing.T) {
	v := MakeBitSet(1000)
	for i := uint(0); i < 1000; i++ {
		if v.Bit(i) != false {
			t.Errorf("Bit %d is set, and it shouldn't be.", i)
		}
	}
}


func TestBitSetAndGet(t *testing.T) {
	v := MakeBitSet(1000)
	v.SetBit(100)
	if v.Bit(100) != true {
		t.Errorf("Bit %d is clear, and it shouldn't be.", 100)
	}
}

func TestLotsOfSetsAndGets(t *testing.T) {
	tot := uint(100000)
	v := MakeBitSet(tot)
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
	v := MakeBitSet(tot)
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
	v := MakeBitSet(64)
	defer func() {
	        if r := recover(); r == nil {
	            t.Error("Out of index error within the next set of bits should have caused a panic")
	        }
	    }()
	v.SetBit(1000)
}

func TestOutOfBoundsOK(t *testing.T) {
	v := MakeBitSet(65)
	defer func() {
	        if r := recover(); r != nil {
	            t.Error("Out of index error within the next set of bits should not caused a panic")
	        }
	    }()
	v.SetBit(66)   
}
