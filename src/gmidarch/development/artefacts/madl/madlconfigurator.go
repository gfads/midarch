package madl

import (
	"gmidarch/development/artefacts/graphs/dot"
	"gmidarch/development/components/component"
	"gmidarch/development/connectors"
	"gmidarch/development/messages"
	"gmidarch/development/repositories/architectural"
	"log"
	"net"
	"reflect"
	"shared"
	"strconv"
	"strings"
)

type MADLConfigurator interface {
	ConfigureEE(*MADL, architectural.ArchitecturalRepository, map[string]messages.EndPoint, MADL)
	Configure(*MADL, architectural.ArchitecturalRepository, map[string]messages.EndPoint)
	configureComponents(*MADL, architectural.ArchitecturalRepository, map[string]messages.EndPoint)
	configureConnectors(*MADL, architectural.ArchitecturalRepository)
	configureInfo(*MADL, architectural.ArchitecturalRepository, map[string]messages.EndPoint)
	checkInterface(interface{}, dot.DOTGraph)
}

type MADLConfiguratorImpl struct{
	madl MADL
}

func NewMADLConfigurator() MADLConfigurator {
	var configurator MADLConfigurator
	configurator = &MADLConfiguratorImpl{}

	return configurator
}

func (confImpl MADLConfiguratorImpl) ConfigureEE(m *MADL, archRepo architectural.ArchitecturalRepository, args map[string]messages.EndPoint, madl MADL) {
	confImpl.madl = madl
	confImpl.Configure(m, archRepo, args)
}

// Configure MADL
func (confImpl MADLConfiguratorImpl) Configure(m *MADL, archRepo architectural.ArchitecturalRepository, args map[string]messages.EndPoint) {

	// Step 1 - Configure Connectors
	confImpl.configureConnectors(m, archRepo)

	// Step 2 - Configure ConnMaps, e.g., c1.e1 = {connector}
	confImpl.configureConnMaps(m)

	// Step 3 - Configure Components
	confImpl.configureComponents(m, archRepo, args)
}

func (confImpl MADLConfiguratorImpl) configureComponents(m *MADL, archRepo architectural.ArchitecturalRepository, args map[string]messages.EndPoint) {
	for i := range m.Components {
		// Step 1 - Configure component's info
		confImpl.configureInfo(m, archRepo, args) // TODO dcruzb: this line should be here? This line don't use components, so probably should be outside the for loop

		// Step 2 - Configure type
		record, _ := archRepo.CompLibrary[m.Components[i].TypeName] // type from repositories
		m.Components[i].Type = record.Type

		// Step 3 - Configure Buffer
		m.Components[i].Buffer = messages.SAMessage{}

		// Step 4 - Configure component's behaviour
		confImpl.configureComponentBehaviour(&m.Components[i], *m, archRepo)

		// Step 5 - Configure component's execution graph
		m.Components[i].Graph = dot.DOTLoaderImpl{}.Create(strings.ToLower(m.Components[i].TypeName + "." + shared.DOT_EXTENSION))
		confImpl.configureComponentGraph(&m.Components[i], *m)

		// Step 6 - Check the compatibility between component's interface and DOT actions
		confImpl.checkInterface(m.Components[i].Type, m.Components[i].Graph)
	}
}

func (confImpl MADLConfiguratorImpl) configureConnectors(m *MADL, archRepo architectural.ArchitecturalRepository) {

	for i := range m.Connectors {

		// Step 1 - Define connector's left arity
		m.Connectors[i].LeftArity = m.CountArity(m.Connectors[i].Id, shared.LEFT_ARITY)

		// Step 2 - Define connector's right arity
		m.Connectors[i].RightArity = m.CountArity(m.Connectors[i].Id, shared.RIGHT_ARITY)

		// Step 3 - Define connector's behaviour (from library)
		confImpl.configureConnectorBehaviour(&m.Connectors[i], *m, archRepo)

		// Step 4 - Create a new connector and assign it to set of connectors
		m.Connectors[i] = connectors.NewConnector(m.Connectors[i].Id, m.Connectors[i].TypeName, m.Connectors[i].Behaviour, m.Connectors[i].LeftArity, m.Connectors[i].RightArity)
	}
}

func (confImpl MADLConfiguratorImpl) configureComponentBehaviour(comp *component.Component, m MADL, archRepo architectural.ArchitecturalRepository) {

	// Step 1 - Load the default behaviour of the component (from library)
	comp.Behaviour = archRepo.CompLibrary[comp.TypeName].Behaviour
	comp.Behaviour = strings.Replace(comp.Behaviour, shared.BEHAVIOUR_ID, strings.ToUpper(comp.Id), 99)

	// Step 2 - Update the default behaviour
	partners := map[string]bool{}
	nPartners := 0
	for a := range m.Attachments {
		if m.Attachments[a].C1.Id == comp.Id || m.Attachments[a].C2.Id == comp.Id {
			key := m.Attachments[a].T.Id
			if _, ok := partners[key]; !ok { // New partner found

				// Step 1 - Increment number of partners, e.g., a partner is a connector
				nPartners++

				// Step 2 - Include new partner on map
				partners[key] = true

				// Step 3 - Define 'e' of the behaviour to be replaced by the new partner
				e := "e" + strconv.Itoa(nPartners)

				// Step 4 - Replace 'e' of the behaviour by the new partner
				comp.Behaviour = strings.Replace(comp.Behaviour, e, key, 99)
			}
		}
	}
}

func (confImpl MADLConfiguratorImpl) configureConnectorBehaviour(conn *connectors.Connector, m MADL, archRepo architectural.ArchitecturalRepository) {

	switch conn.TypeName {
	case shared.ONEWAY:
		confImpl.configureOneway(conn, m)
	case shared.REQUEST_REPLY:
		confImpl.configureRequestreply(conn, m)
	case shared.NTOONE:
		confImpl.configureNtoone(conn, m)
	case shared.ONETON:
		confImpl.configureOneton(conn, m)
	case shared.NTOONEREQREP:
		confImpl.configureNtoonetreqrep(conn, m)
	case shared.ONETONREQREP:
		confImpl.configureOnetonreqrep(conn, m)
	default:
		shared.ErrorHandler(shared.GetFunction(), "Behaviour of connector type '"+conn.TypeName+"'is not supported")
	}
}

func (confImpl MADLConfiguratorImpl) configureOnetonreqrep(conn *connectors.Connector, m MADL) {
	b := strings.ToUpper(conn.Id) + " = "

	for a := range m.Attachments {
		att := m.Attachments[a]
		if att.T.Id == conn.Id {
			c1 := att.C1.Id
			c2 := att.C2.Id
			b = b + shared.INVP + "." + c1 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.INVR + "." + c2 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.TERR + "." + c2 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.TERP + "." + c1 + " " + shared.ACTION_PREFIX
			b = b + " " + strings.ToUpper(conn.Id) + " " + shared.ACTION_PREFIX + " "
		}
	}
	b = b[:strings.LastIndex(b, shared.ACTION_PREFIX)]
	conn.Behaviour = b
}

func (confImpl MADLConfiguratorImpl) configureNtoonetreqrep(conn *connectors.Connector, m MADL) {
	b := strings.ToUpper(conn.Id) + " = "

	for a := range m.Attachments {
		att := m.Attachments[a]
		if att.T.Id == conn.Id {
			c1 := att.C1.Id
			c2 := att.C2.Id
			b = b + shared.INVP + "." + c1 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.INVR + "." + c2 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.TERR + "." + c2 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.TERP + "." + c1 + " " + shared.ACTION_PREFIX
			b = b + " " + strings.ToUpper(conn.Id) + " " + shared.CHOICE + " "
		}
	}
	b = b[:strings.LastIndex(b, shared.CHOICE)]
	conn.Behaviour = b
}

func (confImpl MADLConfiguratorImpl) configureOneway(conn *connectors.Connector, m MADL) {

	behaviour := strings.ToUpper(conn.Id) + " = "

	for a := range m.Attachments {
		att := m.Attachments[a]
		if att.T.Id == conn.Id {
			c1 := att.C1.Id
			c2 := att.C2.Id
			behaviour = behaviour + shared.INVP + "." + c1 + " " + shared.ACTION_PREFIX + " " + shared.INVR + "." + c2 + " " + shared.ACTION_PREFIX + " " + strings.ToUpper(conn.Id)
			break
		}
	}
	conn.Behaviour = behaviour
}

func (confImpl MADLConfiguratorImpl) configureRequestreply(conn *connectors.Connector, m MADL) {

	behaviour := strings.ToUpper(conn.Id) + " = "

	for a := range m.Attachments {
		att := m.Attachments[a]
		if att.T.Id == conn.Id {
			c1 := att.C1.Id
			c2 := att.C2.Id
			behaviour = behaviour + shared.INVP + "." + c1 + " " + shared.ACTION_PREFIX + " " + shared.INVR + "." + c2 + shared.ACTION_PREFIX + " "
			behaviour = behaviour + shared.TERR + "." + c2 + " " + shared.ACTION_PREFIX + " " + shared.TERP + "." + c1 + " "
			behaviour = behaviour + shared.ACTION_PREFIX + strings.ToUpper(conn.Id)
			break
		}
	}
	conn.Behaviour = behaviour
}

func (confImpl MADLConfiguratorImpl) configureNtoone(conn *connectors.Connector, m MADL) {

	b := strings.ToUpper(conn.Id) + " = "

	for a := range m.Attachments {
		att := m.Attachments[a]
		if att.T.Id == conn.Id {
			c1 := att.C1.Id
			c2 := att.C2.Id
			b = b + shared.INVP + "." + c1 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.INVR + "." + c2 + " " + shared.ACTION_PREFIX
			b = b + " " + strings.ToUpper(conn.Id) + shared.CHOICE + " "
		}
	}
	b = b[:strings.LastIndex(b, shared.CHOICE)]
	conn.Behaviour = b
}

func (confImpl MADLConfiguratorImpl) configureOneton(conn *connectors.Connector, m MADL) {

	b := strings.ToUpper(conn.Id) + " = "

	for a := range m.Attachments {
		att := m.Attachments[a]
		if att.T.Id == conn.Id {
			c1 := att.C1.Id
			c2 := att.C2.Id
			b = b + shared.INVP + "." + c1 + " " + shared.ACTION_PREFIX
			b = b + " " + shared.INVR + "." + c2 + " " + shared.ACTION_PREFIX
		}
	}
	b = b + " " + strings.ToUpper(conn.Id)
	conn.Behaviour = b
}

func (MADLConfiguratorImpl) configureComponentGraph(comp *component.Component, m MADL) {

	for e1 := range comp.Graph.EdgesDot {
		for e2 := range comp.Graph.EdgesDot [e1] {
			if shared.IsInternal(comp.Graph.EdgesDot[e1][e2].Action.Name) { // Internal action

				// Configure type of action
				comp.Graph.EdgesDot[e1][e2].Action.IsInternal = true

				// Configure function to be invoked according to action from graph
				comp.Graph.EdgesDot[e1][e2].Action.InternalAction = shared.MyInvoke

			} else { // External action

				// Configure type of action
				comp.Graph.EdgesDot[e1][e2].Action.IsInternal = false

				// Configure connector to be used
				key := comp.Id + "." + comp.Graph.EdgesDot[e1][e2].Action.Name[strings.Index(comp.Graph.EdgesDot[e1][e2].Action.Name, ".")+1:]
				_, ok := m.ConnMaps[key]

				if !ok {
					shared.ErrorHandler(shared.GetFunction(), "Behaviour of component '"+comp.Id+"' has a problem: Action '"+key+"' has not a connector associated!!")
				} else {
					comp.Graph.EdgesDot[e1][e2].Action.Conn = m.ConnMaps[key]
				}

				// Rename external action, e.g., InvP.e1 -> InvP
				tempAction := comp.Graph.EdgesDot[e1][e2].Action.Name[:strings.Index(comp.Graph.EdgesDot[e1][e2].Action.Name, ".")]
				comp.Graph.EdgesDot[e1][e2].Action.Name = tempAction

				// Configure function to be invoked according to action from graph
				switch comp.Graph.EdgesDot[e1][e2].Action.Name {
				case shared.INVR:
					comp.Graph.EdgesDot[e1][e2].Action.ExternalAction = comp.InvR
				case shared.INVP:
					comp.Graph.EdgesDot[e1][e2].Action.ExternalAction = comp.InvP
				case shared.TERR:
					comp.Graph.EdgesDot[e1][e2].Action.ExternalAction = comp.TerR
				case shared.TERP:
					comp.Graph.EdgesDot[e1][e2].Action.ExternalAction = comp.TerP
				default:
					shared.ErrorHandler(shared.GetFunction(), "Action '"+comp.Graph.EdgesDot[e1][e2].Action.Name+"' does not exist in Component type")
				}
			}
		}
	}
}

func (confImpl MADLConfiguratorImpl) configureInfo(m *MADL, archRepo architectural.ArchitecturalRepository, args map[string]messages.EndPoint) {

	// Check if the number of parameters in deployment is one of CRH/SRH components (in order)
	n := 0
	unitIndex := 0
	for i := range m.Components {
		aux := m.Components[i].TypeName
		log.Println("configureInfo-> TypeName:", m.Components[i].TypeName)
		if strings.Contains(aux, "SRH") || strings.Contains(aux, "CRH") {
			n++
		}

		if strings.Contains(aux, "Unit")  {
			log.Println("configureInfo-> Unit madl TypeName:", confImpl.madl.Components[unitIndex].TypeName)

			if strings.Contains(confImpl.madl.Components[unitIndex].TypeName, "SRH") {  // m.Components[i].Type.(adaptive.Unit).UnitId, "SRH") {
				log.Println("configureInfo-> Unit TypeName:", m.Components[i].TypeName)
				log.Println("configureInfo-> Unit Id:", m.Components[i].Id)
				log.Println("configureInfo-> Unit Info", m.Components[i].Info)
				//elementComponent := confImpl.madl.Components[unitIndex].Type.(component.Component) //(*(m.Components[i].Info).([]*interface{})[0]).(*component.Component)
				//fmt.Println("EngineImpl.Execute::comp.Info.([]*interface{})[0]).(component.Component).Info:", (*elementComponent.Info.([]*interface{})[0]).(*component.Component).Info)
				//info := (*comp.Info.([]*interface{})[0]).(component.Component).Info
				//info := (*elementComponent.Info.([]*interface{})[0]).(*component.Component)
				//fmt.Println("EngineImpl.Execute::info is", reflect.TypeOf(info.Type))
				n++
			}
			unitIndex++
		}
	}

	if n != len(args) {
		shared.ErrorHandler(shared.GetFunction(), "SRH/CRH endpoints are missing in 'Deploy'")
	}

	// Configure 'info' of components
	for i := range m.Components {
		if v, ok := args[m.Components[i].Id]; ok {
			if strings.Contains(m.Components[i].TypeName, "SRH") {
				endPoint := messages.EndPoint{Host: v.Host, Port: v.Port}
				conns := []net.Conn{}
				rcvedMsgChan := make(chan messages.ReceivedMessages, shared.MAX_NUMBER_OF_RECEIVED_MESSAGES)
				srhInfo := messages.SRHInfo{EndPoint: endPoint, Conns: conns, RcvedMessages: rcvedMsgChan}
				m.Components[i].Info = srhInfo
			} else if strings.Contains(m.Components[i].TypeName, "CRH") {
				conns := make(map[string]net.Conn, shared.MAX_NUMBER_OF_CONNECTIONS)
				endPoint := messages.EndPoint{Host: v.Host, Port: v.Port}
				m.Components[i].Info = messages.CRHInfo{EndPoint: endPoint, Conns: conns}
			}
		} else { // no info
			m.Components[i].Info = new(interface{})
		}
	}
}

func (MADLConfiguratorImpl) checkInterface(elem interface{}, dot dot.DOTGraph) {
	// Identify dot actions
	dotActions := []string{}
	for e1 := range dot.EdgesDot {
		for e2 := range dot.EdgesDot [e1] {
			edgeTemp := dot.EdgesDot[e1][e2]
			actionNameFDR := edgeTemp.Action.Name
			if shared.IsInternal(actionNameFDR) {
				dotActions = append(dotActions, actionNameFDR)
			}
		}
	}

	//if reflect.TypeOf(elem) == reflect.TypeOf(&adaptive.Unit{}) {
	//	elem.(*adaptive.Unit).UnitId = "Teste"
	//	//elem.(*adaptive.Unit).PrintId()
	//	reflect.ValueOf(elem).MethodByName("PrintId").Call([]reflect.Value{})
	//}

	//fmt.Println("elem is", reflect.TypeOf(elem))
	//fmt.Println("elem kind is", reflect.TypeOf(elem).Kind())
	//fmt.Println("elem kind is", reflect.TypeOf(elem).Elem())
	//fmt.Println("elem.Elem() kind is", reflect.TypeOf(elem).Elem().Kind())
	//
	//fmt.Println("elem value", reflect.ValueOf(elem))
	//fmt.Println("elem value.Elem()", reflect.ValueOf(elem).Elem())
	//fmt.Println("elem methods", reflect.ValueOf(elem).Elem().Type().NumMethod())
	//
	//
	//fmt.Println("elem type name is", reflect.ValueOf(elem).Type().Name())
	//fmt.Println("elem type name is", reflect.TypeOf(elem).Elem().Name())
	//fmt.Println("elem has", reflect.ValueOf(elem).Type().NumMethod(), "methods")

	// Identify interface actions
	interfaceActions := []string{}
	for i := 0; i < reflect.TypeOf(elem).NumMethod(); i++ {
		interfaceActions = append(interfaceActions, reflect.TypeOf(elem).Method(i).Name)
	}

	// Check dot actions
	for i := range dotActions {
		found := false
		for j := range interfaceActions {
			if dotActions[i] == interfaceActions[j] {
				found = true
				break
			}
		}
		if !found {
			shared.ErrorHandler(shared.GetFunction(), "Action '"+dotActions[i]+"' not found in the interface of '"+reflect.TypeOf(elem).String()+"'")
		}
	}
}

func (confImpl MADLConfiguratorImpl) configureConnMaps(m *MADL) {

	// Step 1 - Initialise ConnMaps
	connMaps := make(map[string]connectors.Connector)

	// Step 2 - Populate ConnMaps
	for i := range m.Components {
		comp := m.Components[i]
		partners := map[string]bool{}
		nPartners := 0
		for a := range m.Attachments {
			if m.Attachments[a].C1.Id == comp.Id || m.Attachments[a].C2.Id == comp.Id {
				key := m.Attachments[a].T.Id
				if _, ok := partners[key]; !ok { // New partner found

					// Step 1 - Increment number of partners, e.g., a partner is a connector
					nPartners++

					// Step 2 - Include new partner on map
					partners[key] = true

					// Step 3 - Define 'e' of the behaviour to be replaced by the new partner
					e := "e" + strconv.Itoa(nPartners)

					// Step 4 - Replace 'e' of the behaviour by the new partner
					connMaps[comp.Id+"."+e] = confImpl.getConnector(*m, m.Attachments[a].T.Id)
				}
			}
		}
	}
	m.ConnMaps = connMaps
}

func (MADLConfiguratorImpl) getConnector(m MADL, id string) connectors.Connector {
	var r connectors.Connector
	found := false

	for i := range m.Connectors {
		if m.Connectors[i].Id == id {
			found = true
			r = m.Connectors[i]
			break
		}
	}
	if !found {
		shared.ErrorHandler(shared.GetFunction(), "Connector '"+id+"' does not exist in architecture.")
	}
	return r
}

func (confImpl MADLConfiguratorImpl) hasAttachment(m MADL, c1, t, c2 string) bool {
	found := false

	for i := range m.Attachments {
		if m.Attachments[i].C1.Id == c1 && m.Attachments[i].T.Id == t && m.Attachments[i].C2.Id == c2 {
			found = true
			break
		}
	}
	return found
}
