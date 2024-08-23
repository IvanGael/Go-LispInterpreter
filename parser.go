package main

import (
	"strconv"
)

// Parse reads tokens and constructs a Lisp expression tree
func Parse(tokens []Token) (LispValue, []Token, error) {
	if len(tokens) == 0 {
		return nil, nil, &LispError{Message: "unexpected EOF while reading", Line: 0, Column: 0}
	}

	token := tokens[0]
	tokens = tokens[1:]

	switch token.Type {
	case string(OPEN_BRACKET):
		var elements []LispValue
		for len(tokens) > 0 && tokens[0].Type != string(CLOSE_BRACKET) {
			var elem LispValue
			var err error
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
		return &LispList{Elements: elements}, tokens, nil
	case STRING:
		return &LispString{Value: token.Value}, tokens, nil
	case NUMBER:
		num, _ := strconv.Atoi(token.Value)
		return &LispNumber{Value: num}, tokens, nil
	case FLOAT:
		num, _ := strconv.ParseFloat(token.Value, 64)
		return &LispFloat{Value: num}, tokens, nil
	case BOOLEAN:
		return &LispBoolean{Value: token.Value == "true"}, tokens, nil
	case NIL:
		return &LispNil{}, tokens, nil
	default:
		return &LispAtom{Value: token.Value}, tokens, nil
	}
}
