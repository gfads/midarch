package checker

import (
	"newsolution/gmidarch/development/artefacts/csp"
)

type Checker struct{}

func (Checker) Check(csp csp.CSP) {

	// Use FDR
	fdr := FDRGateway{}
	fdr.Check(csp)
}

