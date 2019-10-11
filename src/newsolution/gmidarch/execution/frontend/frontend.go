package frontend

import (
	"fmt"
	"newsolution/gmidarch/execution/checker"
	"newsolution/gmidarch/execution/creator"
	"newsolution/gmidarch/execution/ee"
	"newsolution/gmidarch/execution/generator"
	"newsolution/gmidarch/execution/loader"
	"newsolution/shared/parameters"
	"os"
)

type FrontEnd struct{}

func (f FrontEnd) Deploy(file string) {

	// Read madl and generate architectural artefacts
	l := loader.Loader{}
	mapp := l.Load(file)

	// Create architecture of the execution environment
	kindOfAdaptability := make([]string, 1)
	kindOfAdaptability[0] = "NONE"
	crt := creator.Creator{}
	mee := crt.Create(mapp,kindOfAdaptability)

	crt.Print(mee)
	crt.Save(mee)

	mee2 := l.Load(mee.Configuration+parameters.MADL_EXTENSION)

	fmt.Println(mee2)

	os.Exit(0)

	// Generate CSPs
	generator := generator.Generator{}
	cspapp := generator.CSP(mapp, false)
	cspee := generator.CSP(mee,true)

	fmt.Println(cspee)

	os.Exit(0)

	generator.SaveCSPFile(cspapp)
	generator.SaveCSPFile(cspee)

	// Check CSPs
	chk := checker.Checker{}
	chk.Check(cspapp)
	chk.Check(cspee)

	os.Exit(0)

	// Check CSP
	//checker := checker.Checker{}
	//

	// Start execution environment
	eeTemp := ee.NewEE(mapp)
	eeTemp.Start()

}
