package main

import (
	"strconv"
	"sync"
)

// a cache for parsed expressions
var (
	parseCache     = make(map[string]LispValue)
	parseCacheLock sync.RWMutex
)

// Parse reads tokens and constructs a Lisp expression tree
func Parse(tokens []Token) (LispValue, []Token, error) {
	if len(tokens) == 0 {
		return nil, nil, &LispError{Message: "unexpected EOF while reading", Line: 0, Column: 0}
	}

	// Check cache for parsed expression
	cacheKey := tokensToString(tokens)
	parseCacheLock.RLock()
	if cachedExpr, ok := parseCache[cacheKey]; ok {
		parseCacheLock.RUnlock()
		return cachedExpr, nil, nil
	}
	parseCacheLock.RUnlock()

	token := tokens[0]
	tokens = tokens[1:]

	var result LispValue
	var err error

	switch token.Type {
	case string(OPEN_BRACKET):
		elements := make([]LispValue, 0, 8)
		for len(tokens) > 0 && tokens[0].Type != string(CLOSE_BRACKET) {
			var elem LispValue
			elem, tokens, err = Parse(tokens)
			if err != nil {
				return nil, nil, err
			}
			elements = append(elements, elem)
		}
		if len(tokens) == 0 {
			return nil, nil, &LispError{Message: "unexpected EOF while reading", Line: token.Line, Column: token.Column}
		}
		tokens = tokens[1:]
		result = &LispList{Elements: elements}
	case STRING:
		result = &LispString{Value: token.Value}
	case NUMBER:
		num, _ := strconv.Atoi(token.Value)
		result = &LispNumber{Value: num}
	case FLOAT:
		num, _ := strconv.ParseFloat(token.Value, 64)
		result = &LispFloat{Value: num}
	case BOOLEAN:
		result = &LispBoolean{Value: token.Value == TRUE}
	case NIL:
		result = &LispNil{}
	default:
		result = &LispAtom{Value: token.Value}
	}

	// Cache the parsed expression
	parseCacheLock.Lock()
	parseCache[cacheKey] = result
	parseCacheLock.Unlock()

	return result, tokens, nil
}
