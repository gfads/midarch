package checker

import (
	"errors"
	"newsolution/shared/parameters"
	"gmidarch/development/framework/configuration/commands"
	"os/exec"
	"strings"
	"gmidarch/development/artefacts/csp"
)

type Checker struct{}

func (Checker) Check(csp csp.CSP) (bool, error) {
	r1 := true
	r2 := *new(error)

	// Check CSP
	isOk, err := check(csp)
	if err != nil {
		r2 = errors.New("EE:: " + err.Error())
		return r1,r2
	}
	if !isOk {
		r1 := isOk
		return r1, r2
	}

	return r1, r2
}

func check(csp csp.CSP) (bool, error) {
	r1 := true
	r2 := *new(error)

	cmdExp := parameters.DIR_FDR + "/" + commands.FDR_COMMAND
	filePath := parameters.DIR_CSP + "/" + csp.CompositionName
	fileName := csp.CompositionName + parameters.CSP_EXTENSION
	inputFile := filePath + "/" + fileName

	out, err := exec.Command(cmdExp, inputFile).Output()
	if err != nil {
		r2 = errors.New("File '" + inputFile + "' has a problem (e.g.,syntax error)")
		return r1, r2
	}
	s := string(out[:])

	if !strings.Contains(s, "Passed") {
		r1 := false
		r2 = errors.New("Deadlock detected")
		return r1, r2
	}
	return r1, r2
}

func (Checker) GenerateDotFiles(csp csp.CSP) (error) {
	r1 := *new(error)

	// Invoke FDR - Generate FDR dot
	// TODO

	return r1
}

/*
func (Checker) LoadDotFiles() {

	// Load dot files
	dotFile := artefacts.DOTFile{}
	dotFiles := map[string]artefacts.DOTFile{}
	dotFiles, err := dotFile.LoadDotFiles(cspMid)
	shared.CheckError(err, "EE")

	// Create dot
	dot := map[string]artefacts.DOT{}
	for i := range dotFiles {
		dot := artefacts.DOT{}
		dot.Create(dotFiles[i])
		dot[i] = dot
	}

	// Create state machines
	stateMachines := map[string]artefacts.GraphExecutable{}
	sm := artefacts.GraphExecutable{}
	for i := range dot {
		sm.Create(dot[i])
		stateMachines[i] = sm
	}
}
*/