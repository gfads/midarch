package fibonacciImpl

type Fibonacci struct{}

func (Fibonacci) F(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return Fibonacci{}.F(n-1) + Fibonacci{}.F(n-2)
	}
}

// Calculate Fibonacci Number based on RPC function signature, where r is the return of the function
func (f Fibonacci) FiboRPC(n int, r *int) error {
	*r = f.F(n)
	return nil
}
