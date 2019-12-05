package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/execution/engine"
	"os"
	"reflect"
	"shared"
	"sync"
)

var allUnitsType sync.Map
var allUnitsGraph sync.Map

type Unit struct {
	UnitId         string
	Behaviour      string
	Graph          graphs.ExecGraph
	ElemOfUnitInfo [] *interface{}
	ElemOfUnit     interface{}
	GraphOfElem    graphs.ExecGraph
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
	allUnitsType.Store(u.UnitId, u.ElemOfUnit)
	allUnitsGraph.Store(u.UnitId, u.GraphOfElem)
}

func (u Unit) I_Execute(msg *messages.SAMessage, info [] *interface{}) {
	var ok bool
	u.ElemOfUnit, ok = allUnitsType.Load(u.UnitId)
	if !ok {
		fmt.Printf("Unit:: Error on acessing the element type")
		os.Exit(0)
	}
	temp, ok := allUnitsGraph.Load(u.UnitId)
	if !ok {
		fmt.Printf("Unit:: Error on acessing the element graph")
		os.Exit(0)
	}
	u.GraphOfElem = temp.(graphs.ExecGraph)
	engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, !shared.EXECUTE_FOREVER)
}

func (u Unit) I_Adaptunit(msg *messages.SAMessage, info [] *interface{}) {
	cmd := msg.Payload.(shared.UnitCommand)

	//fmt.Printf("Unit:: I_Adapt:: %v %v %v\n",reflect.TypeOf(u.ElemOfUnit).Name(),reflect.TypeOf(cmd.Type), u.UnitId)

	if cmd.Cmd != "" {
		unitElemType := reflect.TypeOf(u.ElemOfUnit).Name()
		cmdElemType := reflect.TypeOf(cmd.Type).Name()

		//fmt.Printf("Unit:: Selector:: I_Adaptunit:: [%v] [%v] \n", cmd.Cmd, cmd.Type)

		// Check if the command is to this unit - check by type, i.e., all elements of a given type are adapted
		if unitElemType == cmdElemType {
			if cmd.Cmd == shared.REPLACE_COMPONENT { // TODO
				fmt.Printf("Unit:: **************************** Change happened ****************** \n")
				allUnitsType.Delete(u.UnitId)
				allUnitsType.Store(u.UnitId, cmd.Type) // TODO - Change
				g := u.changeSelector(cmd.Selector)
				allUnitsGraph.Delete(u.UnitId)
				allUnitsGraph.Store(u.UnitId, g)
			} else {
				return
			}
		} else {
			return
		}
	}
}

func (u *Unit) changeSelector(s func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{})) graphs.ExecGraph {
	temp, _ := allUnitsGraph.Load(u.UnitId)

	g := temp.(graphs.ExecGraph)
	for e1 := range g.ExecEdges {
		for e2 := range g.ExecEdges [e1] {
			if g.ExecEdges [e1][e2].Info.IsInternal {
				g.ExecEdges [e1][e2].Info.InternalAction = s
			}
		}
	}
	return g
}
