package csp

import (
	"github.com/gfads/midarch/src/gmidarch/development/artefacts/madl"
	"github.com/gfads/midarch/src/gmidarch/development/components/adaptive"
	"github.com/gfads/midarch/src/gmidarch/development/connectors"
	"github.com/gfads/midarch/src/shared"
	"reflect"
	"strconv"
	"strings"
)

type CompositionProcess struct {
	Components    []string
	Connectors    []string
	SyncPorts     []string
	RenamingPorts map[string][]Renaming
}

type Renaming struct {
	OldName string
	NewName string
}

type CSP struct {
	CompositionName string
	Datatype        []string
	IChannels       []string
	EChannels       []string
	CompProcesses   map[string]string
	ConnProcesses   map[string]string
	Composition     CompositionProcess
	Property        []string
}

func (c *CSP) renameComponentProcesses() {
	compProcesses := map[string]string{}

	for i := range c.CompProcesses {
		compBehaviour := c.CompProcesses[i]
		compId := strings.ToUpper(i)

		subprocesses := strings.Split(compBehaviour, "++")

		if len(subprocesses) > 1 {
			renamedSubprocesses := c.renameSubprocesses(compBehaviour)
			compProcesses[compId] = strings.Replace(renamedSubprocesses, shared.BEHAVIOUR_ID, compId, 99)
		} else {
			compProcesses[compId] = strings.Replace(compBehaviour, shared.BEHAVIOUR_ID, compId, 99)
		}
	}
	c.CompProcesses = compProcesses
}

func (c *CSP) renameConnectorProcesses() {
	connProcesses := map[string]string{}

	for t := range c.ConnProcesses {
		connBehaviour := c.ConnProcesses[t]
		connId := strings.ToUpper(t)
		connProcesses[connId] = strings.Replace(connBehaviour, shared.BEHAVIOUR_ID, connId, 99)
	}
	c.ConnProcesses = connProcesses
}

func (c *CSP) createCompositeProcess() {
	compositionTemp := CompositionProcess{}

	// Load component/connector processes
	for i := range c.CompProcesses {
		compositionTemp.Components = append(compositionTemp.Components, i)
	}
	for i := range c.ConnProcesses {
		compositionTemp.Connectors = append(compositionTemp.Connectors, i)
	}

	//Identify sync ports
	cannonicalNames := map[string]string{}
	for i := range c.EChannels {
		cannonicalName := c.toCanonicalName(c.EChannels[i])
		cannonicalNames[cannonicalName] = cannonicalName
	}
	for i := range cannonicalNames {
		compositionTemp.SyncPorts = append(compositionTemp.SyncPorts, cannonicalNames[i])
	}

	// Renaming port
	eChannels := map[string][]string{}
	for i := range c.ConnProcesses {
		tokens := shared.MyTokenize(c.ConnProcesses[i])
		actions := []string{}
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				actions = append(actions, tokens[j])
			}
			eChannels[i] = actions
		}
	}

	compositionTemp.RenamingPorts = map[string][]Renaming{}
	for i := range eChannels {
		renamingPorts := []Renaming{}
		for j := range eChannels[i] {
			renaming := Renaming{OldName: eChannels[i][j], NewName: c.renameSyncPort(eChannels[i][j], i)}
			renamingPorts = append(renamingPorts, renaming)
		}
		compositionTemp.RenamingPorts[i] = renamingPorts
	}
	c.Composition = compositionTemp
}

func (c *CSP) identifyProcesses(madl madl.MADL) {

	compProcesses := map[string]string{}
	for i := range madl.Components {
		compProcesses[madl.Components[i].Id] = madl.Components[i].Behaviour
	}
	c.CompProcesses = compProcesses

	connProcesses := map[string]string{}
	for i := range madl.Connectors {
		connProcesses[madl.Connectors[i].Id] = madl.Connectors[i].Behaviour
	}
	c.ConnProcesses = connProcesses
}

func (c *CSP) identifyInternalChannels() {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range c.CompProcesses {
		tokens := shared.MyTokenize(c.CompProcesses[i])
		for j := range tokens {
			if shared.IsInternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range c.ConnProcesses {
		tokens := shared.MyTokenize(c.ConnProcesses[i])
		for i := range tokens {
			if shared.IsInternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range r1Temp {
		r1 = append(r1, i)
	}
	c.IChannels = r1
}

func (c *CSP) identifyExternalChannels() {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range c.CompProcesses {
		tokens := shared.MyTokenize(c.CompProcesses[i])
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				iCannonicalAction := c.toCanonicalName(iAction)
				r1Temp[iCannonicalAction] = iCannonicalAction
			}
		}
	}

	for i := range c.ConnProcesses {
		tokens := shared.MyTokenize(c.ConnProcesses[i])

		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				iCannonicalAction := c.toCanonicalName(iAction)
				r1Temp[iCannonicalAction] = iCannonicalAction
			}
		}
	}

	for i := range r1Temp {
		r1 = append(r1, i)
	}
	c.EChannels = r1
}

func (c *CSP) identifyDataTypes() {
	dataTypes := []string{}
	for c := range c.CompProcesses {
		dataTypes = append(dataTypes, c)
	}
	for t := range c.ConnProcesses {
		dataTypes = append(dataTypes, t)
	}
	c.Datatype = dataTypes
}

func (CSP) toCanonicalName(name string) string {
	r1 := ""

	if strings.Contains(name, shared.INVR) {
		r1 = shared.INVR
	}
	if strings.Contains(name, shared.TERR) {
		r1 = shared.TERR
	}
	if strings.Contains(name, shared.INVP) {
		r1 = shared.INVP
	}
	if strings.Contains(name, shared.TERP) {
		r1 = shared.TERP
	}

	if r1 == "" {
		shared.ErrorHandler(shared.GetFunction(), "Channel '"+name+"' has NOT a cannonical name.")
	}

	return r1
}

func (CSP) renameSubprocesses(p string) string {
	subprocesses := strings.Split(p, "++")
	r := ""
	procBaseName := strings.TrimSpace(subprocesses[0][:strings.Index(subprocesses[0], "=")]) // first process

	newProcNames := map[string]string{}
	for i := 1; i < len(subprocesses); i++ {
		procName := strings.TrimSpace(subprocesses[i][:strings.Index(subprocesses[i], "=")])
		newProcNames[procName] = procBaseName + procName
	}

	for i := 0; i < len(subprocesses); i++ {
		for j := range newProcNames {
			subprocesses[i] = strings.Replace(subprocesses[i], j, newProcNames[j], 99)
		}
		r += subprocesses[i] + "\n"
	}

	return r
}

func (CSP) renameSyncPort(action string, processId string) string {
	r1 := ""

	action = action[0:strings.Index(action, ".")]
	switch action {
	case shared.INVP:
		r1 = shared.INVR + "." + strings.ToLower(processId)
	case shared.TERP:
		r1 = shared.TERR + "." + strings.ToLower(processId)
	case shared.INVR:
		r1 = shared.INVP + "." + strings.ToLower(processId)
	case shared.TERR:
		r1 = shared.TERP + "." + strings.ToLower(processId)
	}
	return r1
}

func (c CSP) configureRuntimeBehaviourConnector(conn connectors.Connector) string {
	r1 := ""

	switch conn.TypeName {
	case shared.NTOONE: // TODO
		r1 = "B = "
		for i := 0; i < conn.LeftArity; i++ {
			r1 += " InvP.e" + strconv.Itoa(i+1) + " -> InvR.e" + strconv.Itoa(conn.LeftArity+1) + " -> B [] "
		}
		r1 = r1[:strings.LastIndex(r1, "[]")]
	case shared.ONETON: // TODO
		r1 = "B = InvP.e1"
		for i := 0; i < conn.RightArity; i++ {
			r1 += " -> InvR.e" + strconv.Itoa(i+2)
		}
		r1 += " -> B"
	default:
		shared.ErrorHandler(shared.GetFunction(), "Impossible to define the runtime behaviour of connector type '"+conn.TypeName+"'!!")
	}
	return r1
}

/*
func (c CSP) configureRuntimeBehaviourConnectorOld(madl madl.MADL, connId string) string {
	r1 := ""

	// Define new behaviour
	for i := range madl.Connectors {
		conn := madl.Connectors[i]
		fmt.Println(sharedadaptive.GetFunction(), connId, madl.Connectors[i])
		time.Sleep(1000 * time.Millisecond)
		if strings.ToUpper(connId) == strings.ToUpper(conn.Id) {
			if conn.TypeName == sharedadaptive.ONETON { // TODO
				n := c.countAttachments(madl, conn.Id)
				r1 = c.defineNewBehaviour(n, connectors.NewConnector(conn.Id, sharedadaptive.ONETON, conn.Behaviour, 1, 1), connId) // TODO
				break
			}
		}
	}
	return r1
}
*/

func (CSP) configureRuntimeBehaviourComponent(madl madl.MADL, compId string) string {
	r1 := ""

	// Define new behaviour      - TODO
	for i := range madl.Components {
		comp := madl.Components[i]
		if strings.ToUpper(comp.Id) == strings.ToUpper(compId) {
			if reflect.TypeOf(comp.Type) == reflect.TypeOf(adaptive.Core{}) {
				if strings.ToUpper(madl.Adaptability[0]) == "NONE" { // TODO
					r1 = "B = InvR.e1 -> B"
				} else {
					//r1 = "B = InvR.e1 -> P1 \n P1 = InvP.e2 -> I_Debug -> InvR.e1 -> P1"
					//r1 = "B = InvP.e1 -> I_Debug -> InvR.e2 -> P1"
					r1 = "B = InvP.e1 -> InvR.e2 -> P1"
				}
				break
			}

			if reflect.TypeOf(comp.Type) == reflect.TypeOf(adaptive.Unit{}) {
				if strings.ToUpper(madl.Adaptability[0]) == "NONE" { // TODO
					r1 = "B = I_InitialiseUnit -> P1\n P1 = I_Execute -> P1"
				} else {
					r1 = "B = I_Initialiseunit -> P1 \nP1 = I_Execute -> P1 [] InvP.e1 -> I_AdaptUnit -> P1"
				}
				break
			}
		}
	}
	return r1
}

func (c *CSP) ConfigureProcessBehaviours(madl madl.MADL) {
	for i := range madl.Components {
		// The Component has its behaviour defined at runtime
		if strings.Contains(madl.Components[i].Behaviour, shared.RUNTIME_BEHAVIOUR) {
			madl.Components[i].Behaviour = c.configureRuntimeBehaviourComponent(madl, madl.Components[i].Id)
		}
	}
}

func (CSP) countAttachments(madlGo madl.MADL, connectorId string) int {
	n := 0

	for i := range madlGo.Attachments {
		if madlGo.Attachments[i].T.Id == connectorId {
			n++
		}
	}

	if n == 0 {
		shared.ErrorHandler(shared.GetFunction(), "Connector '"+connectorId+"' not found in Attachments")
	}
	return n
}
