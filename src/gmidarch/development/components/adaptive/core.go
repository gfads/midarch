package adaptive

import (
	"fmt"
	"gmidarch/development/messages"
)

var count int

//@Type: Core
//@Behaviour: Behaviour = RUNTIME
type Core struct {
	//Graph     exec.ExecGraph
}

func NewCore() Core {
	r := new(Core)

	return *r
}

func (Core) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
//	Core{}.I_Debug(msg,info)
	fmt.Println("Core::Selector::msg.Payload", msg.Payload)
}

func (Core) I_Debug(id string, msg *messages.SAMessage, info *interface{}) {
	fmt.Printf("Core:: %v\n",msg.Payload)
}
