package apps

import (
	"fmt"
	"github.com/gfads/midarch/src/gmidarch/development/messages"
	"github.com/gfads/midarch/src/shared"
	"os"
)

// @Type: Client
// @Behaviour: Behaviour = I_Beforesend -> InvR.e1 -> TerR.e1 -> I_Afterreceive -> Behaviour
type Client struct{}

// Calculator client
/*
func (Client) I_Beforesend(id string, msg *messages.SAMessage, info *interface{}) {
	req := new(messages.FunctionalRequest)
	i++
	params := []interface{}{i, i}
	req.Op = "Add"
	req.Params = params

	msg.Payload = req
}
func (Client) I_Afterreceive(id string, msg *messages.SAMessage, info *interface{}) {
	fmt.Println(id, msg.Payload)
	//if i == 1 {
	//	os.Exit(0)
	//}
}
*/
var i = 0

// Naming client

func (Client) I_Beforesend(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	if i == 0 {
		aor := messages.AOR{Host: "localhost", Port: shared.NAMING_PORT, Id: 123456, ProxyName: "calculatorimpl"}
		params := []interface{}{"Calculator", aor}
		req := messages.FunctionalRequest{Op: "Register", Params: params}
		msg.Payload = req
	}
	if i == 1 {
		params := []interface{}{"Calculator"}
		req := messages.FunctionalRequest{Op: "Lookup", Params: params}
		msg.Payload = req
	}
	i++
}
func (Client) I_Afterreceive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	fmt.Println(id, msg.Payload)
	if i == 2 {
		os.Exit(0)
	}
}
