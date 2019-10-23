package frontend

import (
	"newsolution/gmidarch/execution/checker"
	"newsolution/gmidarch/execution/creator"
	"newsolution/gmidarch/execution/ee"
	"newsolution/gmidarch/execution/generator"
	"newsolution/gmidarch/execution/loader"
	"newsolution/injector/versioning"
	"newsolution/shared/parameters"
)

type FrontEnd struct{}

func (f FrontEnd) Deploy(file string) {
	l := loader.Loader{}
	crt := creator.Creator{}
	gen := generator.Generator{}
	chk := checker.Checker{}
	inj := versioning.VersioningInjector{}

	// Read MADL and generate architectural artifacts (App)
	mapp := l.Load(file)

	// Create architecture of the execution environment
	appKindOfAdaptability := make([]string, 1)
	appKindOfAdaptability = mapp.Adaptability
	meeTemp := crt.Create(mapp,appKindOfAdaptability)
	crt.Save(meeTemp)

	// Load MADL and generate architectural artefacts (EE)
	mee := l.Load(meeTemp.Configuration+parameters.MADL_EXTENSION)

	// Configure adaptability of EE - according to the adaptability of the hosted App
	mee.AppAdaptability = mapp.Adaptability

	// Generate & save CSPs
	cspapp := gen.CSP(mapp)
	cspee := gen.CSP(mee)
	gen.SaveCSPFile(cspapp)
	gen.SaveCSPFile(cspee)

	// Check CSPs
	chk.Check(cspapp)
	chk.Check(cspee)

	// Deploy App into EE & start EE
	eeEE := ee.NewEE()
	eeEE.DeployApp(mee,mapp)
	eeEE.Start()

	// Start application only - without execution environment
	//eeApp := ee.NewEE()
	//eeApp.Deploy(mapp)
	//eeApp.Start()

	// Start versioning injector
	inj.Start(mapp.Configuration,"receiver")
}
