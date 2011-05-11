// Copyright 2011 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file tests bit sets
package bitset

import (
	"testing"
)

func TestbBitSetNew(t *testing.T) {
	v := New(10000)
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
	        if r := recover(); r != nil {
	            t.Error("Out of index error within the next set of bits should not have caused a panic")
	        }
	    }()
	v.SetBit(1000)
}

func TestOutOfBoundsOK(t *testing.T) {
	v := New(65)
	defer func() {
	        if r := recover(); r != nil {
	            t.Error("Out of index error within the next set of bits should not caused a panic")
	        }
	    }()
	v.SetBit(66)   
}

func TestMaxSizet(t *testing.T) {
	v := New(1000)
	if v.MaxSize() != 1000 {
		t.Errorf("MaxSize should be 1000, but is %d.", v.MaxSize())
	}
}

func TestSize(t *testing.T) {
	tot := uint(64*4+11) // just some multi unit64 number
	v := New(tot)
	checkLast := true
	for i := uint(0); i < tot; i++ {
		sz := v.Size()
		if sz != i {
			t.Errorf("Size reported as %d, but it should be %d", sz, i)
			checkLast = false
			break
		} 
		v.SetBit(i)
	}
	if checkLast {
		sz := v.Size()
		if sz != tot {
			t.Errorf("After all bits set, size reported as %d, but it should be %d", sz, tot)   
		}
	}
}

// test setting every 3rd bit, just in case something odd is happening
func TestSize2(t *testing.T) {
	tot := uint(64*4+11) // just some multi unit64 number
	v := New(tot)
	for i := uint(0); i < tot; i+=3 {
		sz := v.Size()
		if sz != i/3 {
			t.Errorf("Size reported as %d, but it should be %d", sz, i)
			break
		} 
		v.SetBit(i)
	}
}
