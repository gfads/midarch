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

func (Server) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(Server).I_Process(msg, info)
}

func (Server) OldSelector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}) {

	var f func(*messages.SAMessage, []*interface{})
	switch op {
	case "I_Process":
		f = func(msg *messages.SAMessage, info []*interface{}) {
			elem.(Server).I_Process(msg, info)
		}
	}
	return f
}

func (Server) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	msgTemp := strings.ToUpper(msg.Payload.(string))
	*msg = messages.SAMessage{Payload: msgTemp}
}
