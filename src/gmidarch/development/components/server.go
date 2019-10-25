package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"strings"
)

type Server struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func NewServer() Server {

	// create a new instance of Server
	r := new(Server)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"

	return *r
}

func (Server) I_Process(msg *messages.SAMessage,info [] *interface{}) {
	msgTemp := strings.ToUpper(msg.Payload.(string))
	*msg = messages.SAMessage{Payload: msgTemp}
}
