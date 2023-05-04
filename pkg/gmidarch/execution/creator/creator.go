package creator

import (
	"fmt"
	"github.com/gfads/midarch/pkg/gmidarch/development/artefacts/madl"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/adaptive"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/component"
	"github.com/gfads/midarch/pkg/gmidarch/development/connectors"
	"github.com/gfads/midarch/pkg/shared"
	"reflect"
	"strconv"
	"strings"
)

type Creator struct{}

//func (Creator) CreateNew(mapp madl.MADL, appKindOfAdaptability []string) (madl.MADL) {
//	mEE := madl.MADL{}
//	appIsAdaptive := true
//
//	if len(appKindOfAdaptability) == 1 && appKindOfAdaptability[0] == "NONE" { // TODO
//		appIsAdaptive = false
//	}
//
//	// configuration
//	mEE.Configuration = mapp.Configuration + "_ee"
//
//	// adaptability of app
//	mEE.AppAdaptability = mapp.Adaptability
//
//	// Components
//	comps := []madl.Element{}
//	comps = append(comps, madl.Element{Id: "core", TypeName: reflect.TypeOf(components.Core{}).Name()})
//
//	if appIsAdaptive {
//		comps = append(comps, madl.Element{Id: "monevolutive", TypeName: reflect.TypeOf(components.Monevolutive{}).Name()})
//		comps = append(comps, madl.Element{Id: "monitor", TypeName: reflect.TypeOf(components.Monitor{}).Name()})
//		comps = append(comps, madl.Element{Id: "analyser", TypeName: reflect.TypeOf(components.Analyser{}).Name()})
//		comps = append(comps, madl.Element{Id: "planner", TypeName: reflect.TypeOf(components.Planner{}).Name()})
//		comps = append(comps, madl.Element{Id: "executor", TypeName: reflect.TypeOf(components.Executor{}).Name()})
//	}
//
//	units := []string{}
//	for i := 0; i < len(mapp.Components); i++ {
//		units = append(units, "unit"+strconv.Itoa(i+1))
//	}
//	for i := 0; i < len(units); i++ {
//		comps = append(comps, madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()})
//	}
//
//	// Connectors
//	conns := [] madl.Element{}
//
//	params := make([]interface{}, 1)
//	params[0] = len(units)
//
//	// TODO
//	switch mapp.Configuration {
//	case "senderreceiver":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
//	case "middlewareclient":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params})
//	case "middlewareserver":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params})
//	case "clientserverlocal":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
//	case "clientserver":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
//	case "calculatorlocal":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params})
//	case "midfibonacciserver":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto8{}).Name(), Params: params})
//	case "midfibonacciclient":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto6{}).Name(), Params: params})
//	case "midnamingserver":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params})
//	case "queueingserver":
//		conns = append(conns, madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto9{}).Name(), Params: params})
//
//	default:
//		fmt.Printf("Creator:: Configuration '%v' cannot be executed because 'OnetoN.dot' is not ok!!", mapp.Configuration)
//		os.Exit(0)
//	}
//
//	if appIsAdaptive {
//		conns = append(conns, madl.Element{Id: "t2", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
//		conns = append(conns, madl.Element{Id: "t3", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
//		conns = append(conns, madl.Element{Id: "t4", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
//		conns = append(conns, madl.Element{Id: "t5", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
//		conns = append(conns, madl.Element{Id: "t6", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()})
//	}
//
//	// Attachments
//	atts := []madl.Attachment{}
//
//	if appIsAdaptive {
//		attC1 := madl.Element{Id: "monevolutive", TypeName: reflect.TypeOf(components.Monevolutive{}).Name()}
//		attT := madl.Element{Id: "t2", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
//		attC2 := madl.Element{Id: "monitor", TypeName: reflect.TypeOf(components.Monitor{}).Name()}
//		atts = append(atts, madl.Attachment{attC1, attT, attC2})
//
//		attC1 = madl.Element{Id: "monitor", TypeName: reflect.TypeOf(components.Monitor{}).Name()}
//		attT = madl.Element{Id: "t3", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
//		attC2 = madl.Element{Id: "analyser", TypeName: reflect.TypeOf(components.Analyser{}).Name()}
//		atts = append(atts, madl.Attachment{attC1, attT, attC2})
//
//		attC1 = madl.Element{Id: "analyser", TypeName: reflect.TypeOf(components.Analyser{}).Name()}
//		attT = madl.Element{Id: "t4", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
//		attC2 = madl.Element{Id: "planner", TypeName: reflect.TypeOf(components.Planner{}).Name()}
//		atts = append(atts, madl.Attachment{attC1, attT, attC2})
//
//		attC1 = madl.Element{Id: "planner", TypeName: reflect.TypeOf(components.Planner{}).Name()}
//		attT = madl.Element{Id: "t5", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
//		attC2 = madl.Element{Id: "executor", TypeName: reflect.TypeOf(components.Executor{}).Name()}
//		atts = append(atts, madl.Attachment{attC1, attT, attC2})
//
//		attC1 = madl.Element{Id: "executor", TypeName: reflect.TypeOf(components.Executor{}).Name()}
//		attT = madl.Element{Id: "t6", TypeName: reflect.TypeOf(connectors.Oneway{}).Name()}
//		attC2 = madl.Element{Id: "core", TypeName: reflect.TypeOf(components.Core{}).Name()}
//		atts = append(atts, madl.Attachment{attC1, attT, attC2})
//	}
//
//	for i := 0; i < len(units); i++ {
//		attC1 := madl.Element{Id: "core", TypeName: reflect.TypeOf(components.Core{}).Name()}
//
//		// TODO
//		switch mapp.Configuration {
//		case "senderreceiver":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "middlewareclient":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "middlewareserver":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "clientserverlocal":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "calculatorlocal":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "clientserver":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.OnetoN{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "midfibonacciserver":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto8{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "midfibonacciclient":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto6{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "midnamingserver":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto5{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//		case "queueingserver":
//			attT := madl.Element{Id: "t1", TypeName: reflect.TypeOf(connectors.Oneto9{}).Name(), Params: params}
//			attC2 := madl.Element{Id: units[i], TypeName: reflect.TypeOf(components.Unit{}).Name()}
//			atts = append(atts, madl.Attachment{attC1, attT, attC2})
//
//
//		default:
//			fmt.Printf("Creator:: Configuration '%v' cannot be executed because 'OnetoN.dot' is not ok ", mapp.Configuration)
//			os.Exit(0)
//		}
//	}
//
//	// Adaptability
//	eeKindOfAdaptability := []string{}
//	eeKindOfAdaptability = append(eeKindOfAdaptability, "NONE")
//
//	// configure MADL EE
//	mEE.File = strings.Replace(mapp.File, shared.MADL_EXTENSION, "", 99) + "_ee" + shared.MADL_EXTENSION
//	mEE.Path = mapp.Path
//	mEE.Components = comps
//	mEE.Connectors = conns
//	mEE.Attachments = atts
//	mEE.Adaptability = eeKindOfAdaptability
//
//	return mEE
//}

func (Creator) Create(mapp madl.MADL, appKindOfAdaptability []string) madl.MADL {
	mEE := madl.MADL{}
	appIsAdaptive := true

	if len(appKindOfAdaptability) == 1 && appKindOfAdaptability[0] == shared.NON_ADAPTIVE { // TODO
		appIsAdaptive = false
	}

	// configuration
	mEE.Configuration = mapp.Configuration + "_ee"

	// adaptability of app
	mEE.Adaptability = mapp.Adaptability

	// Components
	comps := []component.Component{}
	comps = append(comps, component.Component{Id: "core", TypeName: reflect.TypeOf(adaptive.Core{}).Name()}) // TODO dcruzb: verify if this component is necessary

	if appIsAdaptive {
		if appKindOfAdaptability[0] == shared.EVOLUTIVE_ADAPTATION {
			comps = append(comps, component.Component{Id: "monevolutive", TypeName: reflect.TypeOf(adaptive.Monevolutive{}).Name()})
		}
		if appKindOfAdaptability[0] == shared.EVOLUTIVE_PROTOCOL_ADAPTATION {
			comps = append(comps, component.Component{Id: "evolutiveprotocol", TypeName: reflect.TypeOf(adaptive.EvolutiveProtocol{}).Name()})
		}
		comps = append(comps, component.Component{Id: "monitor", TypeName: reflect.TypeOf(adaptive.Monitor{}).Name()})
		comps = append(comps, component.Component{Id: "analyser", TypeName: reflect.TypeOf(adaptive.Analyser{}).Name()})
		comps = append(comps, component.Component{Id: "planner", TypeName: reflect.TypeOf(adaptive.Planner{}).Name()})
		comps = append(comps, component.Component{Id: "executor", TypeName: reflect.TypeOf(adaptive.Executor{}).Name()})
	}

	units := []string{}
	for i := 0; i < len(mapp.Components); i++ { //+len(mapp.Connectors)
		units = append(units, "unit"+strconv.Itoa(i+1))
	}
	for i := 0; i < len(units); i++ {
		comps = append(comps, component.Component{Id: units[i], TypeName: reflect.TypeOf(adaptive.Unit{}).Name()})
	}

	// Connectors
	cncts := []connectors.Connector{}

	params := make([]interface{}, 1)
	params[0] = len(units)

	nAttT1 := len(mapp.Components) //+ len(mapp.Connectors)

	behaviour := "B = InvP.e1 -> " //InvR.e2 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> InvR.e7 -> B"
	for i := 0; i < nAttT1; i++ {
		behaviour += fmt.Sprintf("InvR.e'%d' ->", i+2)
	}
	behaviour += "B"

	cncts = append(cncts, connectors.NewConnector(
		"t1",
		shared.ONETON,
		behaviour,
		1, nAttT1))

	//switch nAttT1 {
	//case 3:
	//	cncts = append(cncts, connectors.NewConnector(
	//		"t1",
	//		shared.ONETON,
	//		"B = InvP.e1 -> InvR.e2 -> P2 [] I_Timeout -> P2 \n P2 = InvR.e3 -> P3 [] I_Timeout -> P3 \n P3 = InvR.e4 -> B [] I_Timeout -> B",
	//		1, 3))
	//case 5:
	//	cncts = append(cncts, connectors.Connector{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto5{}).Name(), Params: params})
	//case 6:
	//	cncts = append(cncts, connectors.Connector{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto6{}).Name(), Params: params})
	//case 8:
	//	cncts = append(cncts, connectors.Connector{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto8{}).Name(), Params: params})
	//case 9:
	//	cncts = append(cncts, connectors.Connector{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto9{}).Name(), Params: params})
	//default:
	//	fmt.Printf("Creator:: Configuration '%v' cannot be executed because 'OnetoN.dot' is not ok!!", mapp.Configuration)
	//	os.Exit(0)
	//}

	// Attachments
	atts := []madl.Attachment{}

	if appIsAdaptive {
		cncts = append(cncts, connectors.NewConnector("t2", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1))
		cncts = append(cncts, connectors.NewConnector("t3", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1))
		cncts = append(cncts, connectors.NewConnector("t4", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1))
		cncts = append(cncts, connectors.NewConnector("t5", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1))
		cncts = append(cncts, connectors.NewConnector("t6", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1))

		var attC1 component.Component
		if appKindOfAdaptability[0] == shared.EVOLUTIVE_ADAPTATION {
			attC1 = component.Component{Id: "monevolutive", TypeName: reflect.TypeOf(adaptive.Monevolutive{}).Name()}
		} else {
			attC1 = component.Component{Id: "evolutiveprotocol", TypeName: reflect.TypeOf(adaptive.EvolutiveProtocol{}).Name()}
		}
		attT := connectors.NewConnector("t2", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1)
		attC2 := component.Component{Id: "monitor", TypeName: reflect.TypeOf(adaptive.Monitor{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = component.Component{Id: "monitor", TypeName: reflect.TypeOf(adaptive.Monitor{}).Name()}
		attT = connectors.NewConnector("t3", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1)
		attC2 = component.Component{Id: "analyser", TypeName: reflect.TypeOf(adaptive.Analyser{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = component.Component{Id: "analyser", TypeName: reflect.TypeOf(adaptive.Analyser{}).Name()}
		attT = connectors.NewConnector("t4", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1)
		attC2 = component.Component{Id: "planner", TypeName: reflect.TypeOf(adaptive.Planner{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = component.Component{Id: "planner", TypeName: reflect.TypeOf(adaptive.Planner{}).Name()}
		attT = connectors.NewConnector("t5", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1)
		attC2 = component.Component{Id: "executor", TypeName: reflect.TypeOf(adaptive.Executor{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		attC1 = component.Component{Id: "executor", TypeName: reflect.TypeOf(adaptive.Executor{}).Name()}
		attT = connectors.NewConnector("t6", shared.ONEWAY, "B = InvP.e1 -> InvR.e2 -> B", 1, 1)
		attC2 = component.Component{Id: "core", TypeName: reflect.TypeOf(adaptive.Core{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})
	}

	for i := 0; i < len(units); i++ {
		attC1 := component.Component{Id: "core", TypeName: reflect.TypeOf(adaptive.Core{}).Name()}

		behaviour := "B = InvP.e1 -> " //InvR.e2 -> InvR.e3 -> InvR.e4 -> InvR.e5 -> InvR.e6 -> InvR.e7 -> B"
		for i := 0; i < nAttT1; i++ {
			behaviour += fmt.Sprintf("InvR.e'%d' ->", i+2)
		}
		behaviour += "B"

		attT := connectors.NewConnector(
			"t1",
			shared.ONETON,
			behaviour,
			1, nAttT1)
		attC2 := component.Component{Id: units[i], TypeName: reflect.TypeOf(adaptive.Unit{}).Name()}
		atts = append(atts, madl.Attachment{attC1, attT, attC2})

		//switch nAttT1 {
		//
		//case 3:
		//	attT := component.Component{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto3{}).Name(), Params: params}
		//	attC2 := component.Component{Id: units[i], TypeName: reflect.TypeOf(adaptive.Unit{}).Name()}
		//	atts = append(atts, madl.Attachment{attC1, attT, attC2})
		//case 5:
		//	attT := component.Component{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto5{}).Name(), Params: params}
		//	attC2 := component.Component{Id: units[i], TypeName: reflect.TypeOf(adaptive.Unit{}).Name()}
		//	atts = append(atts, madl.Attachment{attC1, attT, attC2})
		//case 6:
		//	attT := component.Component{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto6{}).Name(), Params: params}
		//	attC2 := component.Component{Id: units[i], TypeName: reflect.TypeOf(adaptive.Unit{}).Name()}
		//	atts = append(atts, madl.Attachment{attC1, attT, attC2})
		//case 8:
		//	attT := component.Component{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto8{}).Name(), Params: params}
		//	attC2 := component.Component{Id: units[i], TypeName: reflect.TypeOf(adaptive.Unit{}).Name()}
		//	atts = append(atts, madl.Attachment{attC1, attT, attC2})
		//case 9:
		//	attT := component.Component{Id: "t1", TypeName: reflect.TypeOf(cncts.Oneto9{}).Name(), Params: params}
		//	attC2 := component.Component{Id: units[i], TypeName: reflect.TypeOf(adaptive.Unit{}).Name()}
		//	atts = append(atts, madl.Attachment{attC1, attT, attC2})
		//default:
		//	fmt.Printf("Creator:: Configuration '%v' cannot be executed because 'OnetoN.dot' is not ok ", mapp.Configuration)
		//	os.Exit(0)
		//}
	}

	// Adaptability
	eeKindOfAdaptability := []string{}
	eeKindOfAdaptability = append(eeKindOfAdaptability, shared.NON_ADAPTIVE)

	// configure MADL EE
	mEE.FileName = strings.Replace(mapp.FileName, shared.MADL_EXTENSION, "", 99) + "_ee" + shared.MADL_EXTENSION
	mEE.Path = mapp.Path
	mEE.Components = comps
	mEE.Connectors = cncts
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
		fmt.Printf("      %v : %v \n", madl.Components[i].Id, madl.Components[i].TypeName)
	}
	fmt.Printf("   Connectors \n")
	for i := range madl.Connectors {
		fmt.Printf("      %v : %v \n", madl.Connectors[i].Id, madl.Connectors[i].TypeName)
	}
	fmt.Printf("   Attachments \n")
	for i := range madl.Attachments {
		fmt.Printf("      %v,%v,%v\n", madl.Attachments[i].C1.Id, madl.Attachments[i].T.Id, madl.Attachments[i].C2.Id)
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
		content = append(content, "      "+m.Components[i].Id+" : "+m.Components[i].TypeName+" \n")
	}
	content = append(content, "\n    Connectors \n")
	for i := range m.Connectors {
		content = append(content, "      "+m.Connectors[i].Id+" : "+m.Connectors[i].TypeName+" \n")
	}
	content = append(content, "\n    Attachments \n")
	for i := range m.Attachments {
		content = append(content, "      "+m.Attachments[i].C1.Id+","+m.Attachments[i].T.Id+","+m.Attachments[i].C2.Id+" \n")
	}

	content = append(content, "\n   Adaptability \n")
	content = append(content, "      "+m.Adaptability[0]+" \n\n") // TODO

	content = append(content, "EndConf \n")

	shared.SaveFile(path, name, ext, content)
}
