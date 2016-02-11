// +build !amd64 appengine

// This is compiled if the popcnt_amd64.go is not compiled

package bitset

// useAsm is a flag used to select the GO or ASM implementation of the popcnt function
var useAsm = false
