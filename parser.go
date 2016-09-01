package main

import (
	"fmt"
	"strings"
)

type parser struct {
	lexer   *lexer
	matched token
	next    token
}

// ParseError is returned if the input cannot be successfuly parsed
type ParseError struct {
	// The original query
	Input string
	// The position where the parsing fails
	Pos int
	// The error message
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Parse error: %s\n%s\n%s^", e.Message, e.Input, strings.Repeat(" ", e.Pos))
}

/*
Parse accepts an input string and returns a list of generators or an error
*/
func ParseGenerators(input string) ([]Generator, error) {
	lexer := newLexer(input)
	return (&parser{
		lexer: lexer,
		next:  lexer.next(),
	}).parseGenerators()
}

/*
Parse accepts an input string and returns a parsed value for a substitution or an error
*/
func ParseSubstValue(input string) (Any, error) {
	lexer := newLexer(input)
	return (&parser{
		lexer: lexer,
		next:  lexer.next(),
	}).parseSubstValue()
}

func (p *parser) parseGenerators() (gen []Generator, err error) {
	defer func() {
		if r := recover(); r != nil {
			gen = nil
			err = ParseError{
				Input:   p.lexer.input,
				Pos:     p.matched.pos,
				Message: fmt.Sprintf("%v", r),
			}
		}
	}()
	gen = p.start()
	if !p.found(tkEOF) {
		p.advance()
		panic("Unexpected input")
	}
	return
}

func (p *parser) parseSubstValue() (res Any, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = nil
			err = ParseError{
				Input:   p.lexer.input,
				Pos:     p.matched.pos,
				Message: fmt.Sprintf("%v", r),
			}
		}
	}()

	switch {
	case p.found(tkLiteral):
		res = p.matched.value
	case p.found(tkRawLiteral):
		rawValue, err := parseRawValue(p.matched.value)
		if err != nil {
			panic("Invalid literal")
		}
		res = rawValue
	case p.found(tkEOF):
		panic("Unexpected end of input")
	default:
		p.advance()
		panic("Unexpected input")
	}

	if !p.found(tkEOF) {
		p.advance()
		panic("Unexpected input")
	}
	return
}

func (p *parser) start() []Generator {
	if p.peek(tkObjStart) || p.peek(tkArrStart) {
		res := []Generator{}
		for {
			switch {
			case p.found(tkObjStart):
				res = append(res, p.obj())
			case p.found(tkArrStart):
				res = append(res, p.arr())
			default:
				return res
			}
		}
	}

	objGen := NewObj()
	for p.found(tkLiteral) {
		field := p.matched.value
		value := p.field(field)
		objGen.Add(field, value)
	}

	return []Generator{objGen}
}

func (p *parser) obj() Generator {
	res := NewObj()
	for p.found(tkLiteral) {
		field := p.matched.value
		value := p.field(field)
		res.Add(field, value)
	}
	p.expect(tkObjEnd)
	return res
}

func (p *parser) arr() Generator {
	res := Arr{}
	for {
		switch {
		case p.found(tkRawLiteral):
			rawValue, err := parseRawValue(p.matched.value)
			if err != nil {
				panic("Invalid literal")
			}
			res.Add(Value{value: rawValue})
		case p.found(tkLiteral):
			if p.peek(tkAssign) || p.peek(tkDot) {
				field := p.matched.value
				value := p.field(field)
				// Add 1-field obj to array
				obj := NewObj()
				obj.Add(field, value)
				res.Add(obj)
			} else {
				res.Add(Value{value: p.matched.value})
			}
		case p.found(tkVar):
			res.Add(Var{varName: p.matched.value})
		case p.found(tkObjStart):
			res.Add(p.obj())
			//Add obj as array elem
		case p.found(tkArrStart):
			res.Add(p.arr())
			// Add array as arr elem
		case p.found(tkArrEnd):
			// return, the array is complete
			return res
		case p.found(tkEOF):
			panic("Unclosed array")
		default:
			p.advance()
			panic("Unexpected input")
		}
	}
}

func (p *parser) field(field string) Generator {

	switch {
	case p.found(tkAssign):
		return p.value()
	case p.found(tkDot):
		p.expect(tkLiteral)
		field := p.matched.value
		value := p.field(field)
		return NewObj().Add(field, value)
	case p.found(tkEOF):
		panic("Unexpected end of input")
	default:
		p.advance()
		panic("Unexpected input")
	}
}

func (p *parser) value() Generator {
	switch {
	case p.found(tkLiteral):
		return Value{value: p.matched.value}
	case p.found(tkVar):
		return Var{varName: p.matched.value}
	case p.found(tkRawLiteral):
		rawValue, err := parseRawValue(p.matched.value)
		if err != nil {
			panic("Invalid literal")
		}
		return Value{value: rawValue}
	case p.found(tkObjStart):
		return p.obj()
	case p.found(tkArrStart):
		return p.arr()
	case p.found(tkEOF):
		panic("Unexpected end of input")
	default:
		p.advance()
		panic("Unexpected input")
	}
}

func (p *parser) expect(class tokenClass) {
	if !p.found(class) {
		p.advance()
		panic(fmt.Sprintf("was expecting %v", class))
	}
}

func (p *parser) peek(class tokenClass) bool {
	return p.next.class == class
}

func (p *parser) found(class tokenClass) bool {
	if p.next.class == class {
		p.matched = p.next
		p.next = p.lexer.next()
		return true
	}
	return false
}

func (p *parser) advance() {
	p.matched = p.next
	p.next = p.lexer.next()
}
