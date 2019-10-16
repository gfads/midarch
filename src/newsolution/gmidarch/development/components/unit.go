package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/messages"
	"newsolution/gmidarch/execution/engine"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
)

type Unit struct {
	Behaviour   string
	Graph       graphs.ExecGraph
	ElemOfUnit  interface{}
	GraphOfElem graphs.ExecGraph
}

func NewUnit() Unit {

	r := new(Unit)
	r.Behaviour = "B = " + parameters.RUNTIME_BEHAVIOUR
	//r.Behaviour = "B = InvP.e1 -> I_Initialiseunit -> P1\n P1 = I_Execute -> P1 [] InvP.e1 -> P1"

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

func (u Unit) I_Initialiseunit(msg *messages.SAMessage, info [] *interface{}) {
}

func (u Unit) I_Execute(msg *messages.SAMessage, info [] *interface{}) {

	engine.Engine{}.Execute(u.ElemOfUnit, u.GraphOfElem, !parameters.EXECUTE_FOREVER)
}

func (Unit) I_Adaptunit(msg *messages.SAMessage, info [] *interface{}) {
	//fmt.Println("Unit:: I_Adaptunit ***********")
}
