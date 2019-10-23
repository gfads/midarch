package components

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/graphs"
	"newsolution/gmidarch/development/messages"
	"newsolution/gmidarch/execution/engine"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"os"
	"reflect"
	"strings"
	"sync"
)

var allGraphs sync.Map

type Unit struct {
	UnitId      string
	Behaviour   string
	Graph       graphs.ExecGraph
	ElemOfUnit  interface{}
	GraphOfElem graphs.ExecGraph
}

func NewUnit() Unit {

	r := new(Unit)
	r.Behaviour = "B = " + parameters.RUNTIME_BEHAVIOUR

	return *r
}

func (u Unit) I_Initialiseunit(msg *messages.SAMessage, info [] *interface{}) {
	allGraphs.Store(u.UnitId, u.GraphOfElem)
}

func (u Unit) I_Execute(msg *messages.SAMessage, info [] *interface{}) {
	temp, _ := allGraphs.Load(u.UnitId)
	u.GraphOfElem = temp.(graphs.ExecGraph)
	engine.Engine{}.Execute(u.ElemOfUnit, u.GraphOfElem, !parameters.EXECUTE_FOREVER)
}

func (u Unit) I_Adaptunit(msg *messages.SAMessage, info [] *interface{}) {
	cmd := msg.Payload.(shared.UnitCommand)

	// Check if the command is to this unit - check by type, i.e., all elements of a given type are adapted
	if cmd.Cmd == "STOP" { // TODO
		s := strings.Split(reflect.TypeOf(u.ElemOfUnit).String(),".")
		unitElemType := s[len(s)-1]
		s = strings.Split(reflect.TypeOf(cmd.Type).String(),".")
		cmdType := s[len(s)-1]

		if unitElemType == cmdType {
			eg := graphs.ExecGraph{}
			newGraph := eg.UpdateGraph(u.GraphOfElem,cmd.Params)
			fmt.Printf("Unit:: %v %v\n", u.GraphOfElem,newGraph)
			os.Exit(0)
			//allGraphs.Store(u.UnitId, u.GraphOfElem)
		}
	} else {
		return
	}
}
