package main

import (
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
	OPEN_BRACKET          = '('
	CLOSE_BRACKET         = ')'
	DOUBLE_QUOTE          = '"'
	EMPTY_STRING          = " "
	DOUBLE_ANTI_SLASH     = '\\'
	NUMBER                = "NUMBER"
	STRING                = "STRING"
	EOF                   = "EOF"
	IDENTIFIER            = "IDENTIFIER"
)

// Tokenize splits the input string into tokens
func Tokenize(input string) []string {
	var tokens []string
	var token strings.Builder
	inString := false
	escapeNext := false
	for _, char := range input {
		switch {
		case unicode.IsSpace(char):
			if !inString && token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			} else if inString {
				token.WriteRune(char)
			}
		case char == OPEN_BRACKET || char == CLOSE_BRACKET:
			if inString {
				token.WriteRune(char)
			} else {
				if token.Len() > 0 {
					tokens = append(tokens, token.String())
					token.Reset()
				}
				tokens = append(tokens, string(char))
			}
		case char == DOUBLE_QUOTE:
			if inString && !escapeNext {
				inString = false
				tokens = append(tokens, "\""+token.String()+"\"")
				token.Reset()
			} else {
				inString = true
			}
			escapeNext = false
		case char == DOUBLE_ANTI_SLASH:
			if inString && !escapeNext {
				escapeNext = true
			} else {
				token.WriteRune(char)
			}
		default:
			token.WriteRune(char)
		}
	}
	if token.Len() > 0 {
		tokens = append(tokens, token.String())
	}
	return tokens
}
