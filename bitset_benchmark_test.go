// Copyright 2014 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file tests bit sets

package bitset

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	sz := 100000
	s := New(uint(sz))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Set(uint(r.Int31n(int32(sz))))
	}
}

func BenchmarkGetTest(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	sz := 100000
	s := New(uint(sz))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Test(uint(r.Int31n(int32(sz))))
	}
}

func BenchmarkSetExpand(b *testing.B) {
	b.StopTimer()
	sz := uint(100000)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var s BitSet
		s.Set(sz)
	}
}

// go test -bench=Count
func BenchmarkCount(b *testing.B) {
	b.StopTimer()
	s := New(100000)
	for i := 0; i < 100000; i += 100 {
		s.Set(uint(i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Count()
	}
}

// go test -bench=Iterate
func BenchmarkIterate(b *testing.B) {
	b.StopTimer()
	s := New(10000)
	for i := 0; i < 10000; i += 3 {
		s.Set(uint(i))
	}
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		c := uint(0)
		for i, e := s.NextSet(0); e; i, e = s.NextSet(i + 1) {
			c++
		}
	}
}

// go test -bench=SparseIterate
func BenchmarkSparseIterate(b *testing.B) {
	b.StopTimer()
	s := New(100000)
	for i := 0; i < 100000; i += 30 {
		s.Set(uint(i))
	}
	b.StartTimer()
	for j := 0; j < b.N; j++ {
		c := uint(0)
		for i, e := s.NextSet(0); e; i, e = s.NextSet(i + 1) {
			c++
		}
	}
}

// go test -bench=BitsetOps
func BenchmarkBitsetOps(b *testing.B) {
	// let's not write into s inside the benchmarks
	s := New(100000)
	for i := 0; i < 100000; i += 100 {
		s.Set(uint(i))
	}
	cpy := s.Clone()

	b.Run("Equal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Equal(cpy)
		}
	})

	b.Run("FlipRange", func(b *testing.B) {
		s = s.Clone()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.FlipRange(0, 100000)
		}
	})

	b.Run("NextSet", func(b *testing.B) {
		s = New(100000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.NextSet(0)
		}
	})

	b.Run("NextClear", func(b *testing.B) {
		s = New(100000)
		s.FlipRange(0, 100000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.NextClear(0)
		}
	})

	b.Run("DifferenceCardinality", func(b *testing.B) {
		empty := New(100000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.DifferenceCardinality(empty)
		}
	})

	b.Run("InPlaceDifference", func(b *testing.B) {
		s = s.Clone()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.InPlaceDifference(cpy)
		}
	})

	b.Run("InPlaceUnion", func(b *testing.B) {
		s = s.Clone()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.InPlaceUnion(cpy)
		}
	})

	b.Run("InPlaceIntersection", func(b *testing.B) {
		s = s.Clone()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.InPlaceIntersection(cpy)
		}
	})

	b.Run("InPlaceSymmetricDifference", func(b *testing.B) {
		s = s.Clone()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s.InPlaceSymmetricDifference(cpy)
		}
	})
}

// go test -bench=LemireCreate
// see http://lemire.me/blog/2016/09/22/swift-versus-java-the-bitset-performance-test/
func BenchmarkLemireCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bitmap := New(0) // we force dynamic memory allocation
		for v := uint(0); v <= 100000000; v += 100 {
			bitmap.Set(v)
		}
	}
}

// go test -bench=LemireCount
// see http://lemire.me/blog/2016/09/22/swift-versus-java-the-bitset-performance-test/
func BenchmarkLemireCount(b *testing.B) {
	bitmap := New(100000000)
	for v := uint(0); v <= 100000000; v += 100 {
		bitmap.Set(v)
	}
	b.ResetTimer()
	sum := uint(0)
	for i := 0; i < b.N; i++ {
		sum += bitmap.Count()
	}
	if sum == 0 { // added just to fool ineffassign
		return
	}
}

// go test -bench=LemireIterate
// see http://lemire.me/blog/2016/09/22/swift-versus-java-the-bitset-performance-test/
func BenchmarkLemireIterate(b *testing.B) {
	bitmap := New(100000000)
	for v := uint(0); v <= 100000000; v += 100 {
		bitmap.Set(v)
	}
	b.ResetTimer()
	sum := uint(0)
	for i := 0; i < b.N; i++ {
		for j, e := bitmap.NextSet(0); e; j, e = bitmap.NextSet(j + 1) {
			sum++
		}
	}
	if sum == 0 { // added just to fool ineffassign
		return
	}
}

// go test -bench=LemireIterateb
// see http://lemire.me/blog/2016/09/22/swift-versus-java-the-bitset-performance-test/
func BenchmarkLemireIterateb(b *testing.B) {
	bitmap := New(100000000)
	for v := uint(0); v <= 100000000; v += 100 {
		bitmap.Set(v)
	}
	b.ResetTimer()
	sum := uint(0)
	for i := 0; i < b.N; i++ {
		for j, e := bitmap.NextSet(0); e; j, e = bitmap.NextSet(j + 1) {
			sum += j
		}
	}

	if sum == 0 { // added just to fool ineffassign
		return
	}
}

// go test -bench=BenchmarkLemireIterateManyb
// see http://lemire.me/blog/2016/09/22/swift-versus-java-the-bitset-performance-test/
func BenchmarkLemireIterateManyb(b *testing.B) {
	bitmap := New(100000000)
	for v := uint(0); v <= 100000000; v += 100 {
		bitmap.Set(v)
	}
	buffer := make([]uint, 256)
	b.ResetTimer()
	sum := uint(0)
	for i := 0; i < b.N; i++ {
		j := uint(0)
		j, buffer = bitmap.NextSetMany(j, buffer)
		for ; len(buffer) > 0; j, buffer = bitmap.NextSetMany(j, buffer) {
			for k := range buffer {
				sum += buffer[k]
			}
			j++
		}
	}

	if sum == 0 { // added just to fool ineffassign
		return
	}
}

func setRnd(bits []uint64, halfings int) {
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := range bits {
		bits[i] = 0xFFFFFFFFFFFFFFFF
		for j := 0; j < halfings; j++ {
			bits[i] &= rnd.Uint64()
		}
	}
}

// go test -bench=BenchmarkFlorianUekermannIterateMany
func BenchmarkFlorianUekermannIterateMany(b *testing.B) {
	var input = make([]uint64, 68)
	setRnd(input, 4)
	var bitmap = From(input)
	buffer := make([]uint, 256)
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		var last, batch = bitmap.NextSetMany(0, buffer)
		for len(batch) > 0 {
			for _, idx := range batch {
				checksum += idx
			}
			last, batch = bitmap.NextSetMany(last+1, batch)
		}
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkFlorianUekermannIterateManyReg(b *testing.B) {
	var input = make([]uint64, 68)
	setRnd(input, 4)
	var bitmap = From(input)
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		for j, e := bitmap.NextSet(0); e; j, e = bitmap.NextSet(j + 1) {
			checksum += j
		}
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

// function provided by FlorianUekermann
func good(set []uint64) (checksum uint) {
	for wordIdx, word := range set {
		var wordIdx = uint(wordIdx * 64)
		for word != 0 {
			var bitIdx = uint(trailingZeroes64(word))
			word ^= 1 << bitIdx
			var index = wordIdx + bitIdx
			checksum += index
		}
	}
	return checksum
}

func BenchmarkFlorianUekermannIterateManyComp(b *testing.B) {
	var input = make([]uint64, 68)
	setRnd(input, 4)
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		checksum += good(input)
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

/////// Mid density

// go test -bench=BenchmarkFlorianUekermannLowDensityIterateMany
func BenchmarkFlorianUekermannLowDensityIterateMany(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 50000; i++ {
		input[rnd.Uint64()%1000000] = 1
	}
	var bitmap = From(input)
	buffer := make([]uint, 256)
	b.ResetTimer()
	var sum = uint(0)
	for i := 0; i < b.N; i++ {
		j := uint(0)
		j, buffer = bitmap.NextSetMany(j, buffer)
		for ; len(buffer) > 0; j, buffer = bitmap.NextSetMany(j, buffer) {
			for k := range buffer {
				sum += buffer[k]
			}
			j++
		}
	}
	if sum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkFlorianUekermannLowDensityIterateManyReg(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 50000; i++ {
		input[rnd.Uint64()%1000000] = 1
	}
	var bitmap = From(input)
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		for j, e := bitmap.NextSet(0); e; j, e = bitmap.NextSet(j + 1) {
			checksum += j
		}
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkFlorianUekermannLowDensityIterateManyComp(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 50000; i++ {
		input[rnd.Uint64()%1000000] = 1
	}
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		checksum += good(input)
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

/////// Mid density

// go test -bench=BenchmarkFlorianUekermannMidDensityIterateMany
func BenchmarkFlorianUekermannMidDensityIterateMany(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 3000000; i++ {
		input[rnd.Uint64()%1000000] |= uint64(1) << (rnd.Uint64() % 64)
	}
	var bitmap = From(input)
	buffer := make([]uint, 256)
	b.ResetTimer()
	sum := uint(0)
	for i := 0; i < b.N; i++ {
		j := uint(0)
		j, buffer = bitmap.NextSetMany(j, buffer)
		for ; len(buffer) > 0; j, buffer = bitmap.NextSetMany(j, buffer) {
			for k := range buffer {
				sum += buffer[k]
			}
			j++
		}
	}

	if sum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkFlorianUekermannMidDensityIterateManyReg(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 3000000; i++ {
		input[rnd.Uint64()%1000000] |= uint64(1) << (rnd.Uint64() % 64)
	}
	var bitmap = From(input)
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		for j, e := bitmap.NextSet(0); e; j, e = bitmap.NextSet(j + 1) {
			checksum += j
		}
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkFlorianUekermannMidDensityIterateManyComp(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 3000000; i++ {
		input[rnd.Uint64()%1000000] |= uint64(1) << (rnd.Uint64() % 64)
	}
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		checksum += good(input)
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

////////// High density

func BenchmarkFlorianUekermannMidStrongDensityIterateMany(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 20000000; i++ {
		input[rnd.Uint64()%1000000] |= uint64(1) << (rnd.Uint64() % 64)
	}
	var bitmap = From(input)
	buffer := make([]uint, 256)
	b.ResetTimer()
	sum := uint(0)
	for i := 0; i < b.N; i++ {
		j := uint(0)
		j, buffer = bitmap.NextSetMany(j, buffer)
		for ; len(buffer) > 0; j, buffer = bitmap.NextSetMany(j, buffer) {
			for k := range buffer {
				sum += buffer[k]
			}
			j++
		}
	}

	if sum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkFlorianUekermannMidStrongDensityIterateManyReg(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 20000000; i++ {
		input[rnd.Uint64()%1000000] |= uint64(1) << (rnd.Uint64() % 64)
	}
	var bitmap = From(input)
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		for j, e := bitmap.NextSet(0); e; j, e = bitmap.NextSet(j + 1) {
			checksum += j
		}
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkFlorianUekermannMidStrongDensityIterateManyComp(b *testing.B) {
	var input = make([]uint64, 1000000)
	var rnd = rand.NewSource(0).(rand.Source64)
	for i := 0; i < 20000000; i++ {
		input[rnd.Uint64()%1000000] |= uint64(1) << (rnd.Uint64() % 64)
	}
	b.ResetTimer()
	var checksum = uint(0)
	for i := 0; i < b.N; i++ {
		checksum += good(input)
	}
	if checksum == 0 { // added just to fool ineffassign
		return
	}
}

func BenchmarkBitsetReadWrite(b *testing.B) {
	s := New(100000)
	for i := 0; i < 100000; i += 100 {
		s.Set(uint(i))
	}
	buffer := bytes.Buffer{}
	temp := New(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.WriteTo(&buffer)
		temp.ReadFrom(&buffer)
		buffer.Reset()
	}
}

func BenchmarkIsSuperSet(b *testing.B) {
	new := func(len int, density float64) *BitSet {
		r := rand.New(rand.NewSource(42))
		bs := New(uint(len))
		for i := 0; i < len; i++ {
			bs.SetTo(uint(i), r.Float64() < density)
		}
		return bs
	}

	bench := func(name string, lenS, lenSS int, density float64, overrideS, overrideSS map[int]bool, f func(*BitSet, *BitSet) bool) {
		s := new(lenS, density)
		ss := new(lenSS, density)

		for i, v := range overrideS {
			s.SetTo(uint(i), v)
		}
		for i, v := range overrideSS {
			ss.SetTo(uint(i), v)
		}

		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = f(ss, s)
			}
		})
	}

	f := func(ss, s *BitSet) bool {
		return ss.IsSuperSet(s)
	}
	fStrict := func(ss, s *BitSet) bool {
		return ss.IsStrictSuperSet(s)
	}

	for _, len := range []int{1, 10, 100, 1000, 10000, 100000} {
		density := 0.5
		bench(fmt.Sprintf("equal, len=%d", len),
			len, len, density, nil, nil, f)
		bench(fmt.Sprintf("equal, len=%d, strict", len),
			len, len, density, nil, nil, fStrict)
	}

	for _, density := range []float64{0, 0.05, 0.2, 0.8, 0.95, 1} {
		len := 10000
		bench(fmt.Sprintf("equal, density=%.2f", density),
			len, len, density, nil, nil, f)
		bench(fmt.Sprintf("equal, density=%.2f, strict", density),
			len, len, density, nil, nil, fStrict)
	}

	for _, diff := range []int{0, 100, 1000, 9999} {
		len := 10000
		density := 0.5
		overrideS := map[int]bool{diff: true}
		overrideSS := map[int]bool{diff: false}
		bench(fmt.Sprintf("subset, len=%d, diff=%d", len, diff),
			len, len, density, overrideS, overrideSS, f)
		bench(fmt.Sprintf("subset, len=%d, diff=%d, strict", len, diff),
			len, len, density, overrideS, overrideSS, fStrict)
	}

	for _, diff := range []int{0, 100, 1000, 9999} {
		len := 10000
		density := 0.5
		overrideS := map[int]bool{diff: false}
		overrideSS := map[int]bool{diff: true}
		bench(fmt.Sprintf("superset, len=%d, diff=%d", len, diff),
			len, len, density, overrideS, overrideSS, f)
		bench(fmt.Sprintf("superset, len=%d, diff=%d, strict", len, diff),
			len, len, density, overrideS, overrideSS, fStrict)
	}
}
