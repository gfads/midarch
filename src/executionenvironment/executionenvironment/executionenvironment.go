package executionenvironment

import (
	"shared/shared"
	"os"
	"verificationtools/fdr"
	"framework/configuration/configuration"
	"framework/messages"
	"graph/execgraph"
	"strings"
	"shared/errors"
	"framework/configuration/commands"
	"strconv"
	"shared/parameters"
	"framework/element"
	"fmt"
	"reflect"
	"framework/libraries"
	"executionenvironment/executionunit"
)

type ExecutionEnvironment struct{}

func (environment ExecutionEnvironment) Deploy(adlFileName string) {

	// Initialize environment
	environment.Initialization()

	// Prepare configuration to be executed
	appConf := environment.PrepareConfiguration(adlFileName)

	// Prepare MAPE-K
	//adaptationManagerConf := environment.PrepareConfiguration("MAPEK.conf")

	// Configure management channels
	managementChannelsApp := environment.ConfigureManagementChannels(appConf)
	//managementChannelsMAPEK := environment.ConfigureManagementChannels(appConf)

	// Start adaptation manager
	//if parameters.IS_CORRECTIVE || parameters.IS_EVOLUTIVE || parameters.IS_PROACTIVE {
	//	go adaptationmanager.AdaptationManager{}.Exec(appConf, managementChannelsApp)
	//	go versioninginjector.InjectAdaptiveEvolution(parameters.PLUGIN_BASE_NAME)
	//}

	// Start App Configuration
	environment.StartConfiguration(appConf, managementChannelsApp)
	//environment.StartConfiguration(adaptationManagerConf, managementChannelsMAPEK)
}

func (environment ExecutionEnvironment) StartConfiguration(conf configuration.Configuration, managementChannels map[string]chan commands.LowLevelCommand) {
	// Start execution units
	for c := range conf.Components {
		go executionunit.ExecutionUnit{}.Exec(conf.Components[c], managementChannels[conf.Components[c].Id])
	}
	for t := range conf.Connectors {
		go executionunit.ExecutionUnit{}.Exec(conf.Connectors[t], managementChannels[conf.Connectors[t].Id])
	}
}

func (environment ExecutionEnvironment) PrepareConfiguration(adlFileName string) configuration.Configuration {

	// Generate Go configuration
	conf := configuration.MapADLIntoGo(adlFileName)

	// Configure structural channels and maps of components/connectors
	environment.ConfigureStructuralChannelsAndMaps(&conf)

	// Check behaviour using FDR
	fdrChecker := new(fdr.FDR)
	ok := fdrChecker.CheckBehaviour(conf)
	if !ok {
		myError := errors.MyError{Source: "Execution Engine", Message: "Configuration has a problem detected by FDR4"}
		myError.ERROR()
	}

	// Generate *.dot files
	// FDR.GenerateFDRGraphs()  // TODO

	// Load graph generated by FDR (*.dot)
	fdrChecker.LoadFDRGraphs(&conf)

	// Generate executable graph
	environment.CreateExecGraphs(&conf)

	// Check if actions and their respective implementations exist
	CheckActionsAndImplementations(conf)

	return conf

}

func (ee ExecutionEnvironment) Initialization() {
	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Perform checks on the library of c
	libraries.CheckLibrary()

	// Show execution parameters
	shared.ShowExecutionParameters(false)
}

func CheckActionsAndImplementations(conf configuration.Configuration) {

	// Check components
	for c := range conf.Components {
		for e1 := range conf.Components[c].ExecGraph.Edges {
			for e2 := range conf.Components[c].ExecGraph.Edges[e1] {
				action := conf.Components[c].ExecGraph.Edges[e1][e2].Action.ActionName
				if shared.IsExternal(action) {
					if action != shared.INVP && action != shared.TERP && action != shared.INVR && action != shared.TERR {
						fmt.Println("EE:: Component '" + conf.Components[c].Id + "' has an invalid action: '" + action)
						os.Exit(0)
					}
				} else {
					if shared.IsInternal(action) {
						st := reflect.TypeOf(conf.Components[c].TypeElem)
						_, ok := st.MethodByName(action)
						if !ok {
							fmt.Println("EE: Component '" + conf.Components[c].Id + "' has an invalid action: '" + action + "'")
							os.Exit(0)
						}

					} else {
						fmt.Println("EE: Component '" + conf.Components[c].Id + "' has an invalid action: '" + action + "'")
						os.Exit(0)
					}
				}
			}
		}
	}
	// Check connectors
	for t := range conf.Connectors {
		for e1 := range conf.Connectors[t].ExecGraph.Edges {
			for e2 := range conf.Connectors[t].ExecGraph.Edges[e1] {
				action := conf.Connectors[t].ExecGraph.Edges[e1][e2].Action.ActionName
				if shared.IsExternal(action) {
					if action != shared.INVP && action != shared.TERP && action != shared.INVR && action != shared.TERR {
						fmt.Println("EE:: Connector '" + conf.Connectors[t].Id + "' has an invalid action: '" + action)
						os.Exit(0)
					}
				} else {
					if shared.IsInternal(action) {
						st := reflect.TypeOf(conf.Connectors[t].TypeElem)
						_, ok := st.MethodByName(action)
						if !ok {
							fmt.Println("EE: Connector '" + conf.Connectors[t].Id + "' has an invalid action: '" + action + "'")
							os.Exit(0)
						}

					} else {
						fmt.Println("EE: Connector '" + conf.Connectors[t].Id + "' has an invalid action: '" + action + "'")
						os.Exit(0)
					}
				}
			}
		}
	}
}

func (ee ExecutionEnvironment) ConfigureManagementChannels(conf configuration.Configuration) map[string]chan commands.LowLevelCommand {
	managementChannels := make(map[string]chan commands.LowLevelCommand)
	for i := range conf.Components {
		id := conf.Components[i].Id
		managementChannels[id] = make(chan commands.LowLevelCommand)
	}
	return managementChannels
}

func (ee ExecutionEnvironment) CreateExecGraphs(conf *configuration.Configuration) {

	// Components
	for c := range conf.Components {
		graph := execgraph.NewGraph(conf.Components[c].FDRGraph.NumNodes)
		eActions := execgraph.Action{}
		var msg messages.SAMessage
		for e1 := range conf.Components[c].FDRGraph.Edges {
			for e2 := range conf.Components[c].FDRGraph.Edges[e1] {
				edgeTemp := conf.Components[c].FDRGraph.Edges[e1][e2]
				actionNameFDR := edgeTemp.Action
				actionNameExec := ""
				if strings.Contains(actionNameFDR, ".") {
					actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
				}
				if shared.IsExternal(actionNameExec) { // External action
					key := conf.Components[c].Id + "." + actionNameFDR
					channel := conf.StructuralChannels[key]
					params := execgraph.Action{}
					switch actionNameExec {
					case "InvR":
						invr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
					case "TerR":
						terr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
					case "InvP":
						invp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
					case "TerP":
						terp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerP, Message: &msg, ActionChannel: &terp, ActionName: actionNameExec}
					}
					mapType := execgraph.Action{}
					mapType = params
					eActions = mapType
				}

				if shared.IsInternal(actionNameFDR) {
					msg := messages.SAMessage{}
					channel := make(chan messages.SAMessage)
					params := execgraph.Action{InternalAction: shared.Invoke, ActionName: actionNameFDR, Message: &msg, ActionChannel: &channel}
					mapType := params
					eActions = mapType
				}
				graph.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
			}
		}
		tempComp := conf.Components[c]
		tempComp.SetExecGraph(graph)
		conf.Components[c] = tempComp
	}

	// Connectors
	for t := range conf.Connectors {
		graph := execgraph.NewGraph(conf.Connectors[t].FDRGraph.NumNodes)
		eActions := execgraph.Action{}
		var msg messages.SAMessage
		for e1 := range conf.Connectors[t].FDRGraph.Edges {
			for e2 := range conf.Connectors[t].FDRGraph.Edges[e1] {
				edgeTemp := conf.Connectors[t].FDRGraph.Edges[e1][e2]
				actionNameFDR := edgeTemp.Action
				actionNameExec := ""
				if strings.Contains(actionNameFDR, ".") {
					actionNameExec = actionNameFDR[:strings.Index(actionNameFDR, ".")]
				}
				if shared.IsExternal(actionNameExec) { // External action
					key := conf.Connectors[t].Id + "." + actionNameFDR
					channel := conf.StructuralChannels[key]
					params := execgraph.Action{}
					switch actionNameExec {
					case "InvR":
						invr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvR, Message: &msg, ActionChannel: &invr, ActionName: actionNameExec}
					case "TerR":
						terr := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerR, Message: &msg, ActionChannel: &terr, ActionName: actionNameExec}
					case "InvP":
						invp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.InvP, Message: &msg, ActionChannel: &invp, ActionName: actionNameExec}
					case "TerP":
						terp := channel
						params = execgraph.Action{ExternalAction: element.Element{}.TerP, Message: &msg, ActionChannel: &terp, ActionName: actionNameExec}
					}
					mapType := execgraph.Action{}
					mapType = params
					eActions = mapType
				}

				if shared.IsInternal(actionNameFDR) {
					msg := messages.SAMessage{}
					params := execgraph.Action{InternalAction: shared.Invoke, ActionName: actionNameFDR, Message: &msg}
					mapType := params
					eActions = mapType
				}
				graph.AddEdge(edgeTemp.From, edgeTemp.To, eActions)
			}
		}
		tempComp := conf.Connectors[t]
		tempComp.SetExecGraph(graph)
		conf.Connectors[t] = tempComp
	}
}

func (ExecutionEnvironment) ConfigureStructuralChannelsAndMaps(conf *configuration.Configuration) {
	structuralChannels := make(map[string]chan messages.SAMessage)

	// Configure structural channels
	for i := range conf.Attachments {
		c1Id := conf.Attachments[i].C1.Id
		c2Id := conf.Attachments[i].C2.Id
		tId := conf.Attachments[i].T.Id

		// c1 -> t
		key01 := c1Id + "." + shared.INVR + "." + tId
		key02 := tId + "." + shared.INVP + "." + c1Id
		key03 := tId + "." + shared.TERP + "." + c1Id
		key04 := c1Id + "." + shared.TERR + "." + tId
		structuralChannels[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]

		// t -> c2
		key01 = tId + "." + shared.INVR + "." + c2Id
		key02 = c2Id + "." + shared.INVP + "." + tId
		key03 = c2Id + "." + shared.TERP + "." + tId
		key04 = tId + "." + shared.TERR + "." + c2Id
		structuralChannels[key01] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan messages.SAMessage, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]
	}
	conf.StructuralChannels = structuralChannels

	// Configure maps
	elemMaps := make(map[string]string)
	partners := make(map[string]string)

	for i := range conf.Attachments {
		c1Id := conf.Attachments[i].C1.Id
		c2Id := conf.Attachments[i].C2.Id
		tId := conf.Attachments[i].T.Id
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
	conf.Maps = elemMaps
}
