package engine

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"reflect"
	"time"
)

type Engine struct{}

func (Engine) Execute(elem interface{}, elemInfo []*interface{}, graph graphs.ExecGraph, executeForever bool) {

	node := 0

	// Execute graph
	for {
		edges := graph.AdjacentEdges(node)
		if len(edges) == 1 {
			if edges[0].Info.IsInternal { // Internal action
				edges[0].Info.InternalAction(elem, elemInfo, edges[0].Info.ActionName, edges[0].Info.Message, edges[0].Info.Info)
			} else { // External action
				edges[0].Info.ExternalAction(edges[0].Info.ActionChannel, edges[0].Info.Message)
			}
			node = edges[0].To
		} else {
			chosen := 0
			choice(elem, elemInfo, &chosen, edges)
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

func choice(elem interface{}, elemInfo [] *interface{}, chosen *int, edges []graphs.ExecEdge) {
	casesInt := make([]reflect.SelectCase, len(edges))
	casesExt := make([]reflect.SelectCase, len(edges), len(edges)+1)
	hasInternalAction := false
	hasExternalAction := false
	var value reflect.Value

	// Assembly select cases
	for i := 0; i < len(edges); i++ {
		if edges[i].Info.IsInternal { // Internal action
			hasInternalAction = true
		    casesInt[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExt[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		} else { // External action (only InvP, TerR)
			hasExternalAction = true
			casesInt[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			if edges[i].Info.ActionName == "InvP" || edges[i].Info.ActionName == "TerR" {
				casesExt[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			} else {
				casesExt[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel), Send:reflect.ValueOf(*edges[i].Info.Message)}
			}
		}
	}

	if hasExternalAction && !hasInternalAction { // Only external actions
		*chosen, value, _ = reflect.Select(casesExt)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	} else if !hasExternalAction && hasInternalAction { // Only internal actions
		for i := range edges {
			if edges[i].Info.IsInternal {
				edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
				go send(edges[i].Info.ActionChannel, *edges[i].Info.Message) // enable choice by sending something to channel
			}
		}
		*chosen, value, _ = reflect.Select(casesInt)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	} else if hasExternalAction && hasInternalAction {
		casesExt = append(casesExt,reflect.SelectCase{Dir: reflect.SelectDefault}) // append default
		*chosen, value, _ = reflect.Select(casesExt)
		if *chosen != len(edges) { // external action executed (as internal actions are not enabled yet)
			if casesExt[*chosen].Dir == reflect.SelectRecv {
				*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
			}
		} else { // No external action executed then try internal
			for i := range edges {
				if edges[i].Info.IsInternal {
					edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
					go send(edges[i].Info.ActionChannel, *edges[i].Info.Message) // enable choice by sending something to channel
				}
			}
			*chosen, value, _ = reflect.Select(casesInt)
			*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		}
	}
}

func timeout(c chan bool, elemInfo []*interface{}) {
	//time.After(10 * time.Second)
	time.Sleep(10 * time.Second) // slow + change
	//time.Sleep(1 * time.Millisecond)  // fast + no change
	c <- true
}

func send(channel *chan messages.SAMessage, msg messages.SAMessage) {
	*channel <- msg
}

/*
func choiceOld2(elem interface{}, elemInfo [] *interface{}, chosen *int, edges []graphs.ExecEdge) {
	casesExternal := make([]reflect.SelectCase, len(edges))
	casesInternal := make([]reflect.SelectCase, len(edges))
	hasInternalAction := false
	hasExternalAction := false
	var value reflect.Value

	// Assembly select cases
	timeoutChan := make(chan bool)
	for i := 0; i < len(edges); i++ {
		if edges[i].Info.IsInternal { // Internal action
			hasInternalAction = true
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timeoutChan)}
			//edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
			//go send(edges[i].Info.ActionChannel, *edges[i].Info.Message) // enable choice by sending something to channel
		} else { // External action
			hasExternalAction = true
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		}
	}

	if hasExternalAction && !hasInternalAction { // Only external actions
		*chosen, value, _ = reflect.Select(casesExternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	} else if !hasExternalAction && hasInternalAction { // Only internal actions
		for i := range edges {
			if edges[i].Info.IsInternal {
				edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
				go send(edges[i].Info.ActionChannel, *edges[i].Info.Message) // enable choice by sending something to channel
			}
		}
		*chosen, value, _ = reflect.Select(casesInternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	} else if hasExternalAction && hasInternalAction { // External & internal actions
		go timeout(timeoutChan,elemInfo)
		*chosen, value, _ = reflect.Select(casesExternal) // external cases first
		if value.Kind().String() != "bool" { // external action executed
			*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		} else { // No external action executed (timeout)
			for i := range edges {
				if edges[i].Info.IsInternal {
					edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
					go send(edges[i].Info.ActionChannel, *edges[i].Info.Message) // enable choice by sending something to channel
				}
			}
			*chosen, value, _ = reflect.Select(casesInternal)
			*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		}
	}
}
*/

/*

func choiceOld1(elem interface{}, elemInfo [] *interface{}, chosen *int, edges []graphs.ExecEdge) {
	casesExternal := make([]reflect.SelectCase, len(edges), shared.NUM_MAX_EDGES_IN_PARALLEL)
	casesInternal := make([]reflect.SelectCase, len(edges), shared.NUM_MAX_EDGES_IN_PARALLEL)
	hasInternalAction := false
	hasExternalAction := false
	var value reflect.Value

	// Assembly select cases
	for i := 0; i < len(edges); i++ {
		if edges[i].Info.IsInternal { // Internal action
			hasInternalAction = true
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
			edges[i].Info.InternalAction(elem, elemInfo, edges[i].Info.ActionName, edges[i].Info.Message, edges[i].Info.Info)
			go send(edges[i].Info.ActionChannel, *edges[i].Info.Message) // enable choice by sending something to channel
		} else { // External action
			hasExternalAction = true
			casesExternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*edges[i].Info.ActionChannel)}
			casesInternal[i] = reflect.SelectCase{Dir: reflect.SelectRecv}
		}
	}

	if hasExternalAction && !hasInternalAction {
		*chosen, value, _ = reflect.Select(casesExternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	} else if hasExternalAction && hasInternalAction {
		casesExternal = append(casesExternal, reflect.SelectCase{Dir: reflect.SelectDefault}) // add default case to external
		*chosen, value, _ = reflect.Select(casesExternal)                                     // select external cases first
		fmt.Printf("Engine:: Choice:: %v \n", *chosen)
		if *chosen != (len(edges)) { // external action executed
			*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		} else { // No external action executed (DEFAULT case)
			*chosen, value, _ = reflect.Select(casesInternal)
			*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
		}
	} else if !hasExternalAction && hasInternalAction {
		*chosen, value, _ = reflect.Select(casesInternal)
		*edges[*chosen].Info.Message = value.Interface().(messages.SAMessage)
	}
}

*/


