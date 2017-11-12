# bitset

## Note: This is a fork of [github.com/willf/bitset](http://github.com/willf/bitset) that provides a read only bitset that can be instantiated from a mmap'd bytes ref and also limits the scope of the API.

*Go language library to map between non-negative integers and boolean values*

[![Master Build Status](https://secure.travis-ci.org/m3db/bitset.png?branch=master)](https://travis-ci.org/m3db/bitset?branch=master)
[![Master Coverage Status](https://coveralls.io/repos/m3db/bitset/badge.svg?branch=master&service=github)](https://coveralls.io/github/m3db/bitset?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/m3db/bitset)](https://goreportcard.com/report/github.com/m3db/bitset)
[![GoDoc](https://godoc.org/github.com/m3db/bitset?status.svg)](http://godoc.org/github.com/m3db/bitset)


## Description

Package bitset implements bitsets, a mapping between non-negative integers and boolean values.
It should be more efficient than map[uint] bool.

It provides methods for setting, clearing and testing individual integers.

BitSets are expanded to the size of the largest set bit; the memory allocation is approximately Max bits, where Max is the largest set bit. BitSets are never shrunk. On creation, a hint can be given for the number of bits that will be used.

### Example use:

```go
package main

import (
	"fmt"
	"math/rand"

	"github.com/m3db/bitset"
)

func main() {
	fmt.Printf("Hello from BitSet!\n")
	var b bitset.BitSet
	// play some Go Fish
	for i := 0; i < 100; i++ {
		card1 := uint(rand.Intn(52))
		card2 := uint(rand.Intn(52))
		b.Set(card1)
		if b.Test(card2) {
			fmt.Println("Go Fish!")
		}
		b.Clear(card1)
	}
}
```

As an alternative to BitSets, one should check out the 'big' package, which provides a (less set-theoretical) view of bitsets.

Godoc documentation is at: https://godoc.org/github.com/m3db/bitset

## Installation

```bash
go get github.com/m3db/bitset
```

## Running all tests

Before committing the code, please check if it passes all tests using (note: this will install some dependencies):
```bash
make qa
```
