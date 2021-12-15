package engine

import (
	"fmt"
	"gmidarch/development/artefacts/graphs/dot"
	"gmidarch/development/components/component"
	"gmidarch/development/messages"
	"math/rand"
	"shared"
)

type Engine interface {
	Execute(*component.Component, bool)
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

func (EngineImpl) Execute(comp *component.Component, executeForever bool) {
	node := 0
	fmt.Println("comp", comp)

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
		} else {
			chosen := choice(comp, edges)
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

func choice(comp *component.Component, edges []dot.DOTEdge) int {
	nEdges := len(edges)
	externalOnly := true
	internalOnly := true
	r := 0

	for e := 0; e < nEdges; e++ {
		if edges[e].Action.IsInternal {
			externalOnly = false
		} else {
			internalOnly = false
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

	if externalOnly {
		shared.ErrorHandler(shared.GetFunction(), "TODO External only in choice")
	}

	if !internalOnly && !externalOnly {
		shared.ErrorHandler(shared.GetFunction(), "TODO External/Internal choices")
	}

	return r
}

func (EngineImpl) Stop()   {}
func (EngineImpl) Resume() {}
