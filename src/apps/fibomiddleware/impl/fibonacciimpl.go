package impl

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
