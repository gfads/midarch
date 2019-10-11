package csp

import (
	"newsolution/shared/parameters"
	"strings"
	"newsolution/shared/shared"
	"errors"
	"reflect"
	"gmidarch/development/framework/connectors"
	"strconv"
	"fmt"
	"os"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/framework/components"
)

type CompositionProcess struct {
	Components    []string
	Connectors    [] string
	SyncPorts     [] string
	RenamingPorts map[string][]Renaming
}

type Renaming struct {
	OldName string
	NewName string
}

type CSP struct {
	CompositionName string
	Datatype        [] string
	IChannels       []string
	EChannels       []string
	CompProcesses   map[string]string
	ConnProcesses   map[string]string
	Composition     CompositionProcess
	Property        []string
}

func (CSP) Create(madlGo madl.MADLGo, maps map[string]string, kindOfMADL int, midAdaptability []string) (CSP, error) {
	r1 := CSP{}
	r2 := *new(error)

	// Solve RUNTIME behaviours
	r2 = r1.ConfigureProcessBehaviours(madlGo, maps, kindOfMADL, midAdaptability)
	if r2 != nil {
			r2 = errors.New("CSP" + r2.Error())
			return r1, r2
		}

	// File name
	r1.CompositionName = madlGo.ConfigurationName

	// Data type
	dataTypes := []string{}
	for c := range madlGo.Components {
		dataTypes = append(dataTypes, madlGo.Components[c].ElemId)
	}
	for t := range madlGo.Connectors {
		dataTypes = append(dataTypes, madlGo.Connectors[t].ElemId)
	}
	r1.Datatype = dataTypes

	// Internal Channels
	r1.IChannels = identifyInternalChannels(madlGo)

	// External Channels
	r1.EChannels = identifyExternalChannels(madlGo)

	// Processes - Components
	compProcesses := map[string]string{}
	for i := range madlGo.Components {
		comp := madlGo.Components[i]
		compId := strings.ToUpper(comp.ElemId)
		subprocesses := strings.Split(comp.CSP, "\n")
		if len(subprocesses) > 1 {
			renamedSubprocesses := renameSubprocesses(comp.CSP)
			compProcesses[compId] = strings.Replace(renamedSubprocesses, "B", compId, 99)
		} else {
			compProcesses[compId] = strings.Replace(comp.CSP, "B", compId, 99)
		}
	}
	r1.CompProcesses = compProcesses

	// Processes - Connectors
	connProcesses := map[string]string{}
	for t := range madlGo.Connectors {
		conn := madlGo.Connectors[t]
		connId := strings.ToUpper(madlGo.Connectors[t].ElemId)
		connProcesses[connId] = strings.Replace(conn.CSP, "B", connId, 99)
	}
	r1.ConnProcesses = connProcesses

	// Composition process - Components/Connectors
	compositionTemp := CompositionProcess{}
	for i := range madlGo.Components {
		compositionTemp.Components = append(compositionTemp.Components, madlGo.Components[i].ElemId)
	}
	for i := range madlGo.Connectors {
		compositionTemp.Connectors = append(compositionTemp.Connectors, madlGo.Connectors[i].ElemId)
	}

	// Composition Process - Sync ports
	cannonicalNames := map[string]string{}
	for i := range r1.EChannels {
		cannonicalName, r2 := toCanonicalName(r1.EChannels[i])
		if r2 != nil {
			r2 = errors.New("CSP:: " + r2.Error())
			return r1, r2
		}
		cannonicalNames[cannonicalName] = cannonicalName
	}
	for i := range cannonicalNames {
		compositionTemp.SyncPorts = append(compositionTemp.SyncPorts, cannonicalNames[i])
	}

	// Composition Process - Renaming ports
	eChannels := map[string][]string{}
	for i := range r1.ConnProcesses {
		tokens := shared.MyTokenize(r1.ConnProcesses[i])
		actions := []string{}
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				actions = append(actions, tokens[j])
			}
			eChannels [i] = actions
		}
	}

	compositionTemp.RenamingPorts = map[string][]Renaming{}
	for i := range eChannels {
		renamingPorts := []Renaming{}
		for j := range eChannels[i] {
			renaming := Renaming{OldName: eChannels[i][j], NewName: r1.RenameSyncPort(eChannels[i][j], i)}
			renamingPorts = append(renamingPorts, renaming)
		}
		compositionTemp.RenamingPorts[i] = renamingPorts
	}
	r1.Composition = compositionTemp

	// Property
	r1.Property = append(r1.Property, strings.Replace(parameters.DEADLOCK_PROPERTY, parameters.CORINGA, madlGo.ConfigurationName, 99))

	return r1, r2
}

func (c *CSP) ConfigureProcessBehaviours(madlGo madl.MADLGo, maps map[string]string, kindOfMADL int, midAdaptability []string) (error) {
	r1 := *new(error)

	// Components
	for i := range madlGo.Components {
		configuredBehaviour := madlGo.Components[i].CSP

		// The Component has its behaviour defined at runtime
		if strings.Contains(configuredBehaviour, parameters.RUNTIME_BEHAVIOUR) {
			configuredBehaviour = updateRuntimeBehaviourComponents(madlGo.Components[i].ElemId, madlGo, kindOfMADL, midAdaptability)
		}

		tokens := strings.Split(configuredBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := strings.ToLower(madlGo.Components[i].ElemId) + "." + strings.ToLower(eX)
				partner, ok := maps[key]
				if !ok {
					r1 = errors.New("Map [" + key + "] of Component " + madlGo.Components[i].ElemId + "  Not FOUND!")
					return r1
				}
				configuredBehaviour = strings.Replace(configuredBehaviour, eX, partner, 99)
			}
		}
		madlGo.Components[i].CSP = configuredBehaviour
	}

	// Connectors
	for i := range madlGo.Connectors {
		configuredBehaviour := madlGo.Connectors[i].CSP

		// The connector has its behaviour defined dynamically
		if strings.Contains(configuredBehaviour, parameters.RUNTIME_BEHAVIOUR) {
			configuredBehaviour = updateRuntimeBehaviourConnectors(madlGo.Connectors[i].ElemId, madlGo)
		}

		tokens := strings.Split(configuredBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := strings.ToLower(madlGo.Connectors[i].ElemId) + "." + strings.ToLower(eX)
				partner, ok := maps[key]
				if !ok {
					r1 = errors.New("Map [" + key + "] of Connectors Not FOUND!")
					return r1
				}
				configuredBehaviour = strings.Replace(configuredBehaviour, eX, partner, 99)
			}
		}
		madlGo.Connectors[i].CSP = configuredBehaviour
	}

	return r1
}

func (CSP) RenameSyncPort(action string, processId string) string {
	r1 := ""

	action = action [0:strings.Index(action, ".")]
	switch action {
	case parameters.INVP:
		r1 = parameters.INVR + "." + strings.ToLower(processId)
	case parameters.TERP:
		r1 = parameters.INVR + "." + strings.ToLower(processId)
	case parameters.INVR:
		r1 = parameters.INVR + "." + strings.ToLower(processId)
	case parameters.TERR:
		r1 = parameters.INVR + "." + strings.ToLower(processId)
	}
	return r1
}

func identifyInternalChannels(madl madl.MADLGo) []string {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range madl.Components {
		tokens := shared.MyTokenize(madl.Components[i].CSP)
		for j := range tokens {
			if shared.IsInternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range madl.Connectors {
		tokens := shared.MyTokenize(madl.Connectors[i].CSP)
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
	return r1
}

func identifyExternalChannels(madl madl.MADLGo) []string {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range madl.Components {
		tokens := shared.MyTokenize(madl.Components[i].CSP)
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				iCannonicalAction, err := toCanonicalName(iAction)
				shared.CheckError(err, "CSP")
				r1Temp[iCannonicalAction] = iCannonicalAction
			}
		}
	}

	for i := range madl.Connectors {
		tokens := shared.MyTokenize(madl.Connectors[i].CSP)

		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				iCannonicalAction, err := toCanonicalName(iAction)
				shared.CheckError(err, "CSP")
				r1Temp[iCannonicalAction] = iCannonicalAction
			}
		}
	}

	for i := range r1Temp {
		r1 = append(r1, i)
	}
	return r1
}

func toCanonicalName(name string) (string, error) {
	r1 := ""
	r2 := *new(error)

	if strings.Contains(name, parameters.INVR) {
		r1 = parameters.INVR
	}
	if strings.Contains(name, parameters.TERR) {
		r1 = parameters.TERR
	}
	if strings.Contains(name, parameters.INVP) {
		r1 = parameters.INVP
	}
	if strings.Contains(name, parameters.TERP) {
		r1 = parameters.TERP
	}

	if r1 == "" {
		r2 = errors.New("Channel '" + name + "' has NOT a cannonical name.")
	}

	return r1, r2
}

func updateRuntimeBehaviourConnectors(connId string, madlGo madl.MADLGo) string {
	r1 := ""

	// Define new behaviour
	for i := range madlGo.Connectors {
		conn := madlGo.Connectors[i]
		if strings.ToUpper(connId) == strings.ToUpper(conn.ElemId) {
			if reflect.TypeOf(conn.ElemType) == reflect.TypeOf(connectors.OneToN{}) {
				n := countAttachments(madlGo, conn.ElemId)
				r1 = defineNewBehaviour(n, connectors.OneToN{}, conn.ElemId)
				break
			}
		}
	}
	return r1
}

func updateRuntimeBehaviourComponents(compId string, madlGo madl.MADLGo, kindOfMADL int, midAdaptability []string) string {
	r1 := ""

	// Define new behaviour
	for i := range madlGo.Components {
		comp := madlGo.Components[i]
		if strings.ToUpper(comp.ElemId) == strings.ToUpper(compId) {
			if reflect.TypeOf(comp.ElemType) == reflect.TypeOf(components.ExecutionEnvironment{}) {
				if (strings.ToUpper(midAdaptability[0]) == "NONE") { // TODO
					r1 = "B = InvR.e1 -> B"
				} else {
					r1 = "B = InvR.e1 -> P1 \nP1 = InvP.e2 -> InvR.e1 -> P1"
				}
				break
			}

			if reflect.TypeOf(comp.ElemType) == reflect.TypeOf(components.ExecutionUnit{}) {
				if strings.ToUpper(midAdaptability[0]) == "NONE" { // TODO
					r1 = "B = I_InitialiseUnit -> P1\n P1 = I_Execute -> P1"
				} else {
					r1 = "B = InvP.e1 -> I_InitialiseUnit -> P1 \nP1 = I_Execute -> P1 [] InvP.e1 -> I_AdaptUnit -> P1"
				}
				break
			}
		}
	}
	return r1
}

func countAttachments(madlGo madl.MADLGo, connectorId string) int {
	n := 0
	for i := range madlGo.Attachments {
		if madlGo.Attachments[i].T.ElemId == connectorId {
			n ++
		}
	}
	return n
}

func defineNewBehaviour(n int, elem interface{}, elemId string) string {
	baseBehaviour := ""

	switch reflect.TypeOf(elem).String() {
	case reflect.TypeOf(connectors.OneToN{}).String():
		baseBehaviour = strings.ToUpper(elemId) + " = InvP.e1"
		for i := 0; i < n; i++ {
			baseBehaviour += " -> InvR.e" + strconv.Itoa(i+2)
		}
		baseBehaviour += " -> " + strings.ToUpper(elemId)
	default:
		fmt.Println("Configuration:: Impossible to define the new behaviour of " + reflect.TypeOf(elem).String())
		os.Exit(0)
	}
	return baseBehaviour
}

func renameSubprocesses(p string) string {
	subprocesses := strings.Split(p, "\n")
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
