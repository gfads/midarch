package manager

import (
	"gmidarch/development/framework/messages"
	"gmidarch/development/creator"
	"errors"
	"gmidarch/development/checker"
	"strings"
	"strconv"
	"newsolution/shared/parameters"
	"gmidarch/development/generator"
	"gmidarch/execution"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/artefacts/csp"
	"gmidarch/development/artefacts/graphs"
	"fmt"
	"os"
	"reflect"
	"gmidarch/development/framework/components"
	"gmidarch/development/framework/element"
	"injector/versioning"
	"newsolution/shared/shared"
	"gmidarch/development/framework/architecturallibrary"
)

type Manager struct {
	MadlMidGo             madl.MADLGo
	CSPMid                csp.CSP
	DotsMid               map[string]csp.DOT
	MadlEEGo              madl.MADLGo
	DotsEE                map[string]csp.DOT
	CSPEE                 csp.CSP
	SMMid                 map[string]graphs.GraphExecutable
	SMEE                  map[string]graphs.GraphExecutable
	MapsMid               map[string]string
	MapsEE                map[string]string
	StructuralChannelsMid map[string]chan messages.SAMessage
	StructuralChannelsEE  map[string]chan messages.SAMessage
}

func NewManager(){}

func (m Manager) Invoke(madlFileName string) (error) {
	r1 := *new(error)

	// Initialize Environment
	m.Initialize()

	// MADLs
	fmt.Println("Manager:: MADL being created...")
	creator := creator.Creator{}
	m.MadlMidGo, m.MadlEEGo, r1 = creator.Create(madlFileName)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	// Create Maps
	fmt.Println("Manager:: Maps being created...")
	m.MapsMid = m.CreateMaps(m.MadlMidGo)
	m.MapsEE = m.CreateMaps(m.MadlEEGo)

	// Create Strcutural Channels
	fmt.Println("Manager:: Channels being created...")
	m.StructuralChannelsMid = m.CreateStructuralChannels(m.MadlMidGo)
	m.StructuralChannelsEE = m.CreateStructuralChannels(m.MadlEEGo)

	// CSP
	fmt.Println("Manager:: CSP being created...")
	generator := generator.Generator{}
	m.CSPMid, r1 = generator.GenerateCSP(m.MadlMidGo,m.MapsMid,parameters.MADL_MID,m.MadlMidGo.Adaptability)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	m.CSPEE, r1 = generator.GenerateCSP(m.MadlEEGo,m.MapsEE,parameters.MADL_EE,m.MadlMidGo.Adaptability)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	r1 = generator.GenerateCSPFile(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	r1 = generator.GenerateCSPFile(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	// Checker
	fmt.Println("Manager:: CSP being checked...")
	checker := checker.Checker{}
	isOk, r1 := checker.Check(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	if !isOk {
		r1 = errors.New("Manager:: CSP specifications have not passed on verification")
		return r1
	}
	isOk, r1 = checker.Check(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	if !isOk {
		r1 = errors.New("Manager:: CSP specifications have not passed on verification")
		return r1
	}

	// Invoke FDR - TODO (after integrating with David's solution)
	r1 = checker.GenerateDotFiles(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}
	r1 = checker.GenerateDotFiles(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager:: " + r1.Error())
		return r1
	}

	// DOTS
	fmt.Println("Manager:: DOTS being created...")
	m.DotsMid, r1 = csp.DOT{}.Create(m.CSPMid)
	if r1 != nil {
		r1 = errors.New("Manager" + r1.Error())
		return r1
	}
	m.DotsEE, r1 = csp.DOT{}.Create(m.CSPEE)
	if r1 != nil {
		r1 = errors.New("Manager" + r1.Error())
	}

	// State Machines
	fmt.Println("Manager:: State Machines being created...")
	m.SMMid = make(map[string]graphs.GraphExecutable)
	for i := range m.DotsMid {
		m.SMMid[strings.ToLower(i)], r1 = execution.Create(m.DotsMid[i],m.StructuralChannelsMid)
		if r1 != nil {
			r1 := errors.New("Manager:: " + r1.Error())
			return r1
		}
	}
	m.SMEE = make(map[string]graphs.GraphExecutable)
	for i := range m.DotsEE {
		m.SMEE[strings.ToLower(i)], r1 = execution.Create(m.DotsEE[i],m.StructuralChannelsEE)
		if r1 != nil {
			r1 := errors.New("Manager:: " + r1.Error())
			return r1
		}
	}

	// Update elements' state machines
	for i := range m.MadlMidGo.Components{
		m.MadlMidGo.Components[i].FDRStateMachine = m.DotsMid[strings.ToUpper(m.MadlMidGo.Components[i].ElemId)].Dotgraph
		m.MadlMidGo.Components[i].GoStateMachine =  m.SMMid[m.MadlMidGo.Components[i].ElemId]
	}
	for i := range m.MadlMidGo.Connectors{
		m.MadlMidGo.Connectors[i].FDRStateMachine = m.DotsMid[strings.ToUpper(m.MadlMidGo.Connectors[i].ElemId)].Dotgraph
		m.MadlMidGo.Connectors[i].GoStateMachine =  m.SMMid[m.MadlMidGo.Connectors[i].ElemId]
	}
	for i := range m.MadlEEGo.Components{
		m.MadlEEGo.Components[i].FDRStateMachine = m.DotsEE[strings.ToUpper(m.MadlEEGo.Components[i].ElemId)].Dotgraph
		m.MadlEEGo.Components[i].GoStateMachine =  m.SMEE[m.MadlEEGo.Components[i].ElemId]
	}
	for i := range m.MadlEEGo.Connectors{
		m.MadlEEGo.Connectors[i].FDRStateMachine = m.DotsEE[strings.ToUpper(m.MadlEEGo.Connectors[i].ElemId)].Dotgraph
		m.MadlEEGo.Connectors[i].GoStateMachine =  m.SMEE[m.MadlEEGo.Connectors[i].ElemId]
	}

	// Update Info of elements
	m.ConfigureInfoEE(&m.MadlEEGo,m.MadlMidGo)
	m.ConfigureInfoApp(&m.MadlMidGo)

	// Deploy App into Units
	fmt.Println("Manager:: Deploying elements into Execution Units...")
	m.ConfigureUnits(&m.MadlEEGo, m.MadlMidGo)

	// Start the execution of the architecture
	core := execution.Core{}
	fmt.Println("Manager:: State Machines being deployed...")
	core.Deploy(m.SMMid,m.MadlMidGo)
	core.Deploy(m.SMEE,m.MadlEEGo)

	// Start Versioning Injector
	versioning.VersioningInjector{}.Start(m.MadlMidGo,parameters.PLUGIN_BASE_NAME)

	return r1
}

func (Manager) CreateMaps(madlGo madl.MADLGo) (map[string]string) {
	r1 := make(map[string]string)

	partners := make(map[string]string)
	for i := 0; i < len(madlGo.Attachments); i++ {
		c1Id := madlGo.Attachments[i].C1.ElemId
		tId := madlGo.Attachments[i].T.ElemId
		c2Id := madlGo.Attachments[i].C2.ElemId

		if !strings.Contains(partners[c1Id], tId) {
			partners[c1Id] += ":" + tId
		}
		if !strings.Contains(partners[tId], c1Id) {
			partners[tId] += ":" + c1Id
		}
		if !strings.Contains(partners[tId], c2Id) {
			partners[tId] += ":" + c2Id
		}
		if !strings.Contains(partners[c2Id], tId) {
			partners[c2Id] += ":" + tId
		}
	}

	for i := range partners {
		p := strings.Split(partners[i], ":")
		c := 1
		for j := range p {
			if p[j] != "" {
				r1[i+".e"+strconv.Itoa(c)] = p[j]
				c++
			}
		}
	}
	return r1
}

func (Manager) CreateStructuralChannels(madlGo madl.MADLGo) (map[string]chan messages.SAMessage) {
	r1 := make(map[string]chan messages.SAMessage)

	// Configure structural channels
	for i := 0; i < len(madlGo.Attachments); i++ {
		c1Id := madlGo.Attachments[i].C1.ElemId
		c2Id := madlGo.Attachments[i].C2.ElemId
		tId := madlGo.Attachments[i].T.ElemId

		// c1 -> t
		key01 := c1Id + "." + parameters.INVR + "." + tId
		key02 := tId + "." + parameters.INVP + "." + c1Id
		key03 := tId + "." + parameters.TERP + "." + c1Id
		key04 := c1Id + "." + parameters.TERR + "." + tId
		r1[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key02] = r1[key01]
		r1[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key04] = r1[key03]

		// t -> c2
		key01 = tId + "." + parameters.INVR + "." + c2Id
		key02 = c2Id + "." + parameters.INVP + "." + tId
		key03 = c2Id + "." + parameters.TERP + "." + tId
		key04 = tId + "." + parameters.TERR + "." + c2Id
		r1[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key02] = r1[key01]
		r1[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		r1[key04] = r1[key03]
	}
	return r1;
}

func (Manager) ConfigureUnits(eeConf *madl.MADLGo, appConf madl.MADLGo) {
	availableUnits := []int{}

	// Identify available units
	for i := 0; i < len(eeConf.Components); i++ {
		if reflect.TypeOf(eeConf.Components[i].ElemType).String() == reflect.TypeOf(components.ExecutionUnit{}).String() {
			availableUnits = append(availableUnits, i)
		}
	}

	// Check if the numbers of units is enough to accomodate the application
	if len(availableUnits) != len(appConf.Components)+len(appConf.Connectors) {
		fmt.Println("Engine:: Available units are not enough to execute the components/connectos. Hence, there is a problem in the configuration")
		os.Exit(0)
	}

	// Associate info to units, e.g., u1.info = c1, u2.info = t, u3.info = c2
	idx := 0
	for c := range appConf.Components {
		elem := eeConf.Components[availableUnits[idx]]
		elem.SetInfo(appConf.Components[c])
		eeConf.Components[availableUnits[idx]] = elem
		idx++
	}

	for t := range appConf.Connectors {
		elem := eeConf.Components[availableUnits[idx]]
		elem.SetInfo(appConf.Connectors[t])
		eeConf.Components[availableUnits[idx]] = elem
		idx++
	}
}

func (Manager) ConfigureInfoEE(eeConf *madl.MADLGo, appConf madl.MADLGo) {
	for c1 := range eeConf.Components {
		componentType := reflect.TypeOf(eeConf.Components[c1].ElemType).String()
		switch  componentType {
		case reflect.TypeOf(components.ExecutionEnvironment{}).String(): // Execution environment
			listOfElements := []element.ElementGo{}
			for c2 := range appConf.Components {
				listOfElements = append(listOfElements, appConf.Components[c2])
			}
			for t := range appConf.Connectors {
				listOfElements = append(listOfElements, appConf.Connectors[t])
			}
			elem := eeConf.Components[c1]
			elem.SetInfo(listOfElements)
			eeConf.Components[c1] = elem
		case reflect.TypeOf(components.MAPEKPlanner{}).String(): // Planner
			plannerInfo := components.MAPEKPlannerInfo{ConfId:eeConf.ConfigurationName,Components:appConf.Components}
			elem := eeConf.Components[c1]
			elem.SetInfo(plannerInfo)
			eeConf.Components[c1] = elem
		case reflect.TypeOf(components.MAPEKMonitorEvolutive{}).String(): // Evolutive Monitor
			elem := eeConf.Components[c1]
			//elem.SetInfo(eeConf.ConfigurationName)
			elem.SetInfo(appConf.ConfigurationName)
			eeConf.Components[c1] = elem
		default:
			elem := eeConf.Components[c1]
			elem.SetInfo(parameters.DEFAULT_INFO)
			eeConf.Components[c1] = elem
		}
	}
}

func (Manager) ConfigureInfoApp(conf *madl.MADLGo) {
	for c := range conf.Components {
		elem := conf.Components[c]
		elem.SetInfo(parameters.DEFAULT_INFO)
		conf.Components[c] = elem
	}
}

func (Manager) Initialize() {
	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Check the CSP behaviours
	architecturallibrary.ArchitecturalLibrary{}.CheckLibrary()

	// Show execution parameters
	shared.ShowExecutionParameters(false)
}