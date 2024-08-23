package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readMultilineInput(scanner *bufio.Scanner) string {
	var input strings.Builder
	parenCount := 0

	for {
		if !scanner.Scan() {
			return input.String()
		}
		line := scanner.Text()
		input.WriteString(line + "\n")

		parenCount += strings.Count(line, "(") - strings.Count(line, ")")

		if parenCount == 0 && input.Len() > 0 {
			return input.String()
		}
	}
}

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

func initEnvironment() Environment {
	env := make(Environment)
	env["t"] = &LispBoolean{Value: true}
	env["nil"] = &LispNil{}
	env["true"] = &LispBoolean{Value: true}
	env["false"] = &LispBoolean{Value: false}
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
	env := initEnvironment()
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
		// REPL mode
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println(" -----+ GoLisp! Type your expressions below. +----- ")
		fmt.Println(" ")
		for {
			fmt.Print("cclisp> ")
			input := readMultilineInput(scanner)
			if input == "" {
				break
			}
			tokens := Tokenize(input)
			expr, _, err := Parse(tokens)
			if err != nil {
				fmt.Println("Error:", err)
				continue
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
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	}
}
