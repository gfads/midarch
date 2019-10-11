package components

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/connectors"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/messages"
	"newsolution/gmidarch/execution/engine"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"os"
	"reflect"
)

type Unit struct {
	CSP   string
	Graph graphs.ExecGraph
}

func NewUnit() Unit {

	// create a new instance of Unit
	r := new(Unit)

	return *r
}

func (u *Unit) ConfigureUnit(elem interface{}, invP *chan messages.SAMessage) {

	// configure the state machine
	u.Graph = *graphs.NewExecGraph(2)
	msg := new(messages.SAMessage)
	actionChannel := make(chan messages.SAMessage)

	info := make([] *interface{}, 2)
	info[0] = new(interface{})
	info[1] = new(interface{})
	info[2] = new(interface{})
	*info[0] = new(messages.SAMessage)
	*info[1] = elem

	newEdgeInfo := graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Execute", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	u.Graph.AddEdge(0, 0, newEdgeInfo)

	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionName: "InvP", ActionType: 2, ActionChannel: invP, Message: msg}
	u.Graph.AddEdge(0, 1, newEdgeInfo)

	actionChannel = make(chan messages.SAMessage)
	info1 := make([]*interface{}, 1)
	info1[0] = new(interface{})
	*info1[0] = new(messages.SAMessage)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_PerformAdaptation", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info1}
	u.Graph.AddEdge(1, 0, newEdgeInfo)
}

func (Unit) I_Execute(msg *messages.SAMessage, elem interface{}) {

	elemTemp := elem

	switch  reflect.TypeOf(elemTemp).String() {
	case "components.Sender":
		elemTemp := elemTemp.(Sender)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Receiver":
		elemTemp := elemTemp.(Receiver)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
/*	case "components.Client":
		elemTemp := elemTemp.(Client)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Server":
		elemTemp := elemTemp.(Server)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Calculatorclient":
		elemTemp := elemTemp.(Calculatorclient)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Calculatorserver":
		elemTemp := elemTemp.(Calculatorserver)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.CRH":
		elemTemp := elemTemp.(CRH)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.SRH":
		elemTemp := elemTemp.(SRH)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Marshaller":
		elemTemp := elemTemp.(Marshaller)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Requestor":
		elemTemp := elemTemp.(Requestor)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.CalculatorProxy":
		elemTemp := elemTemp.(CalculatorProxy)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
	case "components.Invoker":
		elemTemp := elemTemp.(Invoker)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
*/
	case "connectors.Oneway":
		elemTemp := elemTemp.(connectors.Oneway)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
/*	case "connectors.Requestreply":
		elemTemp := elemTemp.(connectors.Requestreply)
		engine.Engine{}.Execute(elemTemp, elemTemp.Graph, !parameters.EXECUTE_FOREVER)
*/	default:
		fmt.Println("Unit:: Architectural element '" + reflect.TypeOf(elemTemp).String() + "' not supported")
		os.Exit(0)
	}
}

func (Unit) I_PerformAdaptation(msg *messages.SAMessage) {
	fmt.Println("Unit:: Perform Adaptation ***********")
}
