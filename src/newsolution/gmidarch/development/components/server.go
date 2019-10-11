package components

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/element"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/shared"
	"strings"
)

type Server struct {
	Behaviour   string
	Graph graphs.ExecGraph
}

func NewServer() Server {

	// create a new instance of Server
	r := new(Server)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (s *Server) Configure(invP, terP *chan messages.SAMessage) Server {

	// configure the state machine
	msg := new(messages.SAMessage)
	s.Graph = *graphs.NewExecGraph(3)

	actionChannel := make(chan messages.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, Message: msg, ActionChannel: invP, ActionType: 2}
	s.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Process", Message: msg, ActionType: 1, ActionChannel: &actionChannel, Info: info}
	s.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, Message: msg, ActionChannel: terP, ActionType: 2}
	s.Graph.AddEdge(2, 0, newEdgeInfo)

	return *s
}

func (Server) I_Process(msg *messages.SAMessage,info [] *interface{}) {
	msgTemp := strings.ToUpper(msg.Payload.(string))
	*msg = messages.SAMessage{Payload: msgTemp}
}
