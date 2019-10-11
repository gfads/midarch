package element

import (
	"gmidarch/development/framework/messages"
	"newsolution/shared/shared"
	"reflect"
	"gmidarch/development/artefacts/graphs"
	"fmt"
	"strings"
)

type ElementGo1 struct {
	ElemId string
	ElemType interface{}
	CSP string
	Info      interface{}
	GoStateMachine graphs.GraphExecutable
	FDRStateMachine graphs.GraphDot
}

// external actions common to all components and connectors
func (ElementGo) InvP1(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (ElementGo) InvR1(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}

func (ElementGo) TerR1(terR *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*terR
}

func (ElementGo) TerP1(terP *chan messages.SAMessage, msg *messages.SAMessage) {
	*terP <- *msg
}

func (ElementGo1) Loop1(elem ElementGo, graph graphs.GraphExecutable) {

	// Execute graph
	node := 0

	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			edge := edges[0]
			if shared.IsInternal(edge.Info.ActionName) {
				r := false
				if strings.Contains(elem.ElemId,"evolutive"){
					fmt.Printf("ElementGo:: %s:: %s %v \n",elem.ElemId,edge.Info.ActionName,edge.Info.ActionChannel)
				}
				edge.Info.InternalAction(elem.ElemType, edge.Info.ActionName, edge.Info.Message, &elem.Info, &r)
			} else {
				if strings.Contains(elem.ElemId,"evolutive"){
					fmt.Printf("ElementGo:: %s:: %s %v \n",elem.ElemId,edge.Info.ActionName,edge.Info.ActionChannel)
				}
				edge.Info.ExternalAction(edge.Info.ActionChannel, edge.Info.Message)
			}
			node = edge.To
		} else {
			msg := messages.SAMessage{}
			chosen := 0
			choice(elem, &msg, &chosen, edges)
			node = edges[chosen].To
			if strings.Contains(elem.ElemId,"evolutive"){
				fmt.Printf("ElementGo:: %s:: %v \n",elem.ElemId,edges[chosen].Info.ActionName)
			}
		}
		if node == 0 {
			break
		}
	}
	return
}

func (e *ElementGo) SetInfo1(info interface{}) {
	e.Info = info
}

func choice1(elem ElementGo, msg *messages.SAMessage, chosen *int, edges []graphs.EdgeExecutable) {
	cases := make([]reflect.SelectCase, len(edges))

	if strings.Contains(elem.ElemId,"evolutive"){
		fmt.Printf("ElementGO:: %v %v\n",edges[0].Info.ActionName, edges[0].Info.Message)
	}

	// Assembly cases
	for i := 0; i < len(edges); i++ {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
		if shared.IsInternal(edges[i].Info.ActionName) {
			// Execute Internal action
			r := false
			edges[i].Info.InternalAction(elem.ElemType, edges[i].Info.ActionName, edges[i].Info.Message, &elem.Info, &r)

			// Update internal channel
			if r {
				go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)
			}
		}
	}

	// Select channel (action)
	var value reflect.Value
	*chosen, value, _ = reflect.Select(cases)
	*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
}

func send1(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}