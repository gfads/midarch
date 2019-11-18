package creator

import (
	"fmt"
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/components"
	"gmidarch/development/connectors"
	"os"
	"reflect"
	"shared"
	"strconv"
	"strings"
)

type Creator struct{}

func (Creator) Create(mapp madl.MADL, appKindOfAdaptability []string) (madl.MADL) {
	mEE := madl.MADL{}
	appIsAdaptive := true

	if len(appKindOfAdaptability) == 1 && appKindOfAdaptability[0] == "NONE" { // TODO
		appIsAdaptive = false
	}

	// configuration
	mEE.Configuration = mapp.Configuration + "_ee"

	// adaptability of app
	mEE.AppAdaptability = mapp.Adaptability

	// Components
	comps := []madl.Element{}
	comps = append(comps, madl.Element{ElemId: "core", TypeName: reflect.TypeOf(components.Core{}).Name()})

	if appIsAdaptive {
		comps = append(comps, madl.Element{ElemId: "monevolutive", TypeName: reflect.TypeOf(components.Monevolutive{}).Name()})
		comps = append(comps, madl.Element{ElemId: "monitor", TypeName: reflect.TypeOf(components.Monitor{}).Name()})
		comps = append(comps, madl.Element{ElemId: "analyser", TypeName: reflect.TypeOf(components.Analyser{}).Name()})
		comps = append(comps, madl.Element{ElemId: "planner", TypeName: reflect.TypeOf(components.Planner{}).Name()})
		comps = append(comps, madl.Element{ElemId: "executor", TypeName: reflect.TypeOf(components.Executor{}).Name()})
	}

	units := []string{}
	for i := 0; i < len(mapp.Components)+len(mapp.Connectors); i++ {
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		comps = append(comps, madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()})
	}

	// Connectors
	conns := [] madl.Element{}

	params := make([]interface{}, 1)
	params[0] = len(units)

	// TODO
	switch mapp.Configuration {
	case "senderreceiver":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
	case "middlewareclient":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params})
	case "middlewareserver":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params})
	case "clientserverlocal":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
	case "clientserver":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
	case "calculatorlocal":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
	case "midfibonacciserver":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto8{}).Name(), Params: params})
	case "midfibonacciclient":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto6{}).Name(), Params: params})
	case "midnamingserver":
		conns = append(conns, madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params})

	default:
		fmt.Printf("Creator:: Configuration '%v' cannot be executed because 'OnetoN.dot' is not ok!!",mapp.Configuration)
		os.Exit(0)
	}

	if appIsAdaptive {
		conns = append(conns, madl.Element{ElemId: "t2", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
		conns = append(conns, madl.Element{ElemId: "t3", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
		conns = append(conns, madl.Element{ElemId: "t4", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
		conns = append(conns, madl.Element{ElemId: "t5", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
		conns = append(conns, madl.Element{ElemId: "t6", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
	}

	// Attachments
	atts := []madl.Attachment{}

	for i := 0; i < len(units); i++ {
		attC1 := madl.Element{ElemId: "core", TypeName: reflect.TypeOf(components.Core{}).Name()}

		// TODO
		switch mapp.Configuration {
		case "senderreceiver":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "middlewareclient":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "middlewareserver":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "clientserverlocal":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "calculatorlocal":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "clientserver":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "midfibonacciserver":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto8{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "midfibonacciclient":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto6{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})
		case "midnamingserver":
			attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params}
			attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
			atts = append(atts, madl.Attachment{attC1, attT, attC2})

		default:
			fmt.Printf("Creator:: Configuration '%v' cannot be executed because 'OnetoN.dot' is not ok ", mapp.Configuration)
			os.Exit(0)
		}

		//attT := madl.Element{ElemId: "t1", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
		//attC2 := madl.Element{ElemId: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
		//atts = append(atts, madl.Attachment{attC1, attT, attC2})
	}

	if appIsAdaptive {
		attC1 := madl.Element{ElemId: "monevolutive", TypeName: reflect.TypeOf(components.Monevolutive{}).Name()}
		attT := madl.Element{ElemId: "t2", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
		attC2 := madl.Element{ElemId: "monitor", TypeName: reflect.TypeOf(components.Monitor{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "monitor", TypeName: reflect.TypeOf(components.Monitor{}).Name()}
		attT = madl.Element{ElemId: "t3", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
		attC2 = madl.Element{ElemId: "analyser", TypeName: reflect.TypeOf(components.Analyser{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "analyser", TypeName: reflect.TypeOf(components.Analyser{}).Name()}
		attT = madl.Element{ElemId: "t4", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
		attC2 = madl.Element{ElemId: "planner", TypeName: reflect.TypeOf(components.Planner{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "planner", TypeName: reflect.TypeOf(components.Planner{}).Name()}
		attT = madl.Element{ElemId: "t5", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
		attC2 = madl.Element{ElemId: "executor", TypeName: reflect.TypeOf(components.Executor{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = madl.Element{ElemId: "executor", TypeName: reflect.TypeOf(components.Executor{}).Name()}
		attT = madl.Element{ElemId: "t6", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
		attC2 = madl.Element{ElemId: "core", TypeName: reflect.TypeOf(components.Core{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})
	}

	// Adaptability
	eeKindOfAdaptability := []string{}
	eeKindOfAdaptability = append(eeKindOfAdaptability, "NONE")

	// configure MADL EE
	mEE.File = strings.Replace(mapp.File, shared.MADL_EXTENSION, "", 99) + "_ee" + shared.MADL_EXTENSION
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

	path := shared.DIR_MADL
	name := m.Configuration
	ext := shared.MADL_EXTENSION

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
