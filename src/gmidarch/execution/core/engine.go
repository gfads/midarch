package engine

import (
	"fmt"
	"gmidarch/development/artefacts/graphs/dot"
	"gmidarch/development/components/component"
	"gmidarch/development/messages"
	"math/rand"
	"reflect"
	"shared"
)

type Engine interface {
	Execute(*component.Component, *bool)
	Stop()
	Resume()
}

type Bufferinfo struct {
	Choice int
	Buffer messages.SAMessage
	Info   interface{}
}

type EngineImpl struct {
	Comp *component.Component
}

func NewEngine() Engine {
	return EngineImpl{}
}

func (EngineImpl) Execute(comp *component.Component, executeForever *bool) {
	node := 0
	fmt.Println("EngineImpl.Execute::Component.Id:", comp.Id)
	if comp.TypeName == "Unit" {
		fmt.Println("EngineImpl.Execute::comp.Info.([]*interface{})[0]).(component.Component).Info:", (*comp.Info.([]*interface{})[0]).(*component.Component).Info)
		//info := (*comp.Info.([]*interface{})[0]).(component.Component).Info
		info := (*comp.Info.([]*interface{})[0]).(*component.Component)
		fmt.Println("EngineImpl.Execute::info is", reflect.TypeOf(info.Type))
		//unit := *comp.Type.(*adaptive.Unit)
		reflect.ValueOf(comp.Type).MethodByName("PrintData").Call([]reflect.Value{})
		//fmt.Println("EngineImpl.Execute::Unit.ElemOfUnit:", unit.ElemOfUnit)
		//fmt.Println("EngineImpl.Execute::Unit.GraphOfElem:", unit.GraphOfElem)
		//fmt.Println("EngineImpl.Execute::Unit.ElemOfUnitInfo:", unit.ElemOfUnitInfo)
	}
	fmt.Println("EngineImpl.Execute::info is", reflect.TypeOf(comp.Type))


	// Execute graph
	for {
		edges := comp.Graph.AdjacentEdges(node)
		if len(edges) == 1 { // Internal action
			if edges[0].Action.IsInternal { // Internal action
				edges[0].Action.InternalAction(comp.Type, comp.Id, edges[0].Action.Name, &comp.Buffer, &comp.Info)
			} else { // External action
				edges[0].Action.ExternalAction(&comp.Buffer, edges[0].Action.Conn, comp.Id, &comp.Info)
			}
			node = edges[0].To // Next node
			fmt.Println("EngineImpl.Execute unique", 0, "comp", comp.Id, "edges[0].To", edges[0].To, "len(edges)", len(edges))
		} else {
			chosen := choice(comp, edges)
			node = edges[chosen].To
			fmt.Println("EngineImpl.Execute chosen", chosen, "comp", comp.Id, "edges[chosen].To", edges[chosen].To, "len(edges)", len(edges))
		}
		if node == 0 {
			if !*executeForever {
				break
			}
		}
	}
	return
}

func choice(comp *component.Component, edges []dot.DOTEdge) int {
	nEdges := len(edges)
	//externalOnly := true
	internalOnly := true
	//r := 0
	cases := make([]reflect.SelectCase, len(edges), len(edges))

	for e := 0; e < nEdges; e++ {
		if edges[e].Action.IsInternal {
			//externalOnly = false
			cases[e] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(edges[e].Action.Conn.Buffer)}
		} else {
			internalOnly = false
			if edges[e].Action.Name == shared.INVP || edges[e].Action.Name == shared.TERR {
				cases[e] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(edges[e].Action.Conn.Buffer)}
			} else {
				cases[e] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(edges[e].Action.Conn.Buffer)} //, Send: reflect.ValueOf(edges[e].Action.)}
			}
		}
	}

	// Only Internal actions in the choice
	if internalOnly {
		start := make(chan bool)
		after := make(chan Bufferinfo)

		// Randomly selects an order for initializing the goroutines
		if rand.Intn(2) == 1 {
			for i := 0; i < nEdges; i++ {
				go func(i int) {
					// Wait for all goroutines
					<-start

					// Make a backup of current buffer/info values
					compBuffer := comp.Buffer
					compInfo := comp.Info

					// Execute action
					edges[i].Action.InternalAction(comp.Type, comp.Id, edges[i].Action.Name, &compBuffer, &compInfo)

					// First action to complete block the others
					after <- Bufferinfo{Choice: i, Buffer: compBuffer, Info: compInfo}
				}(i)
			}
		} else {
			for i := nEdges - 1; i >= 0; i-- {
				go func(i int) {
					// Wait for all go routines
					<-start

					// Make a backup of current buffer/info values
					compBuffer := comp.Buffer
					compInfo := comp.Info

					// Execute action
					edges[i].Action.InternalAction(comp.Type, comp.Id, edges[i].Action.Name, &compBuffer, &compInfo)

					// First action to complete block the others
					after <- Bufferinfo{Choice: i, Buffer: compBuffer, Info: compInfo}
				}(i)
			}
		}

		// Release goroutines to start
		close(start)

		// Wait for completing of a given action
		chosen := <-after

		// Update the values only of the selected operation
		comp.Buffer = chosen.Buffer
		comp.Info = chosen.Info

		return chosen.Choice
	}

	chosen, value, _ := reflect.Select(cases) // External action selection
	//edges[chosen].Action.Conn.Buffer <- value.Interface().(messages.SAMessage)
	comp.Buffer = value.Interface().(messages.SAMessage)
	// Execute action
	if edges[chosen].Action.IsInternal {
		edges[chosen].Action.InternalAction(comp.Type, comp.Id, edges[chosen].Action.Name, &comp.Buffer, &comp.Info)
	} else {
		edges[chosen].Action.ExternalAction(&comp.Buffer, edges[chosen].Action.Conn, comp.Id, &comp.Info)
	}

	return chosen

	//if externalOnly {
	//	shared.ErrorHandler(shared.GetFunction(), "DONE External only in choice")
	//}
	//
	//if !internalOnly && !externalOnly {
	//	shared.ErrorHandler(shared.GetFunction(), "DONE External/Internal choices")
	//}
	//
	//return r
}

func (EngineImpl) Stop()   {}
func (EngineImpl) Resume() {}
