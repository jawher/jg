package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func oneFieldObjGen(field string, value interface{}) Generator {
	return NewObj().Add(field, Value{value: value})
}

func TestParseFieldGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input:    `a=b`,
			expected: oneFieldObjGen("a", "b"),
		},
		{
			input:    `a=42`,
			expected: oneFieldObjGen("a", "42"),
		},
		{
			input:    `a=true`,
			expected: oneFieldObjGen("a", "true"),
		},
		{
			input:    `a=false`,
			expected: oneFieldObjGen("a", "false"),
		},
		{
			input:    `a=null`,
			expected: oneFieldObjGen("a", "null"),
		},

		{
			input:    `a=:42`,
			expected: oneFieldObjGen("a", int64(42)),
		},
		{
			input:    `a=:true`,
			expected: oneFieldObjGen("a", true),
		},
		{
			input:    `a=:false`,
			expected: oneFieldObjGen("a", false),
		},
		{
			input:    `a=:null`,
			expected: oneFieldObjGen("a", nil),
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := Parse(cas.input)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestParseObjectGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input: `a={b=c}`,
			expected: Obj{
				fields: Fields{
					"a": Obj{
						fields: Fields{
							"b": Value{value: "c"},
						},
					},
				},
			},
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := Parse(cas.input)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestParseDotObjectGenerator(t *testing.T) {
	testCases := []struct {
		input    string
		expected Generator
	}{
		{
			input: `a."b.b".c=d`,
			expected: Obj{
				fields: Fields{
					"a": Obj{
						fields: Fields{
							"b.b": Obj{
								fields: Fields{
									"c": Value{value: "d"},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `parent.child1=value1 parent.child2=value2`,
			expected: Obj{
				fields: Fields{
					"parent": Obj{
						fields: Fields{
							"child1": Value{value: "value1"},
							"child2": Value{value: "value2"},
						},
					},
				},
			},
		},
	}

	for _, cas := range testCases {
		t.Logf("Testing input: %s", cas.input)

		ast, err := Parse(cas.input)

		require.NoError(t, err)
		require.Equal(t, []Generator{cas.expected}, ast)
	}
}

func TestComplexParse(t *testing.T) {
	expected := Obj{
		fields: Fields{
			"id":      Value{value: "42"},
			"enabled": Value{value: true},
			"caller": Obj{
				fields: Fields{
					"gender": Obj{
						fields: Fields{
							"code": Value{value: int64(1)},
						},
					},
				},
			},
			"customer": Obj{
				fields: Fields{
					"name": Value{value: "Geralt"},
					"age":  Value{value: "86"},
					"address": Obj{
						fields: Fields{
							"zip": Value{value: "75018"},
						},
					},
				},
			},
		},
	}

	ast, err := Parse(`id = 42 caller.gender.code = :1  customer={name = "Geralt" age  = 86 address.zip = 75018 } enabled = :true`)

	require.NoError(t, err)
	require.Equal(t, []Generator{expected}, ast)
}
