package checker

import (
	"fmt"
	csp2 "gmidarch/development/artefacts/csp"
	"os"
	"os/exec"
	"shared/shared"
	"strings"
)

type FDRGateway struct{}

func (FDRGateway) Check(csp csp2.CSP) {
	cmdExp := shared.DIR_FDR + "/" + shared.FDR_COMMAND
	filePath := shared.DIR_CSP + "/" + csp.CompositionName
	fileName := csp.CompositionName + shared.CSP_EXTENSION
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