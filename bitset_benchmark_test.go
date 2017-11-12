// Copyright 2014 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file tests bit sets

package bitset

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// go test -bench=BenchmarkBitSet
// see http://lemire.me/blog/2016/09/22/swift-versus-java-the-bitset-performance-test/
func BenchmarkBitSetLemireCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		set := NewBitSet(0) // we force dynamic memory allocation
		for v := uint(0); v <= 100000000; v += 100 {
			set.Set(v)
		}
	}
}

func BenchmarkBitSetBinaryWrite(b *testing.B) {
	set := NewBitSet(10000)
	buf := bytes.NewBuffer(make([]byte, len(set.values)*8))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Reset()
		binary.Write(buf, binary.LittleEndian, set.values)
	}
}

func BenchmarkBitSetPutUint64(b *testing.B) {
	set := NewBitSet(10000)
	buf := bytes.NewBuffer(make([]byte, len(set.values)*8))
	b.ResetTimer()

	var single [8]byte
	for i := 0; i < b.N; i++ {
		buf.Reset()
		for _, x := range set.values {
			binary.LittleEndian.PutUint64(single[:], x)
			_, err := buf.Write(single[:])
			if err != nil {
				b.Fatalf(err.Error())
			}
		}
	}
}
