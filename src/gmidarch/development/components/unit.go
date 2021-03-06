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

func (u Unit) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {

	//	fmt.Printf("Unit:: HERE:: %v \n",op)
	switch op[2] {
	case 'E': //"I_Execute":
		elem.(Unit).I_Execute(msg, info, r)
	case 'I': //"I_Initialiseunit":
		elem.(Unit).I_Initialiseunit(msg, info, r)
	case 'A': //"I_Adaptunit":
		elem.(Unit).I_Adaptunit(msg, info, r)
	}
}

func (u Unit) I_Initialiseunit(msg *messages.SAMessage, info [] *interface{}, r *bool) {
	allUnitsType.Store(u.UnitId, u.ElemOfUnit)
	allUnitsGraph.Store(u.UnitId, u.GraphOfElem)
}

func (u Unit) I_Execute(msg *messages.SAMessage, info [] *interface{}, r *bool) {
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
	//engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, !shared.EXECUTE_FOREVER)
	engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, shared.EXECUTE_FOREVER)

	return
}

func (u Unit) I_Adaptunit(msg *messages.SAMessage, info [] *interface{}, r *bool) {
	cmd := msg.Payload.(shared.UnitCommand)

	//fmt.Printf("Unit:: I_Adapt:: %v [%v] %v\n", reflect.TypeOf(u.ElemOfUnit).Name(), cmd.Cmd, u.UnitId)

	if cmd.Cmd != "" {
		unitElemType := reflect.TypeOf(u.ElemOfUnit).Name()
		cmdElemType := reflect.TypeOf(cmd.Type).Name()

		// Check if the command is to this unit - check by type, i.e., all elements of a given type are adapted
		if unitElemType == cmdElemType {
			if cmd.Cmd == shared.REPLACE_COMPONENT { // TODO
				allUnitsType.LoadOrStore(u.UnitId, cmd.Type)
				g := u.changeSelector(cmd.Selector)
				allUnitsGraph.LoadOrStore(u.UnitId, g)
			} else {
				return
			}
		} else {
			return
		}
	}
}

func (u *Unit) changeSelector(s func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool)) graphs.ExecGraph {
	temp, _ := allUnitsGraph.Load(u.UnitId)

	//t1 := time.Now()
	g := temp.(graphs.ExecGraph)
	for e1 := range g.ExecEdges {
		for e2 := range g.ExecEdges [e1] {
			if g.ExecEdges [e1][e2].Info.IsInternal {
				g.ExecEdges [e1][e2].Info.InternalAction = s
			}
		}
	}
	//fmt.Printf("Unit:: %v\n",time.Now().Sub(t1)/1000000.0)
	return g
}
