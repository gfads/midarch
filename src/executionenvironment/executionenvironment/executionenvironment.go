package executionenvironment

import (
	"shared/shared"
	"os"
	"verificationtools/fdr"
	"framework/configuration/configuration"
	"framework/message"
	"graph/execgraph"
	"reflect"
	"strings"
	"shared/errors"
	"framework/configuration/commands"
	"strconv"
	"shared/parameters"
	"executionenvironment/adaptationmanager"
	"executionenvironment/versioninginjector"
	"executionenvironment/executionunit"
	"framework/libraries"
)

type ExecutionEnvironment struct{}

func (ee ExecutionEnvironment) Deploy(adlFileName string) {

	// Load execution parameters
	shared.LoadParameters(os.Args[1:])

	// Perform checks on library
	libraries.CheckLibrary()

	// Generate Go configuration
	conf := configuration.MapADLIntoGo(adlFileName)

	// Configure management channels
	managementChannels := InitializeManagementChannels(conf)

	// Configure structural channels & maps
	ee.ConfigureStructuralChannelsAndMaps(&conf)

	// Check behaviour using FDR
	fdrChecker := new(fdr.FDR)
	ok := fdrChecker.CheckBehaviour(conf)
	if !ok {
		myError := errors.MyError{Source: "Execution Engine", Message: "Configuration has a problem detected by FDR4"}
		myError.ERROR()
	}

	// Load graph generated by FDR (*.dot)
	fdrChecker.LoadFDRGraph(&conf)

	// Generate executable graph
	CreateExecGraph(&conf)

	// Show execution parameters
	shared.ShowExecutionParameters(false)

	// Start engine
	go StartEngine(conf.StateMachine)

	// Start adaptation manager
	if parameters.IS_CORRECTIVE || parameters.IS_EVOLUTIVE || parameters.IS_PROACTIVE {
		go adaptationmanager.AdaptationManager{}.Exec(conf, managementChannels)
		go versioninginjector.InjectAdaptiveEvolution(parameters.PLUGIN_BASE_NAME)
	}

	// Start execution units
	for e := range conf.Components {
		go executionunit.ExecutionUnit{}.Exec(conf.Components[e], conf.StructuralChannels, managementChannels[conf.Components[e].Id])
	}
}

func InitializeManagementChannels(conf configuration.Configuration) map[string]chan commands.LowLevelCommand {
	managementChannels := make(map[string]chan commands.LowLevelCommand)
	for i := range conf.Components {
		id := conf.Components[i].Id
		managementChannels[id] = make(chan commands.LowLevelCommand)
	}
	return managementChannels
}

func CreateExecGraph(conf *configuration.Configuration) {
	graph := execgraph.NewGraph(conf.FDRGraph.NumNodes)
	channels := map[string]chan message.Message{}

	// create channels
	for e1 := range conf.FDRGraph.Edges {
		for e2 := range conf.FDRGraph.Edges[e1] {
			eTemp := conf.FDRGraph.Edges[e1][e2]
			if _, ok := channels[eTemp.Action]; !ok {
				channels[eTemp.Action] = make(chan message.Message)
			}
			graph.AddEdgeX(eTemp.From, eTemp.To, execgraph.Action{Action: eTemp.Action, Channel: channels[eTemp.Action]})
		}
	}

	conf.StateMachine = *graph
	conf.StructuralChannels = channels
}

func StartEngine(g execgraph.Graph) {
	node := 0
	var msg = message.Message{}
	for {
		edges := g.AdjacentEdges(node)
		if len(edges) == 1 { // one edge
			node = edges[0].To
			if shared.IsToElement(edges[0].Action.Action) {
				edges[0].Action.Channel <- msg
			} else {
				msg = <-edges[0].Action.Channel
			}
		} else { // two+ edges
			chosen := 0
			Choice(&msg, &chosen, edges)
			node = edges[chosen].To
		}
	}
}

func Choice(msg *message.Message, chosen *int, edges []execgraph.Edge) {
	cases := make([]reflect.SelectCase, len(edges))
	var value reflect.Value

	for i := 0; i < len(edges); i++ {
		if shared.IsToElement(edges[i].Action.Action) {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.ValueOf(*msg)}
		} else {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(edges[i].Action.Channel), Send: reflect.Value{}}
		}
	}

	*chosen, value, _ = reflect.Select(cases)
	if !shared.IsToElement(edges[*chosen].Action.Action) {
		*msg = value.Interface().(message.Message)
	}
	cases = nil
}

func (ExecutionEnvironment) ConfigureStructuralChannelsAndMaps(conf *configuration.Configuration) {
	structuralChannels := make(map[string]chan message.Message)

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
		structuralChannels[key01] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key04] = structuralChannels[key03]

		// t -> c2
		key01 = tId + "." + shared.INVR + "." + c2Id
		key02 = c2Id + "." + shared.INVP + "." + tId
		key03 = c2Id + "." + shared.TERP + "." + tId
		key04 = tId + "." + shared.TERR + "." + c2Id
		structuralChannels[key01] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
		structuralChannels[key02] = structuralChannels[key01]
		structuralChannels[key03] = make(chan message.Message, parameters.CHAN_BUFFER_SIZE)
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