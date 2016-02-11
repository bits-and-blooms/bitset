// +build amd64,!appengine

package bitset

// *** the following functions are defined in popcnt_amd64.s

//go:noescape

func hasAsm() bool

//go:noescape

func popcntSliceAsm(s []uint64) uint64

//go:noescape

func popcntMaskSliceAsm(s, m []uint64) uint64

//go:noescape

func popcntAndSliceAsm(s, m []uint64) uint64

//go:noescape

func popcntOrSliceAsm(s, m []uint64) uint64

//go:noescape

func popcntXorSliceAsm(s, m []uint64) uint64

// useAsm is a flag used to select the GO or ASM implementation of the popcnt function
var useAsm = hasAsm()
