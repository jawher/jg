package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

type tokenClass string

const (
	tkObjStart   tokenClass = "{"
	tkObjEnd     tokenClass = "}"
	tkArrStart   tokenClass = "["
	tkArrEnd     tokenClass = "]"
	tkDot        tokenClass = "."
	tkLiteral    tokenClass = "LITERAL"
	tkRawLiteral tokenClass = ":LITERAL"
	tkAssign     tokenClass = "="
	tkEOF        tokenClass = "$"
)

type token struct {
	class tokenClass
	value string
	pos   int
}

type lexer struct {
	input string
	start int
	pos   int
	width int
}

func newLexer(input string) *lexer {
	return &lexer{
		input: input,
		pos:   0,
		width: 0,
	}
}

func (lx *lexer) next() token {
	for {
		for lx.peek() == ' ' {
			lx.pop()
		}
		lx.drop()

		r := lx.pop()
		switch {
		case r == eof:
			return lx.emit(tkEOF)
		case r == '{':
			return lx.emit(tkObjStart)
		case r == '}':
			return lx.emit(tkObjEnd)
		case r == '[':
			return lx.emit(tkArrStart)
		case r == ']':
			return lx.emit(tkArrEnd)
		case r == '.':
			return lx.emit(tkDot)
		case r == '=':
			return lx.emit(tkAssign)
		case r == '"':
			return lx.lexString()
		case r == ':':
			lx.drop()
			return lx.lexRawLiteral()
		default:
			return lx.lexLiteral()
		}
	}
}

var notOkInLiteral = []rune{eof, ' ', '\t', '{', '}', '[', ']', '=', '\n', '.'}

func (lx *lexer) lexLiteral() token {
	for notIn(lx.peek(), notOkInLiteral) {
		lx.pop()
	}

	return lx.emit(tkLiteral)
}

var notOkInRawLiteral = []rune{eof, ' ', '\t', '{', '}', '[', ']', '=', '\n'}

func (lx *lexer) lexRawLiteral() token {
	for notIn(lx.peek(), notOkInRawLiteral) {
		lx.pop()
	}

	return lx.emit(tkRawLiteral)
}

func (lx *lexer) lexString() token {
	lx.drop() // get rid of the opening quotes "
	var buffer bytes.Buffer
	for {
		r := lx.pop()
		switch r {
		case eof:
			panic("unclosed string")
		case '\\':
			if lx.peek() == '"' {
				buffer.WriteRune(lx.pop())
			} else {
				buffer.WriteRune('\\')
			}
		case '"':
			return lx.emitV(tkLiteral, buffer.String())
		default:
			buffer.WriteRune(r)
		}
	}

}

func notIn(needle rune, haystack []rune) bool {
	for _, r := range haystack {
		if needle == r {
			return false
		}
	}
	return true
}

const eof = -1

func (l *lexer) pop() rune {
	r, w := l.nextRune()
	l.width = w
	l.pos += w
	return r
}

func (l *lexer) peek() rune {
	r, _ := l.nextRune()
	return r
}

func (l *lexer) nextRune() (rune, int) {
	if l.pos >= len(l.input) {
		return eof, 0
	}
	return utf8.DecodeRuneInString(l.input[l.pos:])
}

func (l *lexer) push() {
	l.pos -= l.width
	l.width = 0
}

func (l *lexer) drop() {
	l.start = l.pos
}

func (l *lexer) matched() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func (l *lexer) emit(class tokenClass) token {
	res := token{
		class: class,
		value: l.input[l.start:l.pos],
		pos:   l.start,
	}
	l.start = l.pos
	return res
}

func (l *lexer) emitV(class tokenClass, v string) token {
	res := token{
		class: class,
		value: v,
		pos:   l.start,
	}
	l.start = l.pos
	return res
}
