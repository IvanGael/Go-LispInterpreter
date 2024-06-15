Simple implementation of lisp interpreter in Go

### Features

- Subset of Lisp: basic arithmetic operations, conditionals, and function definitions.
- Parse Input: A parser to read and tokenize Lisp expressions.
- Evaluate Expressions: An evaluator that recursively processes expressions, handling atoms, function calls, and special forms.
- Define Built-in Functions: Functions like addition, subtraction, multiplication, division and conditionals.
- Support User-Defined Functions: Allow users to define their own functions using defun.
- Support for lambda functions, local variables bindings(let) and logical operations(and, or and not)
- Support for basic list operations (car, cdr, cons, length, and append)
- Create a REPL: Build a Read-Eval-Print Loop (REPL) for interactive use.

### Example Interaction
Arithmetic operations
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


Function definitions and conditionals
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

Local variable bindings
````
> (let ((square (lambda (x) (* x x))))
    (square 5))
25

> (let ((a 10) (b 20))
    (+ a b))
30
> (let ((a 6))
    (if (and (< a 5) (> a 0)) "3 is comprised between 0 and 5" "3 is not comprised between 0 and 5"))
    "3 is not comprised between 0 and 5"
````


Logical operations
````
> (and true true)
true
> (and true false)
false
> (or true false)
true
> (or true true)
true
> (not true)
false
> (not false)
true
````


List operations
````
> (car (list 1 2 3))
1
> (cdr (list 1 2 3))
(2 3)
> (cons 1 (list 2 3))
(1 2 3)
> (length (list 1 2 3 4))
4
> (append (list 1 2) (list 3 4))
(1 2 3 4)
````

### Testing
````
go test -v
````