package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/shared"
	"newsolution/gmidarch/development/impl"
	"fmt"
	"os"
	"newsolution/gmidarch/development/element"
)

type Calculatorserver struct {
	CSP   string
	Graph graphs.ExecGraph
}

func Newcalculatorserver() Calculatorserver {

	// create a new instance of Server
	r := new(Calculatorserver)

	return *r
}

func (c *Calculatorserver) Configure(invP *chan messages.SAMessage, terP *chan messages.SAMessage) {

	// configure the state machine
	c.Graph = *graphs.NewExecGraph(3)

	msg := new(messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	actionChannel := make(chan messages.SAMessage)
	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction:element.Element{}.InvP, ActionType: 2, ActionChannel: invP, Message: msg}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke,ActionName:"I_Process", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionType: 2, ActionChannel: terP, Message: msg}
	c.Graph.AddEdge(2, 0, newEdgeInfo)
}

func (Calculatorserver) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	req := msg.Payload.(shared.Request)
	op := req.Op
	p1 := int(req.Args[0].(float64))
	p2 := int(req.Args[1].(float64))
	r := 0

	switch op {
	case "add":
		r = impl.CalculatorImpl{}.Add(p1,p2)
	case "sub":
		r = impl.CalculatorImpl{}.Sub(p1,p2)
	case "mul":
		r = impl.CalculatorImpl{}.Mul(p1,p2)
	case "div":
		r = impl.CalculatorImpl{}.Div(p1,p2)
	default:
		fmt.Println("Server Calculator:: Operation '" + op + "' not supported!!")
		os.Exit(0)
	}

	*msg = messages.SAMessage{Payload: r}
}