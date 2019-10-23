package components

import (
	"fmt"
	graphs2 "gmidarch/development/artefacts/graphs"
	element2 "gmidarch/development/element"
	messages2 "gmidarch/development/messages"
	"shared/shared"
)

type Client struct {
	Behaviour string
	Graph     graphs2.ExecGraph
}

func NewClient() Client {

	// create a new instance of client
	r := new(Client)
	r.Behaviour = "B = I_Setmessage -> InvR.e1 -> TerR.e1 -> I_Printmessage -> B"

	return *r
}

func (c *Client) Configure(invR, terR *chan messages2.SAMessage) {

	// configure the state machine
	c.Graph = *graphs2.NewExecGraph(4)
	actionChannel := make(chan messages2.SAMessage)

	msg := new(messages2.SAMessage)
	info := make([]*interface{}, 1)
	info[0] = new(interface{})
	*info[0] = msg

	newEdgeInfo := graphs2.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Setmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(0, 1, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.InvR, ActionName: "InvR", ActionType: 2, ActionChannel: invR, Message:msg}
	c.Graph.AddEdge(1, 2, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{ExternalAction: element2.Element{}.TerR, ActionName: "TerR", ActionType: 2, ActionChannel: terR, Message:msg}
	c.Graph.AddEdge(2, 3, newEdgeInfo)
	newEdgeInfo = graphs2.ExecEdgeInfo{InternalAction: shared.Invoke, ActionName: "I_Printmessage", ActionType: 1, ActionChannel: &actionChannel, Message: msg, Info: info}
	c.Graph.AddEdge(3, 0, newEdgeInfo)
}

func (Client) I_Setmessage(msg *messages2.SAMessage, info [] *interface{}) {
	*msg = messages2.SAMessage{Payload: "Hello World from Client"}
}

func (Client) I_Printmessage(msg *messages2.SAMessage, info [] *interface{}) {

	fmt.Println(msg.Payload)
}