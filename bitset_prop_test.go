// Copyright 2014 Will Fitzgerald. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitset

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
)

const (
	bitSetPropTestMaxInitLength  = 4194304
	bitSetPropTestMaxValueLength = bitSetPropTestMaxInitLength * 4
)

func TestBitSetProp(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 8
	properties := gopter.NewProperties(parameters)

	prop := commands.Prop(&commands.ProtoCommands{
		NewSystemUnderTestFunc: func(initialState commands.State) commands.SystemUnderTest {
			return initialState
		},
		DestroySystemUnderTestFunc: func(sut commands.SystemUnderTest) {
			state := sut.(bitSetTestState)
			state.set.ClearAll()
		},
		InitialStateGen: gen.
			UIntRange(0, bitSetPropTestMaxInitLength).
			MapResult(func(r *gopter.GenResult) *gopter.GenResult {
				iface, ok := r.Retrieve()
				if !ok {
					return gopter.NewEmptyResult(reflect.PtrTo(reflect.TypeOf(bitSetTestState{})))
				}
				v, ok := iface.(uint)
				if !ok {
					return gopter.NewEmptyResult(reflect.PtrTo(reflect.TypeOf(bitSetTestState{})))
				}
				state := bitSetTestState{
					set:     NewBitSet(v),
					currSet: make(map[uint]struct{}),
				}
				return gopter.NewGenResult(state, gopter.NoShrinker)
			}),
		GenCommandFunc: func(state commands.State) gopter.Gen {
			return gen.OneGenOf(genBitSetTestCommand,
				genBitSetSetCommand, genBitSetClearAllCommand)
		},
	})
	properties.Property("BitSet Set/Test/ClearAll", prop)
	properties.TestingRun(t)
}

type bitSetTestState struct {
	set     *BitSet
	currSet map[uint]struct{}
}

var genBitSetTestCommand = gen.
	UIntRange(0, bitSetPropTestMaxValueLength).
	Map(func(v uint) commands.Command {
		return &commands.ProtoCommand{
			Name: "Test",
			RunFunc: func(q commands.SystemUnderTest) commands.Result {
				state := q.(bitSetTestState)
				return state.set.Test(v)
			},
			NextStateFunc: func(state commands.State) commands.State {
				return state
			},
			PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
				actual := result.(bool)
				_, expected := state.(bitSetTestState).currSet[v]
				if actual == expected {
					return &gopter.PropResult{Status: gopter.PropTrue}
				}
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  fmt.Errorf("expected %d to be true when testing value", v),
				}
			},
		}
	})

var genBitSetSetCommand = gen.
	UIntRange(0, bitSetPropTestMaxValueLength).
	Map(func(v uint) commands.Command {
		return &commands.ProtoCommand{
			Name: "Set",
			RunFunc: func(q commands.SystemUnderTest) commands.Result {
				state := q.(bitSetTestState)
				state.set.Set(v)
				if !state.set.Test(v) {
					return fmt.Errorf("not set, expected set after setting %v", v)
				}
				return nil
			},
			NextStateFunc: func(state commands.State) commands.State {
				s := state.(bitSetTestState)
				s.currSet[v] = struct{}{}
				return s
			},
			PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
				if result == nil {
					return &gopter.PropResult{Status: gopter.PropTrue}
				}
				return &gopter.PropResult{
					Status: gopter.PropFalse,
					Error:  result.(error),
				}
			},
		}
	})

var genBitSetClearAllCommand = gen.Const(&commands.ProtoCommand{
	Name: "ClearAll",
	RunFunc: func(q commands.SystemUnderTest) commands.Result {
		s := q.(bitSetTestState)
		s.set.ClearAll()
		return nil
	},
	NextStateFunc: func(state commands.State) commands.State {
		s := state.(bitSetTestState)
		for k := range s.currSet {
			delete(s.currSet, k)
		}
		return s
	},
	PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
		s := state.(bitSetTestState)
		for _, val := range s.set.values {
			if val != 0 {
				return &gopter.PropResult{
					Status: gopter.PropError,
					Error:  fmt.Errorf("set expected to be empty after clearAll"),
				}
			}
		}
		return &gopter.PropResult{Status: gopter.PropTrue}
	},
})
