package engine

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"reflect"
	"shared"
)

type Engine struct{}

var hasInternalAction, hasExternalAction bool
var value reflect.Value

func (Engine) Execute(elem interface{}, elemInfo []*interface{}, graph graphs.ExecGraph, executeForever bool) {

	node := 0

	// Execute graph
	for {
		edges := graph.AdjacentEdges(node)
		//fmt.Printf("Engine:: %v\n",edges)
		if len(edges) == 1 {
			if edges[0].Info.IsInternal { // Internal action
				r := true
				edges[0].Info.InternalAction(elem, elemInfo, edges[0].Info.ActionName, edges[0].Info.Message, edges[0].Info.Info, &r)
			} else { // External action
				edges[0].Info.ExternalAction(edges[0].Info.ActionChannel, edges[0].Info.Message)
			}
			node = edges[0].To
		} else {
			chosen := 0
			selectEdge(elem, elemInfo, &chosen, edges)
			node = edges[chosen].To
		}
		if node == 0 {
			if !executeForever {
				break
			}
		}
	}
	return
}

func selectEdge(elem interface{}, elemInfo [] *interface{}, chosen *int, edges []graphs.ExecEdge) {
	casesInt := make([]reflect.SelectCase, len(edges), len(edges))
	casesExt := make([]reflect.SelectCase, len(edges), len(edges)+1)
	hasInternalAction = false
	hasExternalAction = false

	// Assembly select cases
	for i := range edges {
		if edges[i].Info.IsInternal { // Internal actions
			hasInternalAction = true
			casesInt[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExt[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		} else { // External actions
			hasExternalAction = true
			casesInt[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			if edges[i].Info.ActionName == shared.INVP || edges[i].Info.ActionName == shared.TERR {
				casesExt[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			} else {
				casesExt[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel), Send: reflect.ValueOf(*edges[i].Info.Message)}
			}
		}
	}

	// External actions only
	if hasExternalAction && !hasInternalAction {
		*chosen, value, _ = reflect.Select(casesExt) // External action selection
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)

		return // return from this point
	}

	// Internal actions only
	if !hasExternalAction && hasInternalAction {
		for {
			internalEnabled := false
			for i := range edges {
				if edges[i].Info.IsInternal {
					r := true
					edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info, &r)
					if r {
						internalEnabled = true
						go enableInternalAction(edges[i].Info.ActionChannel, *edges[i].Info.Message)
					}
				}
			}
			if internalEnabled {
				break
			}
		}

		*chosen, value, _ = reflect.Select(casesInt) // Internal action selection
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)

		return // return from this point
	}

	// External and internal actions (external actions first and then internal ones)
	if hasExternalAction && hasInternalAction {
		casesExt = append(casesExt, reflect.SelectCase{Dir: reflect.SelectDefault}) // append default case
		idx := 0
		for {
			*chosen, value, _ = reflect.Select(casesExt) // External action selection
			if *chosen != len(edges) {                   // An external action was selected
				//		fmt.Printf("Engine:: External selected:: %v ********\n",idx)
				if casesExt[*chosen].Dir == reflect.SelectRecv { // InvP and TerR
					*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
				}
				return // return from this point
			}

			// If no external action was selected (default case), then select an internal action
			internalEnabled := false
			for i := range edges {
				if edges[i].Info.IsInternal {
					r := true
					edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info, &r)
					if r {
						internalEnabled = true
						go enableInternalAction(edges[i].Info.ActionChannel, *edges[i].Info.Message)
					}
				}
			}
			if internalEnabled {
				break
			}
			idx++
		}

		*chosen, value, _ = reflect.Select(casesInt) // Internal action selection
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	}

	return
}

func enableInternalAction(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}
