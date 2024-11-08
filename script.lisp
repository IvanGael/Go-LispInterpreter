(
  (+ 1 2)
  (% 382 4)
  (let ((fib (lambda (n)
  (if (< n 2)
      n
      (+ (fib (- n 1))
         (fib (- n 2)))))))
    (format t "The 7th number of the Fibonacci sequence is %d" (fib 7)))
  (append (list 1 2) (list 3 4))
  (defun abs (x) (if (< x 0) (- 0 x) x))
  (abs -3)
  (not true)
  (length (list 1 2 3 4))
  (if (> 10 5) "greater" "less")
  (let ((a 10) (b 20)) (* a b))
  (cdr (list 1 2 3))
  (car (list 1 2 3))
)


  



