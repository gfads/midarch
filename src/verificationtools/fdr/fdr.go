package fdr

import (
	"strings"
	"strconv"
	"os"
	"framework/element"
	"os/exec"
	"framework/configuration/configuration"
	"newsolution/shared/shared"
	"newsolution/shared/parameters"
	"newsolution/shared/error"
	"framework/configuration/commands"
	"graph/fdrgraph"
	"log"
	"bufio"
	"framework/architecturallibrary"
	"reflect"
	"fmt"
)

type FDR struct{}

func (fdr FDR) LoadFDRGraphs(conf *configuration.Configuration) {

	// Load component
	dotDir := parameters.DIR_CSP + "/" + strings.Replace(conf.Id, ".confs", "", 99)
	for c := range conf.Components {
		dotFileName := strings.ToUpper(conf.Components[c].Id) + ".dot"
		dotFileName = dotDir + "/" + dotFileName

		fileContent := loadDotFile(dotFileName)
		graph := fdr.CreateFDRGraph(fileContent)
		tempComp := conf.Components[c]
		tempComp.SetFDRGraph(*graph)
		conf.Components[c] = tempComp
	}

	// Load connectors
	for t := range conf.Connectors {
		dotFileName := strings.ToUpper(conf.Connectors[t].Id) + ".dot"
		dotFileName = dotDir + "/" + dotFileName

		fileContent := loadDotFile(dotFileName)
		graph := fdr.CreateFDRGraph(fileContent)
		tempConn := conf.Connectors[t]
		tempConn.SetFDRGraph(*graph)
		conf.Connectors[t] = tempConn
	}
}

func (FDR) CreateFDRGraph(fileContent []string) *fdrgraph.Graph {
	graph := fdrgraph.NewGraph(100)

	for l := range fileContent {
		line := fileContent[l]
		if strings.Contains(line, "->") {
			from, _ := strconv.Atoi(strings.TrimSpace(line[:strings.Index(line, "->")]))
			to, _ := strconv.Atoi(strings.TrimSpace(line[strings.Index(line, "->")+2 : strings.Index(line, "[")]))
			action := line[strings.Index(line, "=")+2 : strings.LastIndex(line, "]")-2]
			graph.AddEdge(from, to, action)
		}
	}
	return graph
}

func loadDotFile(dotFileName string) []string {

	// read file
	fileContent := []string{}
	file, err := os.Open(dotFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	return fileContent
}

func (FDR) SaveCSP(conf configuration.Configuration) {

	// create diretcory if it does not exist
	confDir := parameters.DIR_CSP + "/" + strings.Replace(conf.Id, ".confs", "", 99)
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		os.MkdirAll(confDir, os.ModePerm);
	}

	// create file if it does not exist && truncate otherwise
	fileName := conf.Id + ".artefacts"
	file, err := os.Create(confDir + "/" + fileName)
	if err != nil {
		myError := error.MyError{Source: "FDR", Message: "CSP File not created"}
		myError.ERROR()
	}
	defer file.Close()

	// save data
	_, err = file.WriteString(conf.CSP)
	if err != nil {
		myError := error.MyError{Source: "FDR", Message: "CSP Specification not saved"}
		myError.ERROR()
	}
	err = file.Sync()
	if err != nil {
		myError := error.MyError{Source: "FDR", Message: "CSP File not Synced"}
		myError.ERROR()
	}
	defer file.Close()
}

func (FDR) InvokeFDR(conf configuration.Configuration) bool {
	cmdExp := parameters.DIR_FDR + "/" + commands.FDR_COMMAND
	fileName := conf.Id + ".artefacts"
	dirFile := parameters.DIR_CSP + "/" + strings.Replace(conf.Id, ".confs", "", 99)
	inputFile := dirFile + "/" + fileName

	out, err := exec.Command(cmdExp, inputFile).Output()
	if err != nil {
		fmt.Println(err)
		myError := error.MyError{Source: "FDR", Message: "File '" + inputFile + "' has a problem (e.g.,syntax error)"}
		myError.ERROR()
	}
	s := string(out[:])

	if strings.Contains(s, "Passed") {
		return true
	} else {
		myError := error.MyError{Source: "FDR", Message: "Deadlock detected"}
		myError.ERROR()
		return false
	}
}

func (FDR) CreateCSP(conf configuration.Configuration) string {

	// Configure behaviour expressions of components and conncetores
	configureBehaviours(&conf)

	// general behaviour
	dataTypeExp := createDataTypeExp(conf)
	internalChannelsExp, _ := createInternalChannelExp(conf)
	externalChannelsExp, externalChannels := createExternalChannelExp(conf)
	processesExp, componentProcesses, _ := createProcessExp(conf)
	generalBehaviour := createGeneralBehaviourExp(&conf, externalChannels, componentProcesses)

	// assertion
	assertion := "assert " + conf.Id + " :[deadlock free]"
	csp := dataTypeExp + "\n" + internalChannelsExp + "\n" + externalChannelsExp + "\n" + processesExp + "\n" + generalBehaviour + "\n" + assertion

	return csp
}

func configureBehaviours(conf *configuration.Configuration) {

	// Configure components
	for c := range conf.Components {
		standardBehaviour := architecturallibrary.Repository[reflect.TypeOf(conf.Components[c].TypeElem).String()].CSP
		tokens := strings.Split(standardBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := conf.Components[c].Id + "." + eX
				partner := conf.Maps[key]
				standardBehaviour = strings.Replace(standardBehaviour, eX, partner, 99)
			}
		}
		configuredBehaviour := strings.Replace(standardBehaviour, "B", strings.ToUpper(conf.Components[c].Id), 99)
		conf.Components[c] = element.Element{Id: conf.Components[c].Id, TypeElem: conf.Components[c].TypeElem, CSP: configuredBehaviour}
	}

	// Configure connectors
	for t := range conf.Connectors {
		standardBehaviour := architecturallibrary.Repository[reflect.TypeOf(conf.Connectors[t].TypeElem).String()].CSP
		tokens := strings.Split(standardBehaviour, " ")
		for j := range tokens {
			if shared.IsExternal(tokens[j]) {
				eX := tokens[j][strings.Index(tokens[j], ".")+1:]
				key := conf.Connectors[t].Id + "." + eX
				partner := conf.Maps[key]
				standardBehaviour = strings.Replace(standardBehaviour, eX, partner, 99)
			}
		}
		configuredBehaviour := strings.Replace(standardBehaviour, "B", strings.ToUpper(conf.Connectors[t].Id), 99)
		conf.Connectors[t] = element.Element{Id: conf.Connectors[t].Id, TypeElem: conf.Connectors[t].TypeElem, CSP: configuredBehaviour}
	}
}

func renameSyncPorts(conf *configuration.Configuration, elem element.Element) string {
	id := elem.Id
	tokens := strings.Split(elem.CSP, " ")
	renamingExp := strings.ToUpper(id) + "[["

	for i := range tokens {
		token := strings.TrimSpace(tokens[i])
		if shared.IsExternal(token) {
			action := token[0:strings.Index(token, ".")]
			switch action {
			case parameters.INVP:
				renamingExp += token + " <- " + parameters.INVR + "." + id + ","
			case parameters.TERP:
				renamingExp += token + " <- " + parameters.TERR + "." + id + ","
			case parameters.INVR:
				renamingExp += token + " <- " + parameters.INVP + "." + id + ","
			case parameters.TERR:
				renamingExp += token + " <- " + parameters.TERP + "." + id + ","
			}
		}
	}
	renamingExp = renamingExp[0:strings.LastIndex(renamingExp, ",")] + "]]"
	return renamingExp
}

func createDataTypeExp(conf configuration.Configuration) string {
	dataTypes := make(map[string]string)

	for i := range conf.Components {
		dataTypes [conf.Components[i].Id] = conf.Components[i].Id
	}
	for i := range conf.Connectors {
		dataTypes [conf.Connectors[i].Id] = conf.Connectors[i].Id
	}

	dataTypeExp := "datatype PROCNAMES = "
	for i := range dataTypes {
		dataTypeExp += dataTypes[i] + " | "
	}
	dataTypeExp = dataTypeExp[0:strings.LastIndex(dataTypeExp, "|")]

	return dataTypeExp
}

func myTokenize(s string) [] string {
	tokens := []string{}

	token := ""
	for i := 0; i < len(s); i++ {
		c := s[i : i+1]
		switch c {
		case "=":
			token = ""
		case "-":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case " ":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case "]":
			token = ""
		case ">":
			token = ""
		case "\n":
			token = ""
		case "[":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case "(":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		case ")":
			if strings.TrimSpace(token) != "" {
				tokens = append(tokens, token)
			}
			token = ""
		default:
			token += c
		}
	}
	return tokens
}

func createInternalChannelExp(conf configuration.Configuration) (string, map[string]string) {
	internalChannels := make(map[string]string)

	for i := range conf.Components {
		//tokens := strings.Split(conf.Components[i].CSP, " ")
		tokens := myTokenize(conf.Components[i].CSP)
		for i := range tokens {
			if shared.IsInternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				internalChannels[iAction] = iAction
			}
		}
	}

	for i := range conf.Connectors {
		tokens := myTokenize(conf.Connectors[i].CSP)
		for i := range tokens {
			if shared.IsInternal(tokens[i]) {
				iAction := strings.TrimSpace(tokens[i])
				internalChannels[iAction] = iAction
			}
		}
	}
	internalChannelsExp := "channel "
	for i := range internalChannels {
		internalChannelsExp += internalChannels[i] + ","
	}
	internalChannelsExp = internalChannelsExp[0:strings.LastIndex(internalChannelsExp, ",")]

	return internalChannelsExp, internalChannels
}

func createExternalChannelExp(conf configuration.Configuration) (string, map[string]string) {
	externalChannels := make(map[string]string)

	for i := range conf.Components {
		b := conf.Components[i].CSP
		if strings.Contains(b, parameters.INVR) {
			externalChannels[parameters.INVR] = parameters.INVR
		}
		if strings.Contains(b, parameters.TERR) {
			externalChannels[parameters.TERR] = parameters.TERR
		}
		if strings.Contains(b, parameters.INVP) {
			externalChannels[parameters.INVP] = parameters.INVP
		}
		if strings.Contains(b, parameters.TERP) {
			externalChannels[parameters.TERP] = parameters.TERP
		}
	}
	for i := range conf.Connectors {
		b := conf.Connectors[i].CSP
		if strings.Contains(b, parameters.INVR) {
			externalChannels[parameters.INVR] = parameters.INVR
		}
		if strings.Contains(b, parameters.TERR) {
			externalChannels[parameters.TERR] = parameters.TERR
		}
		if strings.Contains(b, parameters.INVP) {
			externalChannels[parameters.INVP] = parameters.INVP
		}
		if strings.Contains(b, parameters.TERP) {
			externalChannels[parameters.TERP] = parameters.TERP
		}
	}

	externalChannelsExp := "channel "
	for i := range externalChannels {
		externalChannelsExp += externalChannels[i] + ","
	}
	externalChannelsExp = externalChannelsExp[0:strings.LastIndex(externalChannelsExp, ",")] + " : PROCNAMES"

	return externalChannelsExp, externalChannels
}

func createProcessExp(conf configuration.Configuration) (string, map[string]string, map[string]string) {

	componentProcesses := make(map[string]string)
	processesExp := ""
	for i := range conf.Components {
		comp := conf.Components[i]
		id := comp.Id
		subprocesses := strings.Split(comp.CSP, "\n")
		if len(subprocesses) > 1 {
			componentProcesses[strings.ToUpper(id)] = renameSubprocesses(comp.CSP)
		} else {
			componentProcesses[strings.ToUpper(id)] = comp.CSP
		}
		processesExp += componentProcesses[strings.ToUpper(id)] + "\n"
	}
	connectorProcesses := make(map[string]string)
	for i := range conf.Connectors {
		id := conf.Connectors[i].Id
		connectorProcesses[strings.ToUpper(id)] = conf.Connectors[i].CSP
		processesExp += connectorProcesses[strings.ToUpper(id)] + "\n"
	}
	return processesExp, componentProcesses, connectorProcesses
}

func renameSubprocesses(p string) string {
	subprocesses := strings.Split(p, "\n")
	r := ""
	procBaseName := strings.TrimSpace(subprocesses[0][:strings.Index(subprocesses[0], "=")]) // first process

	newProcNames := map[string]string{}
	for i := 1; i < len(subprocesses); i++ {
		procName := strings.TrimSpace(subprocesses[i][:strings.Index(subprocesses[i], "=")])
		newProcNames[procName] = procBaseName + procName
	}

	for i := 0; i < len(subprocesses); i++ {
		for j := range newProcNames {
			subprocesses[i] = strings.Replace(subprocesses[i], j, newProcNames[j], 99)
		}
		r += subprocesses[i] + "\n"
	}
	return r
}

func createGeneralBehaviourExp(conf *configuration.Configuration, externalChannels map[string]string, componentProcesses map[string]string) string {

	// (C1 ||| C2 ||| C3)
	generalBehaviour := conf.Id + " = ("
	for i := range componentProcesses {
		generalBehaviour += i + "|||"
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, "|||")] + ")"

	generalBehaviour += "\n" + "[|{|"
	for i := range externalChannels {
		generalBehaviour += externalChannels[i] + ","
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, ",")] + "|}|]" + "\n"

	generalBehaviour += "("
	for i := range conf.Connectors {
		generalBehaviour += renameSyncPorts(conf, conf.Connectors[i]) + "|||"
	}
	generalBehaviour = generalBehaviour[0:strings.LastIndex(generalBehaviour, "|||")] + ")"

	return generalBehaviour
}
