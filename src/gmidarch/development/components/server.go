package components

import (
	graphs2 "gmidarch/development/artefacts/graphs"
	element2 "gmidarch/development/element"
	messages2 "gmidarch/development/messages"
	"shared/shared"
	"strings"
)

type Server struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewServer() Server {

	// create a new instance of Server
	r := new(Server)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (s *Server) Configure(invP, terP *chan messages2.SAMessage) Server {

	// configure the state machine
	msg := new(messages2.SAMessage)
	s.Graph = *graphs2.NewExecGraph(3)

	actionChannel := make(chan messages2.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Process", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info}
	s.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.TerP, Message: msg, ActionChannel: terP, ActionType: 2}
	s.Graph.AddEdge(2, 0, newEdgeInfo)

	return *s
}

func (Server) I_Process(msg *messages2.SAMessage,info [] *interface{}) {
	msgTemp := strings.ToUpper(msg.Payload.(string))
	*msg = messages2.SAMessage{Payload: msgTemp}
}
