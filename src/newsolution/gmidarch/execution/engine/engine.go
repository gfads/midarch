package engine

import (
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/parameters"
	"reflect"
)

type Engine struct{}

func (Engine) Execute(elem interface{}, graph graphs.ExecGraph, executionMode bool) {

	// Execute graph
	node := 0

	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			edge := edges[0]
			switch edge.Info.ActionType {
			case 1:
				edge.Info.InternalAction(elem, edge.Info.ActionName, edge.Info.Message, edge.Info.Info)
			case 2:
				edge.Info.ExternalAction(edge.Info.ActionChannel, edge.Info.Message)
			}
			node = edge.To
		} else {
			chosen := 0
			choice(elem, &chosen, edges)
			node = edges[chosen].To
		}
		if node == 0 && executionMode != parameters.EXECUTE_FOREVER {
			break
		}
	}
	return
}

func choice(elem interface{}, chosen *int, edges []graphs.ExecEdge) {
	casesExternal := make([]reflect.SelectCase, len(edges)+1)
	casesInternal := make([]reflect.SelectCase, len(edges))

	// Assembly cases
	for i := 0; i < len(edges); i++ {
		switch edges[i].Info.ActionType {
		case 1: // Internal
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			edges[i].Info.InternalAction(elem, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
			go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)

		case 2: // External
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		}
	}
	// default case
	casesExternal[len(edges)] = reflect.SelectCase{Dir: reflect.SelectDefault}

	// select external first
	var value reflect.Value
	*chosen, value, _ = reflect.Select(casesExternal)

	if *chosen != (len(edges)) { // NOT DEFAULT
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		return
	} else {
		*chosen, value, _ = reflect.Select(casesInternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	}
}

func send(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}
