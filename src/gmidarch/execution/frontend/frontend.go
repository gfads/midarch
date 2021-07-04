package frontend

import (
	"fmt"
	"gmidarch/development/artefacts/csp"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/messages"
	"gmidarch/development/repositories/architectural"
	"gmidarch/execution/deployer"
)

type Frontend interface {
	Deploy(string, map[string]messages.EndPoint)
}

type FrontendImpl struct{}

func NewFrontend() Frontend {
	var fe Frontend
	fe = FrontendImpl{}

	return fe
}

func (f FrontendImpl) Deploy(fileName string, args map[string]messages.EndPoint) {

	// Step 1 - Load architectural repositories
	arm := architectural.NewArchitecturalRepositoryManager()
	archRepo := arm.GetRepository()

	// Step 2: Load madl
	fmt.Print("Loading MADL[", fileName, "]...")
	madlLoader := madl.NewMADLLoader()
	madlApp := madlLoader.Load(fileName)
	fmt.Println("ok")

	// Step 3: Syntax check of madl
	fmt.Print("Syntax checking of MADL...")
	madlChecker := madl.NewMADLChecker()
	madlChecker.SyntaxCheck(madlApp)
	fmt.Println("ok")

	// Step 4: Semantic check of madl
	fmt.Print("Semnatic checking of MADL...")
	madlChecker.SemanticCheck(madlApp, archRepo)
	fmt.Println("ok")

	// Step 5: Configure madl
	fmt.Print("Configuring MADL...")
	madlConfigurator := madl.NewMADLConfigurator()
	madlConfigurator.Configure(&madlApp,archRepo,args)
	fmt.Println("ok")

	// Step 6: Generate & save CSP
	fmt.Print("Generating CSP...")
	cspGenerator := csp.NewCSPGenerator()
	cspSpec := cspGenerator.Generate(madlApp)
	cspGenerator.Save(cspSpec)
	fmt.Println("ok")

	// Step 7: Check CSP
	fmt.Print("Checking CSP...")
	checker := csp.NewFDRGateway()
	checker.Check(cspSpec)
	fmt.Println("ok")

	// Step 8: Start execution
	fmt.Print("Starting execution...")
	deployer := deployer.NewDeployer(madlApp)
	//dep.Deploy(madlApp) // Deploy App into EE & start EE
	go deployer.Start()
	fmt.Println("ok")
}
