package main

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
)

func evalMultipleExpressions(env Environment, expressions []LispValue) ([]LispValue, error) {
	results := make([]LispValue, 0, len(expressions))
	for _, expr := range expressions {
		result, err := Eval(env, expr)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

var env Environment

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}

	for key, value := range builtins {
		s = append(s, prompt.Suggest{Text: key, Description: value})
	}

	// Add defined symbols from the environment
	for symbol := range env {
		s = append(s, prompt.Suggest{Text: symbol, Description: "Defined symbol"})
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func executor(input string) {
	tokens := Tokenize(input)
	expr, _, err := Parse(tokens)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if list, ok := expr.(*LispList); ok {
		results, err := evalMultipleExpressions(env, list.Elements)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			for _, result := range results {
				fmt.Println(result)
			}
		}
	} else {
		result, err := Eval(env, expr)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(result)
		}
	}
}

func initEnvironment() Environment {
	env := make(Environment)
	env[T] = &LispBoolean{Value: true}
	env[NIL] = &LispNil{}
	env[TRUE] = &LispBoolean{Value: true}
	env[FALSE] = &LispBoolean{Value: false}
	return env
}

// readFile reads the content of a file and returns it as a string
func readFile(filepath string) (string, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// main
func main() {
	env = initEnvironment()

	if len(os.Args) > 1 {
		// File execution mode
		filepath := os.Args[1]
		content, err := readFile(filepath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		tokens := Tokenize(content)
		expr, _, err := Parse(tokens)
		if err != nil {
			fmt.Println("Error parsing file:", err)
			return
		}
		results, err := evalMultipleExpressions(env, expr.(*LispList).Elements)
		if err != nil {
			fmt.Println("Error evaluating file:", err)
			return
		}
		for _, result := range results {
			fmt.Println(result)
		}
	} else {
		// REPL mode with autocompletion
		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefix("cclisp> "),
			prompt.OptionTitle("CCLisp REPL"),
		)
		p.Run()
	}
}
