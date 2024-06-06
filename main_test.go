package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestTokenize tests the Tokenize function
func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"(define x 10)", []string{"(", "define", "x", "10", ")"}},
		{"(+ 1 2)", []string{"(", "+", "1", "2", ")"}},
		{"(print \"hello world\")", []string{"(", "print", "\"hello world\"", ")"}},
		{"(if (> x 10) \"yes\" \"no\")", []string{"(", "if", "(", ">", "x", "10", ")", "\"yes\"", "\"no\"", ")"}},
	}

	for _, test := range tests {
		result := Tokenize(test.input)
		if !equal(result, test.expected) {
			t.Errorf("Tokenize(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

// TestParse tests the Parse function
func TestParse(t *testing.T) {
	tests := []struct {
		tokens   []string
		expected LispValue
	}{
		{[]string{"10"}, &LispNumber{Value: 10}},
		{[]string{"\"hello\""}, &LispString{Value: "hello"}},
		{[]string{"(", "+", "1", "2", ")"}, &LispList{Elements: []LispValue{&LispAtom{Value: "+"}, &LispNumber{Value: 1}, &LispNumber{Value: 2}}}},
	}

	for _, test := range tests {
		result, _, err := Parse(test.tokens)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("Parse(%v) = %v, %v, want %v", test.tokens, result, err, test.expected)
		}
	}
}

// TestEval tests the Eval function
func TestEval(t *testing.T) {
	env := Environment{
		"x": &LispNumber{Value: 10},
	}

	tests := []struct {
		expr     LispValue
		expected LispValue
	}{
		{&LispNumber{Value: 10}, &LispNumber{Value: 10}},
		{&LispString{Value: "hello"}, &LispString{Value: "hello"}},
		{&LispAtom{Value: "x"}, &LispNumber{Value: 10}},
		{&LispList{Elements: []LispValue{&LispAtom{Value: "+"}, &LispNumber{Value: 1}, &LispNumber{Value: 2}}}, &LispNumber{Value: 3}},
	}

	for _, test := range tests {
		result, err := Eval(env, test.expr)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("Eval(%v) = %v, %v, want %v", test.expr, result, err, test.expected)
		}
	}
}

// TestBuiltinAdd tests the builtinAdd function
func TestBuiltinAdd(t *testing.T) {
	env := Environment{}

	tests := []struct {
		args     []LispValue
		expected LispValue
	}{
		{[]LispValue{&LispNumber{Value: 1}, &LispNumber{Value: 2}}, &LispNumber{Value: 3}},
		{[]LispValue{&LispNumber{Value: 10}, &LispNumber{Value: 20}}, &LispNumber{Value: 30}},
	}

	for _, test := range tests {
		result, err := builtinAdd(env, test.args)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("builtinAdd(%v) = %v, %v, want %v", test.args, result, err, test.expected)
		}
	}
}

// TestBuiltinSub tests the builtinSub function
func TestBuiltinSub(t *testing.T) {
	env := Environment{}

	tests := []struct {
		args     []LispValue
		expected LispValue
	}{
		{[]LispValue{&LispNumber{Value: 10}, &LispNumber{Value: 5}}, &LispNumber{Value: 5}},
		{[]LispValue{&LispNumber{Value: 20}, &LispNumber{Value: 10}, &LispNumber{Value: 5}}, &LispNumber{Value: 5}},
	}

	for _, test := range tests {
		result, err := builtinSub(env, test.args)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("builtinSub(%v) = %v, %v, want %v", test.args, result, err, test.expected)
		}
	}
}

// TestBuiltinMul tests the builtinMul function
func TestBuiltinMul(t *testing.T) {
	env := Environment{}

	tests := []struct {
		args     []LispValue
		expected LispValue
	}{
		{[]LispValue{&LispNumber{Value: 2}, &LispNumber{Value: 3}}, &LispNumber{Value: 6}},
		{[]LispValue{&LispNumber{Value: 4}, &LispNumber{Value: 5}}, &LispNumber{Value: 20}},
	}

	for _, test := range tests {
		result, err := builtinMul(env, test.args)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("builtinMul(%v) = %v, %v, want %v", test.args, result, err, test.expected)
		}
	}
}

// TestBuiltinDiv tests the builtinDiv function
func TestBuiltinDiv(t *testing.T) {
	env := Environment{}

	tests := []struct {
		args     []LispValue
		expected LispValue
		err      string
	}{
		{[]LispValue{&LispNumber{Value: 10}, &LispNumber{Value: 2}}, &LispNumber{Value: 5}, ""},
		{[]LispValue{&LispNumber{Value: 20}, &LispNumber{Value: 5}}, &LispNumber{Value: 4}, ""},
		{[]LispValue{&LispNumber{Value: 10}, &LispNumber{Value: 0}}, nil, "division by zero"},
	}

	for _, test := range tests {
		result, err := builtinDiv(env, test.args)
		if (err != nil && err.Error() != test.err) || (err == nil && !lispValueEqual(result, test.expected)) {
			t.Errorf("builtinDiv(%v) = %v, %v, want %v, %v", test.args, result, err, test.expected, test.err)
		}
	}
}

// TestBuiltinLt tests the builtinLt function
func TestBuiltinLt(t *testing.T) {
	env := Environment{}

	tests := []struct {
		args     []LispValue
		expected LispValue
	}{
		{[]LispValue{&LispNumber{Value: 1}, &LispNumber{Value: 2}}, &LispAtom{Value: "t"}},
		{[]LispValue{&LispNumber{Value: 3}, &LispNumber{Value: 2}}, &LispAtom{Value: "nil"}},
	}

	for _, test := range tests {
		result, err := builtinLt(env, test.args)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("builtinLt(%v) = %v, %v, want %v", test.args, result, err, test.expected)
		}
	}
}

// TestBuiltinGt tests the builtinGt function
func TestBuiltinGt(t *testing.T) {
	env := Environment{}

	tests := []struct {
		args     []LispValue
		expected LispValue
	}{
		{[]LispValue{&LispNumber{Value: 3}, &LispNumber{Value: 2}}, &LispAtom{Value: "t"}},
		{[]LispValue{&LispNumber{Value: 1}, &LispNumber{Value: 2}}, &LispAtom{Value: "nil"}},
	}

	for _, test := range tests {
		result, err := builtinGt(env, test.args)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("builtinGt(%v) = %v, %v, want %v", test.args, result, err, test.expected)
		}
	}
}

// TestBuiltinEq tests the builtinEq function
func TestBuiltinEq(t *testing.T) {
	env := Environment{}

	tests := []struct {
		args     []LispValue
		expected LispValue
	}{
		{[]LispValue{&LispNumber{Value: 2}, &LispNumber{Value: 2}}, &LispAtom{Value: "t"}},
		{[]LispValue{&LispNumber{Value: 2}, &LispNumber{Value: 3}}, &LispAtom{Value: "nil"}},
	}

	for _, test := range tests {
		result, err := builtinEq(env, test.args)
		if err != nil || !lispValueEqual(result, test.expected) {
			t.Errorf("builtinEq(%v) = %v, %v, want %v", test.args, result, err, test.expected)
		}
	}
}

// Helper functions for tests

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func lispValueEqual(a, b LispValue) bool {
	switch x := a.(type) {
	case *LispNumber:
		y, ok := b.(*LispNumber)
		return ok && x.Value == y.Value
	case *LispString:
		y, ok := b.(*LispString)
		return ok && x.Value == y.Value
	case *LispAtom:
		y, ok := b.(*LispAtom)
		return ok && x.Value == y.Value
	case *LispList:
		y, ok := b.(*LispList)
		if !ok || len(x.Elements) != len(y.Elements) {
			return false
		}
		for i := range x.Elements {
			if !lispValueEqual(x.Elements[i], y.Elements[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// mockIO is a simple struct to mock standard input and output
type mockIO struct {
	stdin  *os.File
	stdout *os.File
}

// NewMockIO creates a new mockIO instance with the given input and output
func NewMockIO(input string) *mockIO {
	stdinR, stdinW, _ := os.Pipe()   // Create a pipe for standard input
	stdoutR, stdoutW, _ := os.Pipe() // Create a pipe for standard output
	stdinW.WriteString(input)        // Write input to the pipe
	stdinW.Close()                   // Close the write end of the input pipe
	stdoutW.Close()                  // Close the write end of the output pipe
	return &mockIO{
		stdin:  stdinR,
		stdout: stdoutR,
	}
}

// Stdin returns the mocked standard input
func (m *mockIO) Stdin() *os.File {
	return m.stdin
}

// Stdout returns the mocked standard output
func (m *mockIO) Stdout() *os.File {
	return m.stdout
}

// TestREPL tests the REPL function
func TestREPL(t *testing.T) {
	tests := []struct {
		input           string
		expectedOutputs []string
	}{
		{"(+ 1 2)\n(+ 3 4)\n", []string{"3", "7"}},
		{"(define x 5)\nx\n", []string{"", "5"}},
		{"(exit)\n", []string{""}},
	}

	for _, test := range tests {
		// Create a mockIO instance with the given input
		mock := NewMockIO(test.input)

		// Redirect standard output to a buffer
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Swap standard input with the mock
		oldStdin := os.Stdin
		os.Stdin = mock.Stdin()

		// Run the REPL
		repl()

		// Restore standard input and output
		os.Stdin = oldStdin
		os.Stdout = oldStdout
		w.Close()

		// Read the output from the buffer
		var buf bytes.Buffer
		io.Copy(&buf, r)
		r.Close()
		actualOutputStr := buf.String()

		// Check the outputs
		actualOutputs := strings.Split(strings.TrimSpace(actualOutputStr), "\n")
		for i, expected := range test.expectedOutputs {
			if i >= len(actualOutputs) {
				t.Fatalf("expected output %q, but got no output", expected)
			}
			actual := strings.TrimSpace(actualOutputs[i])

			// Check if the output contains the expected string
			if !strings.Contains(actual, expected) {
				t.Errorf("expected output %q, but got %q", expected, actual)
			}

			// Check if error messages are not present when not expected
			if !strings.HasPrefix(actual, "> Error:") && !strings.HasPrefix(expected, "> Error:") {
				if strings.Contains(actual, "> Error:") {
					t.Errorf("unexpected error message in output: %q", actual)
				}
			}
		}
	}
}
