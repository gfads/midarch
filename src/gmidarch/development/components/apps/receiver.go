package apps

import (
	"fmt"
	"github.com/gfads/midarch/src/gmidarch/development/messages"
)

// @Type: Receiver
// @Behaviour: Behaviour = InvP.e1 -> I_Printmessage -> Behaviour
type Receiver struct{}

func (Receiver) I_Printmessage(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	fmt.Println("From: ", msg.From, "To: ", msg.To, "Message: ", msg.Payload)
}
