// Copyright 2011 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file tests bit sets 

package bitset

import (
	"testing"
)

func TestbBitSetNew(t *testing.T) {
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
