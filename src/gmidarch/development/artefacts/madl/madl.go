package madl

import (
	"fmt"
	"github.com/gfads/midarch/src/gmidarch/development/components/component"
	"github.com/gfads/midarch/src/gmidarch/development/connectors"
	"github.com/gfads/midarch/src/shared"
	"reflect"
	"strings"
)

type MADL struct {
	Path          string
	FileName      string
	Configuration string
	Components    []component.Component
	Connectors    []connectors.Connector
	Attachments   []Attachment
	Adaptability  []string
	//AppAdaptability []string
	ConnMaps   map[string]connectors.Connector
	ActionMaps map[string]string
}

func (m MADL) GetConnector(id string) connectors.Connector {
	found := false
	r1 := connectors.Connector{}

	for i := range m.Connectors {
		if m.Connectors[i].Id == id {
			r1 = m.Connectors[i]
			found = true
		}
	}

	if !found {
		shared.ErrorHandler(shared.GetFunction(), "Connector '"+id+"' does not belong to architecture!!")
	}

	return r1
}

func (m MADL) CountArity(id string, side int) int {
	c := make(map[string]int)
	r := 0
	for i := range m.Attachments {
		if m.Attachments[i].T.Id == id {
			switch side {
			case shared.LEFT_ARITY:
				c[m.Attachments[i].C1.Id+m.Attachments[i].T.Id] = 1
			case shared.RIGHT_ARITY:
				c[m.Attachments[i].T.Id+m.Attachments[i].C2.Id] = 1
			}
		}
	}
	if len(c) == 0 {
		shared.ErrorHandler(shared.GetFunction(), "Impossible to define the arity of connector '"+id+"'")
	} else {
		r = len(c)
	}
	return r
}

func (m MADL) IsInConnectors(e connectors.Connector) bool {
	foundConnector := false

	for i := range m.Connectors {
		if e.Id == m.Connectors[i].Id {
			foundConnector = true
			break
		}
	}
	return foundConnector
}

func (m MADL) IsInComponents(e component.Component) bool {
	foundComponent := false

	for i := range m.Components {
		if e.Id == m.Components[i].Id {
			foundComponent = true
			break
		}
	}
	return foundComponent
}

func (m MADL) IsComponentInAttachments(e component.Component) bool {
	foundComponent := false

	for a := range m.Attachments {
		if m.Attachments[a].C1.Id == e.Id || m.Attachments[a].C2.Id == e.Id {
			foundComponent = true
		}
	}

	return foundComponent
}

func (m MADL) IsConnectorInAttachments(e connectors.Connector) bool {
	foundComponent := false

	for a := range m.Attachments {
		if m.Attachments[a].T.Id == e.Id {
			foundComponent = true
		}
	}
	return foundComponent
}

func CheckMADLCommandInLine(s string, c string) bool {
	r := false

	sUpper := strings.ToUpper(s)

	if strings.Contains(sUpper, c) {
		if !(strings.Index(sUpper, shared.MADL_COMMENT) < strings.Index(sUpper, c)) { // Check comment
			r = true
		}
	}
	return r
}

func (m MADL) PrintComponents() {

	for i := range m.Components {
		fmt.Println(reflect.TypeOf(m.Components[i].Type))
	}
}
