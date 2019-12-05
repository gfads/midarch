package exec

import (
	"fmt"
	"gmidarch/development/artefacts/dot"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/components"
	"gmidarch/development/element"
	"gmidarch/development/messages"
	"os"
	"reflect"
	"shared"
	"strings"
)

type Exec struct{}

func (Exec) Create(id string, elem interface{}, typeName string, dot dot.DOTGraph, maps map[string]string, channels map[string]chan messages.SAMessage) (graphs.ExecGraph) {
	r1 := graphs.NewExecGraph(dot.NumNodes)

	// Check dot actions against elem's interface
	checkInterface(elem, id, dot)

	// initialisation of message and info of a given element
	msg := new(messages.SAMessage)
	*msg = messages.SAMessage{Payload: ""} // TODO
	info := make([]*interface{}, 3, 3)     // 3 can be set any value
	for i := 0; i < 3; i++ {
		info[i] = new(interface{})
		*info[i] = new(interface{})
	}

	for e1 := range dot.EdgesDot {
		for e2 := range dot.EdgesDot [e1] {
			eActions := graphs.ExecEdgeInfo{}
			edgeTemp := dot.EdgesDot[e1][e2]
			actionNameFDR := edgeTemp.Action
			actionNameExec := ""
			if strings.Contains(actionNameFDR, ".") {
				actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
			}
			if shared.IsExternal(actionNameExec) { // External action
				actionNameTemp := strings.Split(actionNameFDR, ".")
				key1 := id + "." + actionNameTemp[1]
				key2 := id + "." + actionNameTemp[0] + "." + maps[key1]
				channel, _ := channels[key2]
				params := graphs.ExecEdgeInfo{}
				switch actionNameExec {
				case shared.INVR:
					invr := channel
					params = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvR, ActionName: "InvR", IsInternal: false, Message: msg, ActionChannel: &invr}
				case shared.TERR:
					terr := channel
					params = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerR, ActionName: "TerR", IsInternal: false, Message: msg, ActionChannel: &terr}
				case shared.INVP:
					invp := channel
					params = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.InvP, ActionName: "InvP", IsInternal: false, Message: msg, ActionChannel: &invp}
				case shared.TERP:
					terp := channel
					params = graphs.ExecEdgeInfo{ExternalAction: element.Element{}.TerP, ActionName: "TerP", IsInternal: false, Message: msg, ActionChannel: &terp}
				}
				mapType := graphs.ExecEdgeInfo{}
				mapType = params
				eActions = mapType
			}

			if shared.IsInternal(actionNameFDR) {
				channel := make(chan messages.SAMessage,shared.CHAN_BUFFER_SIZE)

				// Configure selector of each individual element
				s := components.ConfigureSelector(reflect.TypeOf(elem).Name())

				params := graphs.ExecEdgeInfo{InternalAction: s, ActionName: actionNameFDR, IsInternal: true, ActionChannel: &channel, Message: msg, Info: info}
				mapType := params
				eActions = mapType
			}
			r1.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
		}
	}

	return *r1
}

func checkInterface(elem interface{}, id string, dot dot.DOTGraph) {

	// Identify dot actions
	dotActions := []string{}
	for e1 := range dot.EdgesDot {
		for e2 := range dot.EdgesDot [e1] {
			edgeTemp := dot.EdgesDot[e1][e2]
			actionNameFDR := edgeTemp.Action
			if shared.IsInternal(actionNameFDR) {
				dotActions = append(dotActions, actionNameFDR)
			}
		}
	}

	// Identify interface actions
	interfaceActions := []string{}
	for i := 0; i < reflect.TypeOf(elem).NumMethod(); i++ {
		interfaceActions = append(interfaceActions, reflect.TypeOf(elem).Method(i).Name)
	}

	// Check dot actions
	for i := range dotActions {
		found := false
		for j := range interfaceActions {
			if dotActions[i] == interfaceActions[j] {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Exec:: Action '" + dotActions[i] + "' not found in the interface of '" + reflect.TypeOf(elem).String() + "'")
			os.Exit(0)
		}
	}
}
