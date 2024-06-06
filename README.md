Simple implementation of lisp interpreter in Go

### Features

- Subset of Lisp: basic arithmetic operations, conditionals, and function definitions.
- Parse Input: A parser to read and tokenize Lisp expressions.
- Evaluate Expressions: An evaluator that recursively processes expressions, handling atoms, function calls, and special forms.
- Define Built-in Functions: Functions like addition, subtraction, multiplication, division and conditionals.
- Support User-Defined Functions: Allow users to define their own functions using defun.
- Create a REPL: Build a Read-Eval-Print Loop (REPL) for interactive use.

### Example Interaction
````
> (+ 1 2)
3
> (- 10 4)
6
> (* 3 4)
12
> (/ 8 2)
4
````


````
> (defun square (x) (* x x))
<function>

> (square 4)
16

> (if (= 4 4) "equal" "not equal")
equal

> (if (> 10 5) "greater" "less")
greater

> (defun abs (x) (if (< x 0) (- 0 x) x))
<function>

> (abs -7)
7

> (abs 7)
7

````