package checker

import (
	"fmt"
	"gmidarch/development/artefacts/csp"
	"os"
	"os/exec"
	"shared"
	"strings"
)

type FDRGateway struct{}

func (FDRGateway) Check(csp csp.CSP) {
	cmdExp := shared.DIR_FDR + "/" + shared.FDR_COMMAND
	filePath := shared.DIR_CSP + "/" + csp.CompositionName
	fileName := csp.CompositionName + shared.CSP_EXTENSION
	inputFile := filePath + "/" + fileName

	out, err := exec.Command(cmdExp, inputFile).Output()
	if err != nil {
		if err.Error() == "exit status 127" {
			fmt.Println("CSPGateway:: Problem in the execution of the FDR ( File '"+inputFile+"' ). Error: ", err)
		} else {
			fmt.Println("CSPGateway:: File '"+inputFile+"' has a problem (e.g.,syntax error). Error: ", err)
		}
		os.Exit(0)
	}
	s := string(out[:])

	if !strings.Contains(s, "Passed") {
		fmt.Println("CSPGateway:: File '" + inputFile + "' has a deadlock")
		os.Exit(0)
	}
}