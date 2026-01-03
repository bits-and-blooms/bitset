package bitset

import (
	"math/rand"
	"testing"
)

const (
	// max bit position to avoid excessive memory usage
	maxBitPosition = 100000
	maxBytes       = maxBitPosition / 8
	maxOperations  = 1000
)

// FuzzBasicOps tests basic set, clear, test operations
func FuzzBasicOps(f *testing.F) {
	// Add seed corpus
	f.Add(uint(0), uint(1), true)
	f.Add(uint(63), uint(64), false)
	f.Add(uint(128), uint(256), true)
	f.Add(uint(1000), uint(2000), false)

	f.Fuzz(func(t *testing.T, bit1, bit2 uint, setValue bool) {
		if bit1 > maxBitPosition || bit2 > maxBitPosition {
			return
		}

		b := New(0)

		// Test Set operation
		b.Set(bit1)
		if !b.Test(bit1) {
			t.Errorf("Set(%d) failed: bit should be set", bit1)
		}

		// Test SetTo operation
		b.SetTo(bit2, setValue)
		if b.Test(bit2) != setValue {
			t.Errorf("SetTo(%d, %v) failed: expected %v, got %v", bit2, setValue, setValue, b.Test(bit2))
		}

		// Test Clear operation
		b.Clear(bit1)
		if b.Test(bit1) {
			t.Errorf("Clear(%d) failed: bit should be clear", bit1)
		}

		// Test Flip operation
		originalState := b.Test(bit2)
		b.Flip(bit2)
		if b.Test(bit2) == originalState {
			t.Errorf("Flip(%d) failed: bit state should have changed", bit2)
		}

		// Verify bit operations don't affect other bits unless expected
		if bit1 != bit2 {
			// After clearing bit1 and flipping bit2, bit1 should still be clear
			if b.Test(bit1) {
				t.Errorf("Operations affected unrelated bit %d", bit1)
			}
		}
	})
}

// FuzzRange tests range operations like FlipRange
func FuzzRange(f *testing.F) {
	// Add seed corpus
	f.Add(uint(0), uint(10))
	f.Add(uint(32), uint(96))
	f.Add(uint(63), uint(129))

	f.Fuzz(func(t *testing.T, start, end uint) {
		// Ensure start <= end and reasonable bounds
		if start > end {
			start, end = end, start
		}
		if end > maxBitPosition {
			return
		}

		b := New(0)

		// Test FlipRange
		b.FlipRange(start, end)

		// Verify all bits in range are set (since we started with empty bitset)
		for i := start; i < end; i++ {
			if !b.Test(i) {
				t.Errorf("FlipRange(%d, %d) failed: bit %d should be set", start, end, i)
			}
		}

		// Verify bits outside range are not affected
		if start > 0 && b.Test(start-1) {
			t.Errorf("FlipRange(%d, %d) affected bit outside range: %d", start, end, start-1)
		}
		if end < b.Len() && b.Test(end) {
			t.Errorf("FlipRange(%d, %d) affected bit outside range: %d", start, end, end)
		}

		// Test FlipRange again (should clear all bits in range)
		b.FlipRange(start, end)
		for i := start; i < end; i++ {
			if b.Test(i) {
				t.Errorf("Second FlipRange(%d, %d) failed: bit %d should be clear", start, end, i)
			}
		}
	})
}

// FuzzSetOperations tests set operations like Union, Intersection, Difference
func FuzzSetOperations(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0x0F}, []byte{0xF0})
	f.Add([]byte{0xFF, 0x00}, []byte{0x00, 0xFF})
	f.Add([]byte{0xAA}, []byte{0x55})

	f.Fuzz(func(t *testing.T, data1, data2 []byte) {
		if len(data1) == 0 || len(data2) == 0 || len(data1) > maxBytes || len(data2) > maxBytes {
			return
		}

		// Create bitsets from byte data
		b1 := New(0)
		b2 := New(0)

		for i, byt := range data1 {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b1.Set(uint(i*8 + bit))
				}
			}
		}

		for i, byt := range data2 {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b2.Set(uint(i*8 + bit))
				}
			}
		}

		// Test Union
		union := b1.Union(b2)
		for i := uint(0); i < union.Len(); i++ {
			expected := b1.Test(i) || b2.Test(i)
			if union.Test(i) != expected {
				t.Errorf("Union failed at bit %d: expected %v, got %v", i, expected, union.Test(i))
			}
		}

		// Test Intersection
		intersection := b1.Intersection(b2)
		for i := uint(0); i < max(b1.Len(), b2.Len()); i++ {
			expected := b1.Test(i) && b2.Test(i)
			if intersection.Test(i) != expected {
				t.Errorf("Intersection failed at bit %d: expected %v, got %v", i, expected, intersection.Test(i))
			}
		}

		// Test Difference
		difference := b1.Difference(b2)
		for i := uint(0); i < max(b1.Len(), b2.Len()); i++ {
			expected := b1.Test(i) && !b2.Test(i)
			if difference.Test(i) != expected {
				t.Errorf("Difference failed at bit %d: expected %v, got %v", i, expected, difference.Test(i))
			}
		}

		// Test SymmetricDifference
		symDiff := b1.SymmetricDifference(b2)
		for i := uint(0); i < max(b1.Len(), b2.Len()); i++ {
			expected := b1.Test(i) != b2.Test(i)
			if symDiff.Test(i) != expected {
				t.Errorf("SymmetricDifference failed at bit %d: expected %v, got %v", i, expected, symDiff.Test(i))
			}
		}

		// Test cardinality functions
		if union.Count() != b1.UnionCardinality(b2) {
			t.Errorf("UnionCardinality mismatch: Union().Count()=%d, UnionCardinality()=%d",
				union.Count(), b1.UnionCardinality(b2))
		}

		if intersection.Count() != b1.IntersectionCardinality(b2) {
			t.Errorf("IntersectionCardinality mismatch: Intersection().Count()=%d, IntersectionCardinality()=%d",
				intersection.Count(), b1.IntersectionCardinality(b2))
		}

		if difference.Count() != b1.DifferenceCardinality(b2) {
			t.Errorf("DifferenceCardinality mismatch: Difference().Count()=%d, DifferenceCardinality()=%d",
				difference.Count(), b1.DifferenceCardinality(b2))
		}

		if symDiff.Count() != b1.SymmetricDifferenceCardinality(b2) {
			t.Errorf("SymmetricDifferenceCardinality mismatch: SymmetricDifference().Count()=%d, SymmetricDifferenceCardinality()=%d",
				symDiff.Count(), b1.SymmetricDifferenceCardinality(b2))
		}
	})
}

// FuzzNavigation tests Next/Previous operations
func FuzzNavigation(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0, 5, 10, 63, 64, 128}, uint(7))
	f.Add([]byte{1, 2, 3}, uint(0))
	f.Add([]byte{100, 200}, uint(150))

	f.Fuzz(func(t *testing.T, data []byte, startPos uint) {
		if len(data) == 0 || len(data) > maxBytes {
			return
		}

		b := New(0)
		for i, byt := range data {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b.Set(uint(i*8 + bit))
				}
			}
		}

		if b.Count() > 0 {
			startPos = startPos % b.Count()
		}

		// Test NextSet
		nextBit, found := b.NextSet(startPos)
		if found {
			// Verify the found bit is actually set
			if !b.Test(nextBit) {
				t.Errorf("NextSet(%d) returned unset bit %d", startPos, nextBit)
			}
			// Verify it's the first set bit >= startPos
			for i := startPos; i < nextBit; i++ {
				if b.Test(i) {
					t.Errorf("NextSet(%d) returned %d, but bit %d is set and comes earlier", startPos, nextBit, i)
				}
			}
		}

		// Test PreviousSet
		prevBit, found := b.PreviousSet(startPos)
		if found {
			// Verify the found bit is actually set
			if !b.Test(prevBit) {
				t.Errorf("PreviousSet(%d) returned unset bit %d", startPos, prevBit)
			}
			// Verify it's the last set bit <= startPos
			for i := prevBit + 1; i <= startPos && i < b.Len(); i++ {
				if b.Test(i) {
					t.Errorf("PreviousSet(%d) returned %d, but bit %d is set and comes later", startPos, prevBit, i)
				}
			}
		}

		// Test NextClear
		nextClear, found := b.NextClear(startPos)
		if found {
			// Verify the found bit is actually clear
			if b.Test(nextClear) {
				t.Errorf("NextClear(%d) returned set bit %d", startPos, nextClear)
			}
		}

		// Test PreviousClear
		prevClear, found := b.PreviousClear(startPos)
		if found {
			// Verify the found bit is actually clear
			if b.Test(prevClear) {
				t.Errorf("PreviousClear(%d) returned set bit %d", startPos, prevClear)
			}
		}
	})
}

// FuzzShift tests shift operations
func FuzzShift(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0, 1, 2, 63, 64}, uint(1))
	f.Add([]byte{10, 20, 30}, uint(5))
	f.Add([]byte{100}, uint(32))

	f.Fuzz(func(t *testing.T, data []byte, shiftAmount uint) {
		if len(data) == 0 || len(data) > maxBytes || shiftAmount > maxBitPosition {
			return
		}

		b1 := New(0)
		var rawBits []uint
		for i, byt := range data {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b1.Set(uint(i*8 + bit))
					rawBits = append(rawBits, uint(i*8+bit))
				}
			}
		}
		originalCount := b1.Count()
		b1.ShiftLeft(shiftAmount)

		// After left shift, count should remain the same
		if b1.Count() != originalCount {
			t.Errorf("ShiftLeft(%d) changed bit count: expected %d, got %d", shiftAmount, originalCount, b1.Count())
		}

		// Verify bits are shifted correctly
		for _, originalBit := range rawBits {
			newBit := originalBit + shiftAmount
			if newBit >= originalBit { // no overflow
				if !b1.Test(newBit) {
					t.Errorf("ShiftLeft(%d) failed: bit %d should be set after shifting from %d", shiftAmount, newBit, originalBit)
				}
				if originalBit < shiftAmount && b1.Test(originalBit) {
					t.Errorf("ShiftLeft(%d) failed: original bit %d should be clear", shiftAmount, originalBit)
				}
			}
		}

		// Test ShiftRight
		b2 := New(0)
		for _, bit := range rawBits {
			b2.Set(bit)
		}

		b2.ShiftRight(shiftAmount)

		// Verify bits are shifted correctly
		for _, originalBit := range rawBits {
			if originalBit >= shiftAmount {
				newBit := originalBit - shiftAmount
				if !b2.Test(newBit) {
					t.Errorf("ShiftRight(%d) failed: bit %d should be set after shifting from %d (%v)", shiftAmount, newBit, originalBit, rawBits)
				}
			}
		}
	})
}

// FuzzModification tests insert/delete operations
func FuzzModification(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0, 2, 4, 6}, uint(3))
	f.Add([]byte{10, 20, 30}, uint(15))
	f.Add([]byte{63, 64, 65}, uint(64))

	f.Fuzz(func(t *testing.T, data []byte, modifyPos uint) {
		if len(data) == 0 || len(data) > maxBytes {
			return
		}

		b1 := New(0)
		var rawBits []uint
		for i, byt := range data {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b1.Set(uint(i*8 + bit))
					rawBits = append(rawBits, uint(i*8+bit))
				}
			}
		}

		originalCount := b1.Count()

		b1.InsertAt(modifyPos)

		// After insert, all bits >= modifyPos should be shifted right by 1
		// and the bit count should remain the same
		if b1.Count() != originalCount {
			t.Errorf("InsertAt(%d) changed bit count: expected %d, got %d", modifyPos, originalCount, b1.Count())
		}

		// The inserted position should be clear
		if b1.Test(modifyPos) {
			t.Errorf("InsertAt(%d) failed: inserted position should be clear", modifyPos)
		}

		// Test DeleteAt
		b2 := New(0)
		for _, bit := range rawBits {
			b2.Set(bit)
		}

		// Only delete if the position has a set bit
		if modifyPos < b2.Len() && b2.Test(modifyPos) {
			originalCount := b2.Count()
			b2.DeleteAt(modifyPos)

			// After delete, count should decrease by 1
			if b2.Count() != originalCount-1 {
				t.Errorf("DeleteAt(%d) failed: expected count %d, got %d", modifyPos, originalCount-1, b2.Count())
			}
		}
	})
}

// FuzzCopy tests clone and copy operations
func FuzzCopy(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0, 1, 63, 64, 127, 128})
	f.Add([]byte{10, 100, 200})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > maxBytes {
			return
		}

		// Convert byte data to bit positions
		b := New(0)
		var rawBits []uint
		for i, byt := range data {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b.Set(uint(i*8 + bit))
					rawBits = append(rawBits, uint(i*8+bit))
				}
			}
		}

		// Test Clone
		clone := b.Clone()
		if !b.Equal(clone) {
			t.Errorf("Clone() failed: clone is not equal to original")
		}

		// Modify clone and verify original is unchanged
		if len(rawBits) > 0 {
			testBit := rawBits[0]
			clone.Clear(testBit)
			if !b.Test(testBit) {
				t.Errorf("Clone() failed: modifying clone affected original")
			}
		}

		// Test CopyFull
		dest := New(0)
		b.CopyFull(dest)
		if !b.Equal(dest) {
			t.Errorf("CopyFull() failed: destination is not equal to source")
		}

		// Test Copy with different sizes
		smallDest := New(10)
		count := b.Copy(smallDest)
		if count > 10 && smallDest.Len() != 10 {
			t.Errorf("Copy() to smaller destination failed")
		}
	})
}

// Helper function for max
func max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

// FuzzSerialization tests JSON marshaling/unmarshaling and other serialization
func FuzzSerialization(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0, 1, 63, 64})
	f.Add([]byte{100, 200})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > maxBytes {
			return
		}

		b := New(0)
		for i, byt := range data {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b.Set(uint(i*8 + bit))
				}
			}
		}

		// Test JSON marshaling/unmarshaling
		data, err := b.MarshalJSON()
		if err != nil {
			t.Errorf("MarshalJSON() failed: %v", err)
		}

		b2 := New(0)
		err = b2.UnmarshalJSON(data)
		if err != nil {
			t.Errorf("UnmarshalJSON() failed: %v", err)
		}

		if !b.Equal(b2) {
			t.Errorf("JSON round-trip failed: bitsets are not equal")
		}

		// Test binary marshaling/unmarshaling
		binData, err := b.MarshalBinary()
		if err != nil {
			t.Errorf("MarshalBinary() failed: %v", err)
		}

		b3 := New(0)
		err = b3.UnmarshalBinary(binData)
		if err != nil {
			t.Errorf("UnmarshalBinary() failed: %v", err)
		}

		if !b.Equal(b3) {
			t.Errorf("Binary round-trip failed: bitsets are not equal")
		}

		// Test String representation (doesn't need round-trip, just shouldn't panic)
		_ = b.String()
		_ = b.DumpAsBits()
	})
}

// FuzzRandomOperations performs random sequences of operations
func FuzzRandomOperations(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{1, 2, 3, 4, 5}) // operations: 0=Set, 1=Clear, 2=Flip, 3=Test, 4=Count
	f.Add([]byte{1, 1, 1, 4, 5})
	f.Add([]byte{2, 2, 2, 4, 5})

	f.Fuzz(func(t *testing.T, operations []byte) {
		if len(operations) == 0 || len(operations) > maxOperations {
			return
		}

		b := New(0)
		rng := rand.New(rand.NewSource(int64(len(operations))))

		for i, op := range operations {
			pos := uint(rng.Intn(maxBitPosition))

			switch op % 6 {
			case 0: // Set
				b.Set(pos)
				if !b.Test(pos) {
					t.Errorf("Operation %d: Set(%d) failed", i, pos)
				}

			case 1: // Clear
				b.Clear(pos)
				if b.Test(pos) {
					t.Errorf("Operation %d: Clear(%d) failed", i, pos)
				}

			case 2: // Flip
				before := b.Test(pos)
				b.Flip(pos)
				after := b.Test(pos)
				if before == after {
					t.Errorf("Operation %d: Flip(%d) failed", i, pos)
				}

			case 3: // Test - just ensure it doesn't panic
				_ = b.Test(pos)

			case 4: // Count - ensure it's reasonable
				count := b.Count()
				if count > b.Len() {
					t.Errorf("Operation %d: Count() %d exceeds length %d", i, count, b.Len())
				}

			case 5: // NextSet - just ensure it doesn't panic
				_, _ = b.NextSet(pos)
			}
		}

		// Final consistency check
		manualCount := uint(0)
		for i := uint(0); i < b.Len(); i++ {
			if b.Test(i) {
				manualCount++
			}
		}
		if b.Count() != manualCount {
			t.Errorf("Final consistency check failed: Count()=%d, manual count=%d", b.Count(), manualCount)
		}
	})
}

// FuzzCapacityAndGrowth tests capacity management and memory growth
func FuzzCapacityAndGrowth(f *testing.F) {
	// Add seed corpus for growth testing
	f.Add(uint(0), []byte{1, 2, 3})
	f.Add(uint(100), []byte{10, 20, 30})
	f.Add(uint(1000), []byte{50, 100, 200})

	f.Fuzz(func(t *testing.T, initialCapacity uint, growthPattern []byte) {
		if initialCapacity > maxBitPosition || len(growthPattern) == 0 || len(growthPattern) > maxOperations {
			return
		}

		b := New(initialCapacity)

		// Test growth by setting bits at increasing positions
		var rawBits []uint
		for i, growth := range growthPattern {
			// Calculate next bit position
			nextBit := uint(i*256) + uint(growth)
			if nextBit > maxBitPosition {
				break
			}

			b.Set(nextBit)
			rawBits = append(rawBits, nextBit)

			// Verify length grows appropriately
			if b.Len() <= nextBit {
				t.Errorf("Length %d should be greater than bit position %d", b.Len(), nextBit)
			}

			// Verify the bit is actually set
			if !b.Test(nextBit) {
				t.Errorf("Bit %d should be set", nextBit)
			}
		}

		// Test that all previously set bits are still set
		for _, bit := range rawBits {
			if !b.Test(bit) {
				t.Errorf("Previously set bit %d should still be set", bit)
			}
		}

		// Test count matches number of set bits
		if b.Count() != uint(len(rawBits)) {
			t.Errorf("Count %d should equal number of set bits %d", b.Count(), len(rawBits))
		}

		// Test compact operation
		originalCount := b.Count()
		b.Compact()
		if b.Count() != originalCount {
			t.Errorf("Compact should not change count: expected %d, got %d", originalCount, b.Count())
		}

		// Verify all bits are still set after compact
		for _, bit := range rawBits {
			if !b.Test(bit) {
				t.Errorf("Bit %d should still be set after Compact()", bit)
			}
		}
	})
}

// FuzzIterationConsistency tests that iteration methods return consistent results
func FuzzIterationConsistency(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{1, 5, 10, 20, 50, 100})
	f.Add([]byte{0, 63, 64, 127, 128})
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > maxBytes {
			return
		}

		b := New(0)
		var rawBits []uint
		for i, byt := range data {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b.Set(uint(i*8 + bit))
					rawBits = append(rawBits, uint(i*8+bit))
				}
			}
		}

		if len(rawBits) == 0 {
			return
		}

		// Test NextSetMany method
		buffer := make([]uint, len(rawBits)+10) // Extra space
		_, resultBits := b.NextSetMany(0, buffer)

		// Verify NextSetMany found all bits
		if len(resultBits) != len(rawBits) {
			t.Errorf("NextSetMany found %d bits, expected %d", len(resultBits), len(rawBits))
		}

		// Verify bits are in ascending order
		for i := 1; i < len(resultBits); i++ {
			if resultBits[i] <= resultBits[i-1] {
				t.Errorf("NextSetMany results not in ascending order: %d <= %d", resultBits[i], resultBits[i-1])
			}
		}

		// Test AsSlice method
		sliceBuffer := make([]uint, b.Count())
		resultSlice := b.AsSlice(sliceBuffer)

		// Verify AsSlice matches NextSetMany
		if len(resultSlice) != len(resultBits) {
			t.Errorf("AsSlice length %d doesn't match NextSetMany length %d", len(resultSlice), len(resultBits))
		}

		for i := range resultSlice {
			if i >= len(resultBits) {
				break
			}
			if resultSlice[i] != resultBits[i] {
				t.Errorf("AsSlice[%d]=%d doesn't match NextSetMany[%d]=%d", i, resultSlice[i], i, resultBits[i])
			}
		}

		// Test manual iteration with NextSet
		var manualBits []uint
		for i := uint(0); i < b.Len(); {
			next, found := b.NextSet(i)
			if !found {
				break
			}
			manualBits = append(manualBits, next)
			i = next + 1
		}

		// Verify manual iteration matches other methods
		if len(manualBits) != len(rawBits) {
			t.Errorf("Manual NextSet iteration found %d bits, expected %d", len(manualBits), len(rawBits))
		}

		for i := range manualBits {
			if i >= len(resultBits) {
				break
			}
			if manualBits[i] != resultBits[i] {
				t.Errorf("Manual iteration[%d]=%d doesn't match NextSetMany[%d]=%d", i, manualBits[i], i, resultBits[i])
			}
		}
	})
}

// FuzzStringRepresentations tests string and dump operations
func FuzzStringRepresentations(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{0xFF, 0x00, 0xFF}) // Alternating pattern
	f.Add([]byte{0xAA, 0x55, 0xAA}) // Checkerboard pattern
	f.Add([]byte{0x01, 0x02, 0x04}) // Sparse bits

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 || len(data) > maxBytes {
			return
		}

		b := New(0)
		var setBitsCount uint
		for i, byt := range data {
			for bit := 0; bit < 8; bit++ {
				if byt&(1<<bit) != 0 {
					b.Set(uint(i*8 + bit))
					setBitsCount++
				}
			}
		}

		if setBitsCount == 0 {
			b.Set(0) // Ensure at least one bit is set
		}

		// Test String() method - should not panic
		str := b.String()
		if len(str) == 0 {
			t.Error("String() should return non-empty string for non-empty bitset")
		}

		// Test DumpAsBits() method - should not panic
		bitStr := b.DumpAsBits()
		if len(bitStr) == 0 {
			t.Error("DumpAsBits() should return non-empty string for non-empty bitset")
		}

		// Test that DumpAsBits contains the right number of '1' characters
		oneCount := 0
		for _, char := range bitStr {
			if char == '1' {
				oneCount++
			}
		}

		if uint(oneCount) != b.Count() {
			t.Errorf("DumpAsBits() contains %d '1' characters, but Count() is %d", oneCount, b.Count())
		}

		// Test Bytes() method
		bytes := b.Bytes()
		if len(bytes) == 0 && b.Count() > 0 {
			t.Error("Bytes() should return non-empty slice for non-empty bitset")
		}

		// Create a new bitset from the bytes using From constructor and verify similarity
		if len(bytes) > 0 {
			b2 := From(bytes)
			// Check that they have the same set bits, accounting for potential length differences
			maxLen := b.Len()
			if b2.Len() < maxLen {
				maxLen = b2.Len()
			}
			for i := uint(0); i < maxLen; i++ {
				if b.Test(i) != b2.Test(i) {
					t.Errorf("Bit %d differs between original (%t) and reconstructed (%t)", i, b.Test(i), b2.Test(i))
				}
			}
		}
	})
}
