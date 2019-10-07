package components

import (
	"newsolution/development/artefacts/exec"
	"gmidarch/development/framework/messages"
	"newsolution/shared/shared"
	"newsolution/development/element"
	"fmt"
)

type Client struct {
	CSP      string
	Graph    exec.ExecGraph
}

func NewClient() Client {

	// create a new instance of client
	r := new(Client)

	return *r
}

func (c *Client) Configure(invR, terR *chan messages.SAMessage) {

	// configure the state machine
	c.Graph = *exec.NewExecGraph(4)
	actionChannel := make(chan messages.SAMessage)

	msg := new(interface{})
	args := make([]*interface{}, 1)
	args[0] = new(interface{})
	*args[0] = msg

	newEdgeInfo := exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionName: "InvR", ActionType: 2, ActionChannel: invR, Message:msg}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionName: "TerR", ActionType: 2, ActionChannel: terR, Message:msg}
	c.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = exec.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Printmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Args: args}
	c.Graph.AddEdge(3, 0, newEdgeInfo)
}

func (Client) I_Setmessage(msg *interface{}) {
	*msg = messages.SAMessage{Payload: "Hello World from Client"}
}
func (Client) I_Printmessage(msg *interface{}) {
	temp := *msg

	fmt.Println(temp.(messages.SAMessage).Payload)
}