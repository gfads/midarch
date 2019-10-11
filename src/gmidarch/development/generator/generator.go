package generator

import (
	"errors"
	"newsolution/shared/parameters"
	"strings"
	"newsolution/shared/shared"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/artefacts/csp"
	"fmt"
)

type Generator struct{}

func (g Generator) GenerateCSP(madlGo madl.MADLGo, maps map[string]string, kindOfMADL int, midAdaptability []string) (csp.CSP, error) {
	r1 := csp.CSP{}
	r2 := *new(error)

	// Generate CSP
	r1, err := generateCSP(madlGo, maps, kindOfMADL, midAdaptability)
	if err != nil {
		r2 = errors.New("Generator:: " + err.Error())
	}

	return r1, r2
}

func generateCSP(madlGo madl.MADLGo, maps map[string]string, kindOfMADL int, midAdaptability []string) (csp.CSP, error) {
	r1 := csp.CSP{}
	r2 := *new(error)

	// Generate CSP Mid
	r1, err := csp.CSP{}.Create(madlGo, maps, kindOfMADL,midAdaptability)
	if err != nil {
		r2 = errors.New("Generator:: " + err.Error())
	}

	return r1, r2
}

func (Generator) GenerateCSPFile(cspSpec csp.CSP) (error) {
	r1 := *new(error)

	// Generate File
	cspFile := csp.CSPFile{}
	// File Path
	cspFile.FilePath = parameters.DIR_CSP + "/" + cspSpec.CompositionName

	// File Name
	cspFile.FileName = cspSpec.CompositionName + parameters.CSP_EXTENSION

	fmt.Print("Generator:: ")
	fmt.Println(cspFile.FileName)

	// File content
	dataTypeExp := "datatype PROCNAMES = " + shared.StringComposition(cspSpec.Datatype, "|", true)

	eChannelExp := "channel " + shared.StringComposition(cspSpec.EChannels, ",", false) + " : PROCNAMES"
	iChannelExp := "channel " + shared.StringComposition(cspSpec.IChannels, ",", false)
	processesExp := ""
	for i := range cspSpec.CompProcesses {
		processesExp += cspSpec.CompProcesses[i] + "\n"
	}
	for i := range cspSpec.ConnProcesses {
		processesExp += cspSpec.ConnProcesses[i] + "\n"
	}

	compositionExp := cspSpec.CompositionName + " = (" + strings.ToUpper(shared.StringComposition(cspSpec.Composition.Components, "|||", true)+")")
	compositionExp += "[|{|" + shared.StringComposition(cspSpec.Composition.SyncPorts, ",", false) + "|}|]"

	renamings := []string{}
	conns := []string{}
	for i := range cspSpec.Composition.RenamingPorts {
		for j := range cspSpec.Composition.RenamingPorts[i] {
			r := cspSpec.Composition.RenamingPorts[i][j].OldName + " <- " + cspSpec.Composition.RenamingPorts[i][j].NewName
			renamings = append(renamings, r)
		}
		conns = append(conns, strings.ToUpper(i)+"[["+shared.StringComposition(renamings, ",", false)+"]]")
	}
	compositionExp += "(" + shared.StringComposition(conns, "|||", true) + ")"
	propertyExp := shared.StringComposition(cspSpec.Property, "\n", false)

	cspFile.Content = append(cspFile.Content, dataTypeExp+"\n")
	cspFile.Content = append(cspFile.Content, eChannelExp+"\n")
	cspFile.Content = append(cspFile.Content, iChannelExp+"\n")
	cspFile.Content = append(cspFile.Content, processesExp+"\n")
	cspFile.Content = append(cspFile.Content, compositionExp+"\n")
	cspFile.Content = append(cspFile.Content, propertyExp)

	// Save file
	r1 = cspFile.Save()
	if r1 != nil {
		r1 = errors.New("Generator:: " + r1.Error())
		return r1
	}
	return r1
}
