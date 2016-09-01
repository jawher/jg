package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValue(t *testing.T) {
	cases := []interface{}{
		1,
		1.5,
		nil,
		true,
		"a",
	}

	otherGens := []Generator{
		NewObj(),
		Arr{},
		Value{value: 42},
	}

	for _, cas := range cases {
		v := Value{value: cas}

		require.Equal(t, cas, v.Gen(nil))

		for _, otherGen := range otherGens {
			// Merge of a value wih another generator should always return the other
			require.Equal(t, otherGen, v.Merge(otherGen))
		}
	}
}

func TestVar(t *testing.T) {
	v := Var{varName: "a"}

	require.Equal(t, "b", v.Gen(map[string]Any{"a": "b"}))

	otherGens := []Generator{
		NewObj(),
		Arr{},
		Value{value: 42},
	}
	for _, otherGen := range otherGens {
		// Merge of a var wih another generator should always return the other
		require.Equal(t, otherGen, v.Merge(otherGen))
	}

}

func TestArr(t *testing.T) {
	values := []Any{6, true, " aloha"}

	g := Arr{}
	for _, v := range values {
		g.Add(Value{value: v})
	}

	require.Equal(t, values, g.Gen(nil))

	otherGens := []Generator{
		NewObj(),
		Arr{},
		Value{value: 42},
	}

	for _, otherGen := range otherGens {
		// Merge of an array wih another generator should always return the other
		require.Equal(t, otherGen, g.Merge(otherGen))
	}
}

func TestObj(t *testing.T) {

	fields := map[string]Any{
		"f1": Value{value: 4},
		"f2": Value{value: true},
		"f3": Value{value: "test"},
	}

	g := NewObj()
	for field, value := range fields {
		g.Add(field, Value{value: value})
	}

	require.Equal(t, fields, g.Gen(nil))
}

func TestObjMerge(t *testing.T) {

	obj1 := NewObj().
		Add("only1", Value{value: 111}).
		Add("common",
			NewObj().Add("sub-only1", Value{value: "s1"}).
				Add("sub-common", Value{value: "com1"}),
		)

	obj2 := NewObj().
		Add("only2", Value{value: 222}).
		Add("common",
			NewObj().Add("sub-only2", Value{value: "s2"}).
				Add("sub-common", Value{value: "com2"}),
		)

	expected := NewObj().
		Add("only1", Value{value: 111}).
		Add("only2", Value{value: 222}).
		Add("common",
			NewObj().
				Add("sub-common", Value{value: "com2"}).
				Add("sub-only1", Value{value: "s1"}).
				Add("sub-only2", Value{value: "s2"}),
		)

	actual := obj1.Merge(obj2)

	require.Equal(t, expected, actual)
}
