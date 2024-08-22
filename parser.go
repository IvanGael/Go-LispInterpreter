package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Parse reads tokens and constructs a Lisp expression tree
func Parse(tokens []string) (LispValue, []string, error) {
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("unexpected EOF while reading")
	}

	token := tokens[0]
	tokens = tokens[1:]

	switch token {
	case string(OPEN_BRACKET):
		var elements []LispValue
		for len(tokens) > 0 && tokens[0] != string(CLOSE_BRACKET) {
			var elem LispValue
			var err error
			elem, tokens, err = Parse(tokens)
			if err != nil {
				return nil, nil, err
			}
			elements = append(elements, elem)
		}
		if len(tokens) == 0 {
			return nil, nil, fmt.Errorf("unexpected EOF while reading")
		}
		tokens = tokens[1:]
		return &LispList{Elements: elements}, tokens, nil
	default:
		if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
			return &LispString{Value: token[1 : len(token)-1]}, tokens, nil
		}
		if num, err := strconv.Atoi(token); err == nil {
			return &LispNumber{Value: num}, tokens, nil
		}
		return &LispAtom{Value: token}, tokens, nil
	}
}
