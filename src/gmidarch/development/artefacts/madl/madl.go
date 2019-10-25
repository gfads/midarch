package madl

import (
	"errors"
	"fmt"
	"gmidarch/development/artefacts/dot"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/repositories/architectural"
	"os"
	"reflect"
	"shared"
	"strconv"
	"strings"
)

type MADL struct {
	Path            string
	File            string
	Configuration   string
	Components      []Element
	Connectors      []Element
	Attachments     []Attachment
	Adaptability    []string
	AppAdaptability []string
	Channels        map[string]chan messages.SAMessage
	Maps            map[string]string
}

// Configure executable graph of components
func (m *MADL) ConfigureComponents() {
	lib := new(architectural.ArchitecturalRepository)
	lib.Load()

	for i := range m.Components {
		record, ok := lib.Library[m.Components[i].TypeName]
		if !ok {
			fmt.Printf("MADL:: Component type '%v' not in Library", m.Components[i].TypeName)
			os.Exit(0)
		}
		m.Components[i].Type = record.Type
		m.Components[i].Behaviour = record.Behaviour
		dotgraph := dot.DOT{}.Read(m.Components[i].TypeName + shared.DOT_EXTENSION)
		execgraph := graphs.Exec{}.Create(m.Components[i].ElemId, m.Components[i].Type, m.Components[i].TypeName, dotgraph, m.Maps, m.Channels)

		m.Components[i].Graph = execgraph
	}
}

// Configure executable graph of connectors
func (m *MADL) ConfigureConnectors() {
	lib := new(architectural.ArchitecturalRepository)
	lib.Load()

	for i := range m.Connectors {
		record, ok := lib.Library[m.Connectors[i].TypeName]
		if !ok {
			fmt.Println("MADL:: Connector type '" + m.Connectors[i].TypeName + "'not in Library")
			os.Exit(0)
		}
		m.Connectors[i].Type = record.Type
		m.Connectors[i].Behaviour = record.Behaviour
		dotgraph := dot.DOT{}.Read(m.Connectors[i].TypeName + shared.DOT_EXTENSION)
		execgraph := graphs.Exec{}.Create(m.Connectors[i].ElemId, m.Connectors[i].Type, m.Connectors[i].TypeName, dotgraph, m.Maps, m.Channels)

		m.Connectors[i].Graph = execgraph
	}
}

// Configure channels and Maps - basic to the execution
func (madl *MADL) ConfigureChannelsAndMaps() {
	structuralChannels := make(map[string]chan messages.SAMessage)

	// Configure structural channels
	for i := range madl.Attachments {
		c1Id := madl.Attachments[i].C1.ElemId
		c2Id := madl.Attachments[i].C2.ElemId
		tId := madl.Attachments[i].T.ElemId

		// c1 -> t
		key01 := c1Id + "." + shared.INVR + "." + tId
		key02 := tId + "." + shared.INVP + "." + c1Id
		key03 := tId + "." + shared.TERP + "." + c1Id
		key04 := c1Id + "." + shared.TERR + "." + tId
		structuralChannels[key01] = make(chan messages.SAMessage, shared.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan messages.SAMessage, shared.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]

		// t -> c2
		key01 = tId + "." + shared.INVR + "." + c2Id
		key02 = c2Id + "." + shared.INVP + "." + tId
		key03 = c2Id + "." + shared.TERP + "." + tId
		key04 = tId + "." + shared.TERR + "." + c2Id
		structuralChannels[key01] = make(chan messages.SAMessage, shared.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan messages.SAMessage, shared.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]
	}
	madl.Channels = structuralChannels

	// Configure maps
	elemMaps := make(map[string]string)
	partners := make(map[string]string)

	for i := range madl.Attachments {
		c1Id := madl.Attachments[i].C1.ElemId
		c2Id := madl.Attachments[i].C2.ElemId
		tId := madl.Attachments[i].T.ElemId
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
				elemMaps[i+".e"+strconv.Itoa(c)] = p[j]
				c++
			}
		}
	}
	madl.Maps = elemMaps
}

// Indentify Configuration Name
func (MADL) IdentifyConfigurationName(content []string) (string, error) {
	r1 := ""
	r2 := *new(error)

	for l := range content {
		tempContent := content[l]
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

// Indentify Components
func (MADL) IdentifyComponents(content []string) ([]Element, error) {
	foundComponents := false
	r1 := []Element{}
	r2 := *new(error)

	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "COMPONENTS") {
			foundComponents = true
		} else {
			if foundComponents && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				compId := strings.TrimSpace(temp[0])
				compType := ""
				compType = strings.TrimSpace(temp[1])
				r1 = append(r1, Element{ElemId: compId, TypeName: compType})
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

func (MADL) IdentifyConnectors(content []string) ([]Element, error) {
	foundConnectors := false
	r1 := []Element{}
	r2 := *new(error)

	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "CONNECTORS") {
			foundConnectors = true
		} else {
			if foundConnectors && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ":") {
				temp := strings.Split(tempLine, ":")
				connId := strings.TrimSpace(temp[0])
				connType := strings.TrimSpace(temp[1])
				connTypeName := connType
				r1 = append(r1, Element{ElemId: connId, Type: connType, TypeName: connTypeName})
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

func (MADL) IdentifyAttachments(content []string) ([]Attachment, error) {
	r1 := []Attachment{}
	r2 := *new(error)

	// Identify Attachments
	foundAttachments := false
	for l := range content {
		tempLine := content[l]
		if strings.Contains(strings.ToUpper(tempLine), "ATTACHMENTS") {
			foundAttachments = true
		} else {
			if foundAttachments && !shared.SkipLine(tempLine) && strings.Contains(tempLine, ",") {
				atts := strings.Split(strings.TrimSpace(tempLine), ",")
				c1Temp := strings.TrimSpace(atts[0])
				tTemp := strings.TrimSpace(atts[1])
				c2Temp := strings.TrimSpace(atts[2])

				c1 := Element{ElemId: c1Temp}
				t := Element{ElemId: tTemp}
				c2 := Element{ElemId: c2Temp}

				att := Attachment{c1, t, c2}
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

func (MADL) IdentifyAdaptability(content []string) ([]string, error) {
	r1 := []string{}
	r2 := *new(error)

	foundAdaptability := false
	for l := range content {
		tempLine := content[l]
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

func (m MADL) PrintComponents() {

	for i := range m.Components {
		fmt.Println(reflect.TypeOf(m.Components[i].Type))
	}
}

func (m MADL) Check() (error) {
	r1 := *new(error)

	// Check if all components/connectors were declared
	for a := range m.Attachments {

		if !m.isInComponents(m.Attachments[a].C1) {
			r1 = errors.New("Component '" + m.Attachments[a].C1.ElemId + "' was not Declared!!")
			return r1
		}

		if !m.isInConnectors(m.Attachments[a].T) {
			r1 = errors.New("Connector '" + m.Attachments[a].T.ElemId + "' was not Declared!!")
			return r1

		}
		if !m.isInComponents(m.Attachments[a].C2) {
			r1 = errors.New("Component '" + m.Attachments[a].C2.ElemId + "' was not Declared!!")
			return r1
		}
	}

	// Check if all components/connectors were used
	for c := range m.Components {
		if !m.isComponentInAttachments(m.Components[c]) {
			r1 = errors.New("Component '" + m.Components[c].ElemId + "' declared, but not Used!!")
			return r1
		}
	}

	for t := range m.Connectors {
		if !m.isConnectorInAttachments(m.Connectors[t]) {
			r1 = errors.New("Connector '" + m.Connectors[t].ElemId + "' declared, but not Used!!")
			return r1
		}
	}
	return r1
}

func (m MADL) isInConnectors(e Element) bool {
	foundConnector := false

	for i := range m.Connectors {
		if e.ElemId == m.Connectors[i].ElemId {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

func (m MADL) isInComponents(e Element) bool {
	foundComponent := false

	for i := range m.Components {
		if e.ElemId == m.Components[i].ElemId {
			foundComponent = true
			break
		}
	}
	return foundComponent
}

func (m MADL) isComponentInAttachments(e Element) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].C1.ElemId == e.ElemId || m.Attachments[a].C2.ElemId == e.ElemId) {
			foundComponent = true
		}
	}

	return foundComponent
}

func (m MADL) isConnectorInAttachments(e Element) bool {
	foundComponent := false

	for a := range m.Attachments {
		if (m.Attachments[a].T.ElemId == e.ElemId) {
			foundComponent = true
		}
	}
	return foundComponent
}
