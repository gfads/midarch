package adaptive

import (
	"gmidarch/development/messages"
)

//@Type: Monitor
//@Behaviour: Behaviour = InvP.e1 -> I_Process -> InvR.e2 -> Behaviour
type Monitor struct {}

func (Monitor) I_Process(msg *messages.SAMessage, info [] *interface{}) {
	*msg = messages.SAMessage{Payload: msg.Payload}
}
