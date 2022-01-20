package adaptive

import (
	"fmt"
	"gmidarch/development/artefacts/graphs/exec"
	"gmidarch/development/messages"
	//	"gmidarch/execution/core/engine"
	"os"
	"reflect"
	"shared"
	"sync"
)

var allUnitsType sync.Map
var allUnitsGraph sync.Map

//@Type: Unit
//@Behaviour: Behaviour = RUNTIME
type Unit struct {
	UnitId         string
	Graph          exec.ExecGraph
	ElemOfUnitInfo [] *interface{}
	ElemOfUnit     interface{}
	GraphOfElem    exec.ExecGraph
}

func NewUnit() Unit {
	r := new(Unit)
	//r.Behaviour = "B = " + shared.RUNTIME_BEHAVIOUR

	return *r
}

func (u Unit) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info *interface{}, r *bool) {

	//fmt.Printf("Unit:: HERE:: %v \n",op, msg)
	switch op[2] {
	case 'E': //"I_Execute":
		elem.(Unit).I_Execute(op, msg, info)
	case 'I': //"I_Initialiseunit":
		elem.(Unit).I_Initialiseunit(op, msg, info)
	case 'A': //"I_Adaptunit":
		elem.(Unit).I_Adaptunit(op, msg, info)
	}
}
//msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Initialiseunit(id string, msg *messages.SAMessage, info *interface{}) {
	allUnitsType.Store(u.UnitId, u.ElemOfUnit)
	allUnitsGraph.Store(u.UnitId, u.GraphOfElem)
}
//msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Execute(id string, msg *messages.SAMessage, info *interface{}) {
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

	u.GraphOfElem = temp.(exec.ExecGraph)
	//engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, !shared.EXECUTE_FOREVER)
	//engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, shared.EXECUTE_FOREVER)
	//engine.EngineImpl{}.Execute(u.ElemOfUnit, shared.EXECUTE_FOREVER)

	return
}
//msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Adaptunit(id string, msg *messages.SAMessage, info *interface{}) {
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

func (u *Unit) changeSelector(s func(interface{}, [] *interface{}, string, *messages.SAMessage, []*interface{}, *bool)) exec.ExecGraph {
	temp, _ := allUnitsGraph.Load(u.UnitId)

	//t1 := time.Now()
	g := temp.(exec.ExecGraph)
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