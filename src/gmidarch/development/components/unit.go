package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/execution/engine"
	"reflect"
	"shared"
	"strings"
	"sync"
)

var allUnits sync.Map

type Unit struct {
	UnitId              string
	Behaviour           string
	Graph               graphs.ExecGraph
	ElemOfUnitInfo      [] *interface{}
	ElemOfUnit          interface{}
	GraphOfElem         graphs.ExecGraph
}

func NewUnit() Unit {

	r := new(Unit)
	r.Behaviour = "B = " + shared.RUNTIME_BEHAVIOUR

	return *r
}

func (u Unit) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {

	switch op[2] {
	case 'E': //"I_Execute":
		elem.(Unit).I_Execute(msg, info)
	case 'I': //"I_Initialiseunit":
		elem.(Unit).I_Initialiseunit(msg, info)
	case 'A': //"I_Adaptunit":
		elem.(Unit).I_Adaptunit(msg, info)
	}
}

func (u Unit) I_Initialiseunit(msg *messages.SAMessage, info [] *interface{}) {
	allUnits.Store(u.UnitId, u.ElemOfUnit)
}

func (u Unit) I_Execute(msg *messages.SAMessage, info [] *interface{}) {
	u.ElemOfUnit,_ = allUnits.Load(u.UnitId)
	engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, !shared.EXECUTE_FOREVER)
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
