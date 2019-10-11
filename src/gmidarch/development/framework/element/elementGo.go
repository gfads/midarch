package element

import (
	"gmidarch/development/artefacts/graphs"
	"newsolution/shared/shared"
	"gmidarch/development/framework/messages"
	"reflect"
	"strings"
	"fmt"
)

type ElementGo struct {
	ElemId          string
	ElemType        interface{}
	CSP             string
	Info            interface{}
	GoStateMachine  graphs.GraphExecutable
	FDRStateMachine graphs.GraphDot
}

//func (ElementGo) Loop(elem ElementGo, graph *graphs.GraphExecutable) {
func (ElementGo) Loop(elem *ElementGo) {

	// Execute graph
	node := 0

	for {
		edges := elem.GoStateMachine.AdjacentEdges(node)
		if len(edges) == 1 {
			edge := edges[0]
			if shared.IsInternal(edge.Info.ActionName) {
				r := false
				if strings.Contains(elem.ElemId, "evolutive") {
					fmt.Printf("ElementGO:: BEFORE %v %v \n", edge.Info.ActionName,edge.Info.Message)
				}
				edge.Info.InternalAction(elem.ElemType, edge.Info.ActionName, edge.Info.Message, &elem.Info, &r)
				if strings.Contains(elem.ElemId, "evolutive") {
					fmt.Printf("ElementGO:: AFTER %v %v\n", edge.Info.ActionName,edge.Info.Message)
				}
			} else {
				edge.Info.ExternalAction(edge.Info.ActionChannel, edge.Info.Message)
				if (strings.Contains(edge.Info.ActionName,"InvR") && strings.Contains(elem.ElemId,"evolutive")){
					fmt.Printf("ElementGo::::::::: %v \n",edge.Info.Message)
				}
			}
			node = edge.To
		} else {
			msg := messages.SAMessage{}
			chosen := 0
			choice(*elem, &msg, &chosen, edges)
			node = edges[chosen].To
		}
		if node == 0 {
			break
		}
	}
	return
}

func choice(elem ElementGo, msg *messages.SAMessage, chosen *int, edges []graphs.EdgeExecutable) {
	cases := make([]reflect.SelectCase, len(edges))

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

func send(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}

func receive(channel *chan messages.SAMessage) {
	<-*channel
}

// external actions common to all components and connectors
func (ElementGo) InvP(invP *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*invP
}

func (ElementGo) InvR(invR *chan messages.SAMessage, msg *messages.SAMessage) {
	*invR <- *msg
}

func (ElementGo) TerR(terR *chan messages.SAMessage, msg *messages.SAMessage) {
	*msg = <-*terR
}

func (ElementGo) TerP(terP *chan messages.SAMessage, msg *messages.SAMessage) {
	*terP <- *msg
}

func DefineOldElement(comp interface{}, newElement interface{}) string {
	found := false
	oldElementId := ""
	components := comp.(map[string]ElementGo)

	// TODO check compatibility of old and new elements by type
	for i := range components {
		oldElementType := reflect.TypeOf(components[i].ElemType).String()
		oldElementType = oldElementType[strings.LastIndex(oldElementType, ".")+1 : len(oldElementType)]
		newElementType := reflect.TypeOf(newElement).String()
		newElementType = newElementType[strings.LastIndex(newElementType, ".")+1 : len(newElementType)]
		if oldElementType == newElementType {
			oldElementId = components[i].ElemId
			found = true
		}
	}

	if !found {
		//myError := errors.MyError{Source: "Planner", Message: "New and old components have not COMPATIBLE types"}
		//myError.ERROR()
	}
	return oldElementId
}

func (e *ElementGo) SetInfo(info interface{}) {
	e.Info = info
}
