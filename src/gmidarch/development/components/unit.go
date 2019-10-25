package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/execution/engine"
	"os"
	"reflect"
	"shared"
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
	r.Behaviour = "B = " + shared.RUNTIME_BEHAVIOUR

	return *r
}

func (u Unit) I_Initialiseunit(msg *messages.SAMessage, info [] *interface{}) {
	allUnits.Store(u.UnitId, u.ElemOfUnit)
	fmt.Printf("Unit:: Unit '%v' hosts element '%v' initialised \n",u.UnitId,reflect.TypeOf(u.ElemOfUnit))
}

func (u Unit) I_Execute(msg *messages.SAMessage, info [] *interface{}) {
	newElem, ok := allUnits.Load(u.UnitId)
	if !ok {
		fmt.Printf("Unit:: Element '%v' is not a Unit")
		os.Exit(0)
	}
	u.ElemOfUnit = newElem
	engine.Engine{}.Execute(u.ElemOfUnit, u.GraphOfElem, !shared.EXECUTE_FOREVER)
}

func (u Unit) I_Adaptunit(msg *messages.SAMessage, info [] *interface{}) {
	cmd := msg.Payload.(shared.UnitCommand)

	// Check if the command is to this unit - check by type, i.e., all elements of a given type are adapted
	if cmd.Cmd == shared.REPLACE_COMPONENT { // TODO
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
