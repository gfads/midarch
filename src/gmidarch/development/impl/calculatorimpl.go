package impl

type CalculatorImpl struct {}

func (CalculatorImpl) Add(p1, p2 int) int {
	return p1 + p2
}

func (CalculatorImpl) Sub(p1, p2 int) int {
	return p1 - p2
}

func (CalculatorImpl) Mul(p1, p2 int) int {
	return p1 * p2
}

func (CalculatorImpl) Div(p1, p2 int) int {
	return p1 / p2
}

