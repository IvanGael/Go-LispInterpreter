package main

import (
	"fmt"
)

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
		case FORMAT:
			return builtinFormat(env, args)
		case PLUS:
			return builtinAdd(env, args)
		case MINUS:
			return builtinSub(env, args)
		case STAR:
			return builtinMul(env, args)
		case SLASH:
			return builtinDiv(env, args)
		case LESS_THAN:
			return builtinLt(env, args)
		case LESS_OR_EQUAL_THAN:
			return builtinLtOrEq(env, args)
		case GREATER_THAN:
			return builtinGt(env, args)
		case GREATER_OR_EQUAL_THAN:
			return builtinGtOrEq(env, args)
		case EQUAL:
			return builtinEq(env, args)
		case IF:
			return builtinIf(env, args)
		case DEFUN:
			return builtinDefun(env, args)
		case LAMBDA:
			return builtinLambda(env, args)
		case LET:
			return builtinLet(env, args)
		case AND:
			return builtinAnd(env, args)
		case OR:
			return builtinOr(env, args)
		case NOT:
			return builtinNot(env, args)
		case LIST:
			return builtinList(args)
		case CAR:
			return builtinCar(env, args)
		case CDR:
			return builtinCdr(env, args)
		case CONS:
			return builtinCons(env, args)
		case LENGTH:
			return builtinLength(env, args)
		case APPEND:
			return builtinAppend(env, args)
		default:
			return callFunction(env, fn.Value, args)
		}
	default:
		return nil, fmt.Errorf("unknown expression type: %T", v)
	}
}

// Helper function to convert Lisp values to Go values
func lispValueToGoValue(value LispValue) interface{} {
	switch v := value.(type) {
	case *LispNumber:
		return v.Value
	case *LispString:
		return v.Value
	case *LispAtom:
		return v.Value
	default:
		return v
	}
}

// Built-in function implementations

// builtinFormat is the implementation of the format function
func builtinFormat(env Environment, args []LispValue) (LispValue, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments to format")
	}

	// The second argument is the format string
	formatStr, ok := args[1].(*LispString)
	if !ok {
		return nil, fmt.Errorf("invalid format string: %v", args[1])
	}

	// The remaining arguments are the values to be formatted
	formatArgs := args[2:]

	// Prepare the arguments for fmt.Sprintf
	var sprintfArgs []interface{}
	for _, arg := range formatArgs {
		evalArg, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		sprintfArgs = append(sprintfArgs, lispValueToGoValue(evalArg))
	}

	// Perform the formatting
	formattedStr := fmt.Sprintf(formatStr.Value, sprintfArgs...)

	return &LispString{Value: formattedStr}, nil
}

// builtinAdd is built-in implementation of addition operation
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

// builtinSub is built-in implementation of subtraction operation
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
	diff := number.Value
	for _, arg := range args[1:] {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		number, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to -: %v", val)
		}
		diff -= number.Value
	}
	return &LispNumber{Value: diff}, nil
}

// builtinMul is built-in implementation of multiplication operation
func builtinMul(env Environment, args []LispValue) (LispValue, error) {
	prod := 1
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		number, ok := val.(*LispNumber)
		if !ok {
			return nil, fmt.Errorf("invalid argument to *: %v", val)
		}
		prod *= number.Value
	}
	return &LispNumber{Value: prod}, nil
}

// builtinDiv is built-in implementation of division operation
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

	quot := number.Value
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

		quot /= number.Value
	}

	return &LispNumber{Value: quot}, nil
}

// builtinLt is built-in implementation of less than condition
func builtinLt(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to <")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val2)
	}
	if num1.Value < num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinLtOrEq is built-in implementation of less or equal than condition
func builtinLtOrEq(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to <")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to <: %v", val2)
	}
	if num1.Value <= num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinGt is built-in implementation of greater than condition
func builtinGt(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to >")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val2)
	}
	if num1.Value > num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinGtOrEq is built-in implementation of greater or equal than condition
func builtinGtOrEq(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to >")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok := val1.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val1)
	}
	num2, ok := val2.(*LispNumber)
	if !ok {
		return nil, fmt.Errorf("invalid argument to >: %v", val2)
	}
	if num1.Value >= num2.Value {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinEq is built-in implementation of equal to condition
func builtinEq(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to =")
	}
	val1, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val2, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	num1, ok1 := val1.(*LispNumber)
	num2, ok2 := val2.(*LispNumber)
	if ok1 && ok2 {
		if num1.Value == num2.Value {
			return &LispAtom{Value: "true"}, nil
		}
		return &LispAtom{Value: "false"}, nil
	}
	if val1.String() == val2.String() {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinIf is built-in implementation of if conditional struct
func builtinIf(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("wrong number of arguments to if")
	}
	cond, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	if atom, ok := cond.(*LispAtom); ok && atom.Value == "true" {
		return Eval(env, args[1])
	}
	return Eval(env, args[2])
}

// builtinDefun is built-in implementation of function definition
func builtinDefun(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("wrong number of arguments to defun")
	}
	name, ok := args[0].(*LispAtom)
	if !ok {
		return nil, fmt.Errorf("invalid function name: %v", args[0])
	}
	params, ok := args[1].(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid function parameters: %v", args[1])
	}
	fn := &LispFunction{Name: name, Params: params.Elements, Body: args[2], Env: env}
	env[name.Value] = fn
	return fn, nil
}

// builtinLambda is built-in implementation of lambda function definition
func builtinLambda(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to lambda")
	}
	params, ok := args[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid lambda parameters: %v", args[0])
	}
	return &LispFunction{Params: params.Elements, Body: args[1], Env: env}, nil
}

// builtinLet is built-in implementation of let local variable definition
func builtinLet(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to let")
	}
	bindings, ok := args[0].(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid let bindings: %v", args[0])
	}
	localEnv := make(Environment)
	for key, value := range env {
		localEnv[key] = value
	}
	for _, binding := range bindings.Elements {
		bindList, ok := binding.(*LispList)
		if !ok || len(bindList.Elements) != 2 {
			return nil, fmt.Errorf("invalid let binding: %v", binding)
		}
		key, ok := bindList.Elements[0].(*LispAtom)
		if !ok {
			return nil, fmt.Errorf("invalid let binding key: %v", bindList.Elements[0])
		}
		val, err := Eval(localEnv, bindList.Elements[1])
		if err != nil {
			return nil, err
		}
		localEnv[key.Value] = val
	}
	return Eval(localEnv, args[1])
}

// builtinAnd is built-in implementation of and logical operation
func builtinAnd(env Environment, args []LispValue) (LispValue, error) {
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		if atom, ok := val.(*LispAtom); ok && atom.Value == "false" {
			return &LispAtom{Value: "false"}, nil
		}
	}
	return &LispAtom{Value: "true"}, nil
}

// builtinOr is built-in implementation of or logical operation
func builtinOr(env Environment, args []LispValue) (LispValue, error) {
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		if atom, ok := val.(*LispAtom); ok && atom.Value == "true" {
			return &LispAtom{Value: "true"}, nil
		}
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinNot is built-in implementation of not logical operation
func builtinNot(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to not")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	if atom, ok := val.(*LispAtom); ok && atom.Value == "false" {
		return &LispAtom{Value: "true"}, nil
	}
	return &LispAtom{Value: "false"}, nil
}

// builtinList is built-in implementation of list definition
func builtinList(args []LispValue) (LispValue, error) {
	return &LispList{Elements: args}, nil
}

// builtinCar is built-in implementation of car list operation. It retrieves first element of a list.
func builtinCar(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to car")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to car: %v", val)
	}
	if len(list.Elements) == 0 {
		return nil, fmt.Errorf("car of empty list")
	}
	return list.Elements[0], nil
}

// builtinCdr is built-in implementation of cdr list operation. It retrieves the rest elements of a list.
func builtinCdr(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to cdr")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to cdr: %v", val)
	}
	if len(list.Elements) == 0 {
		return &LispList{Elements: []LispValue{}}, nil
	}
	return &LispList{Elements: list.Elements[1:]}, nil
}

// builtinCons is built-in implementation of cons list operation. It add element to a list.
func builtinCons(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments to cons")
	}
	elem, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	val, err := Eval(env, args[1])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to cons: %v", val)
	}
	return &LispList{Elements: append([]LispValue{elem}, list.Elements...)}, nil
}

// builtinLength is built-in implementation of length list operation. It retrieves the length of a list.
func builtinLength(env Environment, args []LispValue) (LispValue, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments to length")
	}
	val, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	list, ok := val.(*LispList)
	if !ok {
		return nil, fmt.Errorf("invalid argument to length: %v", val)
	}
	return &LispNumber{Value: len(list.Elements)}, nil
}

// builtinAppend is built-in implementation of append list operation. It add a list to another list.
func builtinAppend(env Environment, args []LispValue) (LispValue, error) {
	var result []LispValue
	for _, arg := range args {
		val, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		list, ok := val.(*LispList)
		if !ok {
			return nil, fmt.Errorf("invalid argument to append: %v", val)
		}
		result = append(result, list.Elements...)
	}
	return &LispList{Elements: result}, nil
}

// callFunction calls a user-defined function
func callFunction(env Environment, name string, args []LispValue) (LispValue, error) {
	fn, ok := env[name]
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", name)
	}
	lambda, ok := fn.(*LispFunction)
	if !ok {
		return nil, fmt.Errorf("invalid function: %s", name)
	}
	if len(lambda.Params) != len(args) {
		return nil, fmt.Errorf("wrong number of arguments to %s", name)
	}
	localEnv := make(Environment)
	for key, value := range lambda.Env {
		localEnv[key] = value
	}
	for i, param := range lambda.Params {
		paramName, ok := param.(*LispAtom)
		if !ok {
			return nil, fmt.Errorf("invalid parameter name: %v", param)
		}
		argVal, err := Eval(env, args[i])
		if err != nil {
			return nil, err
		}
		localEnv[paramName.Value] = argVal
	}
	return Eval(localEnv, lambda.Body)
}
