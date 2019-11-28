package exec

import (
	"fmt"
	"gmidarch/development/artefacts/dot"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/components"
	"gmidarch/development/connectors"
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
				var s components.Selector

				switch reflect.TypeOf(elem).Name() {
				case "Sender":
					s = components.NewSender()
				case "Receiver":
					s = components.NewReceiver()
				case "Unit":
					s = components.NewUnit()
				case "Core":
					s = components.NewCore()
				case "Client":
					s = components.NewClient()
				case "Server":
					s = components.NewServer()
				case "Fibonacciserver":
					s = components.Newfibonacciserver()
				case "FibonacciinvokerM":
					s = components.NewFibonacciInvokerM()
				case "SRH":
					s = components.NewSRH()
				case "Fibonacciclient":
					s = components.NewFibonacciclient()
				case "Fibonacciproxy":
					s = components.NewFibonacciproxy()
				case "RequestorM":
					s = components.NewRequestorM()
				case "CRH":
					s = components.NewCRH()
				case "Namingserver":
					s = components.Newnamingserver()
				case "NaminginvokerM":
					s = components.NewnaminginvokerM()
				case "Namingproxy":
					s = components.NewNamingproxy()
				case "Monevolutive":
					s = components.NewMonevolutive()
				case "Monitor":
					s = components.NewMonitor()
				case "Analyser":
					s = components.NewAnalyser()
				case "Planner":
					s = components.NewPlanner()
				case "Executor":
					s = components.NewExecutor()
				case "Oneto8":
					s = connectors.NewOneto8()
				default:
					fmt.Printf("Exec:: Element '%v' not visible in Engine!!\n", reflect.TypeOf(elem).Name())
					os.Exit(0)
				}

				//s = element.Element{}
				params := graphs.ExecEdgeInfo{InternalAction: s.Selector, ActionName: actionNameFDR, IsInternal: true, ActionChannel: &channel, Message: msg, Info: info}
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
