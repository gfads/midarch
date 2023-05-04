package sender

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"strconv"
)

// @Type: Sender
// @Behaviour: Behaviour = I_Setmessage1 -> InvR.e1 -> Behaviour [] I_Setmessage2 -> InvR.e1 -> Behaviour
type Sender struct{}

var idx = 0 // REMOVE

func (s Sender) I_Setmessage1(id string, msg *messages.SAMessage, info *interface{}) {
	msg.Payload = "Message Adapted (type 1 V2) [" + strconv.Itoa(idx) + "]"
	idx++
}

func (s Sender) I_Setmessage2(id string, msg *messages.SAMessage, info *interface{}) {
	msg.Payload = "Message Adapted (type 2 V2) [" + strconv.Itoa(idx) + "]"
	idx++
}

//func (s Sender) GetType() interface{} {
//	fmt.Println("Passou pelo gettype Sender 2")
//
//	return Sender{}
//}
