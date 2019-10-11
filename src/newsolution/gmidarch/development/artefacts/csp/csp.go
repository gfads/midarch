package csp

import (
	"newsolution/gmidarch/development/artefacts/madl"
	"newsolution/gmidarch/development/components"
	"newsolution/gmidarch/development/connectors"
	"newsolution/shared/parameters"
	"strings"
	"newsolution/shared/shared"
	"errors"
	"reflect"
	"strconv"
	"fmt"
	"os"
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

func (c *CSP) ConfigureProcessBehaviours(madl madl.MADL, isEE bool) {

	// Components
	for i := range madl.Components {
		configuredBehaviour := madl.Components[i].Behaviour

		// The Component has its behaviour defined at runtime
		if strings.Contains(configuredBehaviour, parameters.RUNTIME_BEHAVIOUR) {
			configuredBehaviour = updateRuntimeBehaviourComponents(madl.Components[i].ElemId, madl, isEE) // TODO
		}

		tokens := strings.Split(configuredBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := strings.ToLower(madl.Components[i].ElemId) + "." + strings.ToLower(eX)
				partner, ok := madl.Maps[key]
				if !ok {
					fmt.Println("CSP:: Map [" + key + "] of Component " + madl.Components[i].ElemId + "  Not FOUND!")
					os.Exit(0)
				}
				configuredBehaviour = strings.Replace(configuredBehaviour, eX, partner, 99)
			}
		}
		madl.Components[i].Behaviour = configuredBehaviour
	}

	// Connectors
	for i := range madl.Connectors {
		configuredBehaviour := madl.Connectors[i].Behaviour

		// The connector has its behaviour defined dynamically
		if strings.Contains(configuredBehaviour, parameters.RUNTIME_BEHAVIOUR) { // TODO
			configuredBehaviour = updateRuntimeBehaviourConnectors(madl.Connectors[i].ElemId, madl)
		}

		tokens := strings.Split(configuredBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := strings.ToLower(madl.Connectors[i].ElemId) + "." + strings.ToLower(eX)
				partner, ok := madl.Maps[key]
				if !ok {
					fmt.Println("CSP:: Map [" + key + "] of Connectors Not FOUND!")
					os.Exit(0)
				}
				configuredBehaviour = strings.Replace(configuredBehaviour, eX, partner, 99)
			}
		}
		madl.Connectors[i].Behaviour = configuredBehaviour
	}
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

func (CSP) IdentifyInternalChannels(madl madl.MADL) []string {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range madl.Components {
		tokens := shared.MyTokenize(madl.Components[i].Behaviour)
		for j := range tokens {
			if shared.IsInternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				r1Temp[iAction] = iAction
			}
		}
	}

	for i := range madl.Connectors {
		tokens := shared.MyTokenize(madl.Connectors[i].Behaviour)
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

func (c CSP) IdentifyExternalChannels(madl madl.MADL) []string {
	r1 := []string{}
	r1Temp := map[string]string{}

	for i := range madl.Components {
		tokens := shared.MyTokenize(madl.Components[i].Behaviour)
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				iCannonicalAction, err := c.ToCanonicalName(iAction)
				shared.CheckError(err, "CSP")
				r1Temp[iCannonicalAction] = iCannonicalAction
			}
		}
	}

	for i := range madl.Connectors {
		tokens := shared.MyTokenize(madl.Connectors[i].Behaviour)

		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				iAction := strings.TrimSpace(tokens[j])
				iCannonicalAction, err := c.ToCanonicalName(iAction)
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

func (CSP) ToCanonicalName(name string) (string, error) {
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

func updateRuntimeBehaviourConnectors(connId string, madl madl.MADL) string {
	r1 := ""

	// Define new behaviour
	for i := range madl.Connectors {
		conn := madl.Connectors[i]
		if strings.ToUpper(connId) == strings.ToUpper(conn.ElemId) {
			if reflect.TypeOf(conn.Type) == reflect.TypeOf(connectors.OnetoN{}) {
				n := countAttachments(madl, conn.ElemId)
				r1 = defineNewBehaviour(n, connectors.OnetoN{}, conn.ElemId)
				break
			}
		}
	}
	return r1
}

func updateRuntimeBehaviourComponents(compId string, madl madl.MADL, isEE bool) string {
	r1 := ""

	// Define new behaviour
	for i := range madl.Components {
		comp := madl.Components[i]
		if strings.ToUpper(comp.ElemId) == strings.ToUpper(compId) {
			if reflect.TypeOf(comp.Type) == reflect.TypeOf(components.Core{}) {
				if (strings.ToUpper(madl.Adaptability[0]) == "NONE") { // TODO
					r1 = "B = InvR.e1 -> B"
				} else {
					r1 = "B = InvR.e1 -> P1 \nP1 = InvP.e2 -> InvR.e1 -> P1"
				}
				break
			}

			if reflect.TypeOf(comp.Type) == reflect.TypeOf(components.Unit{}) {
				if strings.ToUpper(madl.Adaptability[0]) == "NONE" { // TODO
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

func countAttachments(madlGo madl.MADL, connectorId string) int {
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
	case reflect.TypeOf(connectors.OnetoN{}).String():
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

func (CSP) RenameSubprocesses(p string) string {
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
