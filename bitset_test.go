// Copyright 2014 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file tests bit sets

package bitset

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBitSet(t *testing.T) {
	set := NewBitSet(0)

	total := uint(16 * 64)
	for i := uint(0); i < total; i++ {
		if i%3 == 0 || i%5 == 0 {
			set.Set(i)
		}
	}

	for i := uint(0); i < 2*total; i++ {
		if i%5 == 0 {
			set.Clear(i)
		}
	}

	for i := uint(0); i < total; i++ {
		if i%3 == 0 && i%5 != 0 {
			assert.True(t, set.Test(i))
		} else {
			assert.False(t, set.Test(i))
		}
	}

	set.ClearAll()

	for i := uint(0); i < 2*total; i++ {
		assert.False(t, set.Test(i))
	}
}

func TestReadOnlyBitSet(t *testing.T) {
	set := NewBitSet(0)

	total := uint(16 * 64)
	for i := uint(0); i < total; i++ {
		if i%3 == 0 || i%5 == 0 {
			set.Set(i)
		}
	}

	// Derive read only set
	readOnlySet, err := set.ReadOnly()
	require.NoError(t, err)

	for i := uint(0); i < 2*total; i++ {
		if i < total && (i%3 == 0 || i%5 == 0) {
			assert.True(t, readOnlySet.Test(i))
		} else {
			assert.False(t, readOnlySet.Test(i))
		}
	}

	// Write out read only set and re-verify
	buf := bytes.NewBuffer(nil)
	require.NoError(t, readOnlySet.Write(buf))

	readOnlySet = NewReadOnlyBitSet(buf.Bytes())
	for i := uint(0); i < 2*total; i++ {
		if i < total && (i%3 == 0 || i%5 == 0) {
			assert.True(t, readOnlySet.Test(i))
		} else {
			assert.False(t, readOnlySet.Test(i))
		}
	}
}
