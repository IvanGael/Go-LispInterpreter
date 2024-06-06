package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// LispValue represents a value in our Lisp interpreter
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
	sb.WriteString("(")
	for i, elem := range l.Elements {
		sb.WriteString(elem.String())
		if i < len(l.Elements)-1 {
			sb.WriteString(" ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

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
		case char == '(' || char == ')':
			if inString {
				token.WriteRune(char)
			} else {
				if token.Len() > 0 {
					tokens = append(tokens, token.String())
					token.Reset()
				}
				tokens = append(tokens, string(char))
			}
		case char == '"':
			if inString && !escapeNext {
				inString = false
				tokens = append(tokens, "\""+token.String()+"\"")
				token.Reset()
			} else {
				inString = true
			}
			escapeNext = false
		case char == '\\':
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

// Parse reads tokens and constructs a Lisp expression tree
func Parse(tokens []string) (LispValue, []string, error) {
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("unexpected EOF while reading")
	}

	token := tokens[0]
	tokens = tokens[1:]

	switch token {
	case "(":
		var elements []LispValue
		for len(tokens) > 0 && tokens[0] != ")" {
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
	case ")":
		return nil, nil, fmt.Errorf("unexpected )")
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

// Environment represents the mapping of symbols to their values
type Environment map[string]LispValue

// Eval evaluates a Lisp expression in the given environment
func Eval(env Environment, expr LispValue) (LispValue, error) {
	switch v := expr.(type) {
	case *LispAtom:
		if val, ok := env[v.Value]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("unbound symbol: %s", v.Value)
	case *LispNumber:
		return v, nil
	case *LispString:
		return v, nil
	case *LispList:
		if len(v.Elements) == 0 {
			return v, nil
		}
		fn, ok := v.Elements[0].(*LispAtom)
		if !ok {
			return nil, fmt.Errorf("invalid function call: %v", v.Elements[0])
		}
		args := v.Elements[1:]
		switch fn.Value {
		case "+":
			return builtinAdd(env, args)
		case "-":
			return builtinSub(env, args)
		case "*":
			return builtinMul(env, args)
		case "/":
			return builtinDiv(env, args)
		case "<":
			return builtinLt(env, args)
		case ">":
			return builtinGt(env, args)
		case "=":
			return builtinEq(env, args)
		case "if":
			return builtinIf(env, args)
		case "defun":
			return builtinDefun(env, args)
		default:
			return callFunction(env, fn.Value, args)
		}
	default:
		return nil, fmt.Errorf("unknown expression type: %T", v)
	}
}

// Built-in function implementations

func builtinAdd(env Environment, args []LispValue) (LispValue, error) {
	sum := 0
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		number, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to +: %v", val)
		}
		sum += number.Value
	}
	return &LispNumber{Value: sum}, nil
}

func builtinSub(env Environment, args []LispValue) (LispValue, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("wrong number of arguments to -")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	number, ok := val.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to -: %v", val)
	}
	result := number.Value
	for _, arg := range args[1:] {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		number, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to -: %v", val)
		}
		result -= number.Value
	}
	return &LispNumber{Value: result}, nil
}

func builtinMul(env Environment, args []LispValue) (LispValue, error) {
	product := 1
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		number, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to *: %v", val)
		}
		product *= number.Value
	}
	return &LispNumber{Value: product}, nil
}

func builtinDiv(env Environment, args []LispValue) (LispValue, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("wrong number of arguments to /")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	number, ok := val.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to /: %v", val)
	}
	result := number.Value
	for _, arg := range args[1:] {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		number, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to /: %v", val)
		}
		if number.Value == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result /= number.Value
	}
	return &LispNumber{Value: result}, nil
}

func builtinLt(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to <")
	}
	left, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	right, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	leftNum, ok := left.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", left)
	}
	rightNum, ok := right.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", right)
	}
	if leftNum.Value < rightNum.Value {
		return &LispAtom{Value: "t"}, nil
	}
	return &LispAtom{Value: "nil"}, nil
}

func builtinGt(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to >")
	}
	left, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	right, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	leftNum, ok := left.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", left)
	}
	rightNum, ok := right.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", right)
	}
	if leftNum.Value > rightNum.Value {
		return &LispAtom{Value: "t"}, nil
	}
	return &LispAtom{Value: "nil"}, nil
}

func builtinEq(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to =")
	}
	left, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	right, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	leftNum, ok := left.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to =: %v", left)
	}
	rightNum, ok := right.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to =: %v", right)
	}
	if leftNum.Value == rightNum.Value {
		return &LispAtom{Value: "t"}, nil
	}
	return &LispAtom{Value: "nil"}, nil
}

func builtinIf(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("wrong number of arguments to if")
	}
	cond, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	condAtom, ok := cond.(*LispAtom)
	if !ok {
		return nil, fmt.Errorf("invalid condition in if: %v", cond)
	}
	if condAtom.Value != "nil" {
		return Eval(env, args[1])
	}
	return Eval(env, args[2])
}

func builtinDefun(env Environment, args []LispValue) (LispValue, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("wrong number of arguments to defun")
	}
	nameAtom, ok := args[0].(*LispAtom)
	if !ok {
		return nil, fmt.Errorf("invalid function name: %v", args[0])
	}
	paramsList, ok := args[1].(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid parameter list: %v", args[1])
	}
	body := args[2]
	env[nameAtom.Value] = &LispFunction{
		Params: paramsList.Elements,
		Body:   body,
	}
	return nameAtom, nil
}

// LispFunction represents a user-defined function
type LispFunction struct {
	Params []LispValue
	Body   LispValue
}

func (f *LispFunction) String() string {
	return "<function>"
}

func callFunction(env Environment, name string, args []LispValue) (LispValue, error) {
	fn, ok := env[name].(*LispFunction)
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", name)
	}
	if len(args) != len(fn.Params) {
		return nil, fmt.Errorf("wrong number of arguments to %s", name)
	}
	newEnv := make(Environment)
	for key, val := range env {
		newEnv[key] = val
	}
	for i, param := range fn.Params {
		paramAtom, ok := param.(*LispAtom)
		if !ok {
			return nil, fmt.Errorf("invalid parameter: %v", param)
		}
		argValue, err := Eval(env, args[i])
		if err != nil {
			return nil, err
		}
		newEnv[paramAtom.Value] = argValue
	}
	return Eval(newEnv, fn.Body)
}

// REPL implementation

func repl() {
	env := make(Environment)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		tokens := Tokenize(input)
		ast, _, err := Parse(tokens)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		result, err := Eval(env, ast)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println(result)
	}
}

func main() {
	repl()
}
