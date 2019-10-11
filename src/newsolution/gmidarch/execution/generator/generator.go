package generator

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/csp"
	"newsolution/gmidarch/development/artefacts/madl"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"os"
	"strings"
)

type Generator struct{}

func (g Generator) CSP(madl madl.MADL, isEE bool) (csp.CSP) {
	c := csp.CSP{}

	// Solve RUNTIME behaviours
	c.ConfigureProcessBehaviours(madl, isEE)

	// File name
	c.CompositionName = madl.Configuration

	// CSP Data types
	dataTypes := []string{}
	for c := range madl.Components {
		dataTypes = append(dataTypes, madl.Components[c].ElemId)
	}
	for t := range madl.Connectors {
		dataTypes = append(dataTypes, madl.Connectors[t].ElemId)
	}
	c.Datatype = dataTypes

	// Internal Channels
	c.IChannels = c.IdentifyInternalChannels(madl)

	// External Channels
	c.EChannels = c.IdentifyExternalChannels(madl)

	// Processes - Components
	compProcesses := map[string]string{}
	for i := range madl.Components {
		comp := madl.Components[i]
		compId := strings.ToUpper(comp.ElemId)
		subprocesses := strings.Split(comp.Behaviour, "\n")
		if len(subprocesses) > 1 {
			renamedSubprocesses := c.RenameSubprocesses(comp.Behaviour)
			compProcesses[compId] = strings.Replace(renamedSubprocesses, "B", compId, 99)
		} else {
			compProcesses[compId] = strings.Replace(comp.Behaviour, "B", compId, 99)
		}
	}
	c.CompProcesses = compProcesses

	// Processes - Connectors
	connProcesses := map[string]string{}
	for t := range madl.Connectors {
		conn := madl.Connectors[t]
		connId := strings.ToUpper(madl.Connectors[t].ElemId)
		connProcesses[connId] = strings.Replace(conn.Behaviour, "B", connId, 99)
	}
	c.ConnProcesses = connProcesses

	// Composition process - Components/Connectors
	compositionTemp := csp.CompositionProcess{}
	for i := range madl.Components {
		compositionTemp.Components = append(compositionTemp.Components, madl.Components[i].ElemId)
	}
	for i := range madl.Connectors {
		compositionTemp.Connectors = append(compositionTemp.Connectors, madl.Connectors[i].ElemId)
	}

	// Composition Process - Sync ports
	cannonicalNames := map[string]string{}
	for i := range c.EChannels {
		cannonicalName, r2 := c.ToCanonicalName(c.EChannels[i])
		if r2 != nil {
			fmt.Println("CSP:: " + r2.Error())
			os.Exit(0)
		}
		cannonicalNames[cannonicalName] = cannonicalName
	}
	for i := range cannonicalNames {
		compositionTemp.SyncPorts = append(compositionTemp.SyncPorts, cannonicalNames[i])
	}

	// Composition Process - Renaming ports
	eChannels := map[string][]string{}
	for i := range c.ConnProcesses {
		tokens := shared.MyTokenize(c.ConnProcesses[i])
		actions := []string{}
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				actions = append(actions, tokens[j])
			}
			eChannels [i] = actions
		}
	}

	compositionTemp.RenamingPorts = map[string][]csp.Renaming{}
	for i := range eChannels {
		renamingPorts := []csp.Renaming{}
		for j := range eChannels[i] {
			renaming := csp.Renaming{OldName: eChannels[i][j], NewName: c.RenameSyncPort(eChannels[i][j], i)}
			renamingPorts = append(renamingPorts, renaming)
		}
		compositionTemp.RenamingPorts[i] = renamingPorts
	}
	c.Composition = compositionTemp

	// Property
	c.Property = append(c.Property, strings.Replace(parameters.DEADLOCK_PROPERTY, parameters.CORINGA, madl.Configuration, 99))

	return c
}

func (Generator) SaveCSPFile(c csp.CSP) {

	path := parameters.DIR_CSP + "/" + c.CompositionName
	file := c.CompositionName + parameters.CSP_EXTENSION

	// File content
	dataTypeExp := "datatype PROCNAMES = " + shared.StringComposition(c.Datatype, "|", true)

	eChannelExp := "channel " + shared.StringComposition(c.EChannels, ",", false) + " : PROCNAMES"
	iChannelExp := "channel " + shared.StringComposition(c.IChannels, ",", false)
	processesExp := ""
	for i := range c.CompProcesses {
		processesExp += c.CompProcesses[i] + "\n"
	}
	for i := range c.ConnProcesses {
		processesExp += c.ConnProcesses[i] + "\n"
	}

	compositionExp := c.CompositionName + " = (" + strings.ToUpper(shared.StringComposition(c.Composition.Components, "|||", true)+")")
	compositionExp += "[|{|" + shared.StringComposition(c.Composition.SyncPorts, ",", false) + "|}|]"

	renamings := []string{}
	conns := []string{}
	for i := range c.Composition.RenamingPorts {
		for j := range c.Composition.RenamingPorts[i] {
			r := c.Composition.RenamingPorts[i][j].OldName + " <- " + c.Composition.RenamingPorts[i][j].NewName
			renamings = append(renamings, r)
		}
		conns = append(conns, strings.ToUpper(i)+"[["+shared.StringComposition(renamings, ",", false)+"]]")
	}
	compositionExp += "(" + shared.StringComposition(conns, "|||", true) + ")"
	propertyExp := shared.StringComposition(c.Property, "\n", false)

	content := []string{}
	content = append(content, dataTypeExp+"\n")
	content = append(content, eChannelExp+"\n")
	content = append(content, iChannelExp+"\n")
	content = append(content, processesExp+"\n")
	content = append(content, compositionExp+"\n")
	content = append(content, propertyExp)

	// Save file
	shared.SaveFile(path, file, parameters.CSP_EXTENSION, content)
}