package components

import (
	"gmidarch/development/framework/messages"
	"gmidarch/development/framework/element"
	"gmidarch/shared/shared"
)

type ExecutionUnit struct {
}

var unitMsg messages.SAMessage

// Unit initialization
func (unit ExecutionUnit) I_InitialiseUnit(msg *messages.SAMessage, info *interface{}, r *bool) {
	elem := (*info).(element.ElementGo)
	unitMsg = messages.SAMessage{}
	for e1 := range elem.GoStateMachine.EdgesExecutable {
		for e2 := range elem.GoStateMachine.EdgesExecutable[e1] {
			elem.GoStateMachine.EdgesExecutable[e1][e2].Info.Message = &unitMsg
		}
	}
	*r = true
}

// Unit execution
func (unit ExecutionUnit) I_Execute(msg *messages.SAMessage, info *interface{}, r *bool) {
	elem := (*info).(element.ElementGo)

	//shared.Invoke(elem, "Loop", elem, &elem.GoStateMachine)
	shared.Invoke(elem, "Loop", elem, elem.GoStateMachine)
	*r = true
}

// Unit adaptation
func (unit ExecutionUnit) I_AdaptUnit(msg *messages.SAMessage, info *interface{}, r *bool) {
}
