package engine

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"reflect"
	"shared"
)

type Engine struct{}

func (Engine) Execute(elem interface{}, elemInfo []*interface{}, graph graphs.ExecGraph, executionMode bool) {

	node := 0

	// Execute graph
	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			if edges[0].Info.IsInternal { // Internal action
				//fmt.Printf("Engine:: %v :: %v :: BEFORE\n",reflect.TypeOf(elem),edges[0].Info.ActionName)
				edges[0].Info.InternalAction(elem, elemInfo, edges[0].Info.ActionName, edges[0].Info.Message, edges[0].Info.Info)
				//fmt.Printf("Engine:: %v :: %v :: AFTER\n",reflect.TypeOf(elem),edges[0].Info.ActionName)
			} else { // External action
				//fmt.Printf("Engine:: %v :: %v :: BEFORE\n",reflect.TypeOf(elem),edges[0].Info.ActionName)
				edges[0].Info.ExternalAction(edges[0].Info.ActionChannel, edges[0].Info.Message)
				//fmt.Printf("Engine:: %v :: %v :: AFTER\n",reflect.TypeOf(elem),edges[0].Info.ActionName)
			}
			node = edges[0].To
		} else {
			chosen := 0
			choice(elem, elemInfo, &chosen, edges)
			node = edges[chosen].To
		}
		if node == 0 {
			if !shared.EXECUTE_FOREVER {
				break
			}
		}
	}
	return
}

func choice(elem interface{}, elemInfo [] *interface{}, chosen *int, edges []graphs.ExecEdge) {
	casesExternal := make([]reflect.SelectCase, len(edges)+1, shared.NUM_MAX_EDGES)
	casesInternal := make([]reflect.SelectCase, len(edges), shared.NUM_MAX_EDGES)

	// Assembly cases
	//for i := 0; i < len(edges); i++ {
	for i := range edges {
		if edges[i].Info.IsInternal { // Internal action
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
			go send(edges[i].Info.ActionChannel, *edges[i].Info.Message)
		} else { // External action
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		}
	}

	// add default case to choice options
	casesExternal[len(edges)] = reflect.SelectCase{Dir: reflect.SelectDefault}

	// select external first
	var value reflect.Value
	*chosen, value, _ = reflect.Select(casesExternal)

	if *chosen != (len(edges)) { // NOT DEFAULT, i.e., no external action executed
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	} else { // DEFAULT
		*chosen, value, _ = reflect.Select(casesInternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	}
}

func send(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}
