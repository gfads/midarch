package adaptive

import (
	"github.com/gfads/midarch/src/gmidarch/development/messages"
)

// @Type: Monitor
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Monitor struct{}

func (Monitor) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("Monitor.I_Process::Payload:", msg.Payload)
	*msg = messages.SAMessage{Payload: msg.Payload}
}
