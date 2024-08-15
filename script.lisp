(
  (+ 1 2)
  (let ((fib (lambda (n)
  (if (< n 2)
      n
      (+ (fib (- n 1))
         (fib (- n 2)))))))
    (format t "The 7th number of the Fibonacci sequence is %d" (fib 7)))
  (append (list 1 2) (list 3 4))
)


  



