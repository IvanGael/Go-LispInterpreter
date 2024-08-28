package main

import (
	"strconv"
	"strings"
	"unicode"
)

// Token types
const (
	FORMAT                = "format"
	PLUS                  = "+"
	MINUS                 = "-"
	STAR                  = "*"
	SLASH                 = "/"
	PERCENT               = "%"
	LESS_THAN             = "<"
	LESS_OR_EQUAL_THAN    = "<="
	GREATER_THAN          = ">"
	GREATER_OR_EQUAL_THAN = ">="
	EQUAL                 = "="
	IF                    = "if"
	DEFUN                 = "defun"
	LAMBDA                = "lambda"
	LET                   = "let"
	AND                   = "and"
	OR                    = "or"
	NOT                   = "not"
	LIST                  = "list"
	CAR                   = "car"
	CDR                   = "cdr"
	CONS                  = "cons"
	LENGTH                = "length"
	APPEND                = "append"
	POW                   = "pow"
	SQRT                  = "sqrt"
	CONCAT                = "concat"
	SUBSTRING             = "substring"
	IS_NUMBER             = "isNumber"
	IS_STRING             = "isString"
	READ                  = "read"
	PRINT                 = "print"
	OPEN_BRACKET          = '('
	CLOSE_BRACKET         = ')'
	DOUBLE_QUOTE          = '"'
	EMPTY_STRING          = " "
	DOUBLE_ANTI_SLASH     = '\\'
	ANTI_SLASH_N          = '\n'
	DOT                   = "."
	TRUE                  = "true"
	FALSE                 = "false"
	NIL                   = "nil"
	T                     = "t"
	NUMBER                = "NUMBER"
	FLOAT                 = "FLOAT"
	STRING                = "STRING"
	EOF                   = "EOF"
	IDENTIFIER            = "IDENTIFIER"
	BOOLEAN               = "BOOLEAN"
	FUNCTION              = "FUNCTION"
)

type Token struct {
	Type   string
	Value  string
	Line   int
	Column int
}

// Tokenize splits the input string into tokens
func Tokenize(input string) []Token {
	var tokens []Token
	var token strings.Builder
	inString := false
	escapeNext := false
	line, column := 1, 1

	for _, char := range input {
		switch {
		case unicode.IsSpace(char):
			if !inString && token.Len() > 0 {
				tokens = append(tokens, createToken(token.String(), line, column-token.Len()))
				token.Reset()
			} else if inString {
				token.WriteRune(char)
			}
			if char == ANTI_SLASH_N {
				line++
				column = 1
			} else {
				column++
			}
		case char == OPEN_BRACKET || char == CLOSE_BRACKET:
			if inString {
				token.WriteRune(char)
			} else {
				if token.Len() > 0 {
					tokens = append(tokens, createToken(token.String(), line, column-token.Len()))
					token.Reset()
				}
				tokens = append(tokens, Token{Type: string(char), Value: string(char), Line: line, Column: column})
			}
			column++
		case char == DOUBLE_QUOTE:
			if inString && !escapeNext {
				inString = false
				tokens = append(tokens, Token{Type: STRING, Value: token.String(), Line: line, Column: column - token.Len()})
				token.Reset()
			} else {
				inString = true
			}
			escapeNext = false
			column++
		case char == DOUBLE_ANTI_SLASH:
			if inString && !escapeNext {
				escapeNext = true
			} else {
				token.WriteRune(char)
			}
			column++
		default:
			token.WriteRune(char)
			column++
		}
	}

	if token.Len() > 0 {
		tokens = append(tokens, createToken(token.String(), line, column-token.Len()))
	}

	return tokens
}

func createToken(value string, line, column int) Token {
	tokenType := IDENTIFIER
	switch value {
	case TRUE, FALSE:
		tokenType = BOOLEAN
	case NIL:
		tokenType = NIL
	default:
		if _, err := strconv.ParseFloat(value, 64); err == nil {
			if strings.Contains(value, DOT) {
				tokenType = FLOAT
			} else {
				tokenType = NUMBER
			}
		}
	}
	return Token{Type: tokenType, Value: value, Line: line, Column: column}
}
