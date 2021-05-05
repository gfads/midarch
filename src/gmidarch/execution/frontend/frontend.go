package frontend

import (
	"fmt"
	"gmidarch/execution/creator"
	"gmidarch/execution/deployer"
	"gmidarch/execution/generator"
	"gmidarch/execution/loader"
	"os"
	"shared"
)

type FrontEnd struct{}

func (f FrontEnd) Deploy(file string) {
	l := loader.Loader{}
	crt := creator.Creator{}
	gen := generator.Generator{}
	//chk := checker.Checker{}
	dep := deployer.NewEE()

	// Read MADL and generate architectural artifacts (App)
	mapp := l.Load(file)

	switch mapp.Adaptability[0] { // TODO
	case shared.NON_ADAPTIVE:
		cspapp := gen.CSP(mapp) // Generate CSP
		gen.SaveCSPFile(cspapp) // Save CSP

		//chk.Check(cspapp)       // Check CSP

		eeApp := deployer.NewEE()
		eeApp.Deploy(mapp) // Deploy app
		eeApp.Start()      // Start app

	case shared.EVOLUTIVE_ADAPTATION:
		appKindOfAdaptability := make([]string, 1, 1)
		appKindOfAdaptability = mapp.Adaptability
		meeTemp := crt.Create(mapp, appKindOfAdaptability)           // Create architecture of the execution environment

		crt.Save(meeTemp)                                            // Save architecture of ee

		mee := l.Load(meeTemp.Configuration + shared.MADL_EXTENSION) // Load MADL and generate architectural artefacts (EE)
		mee.AppAdaptability = mapp.Adaptability                      // Configure adaptability of EE's App

		cspee := gen.CSP(mee)
		gen.SaveCSPFile(cspee)    		// Generate & save CSPs

		// Check CSPs
		//chk.Check(cspee)  // TODO think about as it takes a long time and may be correct by construction

		dep.DeployApp(mee, mapp) // Deploy App into EE & start EE
		dep.Start()
	default:
		fmt.Printf("Frontend:: Something wrong with the adaptability of %v\n",mapp.Configuration)
		os.Exit(0)
	}
}
