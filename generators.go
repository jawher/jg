package main

import (
	"fmt"
	"strconv"
)

type Generator interface {
	Gen(map[string]Any) Any
	Merge(Generator) Generator
}

type Any interface{}

type Value struct {
	value Any
}

func (v Value) Gen(substs map[string]Any) Any {
	return v.value
}

func (v Value) Merge(g Generator) Generator {
	// Values cannot be merged. Return the new one
	return g
}

type Var struct {
	varName string
}

func (v Var) Gen(substs map[string]Any) Any {
	return substs[v.varName]
}

func (v Var) Merge(g Generator) Generator {
	// Vars cannot be merged. Return the new one
	return g
}

type Fields map[string]Generator

type Obj struct {
	fields Fields
}

func NewObj() Obj {
	return Obj{
		fields: Fields{},
	}
}

func (obj Obj) Merge(g Generator) Generator {
	switch g := g.(type) {
	case Obj:
		// Objects can be merged together
		res := NewObj()
		for f, v := range obj.fields {
			res.Add(f, v)
		}
		for f, v := range g.fields {
			res.Add(f, v)
		}
		return res
	default:
		// other types, less so, return the new one
		return g
	}
}

func (obj Obj) Gen(substs map[string]Any) Any {
	res := map[string]Any{}
	for field, valueGen := range obj.fields {
		res[field] = valueGen.Gen(substs)
	}
	return res
}

func (obj Obj) Add(field string, value Generator) Obj {
	existingGenerator, found := obj.fields[field]
	if found {
		value = existingGenerator.Merge(value)
	}
	obj.fields[field] = value
	return obj
}

type Arr []Generator

func (arr Arr) Merge(g Generator) Generator {
	// arrays can' t be merged with other generators
	return g
}

func (arr Arr) Gen(substs map[string]Any) Any {
	res := make([]Any, len(arr))
	for idx, elemGen := range arr {
		res[idx] = elemGen.Gen(substs)
	}
	return res
}

func (arr *Arr) Add(g Generator) Generator {
	*arr = append(*arr, g)
	return *arr
}

func parseRawValue(value string) (Any, error) {
	switch {
	case value == "true":
		return true, nil
	case value == "false":
		return false, nil
	case value == "null":
		return nil, nil
	default:
		var v Any
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			v, err = strconv.ParseFloat(value, 64)
		}
		if err != nil {
			return nil, fmt.Errorf("Invalid raw literal %q: isn't any of true, false, null or a numeric", value)
		}
		return v, nil
	}
}
