package apps

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"strconv"
	"time"
)

// @Type: Sender
// @Behaviour: Behaviour = I_Setmessage1 -> InvR.e1 -> Behaviour [] I_Setmessage2 -> InvR.e1 -> Behaviour
type Sender struct{}

var idx = 0 // REMOVE

func (s Sender) I_Setmessage1(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	time.Sleep(20 * time.Second)
	msg.Payload = "Message (type 1) [" + strconv.Itoa(idx) + "]"
	idx++
}

func (s Sender) I_Setmessage2(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	time.Sleep(20 * time.Second)
	msg.Payload = "Message (type 2) [" + strconv.Itoa(idx) + "]"
	idx++
	//shared.ExecuteForever = false
}
