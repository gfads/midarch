package csp

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/artefacts/madl"
	"github.com/gfads/midarch/pkg/shared"
	"strings"
)

type CSPGenerator interface {
	Generate(madl.MADL) CSP
	Save(CSP)
}

type CSPGeneratorImpl struct {
}

func NewCSPGenerator() CSPGenerator {
	return CSPGeneratorImpl{}
}

func (g CSPGeneratorImpl) Generate(madl madl.MADL) CSP {
	csp := CSP{}

	// File name
	csp.CompositionName = madl.Configuration

	// Update Runtime Behaviours
	csp.ConfigureProcessBehaviours(madl)

	// Step 1 - Identify processes
	csp.identifyProcesses(madl)

	// Step 2 - Identify Data types
	csp.identifyDataTypes()

	// Step 3 - identify Internal Channels
	csp.identifyInternalChannels()

	// Step 4 - External Channels
	csp.identifyExternalChannels()

	// Step 5 - Rename component processes
	csp.renameComponentProcesses()

	// Step 6 - Rename connector processes
	csp.renameConnectorProcesses()

	// Step 7 - Create Composition process
	csp.createCompositeProcess()

	// Step 8 - Include Properties
	csp.Property = append(csp.Property, strings.Replace(shared.DEADLOCK_PROPERTY, shared.CORINGA, csp.CompositionName, 99))

	return csp
}

func (CSPGeneratorImpl) Save(c CSP) {

	path := shared.DIR_CSP + "/" + c.CompositionName
	file := c.CompositionName + "." + shared.CSP_EXTENSION

	// Data type
	dataTypeExp := ""
	if len(c.Datatype) > 0 {
		dataTypeExp = "datatype PROCNAMES = " + shared.StringComposition(c.Datatype, "|", true)
	}

	// External channels
	eChannelExp := ""
	if len(c.EChannels) > 0 {
		eChannelExp = "channel " + shared.StringComposition(c.EChannels, ",", false) + " : PROCNAMES"
	}

	// Internal channels
	iChannelExp := ""
	if len(c.IChannels) > 0 {
		iChannelExp = "channel " + shared.StringComposition(c.IChannels, ",", false)
	}

	processesExp := ""
	for i := range c.CompProcesses {
		processesExp += c.CompProcesses[i] + "\n"
	}
	for i := range c.ConnProcesses {
		processesExp += c.ConnProcesses[i] + "\n"
	}

	compositionExp := ""
	if len(c.Composition.Components) > 0 {
		compositionExp = c.CompositionName + " = (" + strings.ToUpper(shared.StringComposition(c.Composition.Components, "|||", true)+")")
	}

	if len(c.Composition.SyncPorts) > 0 {
		compositionExp += "[|{|" + shared.StringComposition(c.Composition.SyncPorts, ",", false) + "|}|]"
	}

	renamings := []string{}
	conns := []string{}
	for i := range c.Composition.RenamingPorts {
		for j := range c.Composition.RenamingPorts[i] {
			r := c.Composition.RenamingPorts[i][j].OldName + " <- " + c.Composition.RenamingPorts[i][j].NewName
			renamings = append(renamings, r)
		}
		conns = append(conns, strings.ToUpper(i)+"[["+shared.StringComposition(renamings, ",", false)+"]]")
	}

	if len(conns) > 0 {
		compositionExp += "(" + shared.StringComposition(conns, "|||", true) + ")"
	}

	propertyExp := ""
	if len(c.Property) > 0 {
		propertyExp = shared.StringComposition(c.Property, "\n", false)
	}

	content := []string{}
	content = append(content, dataTypeExp+"\n")
	content = append(content, eChannelExp+"\n")
	content = append(content, iChannelExp+"\n")
	content = append(content, processesExp+"\n")
	content = append(content, compositionExp+"\n")
	content = append(content, propertyExp)

	// Save file
	shared.SaveFile(path, file, shared.CSP_EXTENSION, content)
}
