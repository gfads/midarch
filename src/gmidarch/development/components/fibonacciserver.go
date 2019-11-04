package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"shared"
)

type Fibonacciserver struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Newfibonacciserver() Fibonacciserver {

	r := new(Fibonacciserver)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (Fibonacciserver) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)
	op := req.Op
	p1 := req.Args[0].(int)
	r := 0

	switch op {
	case "fibo":
		r = fibo(p1) // TODO
	default:
		fmt.Printf("Fibonacci Server:: Operation '%v' not supported!!",op)
		os.Exit(0)
	}
	*msg = messages.SAMessage{Payload: r}
}

func fibo(n int) int {
	return f(n)
}

func f(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return f(n-1) + f(n-2)
	}
}
