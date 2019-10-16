package creator

import (
	"fmt"
	"newsolution/gmidarch/development/artefacts/madl"
	"newsolution/shared/parameters"
	"newsolution/shared/shared"
	"strconv"
	"strings"
)

type Creator struct{}

func (Creator) Create(mapp madl.MADL, appKindOfAdaptability []string) (madl.MADL) {
	mEE := madl.MADL{}
	appIsAdaptive := true

	if len(appKindOfAdaptability) == 1 && appKindOfAdaptability[0] == "NONE" {
		appIsAdaptive = false
	}

	// configuration
	mEE.Configuration = mapp.Configuration + "_EE"

	// adaptability of app
	mEE.AppAdaptability = mapp.Adaptability

	// Components
	comps := []madl.Element{}
	comps = append(comps, madl.Element{ElemId: "core", TypeName: "Core"})

	if appIsAdaptive {
		comps = append(comps, madl.Element{ElemId: "monevolutive", TypeName: "Monevolutive"}) //TODO
		comps = append(comps, madl.Element{ElemId: "monitor", TypeName: "Monitor"})
		comps = append(comps, madl.Element{ElemId: "analyser", TypeName: "Analyser"})
		comps = append(comps, madl.Element{ElemId: "planner", TypeName: "Planner"})
		comps = append(comps, madl.Element{ElemId: "executor", TypeName: "Executor"})
	}

	units := []string{}
	for i := 0; i < len(mapp.Components)+len(mapp.Connectors); i++ {
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		comps = append(comps, madl.Element{ElemId: units[i], TypeName: "Unit"})
	}

	// Connectors
	conns := [] madl.Element{}

	params := make([]interface{},1)
	params[0] = len(units)
	conns = append(conns, madl.Element{ElemId: "t1", TypeName: "OnetoN", Params:params}) // TODO

	if appIsAdaptive {
		conns = append(conns, madl.Element{ElemId: "t2", TypeName: "Oneway"})
		conns = append(conns, madl.Element{ElemId: "t3", TypeName: "Oneway"})
		conns = append(conns, madl.Element{ElemId: "t4", TypeName: "Oneway"})
		conns = append(conns, madl.Element{ElemId: "t5", TypeName: "Oneway"})
		conns = append(conns, madl.Element{ElemId: "t6", TypeName: "Oneway"})
	}

	// Attachments
	atts := []madl.Attachment{}

	for i := 0; i < len(units); i++ {
		attC1 := madl.Element{ElemId: "core", TypeName: "Core"}
		attT := madl.Element{ElemId: "t1", TypeName: "OnetoN"}
		attC2 := madl.Element{ElemId: units[i], TypeName: "ExecutionUnit"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})
	}

	if appIsAdaptive {
		attC1 := madl.Element{ElemId: "monevolutive", TypeName: "Monevolutive"}
		attT := madl.Element{ElemId: "t2", TypeName: "Oneway"}
		attC2 := madl.Element{ElemId: "monitor", TypeName: "Monitor"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "monitor", TypeName: "Monitor"}
		attT = madl.Element{ElemId: "t3", TypeName: "Oneway"}
		attC2 = madl.Element{ElemId: "analyser", TypeName: "Analyser"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "analyser", TypeName: "Analyser"}
		attT = madl.Element{ElemId: "t4", TypeName: "OneWay"}
		attC2 = madl.Element{ElemId: "planner", TypeName: "Planner"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "planner", TypeName: "Planner"}
		attT = madl.Element{ElemId: "t5", TypeName: "Oneway"}
		attC2 = madl.Element{ElemId: "executor", TypeName: "Executor"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "executor", TypeName: "Executor"}
		attT = madl.Element{ElemId: "t6", TypeName: "Oneway"}
		attC2 = madl.Element{ElemId: "core", TypeName: "Core"}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})
	}

	// Adaptability
	eeKindOfAdaptability := []string{}
	eeKindOfAdaptability = append(eeKindOfAdaptability, "NONE")

	// configure MADL EE
	mEE.File = strings.Replace(mapp.File, parameters.MADL_EXTENSION, "", 99) + "_EE" + parameters.MADL_EXTENSION
	mEE.Path = mapp.Path
	mEE.Components = comps
	mEE.Connectors = conns
	mEE.Attachments = atts
	mEE.Adaptability = eeKindOfAdaptability

	return mEE
}

func (Creator) Print(madl madl.MADL) {

	// Configuration
	fmt.Printf("Configuration %v := \n\n", madl.Configuration)

	// Components
	fmt.Printf("   Components \n")
	for i := range madl.Components {
		fmt.Printf("      %v : %v \n", madl.Components[i].ElemId, madl.Components[i].TypeName)
	}
	fmt.Printf("   Connectors \n")
	for i := range madl.Connectors {
		fmt.Printf("      %v : %v \n", madl.Connectors[i].ElemId, madl.Connectors[i].TypeName)
	}
	fmt.Printf("   Attachments \n")
	for i := range madl.Attachments {
		fmt.Printf("      %v,%v,%v\n", madl.Attachments[i].C1.ElemId, madl.Attachments[i].T.ElemId, madl.Attachments[i].C2.ElemId)
	}

	fmt.Printf("\n   Adaptability \n")
	fmt.Printf("      %v \n\n", madl.Adaptability[0]) // TODO
	fmt.Printf("EndConf \n")
}

func (Creator) Save(m madl.MADL) {
	content := []string{}

	path := parameters.DIR_MADL
	name := m.Configuration
	ext := parameters.MADL_EXTENSION

	// Configuration
	content = append(content, "Configuration "+m.Configuration+" := \n\n")

	// Components
	content = append(content, "   Components \n")
	for i := range m.Components {
		content = append(content, "      "+m.Components[i].ElemId+" : "+m.Components[i].TypeName+" \n")
	}
	content = append(content, "\n    Connectors \n")
	for i := range m.Connectors {
		content = append(content, "      "+m.Connectors[i].ElemId+" : "+m.Connectors[i].TypeName+" \n")
	}
	content = append(content, "\n    Attachments \n")
	for i := range m.Attachments {
		content = append(content, "      "+m.Attachments[i].C1.ElemId+","+m.Attachments[i].T.ElemId+","+m.Attachments[i].C2.ElemId+" \n")
	}

	content = append(content, "\n   Adaptability \n")
	content = append(content, "      "+m.Adaptability[0]+" \n\n") // TODO

	content = append(content, "EndConf \n")

	shared.SaveFile(path, name, ext, content)
}
