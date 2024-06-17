package adaptive

import (
	"fmt"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/artefacts/graphs/dot"
	"github.com/gfads/midarch/pkg/gmidarch/development/artefacts/graphs/exec"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/component"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/middleware"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	engine "github.com/gfads/midarch/pkg/gmidarch/execution/core"
	"github.com/gfads/midarch/pkg/shared/lib"

	//	"github.com/gfads/midarch/src/gmidarch/execution/core/engine"
	"os"
	"reflect"
	"sync"

	"github.com/gfads/midarch/pkg/shared"
)

var allUnitsType sync.Map
var allUnitsGraph sync.Map

// @Type: Unit
// @Behaviour: Behaviour = RUNTIME
type Unit struct {
	UnitId         string
	Graph          dot.DOTGraph
	ElemOfUnitInfo interface{} //[] *
	ElemOfUnit     interface{}
	GraphOfElem    dot.DOTGraph
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

// msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Initialiseunit(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	allUnitsType.Store(u.UnitId, u.ElemOfUnit)
	allUnitsGraph.Store(u.UnitId, u.GraphOfElem)
}

// msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Execute(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	// lib.PrintlnDebug("-----------------------------------------> Unit.I_Execute::", u.UnitId, "::TypeName:", (*(*info).([]*interface{})[0]).(*component.Component).TypeName, "::msg.Payload", msg.Payload, "::info:", info)
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
	//fmt.Println("Unit.I_Execute::", u.UnitId, "::info:", elementComponent)
	//fmt.Println("Unit.I_Execute::elementComponent is", reflect.TypeOf(elementComponent))
	//fmt.Println("Unit.I_Execute::elementComponent kind is", reflect.TypeOf(elementComponent).Kind())

	//engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, !shared.EXECUTE_FOREVER)
	//engine.Engine{}.Execute(u.ElemOfUnit, u.ElemOfUnitInfo, u.GraphOfElem, shared.EXECUTE_FOREVER)
	//engine.EngineImpl{}.Execute(u.ElemOfUnit.(*component.Component), shared.EXECUTE_FOREVER)
	//fmt.Println(">>>>>>>><<<<<<<<<<<<<>>>>>>>>>>>><<<<<<<<< Unit:", u.UnitId, "TypeName:", elementComponent.TypeName, "executing:", elementComponent.Executing)
	if elementComponent.Executing == nil || !*elementComponent.Executing {
		// lib.PrintlnDebug("Will execute elementComponent.TypeName:", elementComponent.TypeName)
		var executeForever = true
		elementComponent.ExecuteForever = &executeForever
		//fmt.Println("Setará executeforever:", elementComponent.TypeName)
		if strings.Contains(elementComponent.TypeName, "SRH") {
			//fmt.Println("Setou executeforever")
			//log.Println("Setou executeforever")
			infoTemp := elementComponent.Info
			srhInfo := infoTemp.(*messages.SRHInfo)
			srhInfo.ExecuteForever = elementComponent.ExecuteForever
		}
		go engine.EngineImpl{}.Execute(elementComponent, elementComponent.ExecuteForever) //&shared.ExecuteForever) //shared.EXECUTE_FOREVER)
	} // TODO dcruzb: add sleep no else do executing para melhorar consumo de CPU

	return
}

// msg *messages.SAMessage, info [] *interface{}, r *bool
func (u Unit) I_Adaptunit(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("-----------------------------------------> Unit.I_Adaptunit::", u.UnitId, "::TypeName:",(*(*info).([]*interface{})[0]).(*component.Component).TypeName,"::msg.Payload", msg.Payload, "::info:", info)
	cmd := shared.UnitCommand{}
	if msg.Payload != nil {
		cmd = msg.Payload.(shared.UnitCommand)
	} else {
		//fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::msg.Payload->nil")
	}

	// fmt.Printf("Unit:: I_Adapt:: %v [%v] %v\n", reflect.TypeOf(u.ElemOfUnit).Name(), cmd.Cmd, u.UnitId)

	if cmd.Cmd != "" && cmd.Cmd != "Nothing" {
		elementComponent := (*(*info).([]*interface{})[0]).(*component.Component)
		// lib.PrintlnInfo("--------------Unit.I_Adaptunit::", u.UnitId, ":: Adapt ---->", elementComponent.TypeName)
		unitElemType := elementComponent.TypeName //reflect.TypeOf(u.ElemOfUnit).Name()
		cmdElemType := reflect.ValueOf(cmd.Type).Elem().Type().Name()
		//log.Println("")
		//log.Println("")
		// lib.PrintlnInfo("--------------Unit.I_Adaptunit::", u.UnitId, ":: Adapt to ---->", cmdElemType)
		//log.Println("")
		//log.Println("")

		// Check if the command is to this unit - check by type, i.e., all elements of a given type are adapted
		if shared.CompatibleComponents(strings.ToUpper(unitElemType), strings.ToUpper(cmdElemType)) {
			if cmd.Cmd == shared.REPLACE_COMPONENT { // TODO
				//log.Println("")
				//log.Println("")
				//log.Println("")
				//log.Println("")
				//log.Println("")
				//log.Println("")
				//allUnitsType.LoadOrStore(u.UnitId, cmd.Type)
				//g := u.changeSelector(cmd.Selector)
				//allUnitsGraph.LoadOrStore(u.UnitId, g)
				// lib.PrintlnInfo("--------------Unit.I_Adaptunit::unitElemType(from)", unitElemType, ":: cmdElemType(to)", cmdElemType)
				//fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::Unit.Type", cmd.Type)
				//fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::Unit.Type is", reflect.TypeOf(cmd.Type))

				//fmt.Println("Unit.I_Adaptunit::", u.UnitId, "::info:", elementComponent)
				var adaptTo string
				if strings.Contains(cmdElemType, "TCP") {
					adaptTo = "tcp"
					// lib.PrintlnInfo("****** Adapt to TCP")
				} else if strings.Contains(cmdElemType, "UDP") {
					adaptTo = "udp"
					// lib.PrintlnInfo("****** Adapt to UDP")
				} else if strings.Contains(cmdElemType, "TLS") {
					adaptTo = "tls"
					// lib.PrintlnInfo("****** Adapt to TLS")
				} else if strings.Contains(cmdElemType, "QUIC") {
					adaptTo = "quic"
					// lib.PrintlnInfo("****** Adapt to QUIC")
				} else if strings.Contains(cmdElemType, "RPC") {
					adaptTo = "rpc"
					// lib.PrintlnInfo("****** Adapt to RPC")
				} else if strings.Contains(cmdElemType, "HTTP2") {
					adaptTo = "http2"
					// lib.PrintlnInfo("****** Adapt to HTTP2")
				} else if strings.Contains(cmdElemType, "HTTPS") {
					adaptTo = "https"
					// lib.PrintlnInfo("****** Adapt to HTTPS")
				} else if strings.Contains(cmdElemType, "HTTP") {
					adaptTo = "http"
					// lib.PrintlnInfo("****** Adapt to HTTP")
				}

				var adaptFrom string
				if strings.Contains(unitElemType, "TCP") {
					adaptFrom = "tcp"
					// lib.PrintlnInfo("****** Adapt from TCP")
				} else if strings.Contains(unitElemType, "UDP") {
					adaptFrom = "udp"
					// lib.PrintlnInfo("****** Adapt from UDP")
				} else if strings.Contains(unitElemType, "TLS") {
					adaptFrom = "tls"
					// lib.PrintlnInfo("****** Adapt from TLS")
				} else if strings.Contains(unitElemType, "QUIC") {
					adaptFrom = "quic"
					// lib.PrintlnInfo("****** Adapt from QUIC")
				} else if strings.Contains(unitElemType, "RPC") {
					adaptFrom = "rpc"
					// lib.PrintlnInfo("****** Adapt from RPC")
				} else if strings.Contains(unitElemType, "HTTP2") {
					adaptFrom = "http2"
					// lib.PrintlnInfo("****** Adapt from HTTP2")
				} else if strings.Contains(unitElemType, "HTTPS") {
					adaptFrom = "https"
					// lib.PrintlnInfo("****** Adapt from HTTPS")
				} else if strings.Contains(unitElemType, "HTTP") {
					adaptFrom = "http"
					// lib.PrintlnInfo("****** Adapt from HTTP")
				}

				isSRH := strings.Contains(cmdElemType, "SRH")
				isCRH := strings.Contains(cmdElemType, "CRH")

				if isSRH {
					reset := false

					infoTemp := elementComponent.Info
					srhInfo := infoTemp.(*messages.SRHInfo)
					for idx, client := range srhInfo.Clients { // TODO dcruzb: probably no protocol use srhInfo.Clients anymore. Verify and remove
						//fmt.Println("Vai adaptar")
						// if Client from Connection Pool have a client connected
						if client.Ip != "" {
							//fmt.Println("Vai adaptar: IP:", client.Ip)
							if (adaptFrom == "udp" && client.UDPConnection == nil) ||
								(adaptFrom == "tcp" && client.Connection == nil) ||
								(adaptFrom == "tls" && client.Connection == nil) ||
								(adaptFrom == "quic" && client.QUICStream == nil) {
								//fmt.Println("Vai adaptar: pulou sem conexão")
								continue
							}
							//fmt.Println("Vai adaptar: entrou AdaptId:", client.AdaptId)
							client.AdaptId = idx
							miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{adaptTo, client.AdaptId}, client.AdaptId) // idx is the Connection ID
							msg := &messages.SAMessage{}
							msg.ToAddr = client.Ip
							//log.Println("msg.ToAddr:", msg.ToAddr)
							msg.Payload = middleware.Gobmarshaller{}.Marshall(miopPacket)
							// Coordinate the protocol change
							shared.MyInvoke(elementComponent.Type, elementComponent.Id, "I_Send", msg, &elementComponent.Info, &reset)
						}
					}
					if srhInfo.Protocol != nil {
						for idx, client := range srhInfo.Protocol.GetClients() {
							// fmt.Println("Vai adaptar IP:", (*client).Address())
							// if Client from Connection Pool have a client connected
							if adaptFrom == "rpc" || adaptFrom == "http" || adaptFrom == "https" || adaptFrom == "http2" || (*client).Address() != "" {
								// fmt.Println("Vai adaptar: IP:", (*client).Address())
								// if (adaptFrom == "udp" && client.UDPConnection == nil) ||
								// 	(adaptFrom == "tcp" && client.Connection == nil) ||
								// 	(adaptFrom == "quic" && client.QUICStream == nil) {
								// 	//fmt.Println("Vai adaptar: pulou sem conexão")
								// 	continue
								// }
								// fmt.Println("Vai adaptar: entrou AdaptId:", (*client).AdaptId())
								(*client).SetAdaptId(idx)
								miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{adaptTo, (*client).AdaptId()}, (*client).AdaptId()) // idx is the Connection ID
								msg := &messages.SAMessage{}
								msg.ToAddr = (*client).Address()
								// log.Println("msg.ToAddr:", msg.ToAddr)
								msg.Payload = middleware.Gobmarshaller{}.Marshall(miopPacket)
								// Coordinate the protocol change
								shared.MyInvoke(elementComponent.Type, elementComponent.Id, "I_Send", msg, &elementComponent.Info, &reset)
							}
						}
					}
					time.Sleep(200 * time.Millisecond)
				} // else if isCRH {
				//time.Sleep(10 * time.Second)
				//fmt.Println("Unit.I_Adaptunit:: 10 seconds passed", u.UnitId, "::info:", elementComponent)
				//cmd.Type = shared.GetComponentTypeByNameFromRAM(unitElemType)
				//fmt.Println("unitElemType", unitElemType, "cmd.Type", cmd.Type)
				//shared.ErrorHandler(shared.GetFunction(), "Teste")
				//}

				*elementComponent.ExecuteForever = false
				for *elementComponent.Executing == true {
					// lib.PrintlnInfo("Awaiting to stop executing")
					time.Sleep(200 * time.Millisecond)
				}
				// lib.PrintlnInfo("Execution stopped")
				// lib.PrintlnInfo("****** elementComponent.TypeName:", elementComponent.TypeName)
				// lib.PrintlnInfo("****** cmdElemType:", cmdElemType)
				lib.PrintlnInfo("****** adapt From:", adaptFrom, "=> To:", adaptTo)
				//time.Sleep(6 * time.Second)
				elementComponent.Type = cmd.Type
				elementComponent.TypeName = cmdElemType

				if isCRH {
					//time.Sleep(2000 * time.Millisecond)
					//fmt.Println("Unit.I_Adaptunit:: 2 seconds passed", u.UnitId) //, "::info:", elementComponent)
					// lib.PrintlnInfo("Will close CRH connections")
					infoTemp := elementComponent.Info
					crhInfo := infoTemp.(messages.CRHInfo)
					for idx, protocol := range crhInfo.Protocols {
						if protocol != nil {
							protocol.CloseConnection()
						}
						crhInfo.Protocols[idx] = nil
						delete(crhInfo.Protocols, idx)
					}

					for _, conn := range crhInfo.Conns {
						conn.Close()
					}
					// if adaptTo != "quic" {
					for _, conn := range crhInfo.QuicConns {
						conn.CloseWithError(0, "Close to adapt to "+adaptTo)
					}
					for _, stream := range crhInfo.QuicStreams {
						stream.Close()
					}
					// }
					// lib.PrintlnInfo("CRH connections closed")
					//shared.MyInvoke(elementComponent.Type, elementComponent.Id, "I_Process", msg, &elementComponent.Info, reset)
				} else if isSRH {
					infoTemp := elementComponent.Info
					srhInfo := infoTemp.(*messages.SRHInfo)
					for len(srhInfo.Clients) > 0 {
						// lib.PrintlnInfo("Will initialize:", srhInfo.Clients[len(srhInfo.Clients)-1].Ip)
						tmpClient := srhInfo.Clients[len(srhInfo.Clients)-1]
						srhInfo.Clients = messages.Remove(srhInfo.Clients, len(srhInfo.Clients)-1)
						tmpClient.Initialize()
						// lib.PrintlnInfo("Initialized")
					}
					if srhInfo.Protocol != nil {
						// lib.PrintlnInfo("Will stop server")
						srhInfo.Protocol.StopServer()
						// lib.PrintlnInfo("Server stoped")
						srhInfo.Protocol = nil
					}
				}

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
		//fmt.Println("Unit::msg.Payload.Cmd->empty")
	}
}

func (u *Unit) changeSelector(s func(interface{}, []*interface{}, string, *messages.SAMessage, []*interface{}, *bool)) exec.ExecGraph {
	//fmt.Println("-----------------------------------------> Unit::changeSelector")

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
