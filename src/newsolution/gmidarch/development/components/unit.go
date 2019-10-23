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

var allUnits sync.Map

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
	allUnits.Store(u.UnitId, u.ElemOfUnit)
}

func (u Unit) I_Execute(msg *messages.SAMessage, info [] *interface{}) {
	newElem, ok := allUnits.Load(u.UnitId)
	if !ok {
		fmt.Printf("Unit:: Element '%v' is not a Unit")
		os.Exit(0)
	}
	u.ElemOfUnit = newElem
	engine.Engine{}.Execute(u.ElemOfUnit, u.GraphOfElem, !parameters.EXECUTE_FOREVER)
}

func (u Unit) I_Adaptunit(msg *messages.SAMessage, info [] *interface{}) {
	cmd := msg.Payload.(shared.UnitCommand)

	// Check if the command is to this unit - check by type, i.e., all elements of a given type are adapted
	if cmd.Cmd == parameters.REPLACE_COMPONENT { // TODO
		s := strings.Split(reflect.TypeOf(u.ElemOfUnit).String(), ".")
		unitElemType := s[len(s)-1]
		s = strings.Split(reflect.TypeOf(cmd.Type).String(), ".")
		cmdElemType := s[len(s)-1]

		if unitElemType == cmdElemType {
			allUnits.Delete(u.UnitId)
			allUnits.Store(u.UnitId, cmd.Type)
		}
	} else {
		return
	}
}
