package adaptive

import (
	"fmt"
	"gmidarch/development/artefacts/graphs/dot"
	"gmidarch/development/artefacts/graphs/exec"
	"gmidarch/development/components/component"
	"gmidarch/development/components/middleware"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	"gmidarch/execution/core"
	"log"
	"strings"
	"time"

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
	UnitId         	string
	Graph          	dot.DOTGraph
	ElemOfUnitInfo 	interface{} //[] *
	ElemOfUnit     	interface{}
	GraphOfElem    	dot.DOTGraph
}

func NewUnit() Unit {
	r := new(Unit)
	//r.Behaviour = "B = " + shared.RUNTIME_BEHAVIOUR

	return *r
}
func (u Unit) PrintId() {
	fmt.Println("Id:", u.UnitId)
}

func (u Unit) PrintData() {
	fmt.Println("Unit.PrintData::Unit.Id:", u.UnitId)
	fmt.Println("Unit.PrintData::Unit.ElemOfUnit:", u.ElemOfUnit)
	fmt.Println("Unit.PrintData::Unit.GraphOfElem:", u.GraphOfElem)
	fmt.Println("Unit.PrintData::Unit.ElemOfUnitInfo:", u.ElemOfUnitInfo)
}

//func (u Unit) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info *interface{}, r *bool) {
//
//	//fmt.Printf("Unit:: HERE:: %v \n",op, msg)
//	switch op[2] {
//	case 'E': //"I_Execute":
//		elem.(Unit).I_Execute(op, msg, info)
//	case 'I': //"I_Initialiseunit":
//		elem.(Unit).I_Initialiseunit(op, msg, info)
//	case 'A': //"I_Adaptunit":
//		elem.(Unit).I_Adaptunit(op, msg, info)
//	}
//}

//msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Initialiseunit(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	allUnitsType.Store(u.UnitId, u.ElemOfUnit)
	allUnitsGraph.Store(u.UnitId, u.GraphOfElem)
}

//msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Execute(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	fmt.Println("-----------------------------------------> Unit.I_Execute::", u.UnitId, "::TypeName:",(*(*info).([]*interface{})[0]).(*component.Component).TypeName,"::msg.Payload", msg.Payload, "::info:", info)
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

	u.GraphOfElem = temp.(dot.DOTGraph)

	//fmt.Println("Unit.I_Execute::ElemOfUnit is", reflect.TypeOf(u.ElemOfUnit))
	//fmt.Println("Unit.I_Execute::ElemOfUnit kind is", reflect.TypeOf(u.ElemOfUnit).Kind())
	//(*d.Madl.Components[i].Info.([]*interface{})[0]).(component.Component)

	elementComponent := (*(*info).([]*interface{})[0]).(*component.Component)
	fmt.Println("Unit.I_Execute::", u.UnitId, "::info:", elementComponent)
	//fmt.Println("Unit.I_Execute::elementComponent is", reflect.TypeOf(elementComponent))
	//fmt.Println("Unit.I_Execute::elementComponent kind is", reflect.TypeOf(elementComponent).Kind())

	//engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, !shared.EXECUTE_FOREVER)
	//engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, shared.EXECUTE_FOREVER)
	//engine.EngineImpl{}.Execute(u.ElemOfUnit.(*component.Component), shared.EXECUTE_FOREVER)
	fmt.Println(">>>>>>>><<<<<<<<<<<<<>>>>>>>>>>>><<<<<<<<< Unit:", u.UnitId, "TypeName:", elementComponent.TypeName, "executing:", elementComponent.Executing)
	if !elementComponent.Executing {
		elementComponent.ExecuteForever = true
		go engine.EngineImpl{}.Execute(elementComponent, &elementComponent.ExecuteForever) //&shared.ExecuteForever) //shared.EXECUTE_FOREVER)
	}

	return
}

//msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Adaptunit(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	fmt.Println("-----------------------------------------> Unit.I_Adaptunit::", u.UnitId, "::TypeName:",(*(*info).([]*interface{})[0]).(*component.Component).TypeName,"::msg.Payload", msg.Payload, "::info:", info)
	cmd := shared.UnitCommand{}
	if msg.Payload != nil {
		cmd = msg.Payload.(shared.UnitCommand)
	} else {
		fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::msg.Payload->nil")
	}


	//fmt.Printf("Unit:: I_Adapt:: %v [%v] %v\n", reflect.TypeOf(u.ElemOfUnit).Name(), cmd.Cmd, u.UnitId)

	if cmd.Cmd != "" && cmd.Cmd != "Nothing" {
		elementComponent := (*(*info).([]*interface{})[0]).(*component.Component)
		unitElemType := elementComponent.TypeName //reflect.TypeOf(u.ElemOfUnit).Name()
		cmdElemType := reflect.TypeOf(cmd.Type).Name()
		//log.Println("")
		//log.Println("")
		//log.Println("--------------Unit.I_Adaptunit::", u.UnitId, ":: Adapt to ---->", cmdElemType)
		//log.Println("")
		//log.Println("")

		// Check if the command is to this unit - check by type, i.e., all elements of a given type are adapted
		if unitElemType == cmdElemType {
			if cmd.Cmd == shared.REPLACE_COMPONENT { // TODO
				log.Println("")
				log.Println("")
				log.Println("")
				log.Println("")
				log.Println("")
				log.Println("")
				//allUnitsType.LoadOrStore(u.UnitId, cmd.Type)
				//g := u.changeSelector(cmd.Selector)
				//allUnitsGraph.LoadOrStore(u.UnitId, g)
				log.Println("--------------Unit.I_Adaptunit::unitElemType", unitElemType, ":: cmdElemType", cmdElemType)
				//fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::Unit.Type", cmd.Type)
				fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::Unit.Type is", reflect.TypeOf(cmd.Type))

				fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::info:", elementComponent)
				if strings.Contains(unitElemType, "SRH") {
					reset := false

					infoTemp := elementComponent.Info
					srhInfo := infoTemp.(*messages.SRHInfo)
					for _, client := range srhInfo.Clients {
						// if Client from Connection Pool have a client connected
						if client.Ip != "" {
							miop := miop.CreateReqPacket("ChangeProtocol", []interface{}{"tcp"})
							msg := &messages.SAMessage{}
							msg.ToAddr = client.Ip
							log.Println("msg.ToAddr:", msg.ToAddr)
							msg.Payload = middleware.Jsonmarshaller{}.Marshall(miop)
							// Coordinate the protocol change
							shared.MyInvoke(elementComponent.Type, elementComponent.Id, "I_Send", msg, &elementComponent.Info, &reset)
						}
					}
				}

				elementComponent.ExecuteForever = false
				for elementComponent.Executing == true {
					time.Sleep(200 * time.Millisecond)
				}
				elementComponent.Type = cmd.Type

				//infoTemp := make([]*interface{}, 1)
				//infoTemp[0] = new(interface{})
				//*infoTemp[0] = component.Component{Id: u.UnitId, TypeName: reflect.TypeOf(cmd.Type).Name()} //cmd.Type //components[idx]

				//fmt.Println("NewEEDeployer::Unit.Graph", components[idx].Graph)
				//fmt.Println("NewEEDeployer::Unit.Info", components[idx].Info)

				//u.Info = infoTemp //TODO dcruzb: tem que fazer o cara que tem esta unit mudar o info dele para conter o componente a ser criado baseado no tipo que foi recebido no cmd.params
			} else {
				return
			}
		} else {
			return
		}
	} else {
		fmt.Println("Unit::msg.Payload.Cmd->empty")
	}
}

func (u *Unit) changeSelector(s func(interface{}, []*interface{}, string, *messages.SAMessage, []*interface{}, *bool)) exec.ExecGraph {
	fmt.Println("-----------------------------------------> Unit::changeSelector")

	temp, _ := allUnitsGraph.Load(u.UnitId)

	//t1 := time.Now()
	g := temp.(exec.ExecGraph)
	for e1 := range g.ExecEdges {
		for e2 := range g.ExecEdges[e1] {
			if g.ExecEdges[e1][e2].Info.IsInternal { // TODO dcruzb: it is needed te compare the action name too, otherwise it will change all the actions to the last one
				g.ExecEdges[e1][e2].Info.InternalAction = s
			}
		}
	}
	//fmt.Printf("Unit:: %v\n",time.Now().Sub(t1)/1000000.0)
	return g
}
