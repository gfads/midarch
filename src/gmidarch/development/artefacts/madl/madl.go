package madl

import (
	"newsolution/shared/shared"
	"strings"
	"errors"
	"strconv"
	"newsolution/shared/parameters"
	"gmidarch/development/framework/element"
	"gmidarch/development/framework/configuration/attachments"
	"fmt"
	"os"
	"bufio"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/framework/architecturallibrary"
)

type Element struct {
	ElemId         string
	TypeName       string
	Type           interface{}
	CSP            string
	Info           interface{}
	GoStateMachine graphs.GraphExecutable
}

type Attachment struct {
	C1 Element
	T  Element
	C2 Element
}

type MADL struct {
	Path          string
	File          string
	Content       []string
	Configuration string
	Components    []Element
	Connectors    []Element
	Attachments   []Attachment
	Adaptability  []string
}

func (m *MADLGo) Create(madl MADL) error {
	r1 := *new(error)
	lib := architecturallibrary.ArchitecturalLibrary{}

	// Load architectural library
	err := lib.Load()
	if err != nil {
		r1 = errors.New("MADLGO:: " + err.Error())
		return r1
	}

	// Configuration
	m.ConfigurationName = madl.ConfigurationName

	// Components
	comps := []element.ElementGo{}
	for c := range madl.Components {
		compMadl := madl.Components[c]
		err, compRecord := lib.GetRecord(compMadl.ElemType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		compGoTemp := element.ElementGo{ElemId: compMadl.ElemId, ElemType: compRecord.Go, CSP: compRecord.CSP}
		comps = append(comps, compGoTemp)
	}
	m.Components = comps

	// Connectors
	conns := []element.ElementGo{}
	for c := range madl.Connectors {
		connMadl := madl.Connectors[c]
		err, connRecord := lib.GetRecord(connMadl.ElemType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		connGoTemp := element.ElementGo{ElemId: connMadl.ElemId, ElemType: connRecord.Go, CSP: connRecord.CSP}
		conns = append(conns, connGoTemp)
	}
	m.Connectors = conns

	// Attachments
	atts := []attachments.AttachmentGo{}
	for a := 0; a < len(madl.Attachments); a++ {
		attMadl := madl.Attachments[a]

		c1Type := attMadl.C1.ElemType
		err, c1Record := lib.GetRecord(c1Type)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		c1Go := element.ElementGo{ElemId: attMadl.C1.ElemId, ElemType: c1Record.Go, CSP: c1Record.CSP}

		tType := attMadl.T.ElemType
		err, tRecord := lib.GetRecord(tType)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		tGo := element.ElementGo{ElemId: attMadl.T.ElemId, ElemType: tRecord.Go, CSP: tRecord.CSP}

		c2Type := attMadl.C2.ElemType
		err, c2Record := lib.GetRecord(c2Type)
		if err != nil {
			r1 = errors.New("MADLGO:: " + err.Error())
			return r1
		}
		c2Go := element.ElementGo{ElemId: attMadl.C2.ElemId, ElemType: c2Record.Go, CSP: c2Record.CSP}

		atts = append(atts, attachments.AttachmentGo{c1Go, tGo, c2Go})
	}
	m.Attachments = atts

	// Adaptability
	m.Adaptability = madl.Adaptability

	return r1
}

func (m *MADL) Read(file string) {

	// Check file name
	err := m.CheckFileName(file)
	shared.CheckError(err, "MADL")

	// configure r
	m.File = file
	m.Path = parameters.DIR_CONF

	fullPathAdlFileName := m.Path + "/" + m.File

	// read file
	fileContent := []string{}
	fileTemp, err := os.Open(fullPathAdlFileName)
	shared.CheckError(err, "MADL")
	defer fileTemp.Close()

	scanner := bufio.NewScanner(fileTemp)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// configure r
	m.Content = fileContent
}

func (MADL) CheckFileName(fileName string) error {
	r := *new(error)

	len := len(fileName)

	if len <= 5 {
		r = errors.New("File Name Invalid")
	} else {
		if fileName[len-5:] != parameters.MADL_EXTENSION {
			r = errors.New("Invalid extension of '" + fileName + "'")
		} else {
			r = nil
		}
	}

	return r
}

func (m *MADL) Create(madlFile MADLFile) error {
	r1 := *new(error)

	m.File = madlFile

	// Configuration
	configurationName, err := m.IdentifyConfigurationName()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Configuration = configurationName

	// Components
	components, err := m.IdentifyComponents()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Components = components

	// Connectors
	connectors, err := m.IdentifyConnectors()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Connectors = connectors

	// Attachments
	attachments, err := m.IdentifyAttachments()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Attachments = attachments
	m.SetAttachmentTypes()

	adaptability, err := m.IdentifyAdaptability()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}
	m.Adaptability = adaptability

	err = m.Check()
	if (err != nil) {
		r1 = errors.New("MADL: " + err.Error())
		return r1
	}

	return r1
}

func (m MADL) IdentifyConfigurationName() (string, error) {
	r1 := ""
	r2 := *new(error)

	for l := range m.Content {
		tempContent := m.Content[l]
		if strings.Contains(strings.ToUpper(tempContent), "CONFIGURATION") {
			temp := strings.Split(tempContent, " ")
			r1 = strings.TrimSpace(temp[1])
		}
	}
	if r1 == "" {
		r2 = errors.New("Configuration name not defined.")
	}
	return r1, r2
}

func (m MADL) IdentifyComponents() ([]element.ElementMADL, error) {
	foundComponents := false
	r1 := []element.ElementMADL{}
	r2 := *new(error)

	for l := range m.Content {
		tempLine := m.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "COMPONENTS") {
			foundComponents = true
		} else {
			if foundComponents && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				compId := strings.TrimSpace(temp[0])
				compType := ""
				compType = strings.TrimSpace(temp[1])
				r1 = append(r1, element.ElementMADL{ElemId: compId, ElemType: compType})
			} else {
				if foundComponents && !shared.SkipLine(tempLine) && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if len(r1) == 0 {
		r2 = errors.New("MADL:: 'Components' not well formed.")
	}

	return r1, r2
}

func (m MADL) IdentifyConnectors() ([]element.ElementMADL, error) {
	foundConnectors := false
	r1 := []element.ElementMADL{}
	r2 := *new(error)

	for l := range m.Content {
		tempLine := m.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "CONNECTORS") {
			foundConnectors = true
		} else {
			if foundConnectors && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				connId := strings.TrimSpace(temp[0])
				connType := strings.TrimSpace(temp[1])
				r1 = append(r1, element.ElementMADL{ElemId: connId, ElemType: connType})
			} else {
				if foundConnectors && tempLine != "" && !strings.Contains(tempLine, ":") {
					break
				}
			}
		}
	}

	if len(r1) == 0 {
		r2 = errors.New("MADL:: 'Connectors' not well formed.")
	}

	return r1, r2
}

func (m MADL) IdentifyAttachments() ([]attachments.AttachmentMADL, error) {
	r1 := []attachments.AttachmentMADL{}
	r2 := *new(error)

	// Identify Attachments
	foundAttachments := false
	for l := range m.Content {
		tempLine := m.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ATTACHMENTS") {
			foundAttachments = true
		} else {
			if foundAttachments && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ",") {
				atts := strings.Split(strings.TrimSpace(tempLine), ",")
				c1Temp := atts[0]
				tTemp := atts[1]
				c2Temp := atts[2]

				c1 := element.ElementMADL{ElemId: c1Temp}
				t := element.ElementMADL{ElemId: tTemp}
				c2 := element.ElementMADL{ElemId: c2Temp}

				att := attachments.AttachmentMADL{c1, t, c2}
				r1 = append(r1, att)
			} else {
				if foundAttachments && tempLine != "" && !strings.Contains(tempLine, ",") {
					break
				}
			}
		}
	}

	if len(r1) == 0 {
		r2 = errors.New("MADL:: 'Attachments' not well formed.")
	}

	return r1, r2
}

func (m MADL) IdentifyAdaptability() ([]string, error) {
	r1 := []string{}
	r2 := *new(error)

	foundAdaptability := false
	for l := range m.Content {
		tempLine := m.Content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ADAPTABILITY") {
			foundAdaptability = true
		} else {
			if foundAdaptability && !shared.SkipLine(tempLine) && shared.IsAdaptationType(tempLine) {
				r1 = append(r1, strings.ToUpper(strings.TrimSpace(tempLine)))
			} else {
				if foundAdaptability && !shared.SkipLine(tempLine) && !shared.IsAdaptationType(tempLine) {
					break
				}
			}
		}
	}

	if !foundAdaptability || len(r1) == 0 {
		r2 = errors.New("'Adaptability' NOT well defined!")
	}

	return r1, r2
}

func (m *MADL) SetAttachmentTypes() {
	for a := range m.Attachments {
		c1 := m.Attachments[a].C1.ElemId
		c1Type, err := m.GetType(c1)
		shared.CheckError(err, "MADL")
		t := m.Attachments[a].T.ElemId
		tType, err := m.GetType(t)
		shared.CheckError(err, "MADL")
		c2 := m.Attachments[a].C2.ElemId
		c2Type, err := m.GetType(c2)
		shared.CheckError(err, "MADL")

		tempAttachment := attachments.AttachmentMADL{element.ElementMADL{ElemId: c1, ElemType: c1Type}, element.ElementMADL{ElemId: t, ElemType: tType}, element.ElementMADL{c2, c2Type}}
		m.Attachments[a] = tempAttachment
	}
}

func (m MADL) GetType(elemId string) (string, error) {
	r1 := ""
	r2 := *new(error)
	found := false

	for c := range m.Components {
		if m.Components[c].ElemId == elemId {
			found = true
			r1 = m.Components[c].ElemType
			break
		}
	}

	if !found {
		for c := range m.Connectors {
			if m.Connectors[c].ElemId == elemId {
				found = true
				r1 = m.Connectors[c].ElemType
				break
			}
		}
	}

	if !found {
		r2 = errors.New("MADL:: Type of element '" + elemId + "' not found.")
	}
	return r1, r2
}

func (m MADL) CreateEE(kindOfAdaptability []string) (MADL, error) {
	r1 := MADL{}
	r2 := *new(error)
	isAdaptive := true

	if len(kindOfAdaptability) == 1 && kindOfAdaptability[0] == "NONE" {
		isAdaptive = false
	}

	// configuration
	r1.Configuration = m.Configuration + "_EE"

	// Components
	comps := []element.ElementMADL{}

	comps = append(comps, element.ElementMADL{"ee", "ExecutionEnvironment"})

	if isAdaptive {
		comps = append(comps, element.ElementMADL{"monitorevolutive", "MAPEKMonitorEvolutive"}) //TODO
		comps = append(comps, element.ElementMADL{"mapekmonitor", "MAPEKMonitor"})
		comps = append(comps, element.ElementMADL{"analyser", "MAPEKAnalyser"})
		comps = append(comps, element.ElementMADL{"planner", "MAPEKPlanner"})
		comps = append(comps, element.ElementMADL{"executor", "MAPEKExecutor"})
	}

	units := []string{}
	for i := 0; i < len(m.Components)+len(m.Connectors); i++ {
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		comps = append(comps, element.ElementMADL{units[i], "ExecutionUnit"})
	}

	// Connectors
	conns := [] element.ElementMADL{}

	conns = append(conns, element.ElementMADL{"t1", "OneToN"})

	if isAdaptive {
		conns = append(conns, element.ElementMADL{"t2", "OneWay"})
		conns = append(conns, element.ElementMADL{"t3", "OneWay"})
		conns = append(conns, element.ElementMADL{"t4", "OneWay"})
		conns = append(conns, element.ElementMADL{"t5", "OneWay"})
		conns = append(conns, element.ElementMADL{"t6", "OneWay"})
	}

	// Attachments
	atts := []attachments.AttachmentMADL{}

	for i := 0; i < len(units); i++ {
		attC1 := element.ElementMADL{"ee", "ExecutionEnvironment"}
		attT := element.ElementMADL{"t1", "OneToN"}
		attC2 := element.ElementMADL{units[i], "ExecutionUnit"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})
	}

	if isAdaptive {
		attC1 := element.ElementMADL{"monitorevolutive", "MAPEKMonitorEvolutive"}
		attT := element.ElementMADL{"t2", "OneWay"}
		attC2 := element.ElementMADL{"mapekmonitor", "MAPEKMonitor"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = element.ElementMADL{"mapekmonitor", "MAPEKMonitor"}
		attT = element.ElementMADL{"t3", "OneWay"}
		attC2 = element.ElementMADL{"analyser", "MAPEKAnalyser"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = element.ElementMADL{"analyser", "MAPEKAnalyser"}
		attT = element.ElementMADL{"t4", "OneWay"}
		attC2 = element.ElementMADL{"planner", "MAPEKPlanner"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = element.ElementMADL{"planner", "MAPEKPlanner"}
		attT = element.ElementMADL{"t5", "OneWay"}
		attC2 = element.ElementMADL{"executor", "MAPEKExecutor"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})

		attC1 = element.ElementMADL{"executor", "MAPEKExecutor"}
		attT = element.ElementMADL{"t6", "OneWay"}
		attC2 = element.ElementMADL{"ee", "ExecutionEnvironment"}
		atts = append(atts, attachments.AttachmentMADL{attC1, attT, attC2})
	}

	// Adaptability
	adaptability := []string{}
	adaptability = append(adaptability, "NONE") // TODO

	// configure MADL EE
	r1.File = strings.Replace(m.File, parameters.MADL_EXTENSION, "", 99) + "_EE" + parameters.MADL_EXTENSION
	r1.Path = m.Path
	r1.Components = comps
	r1.Connectors = conns
	r1.Attachments = atts
	r1.Adaptability = adaptability

	return r1, r2
}

func (m MADL) Print() {
	lines := []string{}

	// Configuration
	lines = append(lines, "Configuration "+m.Configuration+" := "+"\n")

	// Components
	lines = append(lines, "Components"+"\n")
	for i := range m.Components {
		lines = append(lines, "   "+m.Components[i].ElemId+" : "+m.Components[i].ElemType+"\n")
	}

	// Connectors
	lines = append(lines, "Connectors"+"\n")
	for i := range m.Connectors {
		lines = append(lines, "   "+m.Connectors[i].ElemId+" : "+m.Connectors[i].ElemType+"\n")
	}

	// Attachments
	lines = append(lines, "Attachments"+"\n")
	for i := range m.Attachments {
		lines = append(lines, "   "+m.Attachments[i].C1.ElemId+","+m.Attachments[i].T.ElemId+","+m.Attachments[i].C2.ElemId+"\n")
	}

	// Adaptability
	lines = append(lines, "Adaptability"+"\n")
	lines = append(lines, "   "+m.Adaptability[0]+"\n") // TODO

	lines = append(lines, "EndConf")

	for i := range lines {
		fmt.Print(lines[i])
	}
}

func (m MADL) Check() (error) {
	r1 := *new(error)

	// Check if all components/connectors were declared
	for a := range m.Attachments {

		if !m.IsInComponents(m.Attachments[a].C1) {
			r1 = errors.New("Component '" + m.Attachments[a].C1.ElemId + "' was not Declared!!")
			return r1
		}

		if !m.IsInConnectors(m.Attachments[a].T) {
			r1 = errors.New("Connector '" + m.Attachments[a].T.ElemId + "' was not Declared!!")
			return r1

		}
		if !m.IsInComponents(m.Attachments[a].C2) {
			r1 = errors.New("Component '" + m.Attachments[a].C2.ElemId + "' was not Declared!!")
			return r1
		}
	}

	// Check if all components/connectors were used
	for c := range m.Components {
		if !m.IsComponentInAttachments(m.Components[c]) {
			r1 = errors.New("Component '" + m.Components[c].ElemId + "' declared, but not Used!!")
			return r1
		}
	}

	for t := range m.Connectors {
		if !m.IsConnectorInAttachments(m.Connectors[t]) {
			r1 = errors.New("Connector '" + m.Connectors[t].ElemId + "' declared, but not Used!!")
			return r1
		}
	}
	return r1
}

func (m MADL) IsInConnectors(e element.ElementMADL) bool {
	foundConnector := false

	for i := range m.Connectors {
		if e.ElemId == m.Connectors[i].ElemId {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

func (m MADL) IsInComponents(e element.ElementMADL) bool {
	foundComponent := false

	for i := range m.Components {
		if e.ElemId == m.Components[i].ElemId {
			foundComponent = true
			break
		}
	}
	return foundComponent
}

func (m MADL) IsComponentInAttachments(e element.ElementMADL) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].C1.ElemId == e.ElemId || m.Attachments[a].C2.ElemId == e.ElemId) {
			foundComponent = true
		}
	}

	return foundComponent
}

func (m MADL) IsConnectorInAttachments(e element.ElementMADL) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].T.ElemId == e.ElemId) {
			foundComponent = true
		}
	}
	return foundComponent
}
