package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexerSingleTokens(t *testing.T) {
	cases := []struct {
		input    string
		expected token
	}{
		{"a", token{class: tkLiteral, value: "a", pos: 0}},
		{"$a", token{class: tkVar, value: "a", pos: 1}},
		{"a/b", token{class: tkLiteral, value: "a/b", pos: 0}},
		{`"a b"`, token{class: tkLiteral, value: "a b", pos: 1}},
		{`"a.b"`, token{class: tkLiteral, value: "a.b", pos: 1}},
		{`"a\"b"`, token{class: tkLiteral, value: `a"b`, pos: 1}},
		{`":a"`, token{class: tkLiteral, value: `:a`, pos: 1}},

		{":a", token{class: tkRawLiteral, value: "a", pos: 1}},
		{":true", token{class: tkRawLiteral, value: "true", pos: 1}},
		{":false", token{class: tkRawLiteral, value: "false", pos: 1}},
		{":null", token{class: tkRawLiteral, value: "null", pos: 1}},
		{":1.2345", token{class: tkRawLiteral, value: "1.2345", pos: 1}},
		{":-1.2345e-4", token{class: tkRawLiteral, value: "-1.2345e-4", pos: 1}},

		{"=", token{class: tkAssign, value: "=", pos: 0}},
		{"{", token{class: tkObjStart, value: "{", pos: 0}},
		{"}", token{class: tkObjEnd, value: "}", pos: 0}},
		{"[", token{class: tkArrStart, value: "[", pos: 0}},
		{"]", token{class: tkArrEnd, value: "]", pos: 0}},
		{".", token{class: tkDot, value: ".", pos: 0}},
		{"", token{class: tkEOF, value: "", pos: 0}},
	}

	for _, cas := range cases {
		t.Logf("Testing input: %s", cas.input)
		lx := newLexer(cas.input)
		require.Equal(t, cas.expected, lx.next())
	}
}

func TestLexerMultiTokens(t *testing.T) {
	cases := []struct {
		input    string
		expected []token
	}{
		{"lit1 lit2", []token{
			{class: tkLiteral, value: "lit1", pos: 0},
			{class: tkLiteral, value: "lit2", pos: 5},
			{class: tkEOF, value: "", pos: 9},
		}},
		{"lit1.lit2=", []token{
			{class: tkLiteral, value: "lit1", pos: 0},
			{class: tkDot, value: ".", pos: 4},
			{class: tkLiteral, value: "lit2", pos: 5},
			{class: tkAssign, value: "=", pos: 9},
			{class: tkEOF, value: "", pos: 10},
		}},
		{`"lit1.lit2".lit3=`, []token{
			{class: tkLiteral, value: "lit1.lit2", pos: 1},
			{class: tkDot, value: ".", pos: 11},
			{class: tkLiteral, value: "lit3", pos: 12},
			{class: tkAssign, value: "=", pos: 16},
			{class: tkEOF, value: "", pos: 17},
		}},
		{"lit1=lit2", []token{
			{class: tkLiteral, value: "lit1", pos: 0},
			{class: tkAssign, value: "=", pos: 4},
			{class: tkLiteral, value: "lit2", pos: 5},
			{class: tkEOF, value: "", pos: 9},
		}},
		{`"a=b"=c`, []token{
			{class: tkLiteral, value: "a=b", pos: 1},
			{class: tkAssign, value: "=", pos: 5},
			{class: tkLiteral, value: "c", pos: 6},
			{class: tkEOF, value: "", pos: 7},
		}},
	}

	for _, cas := range cases {
		t.Logf("Testing input: %s", cas.input)
		lx := newLexer(cas.input)
		for _, expected := range cas.expected {
			require.Equal(t, expected, lx.next())
		}

		next := lx.next()
		require.True(t, next.class == tkEOF, "Lexer returned more tokens than expected: %v", next)
	}
}
