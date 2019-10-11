package execution

import (
	"gmidarch/development/framework/element"
	"strings"
	"newsolution/shared/shared"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/artefacts/csp"
	"gmidarch/development/framework/messages"
	"newsolution/shared/parameters"
)

type Core struct{}

func (ee Core) Deploy(stateMachines map[string]graphs.GraphExecutable, madlGo madl.MADLGo) error {
	r1 := *new(error)

	for i := range madlGo.Components {
		elemId := madlGo.Components[i].ElemId
		madlGo.Components[i].GoStateMachine = stateMachines[elemId]
		go startElement(madlGo.Components[i])
	}

	for i := range madlGo.Connectors {
		elemId := madlGo.Connectors[i].ElemId
		madlGo.Connectors[i].GoStateMachine = stateMachines[elemId]
		go startElement(madlGo.Connectors[i])
	}
	return r1
}

func startElement(elem element.ElementGo) {
	for {
		//shared.Invoke(elem, "Loop", elem, &elem.GoStateMachine)
		shared.Invoke(elem, "Loop", &elem)
	}
}

func Create(dot csp.DOT, sc map[string]chan messages.SAMessage) (graphs.GraphExecutable, error) {
	r1 := graphs.NewGraph(dot.Dotgraph.NumNodes)
	r2 := *new(error)

	elemId := strings.Replace(dot.SourceDotFile.FileName, parameters.DOT_EXTENSION, "", 99)

	eActions := graphs.EdgeExecutableInfo{}
	var msg messages.SAMessage
	for e1 := range dot.Dotgraph.EdgesDot {
		for e2 := range dot.Dotgraph.EdgesDot [e1] {
			edgeTemp := dot.Dotgraph.EdgesDot[e1][e2]
			actionNameFDR := edgeTemp.Action
			actionNameExec := ""
			if strings.Contains(actionNameFDR, ".") {
				actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
			}
			if shared.IsExternal(actionNameExec) { // External action
				key := strings.ToLower(elemId) + "." + actionNameFDR
				channel, _ := sc[key]
				params := graphs.EdgeExecutableInfo{}
				switch actionNameExec {
				case parameters.INVR:
					invr := channel
					params = graphs.EdgeExecutableInfo{ExternalAction: element.ElementGo{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
				case parameters.TERR:
					terr := channel
					params = graphs.EdgeExecutableInfo{ExternalAction: element.ElementGo{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
				case parameters.INVP:
					invp := channel
					params = graphs.EdgeExecutableInfo{ExternalAction: element.ElementGo{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
				case parameters.TERP:
					terp := channel
					params = graphs.EdgeExecutableInfo{ExternalAction: element.ElementGo{}.TerP, Message: &msg, ActionChannel: &terp, ActionName: actionNameExec}
				}
				mapType := graphs.EdgeExecutableInfo{}
				mapType = params
				eActions = mapType
			}

			if shared.IsInternal(actionNameFDR) {
				msg := messages.SAMessage{}
				channel := make(chan messages.SAMessage)
				params := graphs.EdgeExecutableInfo{InternalAction: shared.Invoke, ActionName: actionNameFDR, Message: &msg, ActionChannel: &channel}
				mapType := params
				eActions = mapType
			}
			r1.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
		}
	}

	return *r1, r2
}
