package checker

import (
	"fmt"
	"gmidarch/development/framework/configuration/commands"
	"newsolution/gmidarch/development/artefacts/csp"
	"newsolution/shared/parameters"
	"os"
	"os/exec"
	"strings"
)

type FDRGateway struct{}

func (FDRGateway) Check(csp csp.CSP) {
	cmdExp := parameters.DIR_FDR + "/" + commands.FDR_COMMAND
	filePath := parameters.DIR_CSP + "/" + csp.CompositionName
	fileName := csp.CompositionName + parameters.CSP_EXTENSION
	inputFile := filePath + "/" + fileName

	out, err := exec.Command(cmdExp, inputFile).Output()
	if err != nil {
		fmt.Println("CSPGateway:: File '" + inputFile + "' has a problem (e.g.,syntax error)")
		os.Exit(0)
	}
	s := string(out[:])

	if !strings.Contains(s, "Passed") {
		fmt.Println("CSPGateway:: File '" + inputFile + "' has a deadlock")
		os.Exit(0)
	}
}