package bitset

import "testing"

func TestExtractToReusedDestination(t *testing.T) {
	for _, tc := range []struct {
		name           string
		src, mask, dst *BitSet
		alias          string
	}{
		{"single word", New(64).Set(1), New(64).Set(1).Set(4), New(128).Set(0).Set(1).Set(80), ""},
		{"zero source", New(128), New(128).SetAll(), New(192).SetAll(), ""},
		{"full words", New(128).Set(0).Set(70), New(128).SetAll(), New(192).SetAll(), ""},
		{"cross word", New(128).Set(1).Set(70), extractReuseMask(), New(192).SetAll(), ""},
		{"source destination", New(128).Set(1).Set(70).Set(127), extractReuseMask(), nil, "source"},
		{"mask destination", New(128).Set(1).Set(70), extractReuseMask(), nil, "mask"},
		{"short source", New(64).Set(1), New(256).Set(1).Set(70).Set(200), New(256).SetAll(), ""},
		{"short source zeros", New(64).Set(0), New(256).Set(1).Set(70).Set(200), New(256).SetAll(), ""},
		{"alias growth", New(64).SetAll(), extractGrowthMask(), nil, "source"},
		{"nonalias growth", New(64).SetAll(), extractGrowthMask(), New(64).SetAll(), ""},
		{"zero logical source with words", FromWithLength(0, []uint64{0x2}), New(64).Set(1).Set(4), New(128).SetAll(), ""},
		{"extra source words", FromWithLength(1, []uint64{1, 1 << 6}), New(128).Set(0).Set(70), New(128).SetAll(), ""},
		{"zero logical mask with words", New(64).Set(1), FromWithLength(0, []uint64{0x12}), New(128).SetAll(), ""},
		{"extra mask words", New(128).Set(70), FromWithLength(1, []uint64{1, 1 << 6}), New(128).SetAll(), ""},
		{"source alias shrink", FromWithLength(1, []uint64{0, 1}), From([]uint64{1, 1}), nil, "source"},
		{"source nonalias shrink", FromWithLength(1, []uint64{0, 1}), From([]uint64{1, 1}), FromWithLength(1, []uint64{0, 1}), ""},
		{"mask alias shrink", From([]uint64{1, 0}), FromWithLength(0, []uint64{1, 1}), nil, "mask"},
		{"mask nonalias shrink", From([]uint64{1, 0}), FromWithLength(0, []uint64{1, 1}), FromWithLength(0, []uint64{1, 1}), ""},
		{"both alias shrink", FromWithLength(1, []uint64{1, 1}), nil, nil, "both"},
		{"both nonalias shrink", FromWithLength(1, []uint64{1, 1}), FromWithLength(1, []uint64{1, 1}), FromWithLength(1, []uint64{1, 1}), ""},
		{"empty source", New(0), extractGrowthMask(), New(256).SetAll(), ""},
		{"empty mask", New(128).Set(1), New(128), New(192).SetAll(), ""},
		{"grow destination", New(128).Set(1).Set(70), extractReuseMask(), New(0), ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dst := tc.dst
			if tc.alias == "both" {
				tc.mask, dst = tc.src, tc.src
			}
			if tc.alias == "source" {
				dst = tc.src
			}
			if tc.alias == "mask" {
				dst = tc.mask
			}
			// Preserve all backing words, including extra words supplied with FromWithLength.
			src := From(append([]uint64(nil), tc.src.Words()...))
			mask := From(append([]uint64(nil), tc.mask.Words()...))
			want := FromWithLength(dst.Len(), append([]uint64(nil), dst.Words()...))
			if src.Len() != 0 && mask.Count() > want.Len() {
				want.Set(mask.Count() - 1).Clear(mask.Count() - 1)
			}
			var out uint
			// The legacy operation truncates the mask at the source's final word.
			end := src.Len()
			for in := uint(0); in < mask.Len() && in < end; in++ {
				if mask.Test(in) {
					want.SetTo(out, src.Test(in))
					out++
				}
			}
			tc.src.ExtractTo(tc.mask, dst)
			if !dst.Equal(want) {
				t.Errorf("got %v, want %v", dst, want)
			}
			if dst.Len() != want.Len() {
				t.Errorf("length got %d, want %d", dst.Len(), want.Len())
			}
		})
	}
}

func extractReuseMask() *BitSet {
	mask := New(128)
	for i := uint(0); i < 128; i++ {
		if i%3 != 0 {
			mask.Set(i)
		}
	}
	return mask
}

func extractGrowthMask() *BitSet {
	mask := New(256).Set(1)
	for i := uint(64); i < 256; i++ {
		mask.Set(i)
	}
	return mask
}
