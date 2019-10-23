package components

import (
	"fmt"
	graphs2 "gmidarch/development/artefacts/graphs"
	element2 "gmidarch/development/element"
	impl2 "gmidarch/development/impl"
	messages2 "gmidarch/development/messages"
	"os"
	"shared/shared"
)

type Calculatorserver struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func Newcalculatorserver() Calculatorserver {

	// create a new instance of Server
	r := new(Calculatorserver)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (c *Calculatorserver) Configure(invP *chan messages2.SAMessage, terP *chan messages2.SAMessage) {

	// configure the state machine
	c.Graph = *graphs2.NewExecGraph(3)

	msg := new(messages2.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	actionChannel := make(chan messages2.SAMessage)
	newEdgeInfo := graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_Process", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	c.Graph.AddEdge(2, 0, newEdgeInfo)
}

func (Calculatorserver) I_Process(msg *messages2.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)
	op := req.Op
	p1 := int(req.Args[0].(float64))
	p2 := int(req.Args[1].(float64))
	r := 0

	switch op {
	case "add":
		r = impl2.CalculatorImpl{}.Add(p1,p2)
	case "sub":
		r = impl2.CalculatorImpl{}.Sub(p1,p2)
	case "mul":
		r = impl2.CalculatorImpl{}.Mul(p1,p2)
	case "div":
		r = impl2.CalculatorImpl{}.Div(p1,p2)
	default:
		fmt.Println("Server Calculator:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}

	*msg = messages2.SAMessage{Payload: r}
}