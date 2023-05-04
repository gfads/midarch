package csp

import (
	"github.com/gfads/midarch/pkg/shared"
	"os/exec"
	"strings"
)

type FDRGateway interface {
	Check(CSP)
}

type FDRGatewayImpl struct{}

func NewFDRGateway() FDRGateway {
	return FDRGatewayImpl{}
}
func (FDRGatewayImpl) Check(csp CSP) {
	cmdExp := shared.DIR_FDR + "/" + shared.FDR_COMMAND
	filePath := shared.DIR_CSP + "/" + csp.CompositionName
	fileName := csp.CompositionName + "." + shared.CSP_EXTENSION
	inputFile := filePath + "/" + fileName

	out, err := exec.Command(cmdExp, inputFile).Output()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "File '"+inputFile+"' has a problem (e.g.,syntax error)")
	}
	s := string(out[:])

	if !strings.Contains(s, "Passed") {
		shared.ErrorHandler(shared.GetFunction(), "File '"+inputFile+"' has a deadlock")
	}
}
