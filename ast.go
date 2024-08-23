package main

import (
	"strconv"
	"strings"
)

// LispValue represents a value
type LispValue interface {
	String() string
}

// LispAtom represents an atomic value (symbol)
type LispAtom struct {
	Value string
}

func (a *LispAtom) String() string {
	return a.Value
}

// LispNumber represents a numeric value
type LispNumber struct {
	Value int
}

func (n *LispNumber) String() string {
	return strconv.Itoa(n.Value)
}

// LispFloat represents a float value
type LispFloat struct {
	Value float64
}

func (f *LispFloat) String() string {
	return strconv.FormatFloat(f.Value, 'f', -1, 64)
}

// LispString represents a string value
type LispString struct {
	Value string
}

func (s *LispString) String() string {
	return "\"" + s.Value + "\""
}

// LispList represents a list of Lisp values
type LispList struct {
	Elements []LispValue
}

func (l *LispList) String() string {
	var sb strings.Builder
	sb.WriteString(string(OPEN_BRACKET))
	for i, elem := range l.Elements {
		sb.WriteString(elem.String())
		if i < len(l.Elements)-1 {
			sb.WriteString(EMPTY_STRING)
		}
	}
	sb.WriteString(string(CLOSE_BRACKET))
	return sb.String()
}

// LispFunction represents a user-defined function
type LispFunction struct {
	Name   *LispAtom
	Params []LispValue
	Body   LispValue
	Env    Environment
}

func (f *LispFunction) String() string {
	if f.Name != nil {
		return strings.ToUpper(f.Name.Value)
	}
	return "FUNCTION"
}

// LispBoolean represents a boolean value
type LispBoolean struct {
	Value bool
}

func (b *LispBoolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// LispNil represents a nil/null value
type LispNil struct{}

func (n *LispNil) String() string {
	return "nil"
}
