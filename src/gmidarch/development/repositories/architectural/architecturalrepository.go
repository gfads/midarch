package architectural

import (
	"gmidarch/development/components/adaptive"
	"gmidarch/development/components/apps"
	"gmidarch/development/components/component"
	"gmidarch/development/components/middleware"
	"gmidarch/development/components/proxies/calculatorproxy"
	"gmidarch/development/components/proxies/fibonacciProxy"
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/connectors"
	"gmidarch/development/messages"
	"io/ioutil"
	"shared"
	"strconv"
)

type ArchitecturalRepository struct {
	CompLibrary map[string]*component.Component
	ConnLibrary map[string]connectors.Connector
}

// Set of existing Components
var SetOfComponentTypesRAM = map[string]interface{} {
	"Client":            	&apps.Client{},
	"Server":            	&apps.Server{},
	"Sender":            	&apps.Sender{},
	"Receiver":				&apps.Receiver{},
	"Core":					&adaptive.Core{},
	"Unit":					&adaptive.Unit{},
	"Monevolutive":			&adaptive.Monevolutive{},
	"EvolutiveProtocol":	&adaptive.EvolutiveProtocol{},
	"Monitor":				&adaptive.Monitor{},
	"Planner":				&adaptive.Planner{},
	"Executor":				&adaptive.Executor{},
	"Analyser":				&adaptive.Analyser{},
	"Jsonmarshaller":		&middleware.Jsonmarshaller{},
	"Gobmarshaller":     	&middleware.Gobmarshaller{},
	"CRHTCP":            	&middleware.CRHTCP{},
	"SRHTCP":            	&middleware.SRHTCP{},
	"CRHUDP":            	&middleware.CRHUDP{},
	"SRHUDP":            	&middleware.SRHUDP{},
	"CRHTLS":            	&middleware.CRHTLS{},
	"SRHTLS":            	&middleware.SRHTLS{},
	"Calculatorinvoker": 	&middleware.Calculatorinvoker{},
	"FibonacciInvoker": 	&middleware.FibonacciInvoker{},
	"Requestor":         	&middleware.Requestor{},
	"Naminginvoker":     	&middleware.Naminginvoker{},
	"Namingserver":      	&middleware.Namingserver{},
	"Namingproxy":       	&namingproxy.Namingproxy{},
	"Calculatorproxy":   	&calculatorproxy.Calculatorproxy{},
	"FibonacciProxy":   	&fibonacciProxy.FibonacciProxy{}}

// Set of existing Connectors
var SetOfConnectorTypesRAM = map[string]connectors.Connector{
	shared.ONEWAY: {TypeName: shared.ONEWAY, Behaviour: "B = InvP.e1 -> InvR.e2 -> B", DefaultRightArity: 1,
		DefaultLeftArity: 1},
	shared.REQUEST_REPLY: {TypeName: shared.REQUEST_REPLY, Behaviour: "B = InvP.e1 -> InvR.e2 -> TerR.e2 -> TerP.e1 -> B", DefaultRightArity: 1,
		DefaultLeftArity: 1},
	shared.ONETON:       {Behaviour: shared.RUNTIME_BEHAVIOUR, DefaultLeftArity: 1, DefaultRightArity: shared.MAX_RIGHT_ARITY},
	shared.NTOONE:       {TypeName: shared.NTOONE, Behaviour: shared.RUNTIME_BEHAVIOUR, DefaultLeftArity: shared.MAX_LEFT_ARITY, DefaultRightArity: 1},
	shared.NTOONEREQREP: {TypeName: shared.NTOONEREQREP, Behaviour: shared.RUNTIME_BEHAVIOUR, DefaultLeftArity: shared.MAX_LEFT_ARITY, DefaultRightArity: 1},
	shared.ONETONREQREP: {TypeName: shared.ONETONREQREP, Behaviour: shared.RUNTIME_BEHAVIOUR, DefaultLeftArity: 1, DefaultRightArity: shared.MAX_RIGHT_ARITY}}

func LoadArchitecturalRepository() ArchitecturalRepository {
	r := ArchitecturalRepository{}
	shared.ArchitecturalComponentTypes = SetOfComponentTypesRAM

	// Initialize repositories
	r.CompLibrary = make(map[string]*component.Component)
	r.ConnLibrary = make(map[string]connectors.Connector)

	// Read type and behaviour from actual files
	SetOfComponentTypesFile := ReadComponentTypesFromDisk()

	// Check the consistency of RAM/File repositories
	if len(SetOfComponentTypesRAM) != len(SetOfComponentTypesFile) {
		shared.ErrorHandler(shared.GetFunction(),
			"The set of components in RAM("+strconv.Itoa(len(SetOfComponentTypesRAM))+") and Disk("+strconv.Itoa(len(SetOfComponentTypesFile))+") are inconsistent!!")
	}

	// Store components on the architectural repositories
	for i := range SetOfComponentTypesFile {
		newComp := &component.Component{}
		typeComp, ok := SetOfComponentTypesRAM[i]
		if !ok {
			shared.ErrorHandler(shared.GetFunction(), "Component '"+i+"' is in File, but not in RAM!!")
		}
		newComp.Type = typeComp
		newComp.TypeName = i                           // From file
		newComp.Behaviour = SetOfComponentTypesFile[i] // From file
		newComp.Buffer = messages.SAMessage{}          // Initialisation
		r.CompLibrary[i] = newComp
	}

	// Store connectors on the architectural repositories
	for i := range SetOfConnectorTypesRAM {
		r.ConnLibrary[i] = SetOfConnectorTypesRAM[i]
	}
	return r
}

func ReadComponentTypesFromDisk() map[string]string {
	compLibrary := map[string]string{}

	// Identify adaptive components
	adaptiveFiles, err1 := ioutil.ReadDir(shared.DIR_ADAPTIVE_COMPONENTS)
	if err1 != nil {
		shared.ErrorHandler(shared.GetFunction(), err1.Error())
	}

	for file := range adaptiveFiles {
		fullPathName := shared.DIR_ADAPTIVE_COMPONENTS + "/" + adaptiveFiles[file].Name()
		typeName, behaviour := shared.GetTypeAndBehaviour(fullPathName)
		compLibrary[typeName] = behaviour
	}

	// Identify application components
	appFiles, err1 := ioutil.ReadDir(shared.DIR_APP_COMPONENTS)
	if err1 != nil {
		shared.ErrorHandler(shared.GetFunction(), err1.Error())
	}

	for file := range appFiles {
		fullPathName := shared.DIR_APP_COMPONENTS + "/" + appFiles[file].Name()
		typeName, behaviour := shared.GetTypeAndBehaviour(fullPathName)
		compLibrary[typeName] = behaviour
	}

	// Identify middleware components
	midFiles, err1 := ioutil.ReadDir(shared.DIR_MIDDLEWARE_COMPONENTS)
	if err1 != nil {
		shared.ErrorHandler(shared.GetFunction(), err1.Error())
	}

	for file := range midFiles {
		fullPathName := shared.DIR_MIDDLEWARE_COMPONENTS + "/" + midFiles[file].Name()
		typeName, behaviour := shared.GetTypeAndBehaviour(fullPathName)
		compLibrary[typeName] = behaviour
	}

	// Identify proxy components (proxies dir has folders)
	proxiesFolders, err1 := ioutil.ReadDir(shared.DIR_PROXIES_COMPONENTS)
	if err1 != nil {
		shared.ErrorHandler(shared.GetFunction(), err1.Error())
	}

	for folder := range proxiesFolders {
		temp, err1 := ioutil.ReadDir(shared.DIR_PROXIES_COMPONENTS + "/" + proxiesFolders[folder].Name())
		if err1 != nil {
			shared.ErrorHandler(shared.GetFunction(), err1.Error())
		}
		for file := range temp {
			fullPathName := shared.DIR_PROXIES_COMPONENTS + "/" + proxiesFolders[folder].Name() + "/" + temp[file].Name()
			typeName, behaviour := shared.GetTypeAndBehaviour(fullPathName)
			compLibrary[typeName] = behaviour
		}
	}

	return compLibrary
}