package shared

import (
	"bufio"
	"fmt"
	"gmidarch/development/messages"
	"log"
	"os"
	"plugin"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var ExecuteForever = true

// File extensions
const MADL_EXTENSION = "madl"
const DOT_EXTENSION = "dot"

// Graphs
const PREFIX_INTERNAL_ACTION = "I_"
const MAXIMUM_GRAPH_SIZE = 15
const NUM_MAX_NODES int = 50
const EXECUTE_FOREVER = true

// Directories
var DIR_GO = LocalizegGo() + "/bin"
var DIR_BASE = LocalizegMidArch()
var DIR_MADL = DIR_BASE + "/src/apps/artefacts/madls"
var DIR_DOT = DIR_BASE + "/src/gmidarch/development/repositories/dot"
var DIR_ADAPTIVE_COMPONENTS = DIR_BASE + "/src/gmidarch/development/components/adaptive"
var DIR_APP_COMPONENTS = DIR_BASE + "/src/gmidarch/development/components/apps"
var DIR_MIDDLEWARE_COMPONENTS = DIR_BASE + "/src/gmidarch/development/components/middleware"
var DIR_PROXIES_COMPONENTS = DIR_BASE + "/src/gmidarch/development/components/proxies"
var DIR_PLUGINS = DIR_BASE + "/src/gmidarch/execution/repositories/plugins"
var DIR_PLUGINS_SOURCE = DIR_BASE + "/src/gmidarch/development/repositories/plugins"
var DIR_PLUGINS_IMPORT = "gmidarch/development/repositories/plugins"
var AdaptId = -1
var LocalAddr = ""
var ArchitecturalComponentTypes = map[string]interface{}{}
var Adaptability []string

// MADL
const MADL_COMMENT = "//"

// Connectors
const ONEWAY = "Oneway"
const REQUEST_REPLY = "Requestreply"
const ONETON = "Oneton"
const NTOONE = "Ntoone"
const NTOONEREQREP = "Ntoonereqrep"
const ONETONREQREP = "Onetonreqrep"
const LEFT_ARITY = 1
const RIGHT_ARITY = 2
const MAX_LEFT_ARITY = 99
const MAX_RIGHT_ARITY = 99

// Port names
const INVP = "InvP"
const TERP = "TerP"
const INVR = "InvR"
const TERR = "TerR"

// Network setups
const NAMING_PORT = "1313"
const NAMING_HOST = "namingserver"
const CALCULATOR_HOST = "server"
const CALCULATOR_PORT = "1314"
const FIBONACCI_PORT = "1315"
const QUEUEING_PORT = "1316"
const MAX_NUMBER_OF_CONNECTIONS = 10
const SIZE_OF_MESSAGE_SIZE = 4
const NUM_MAX_MESSAGE_BYTES int = 1024
const MAX_NUMBER_OF_RECEIVED_MESSAGES = 300 // messages received and not processed by srh

const ATTEMPTS_TO_OPEN_A_PLUGIN = 1000

// Evolution
const FIRST_MONITOR_TIME time.Duration = 5 * time.Second
const MONITOR_TIME time.Duration = 100 * time.Millisecond

var INJECTION_TIME time.Duration

var SetOfPorts = map[string]string{
	"NAMING_PORT":     NAMING_PORT,
	"CALCULATOR_PORT": CALCULATOR_PORT,
	"FIBONACCI_PORT":  FIBONACCI_PORT,
	"QUEUEING_PORT":   QUEUEING_PORT}

var AdaptationTypes = map[string]string{
	"EVOLUTIVE": "EVOLUTIVE",
	"EVOLUTIVE_PROTOCOL": "EVOLUTIVE_PROTOCOL",
	"NONE":      "NONE"}

type MonitoredEvolutiveData []string // used in channel Monitor -> Analyser (Evolutive)

type EvolutiveAnalysisResult struct {
	NeedAdaptation         bool
	MonitoredEvolutiveData MonitoredEvolutiveData
}

type AdaptationPlan struct {
	Operations []string
	Params     map[string][]string
}

type UnitCommand struct {
	Cmd      string
	Params   plugin.Plugin
	Type     interface{}
	Selector func(interface{}, []*interface{}, string, *messages.SAMessage, []*interface{}, *bool)
}

const EVOLUTIVE_ADAPTATION string = "EVOLUTIVE"
const EVOLUTIVE_PROTOCOL_ADAPTATION string = "EVOLUTIVE_PROTOCOL"
const NON_ADAPTIVE string = "NONE"

const REPLACE_COMPONENT = "REPLACE_COMPONENT"

// CSP
const RUNTIME_BEHAVIOUR = "RUNTIME"
const CSP_EXTENSION = "csp"
const CORINGA = "XXX"
const DEADLOCK_PROPERTY = "assert " + CORINGA + " :[deadlock free]"
const BEHAVIOUR_TAG = "//@Behaviour:"
const TYPE_TAG = "//@Type:"
const ACTION_PREFIX = "->"
const CHOICE = "[]"
const BEHAVIOUR_ID = "Behaviour"

var DIR_FDR = LocalizegFDR()
var DIR_CSP = DIR_BASE + "/src/apps/artefacts/csp"

const FDR_COMMAND = "refines"

// Utility functions
func MyInvoke(compType interface{}, compId string, op string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	inputs := make([]reflect.Value, 4)

	inputs[0] = reflect.ValueOf(compId)
	inputs[1] = reflect.ValueOf(msg)
	inputs[2] = reflect.ValueOf(info)
	inputs[3] = reflect.ValueOf(reset)

	//fmt.Println("MyInvoke( compId:", compId, "- msg:", msg, "- info:", info, "- Method Name:", op, ")")
	reflect.ValueOf(compType).MethodByName(op).Call(inputs)
}

func ErrorHandler(f string, msg string) {
	fmt.Println(f + "::" + msg)
	os.Exit(1)
}

func GetFunction() string {
	fpcs := make([]uintptr, 1)

	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)
	if n == 0 {
		fmt.Println("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		fmt.Println("MSG CALLER WAS NIL")
	}

	// Print the file name and line number
	//fmt.Println(caller.FileLine(fpcs[0]-1))

	// Print the name of the function
	//fmt.Println(caller.Name())

	return caller.Name()
}

func LocalizegGo() string {
	r := ""
	found := false

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == "GOROOT" {
			r = pair[1]
			found = true
		}
	}

	if !found {
		fmt.Println("Shared:: Error:: OS Environment variable 'GOROOT' not configured\n")
		os.Exit(1)
	}
	return r
}

func LocalizegMidArch() string {
	r := ""
	found := false

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == "GMIDARCH" {
			r = pair[1]
			found = true
		}
	}

	if !found {
		ErrorHandler(GetFunction(), "OS Environment variable 'GMIDARCH' not configured")
	}
	return r
}

func LocalizegFDR() string {
	r := ""
	found := false

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == "FDR4" {
			r = pair[1]
			found = true
		}
	}

	if !found {
		fmt.Println("Shared:: Error:: OS Environment variable 'FDR4' not configured\n")
		os.Exit(1)
	}
	return r
}

func SkipLine(line string) bool {

	if line == "" || strings.TrimSpace(line)[:2] == MADL_COMMENT {
		return true
	} else {
		return false
	}
}

func CheckFileName(fileName, fileExtension string) {
	if len(fileName) <= len(fileExtension)+1 {
		ErrorHandler(GetFunction(), "File Name '"+fileName+"'+Invalid")
	} else {

		switch fileExtension {
		case MADL_EXTENSION:
			if fileName[len(fileName)-4:] != MADL_EXTENSION {
				ErrorHandler(GetFunction(), "Invalid extension of '"+fileName+"'")
			}
		case DOT_EXTENSION:
			if fileName[len(fileName)-3:] != DOT_EXTENSION {
				ErrorHandler(GetFunction(), "Invalid extension of '"+fileName+"'")
			}
		default:
		}
	}
}

func IsInternal(action string) bool {

	if len(action) <= 2 {
		return false
	}
	if action[0:2] == PREFIX_INTERNAL_ACTION {
		return true
	}
	return false
}

func IsExternal(action string) bool {
	r := false

	if len(action) >= 2 {
		action := strings.TrimSpace(action)
		if strings.Contains(action, ".") {
			action = action[:strings.Index(action, ".")]
		}

		if action == INVP || action == TERP || action == INVR || action == TERR {
			r = true
		} else {
			r = false
		}
	} else {
		r = false
	}
	return r
}

func MyTokenize(s string) []string {
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

func StringComposition(e []string, sep string, hasSpace bool) string {
	r1 := ""

	for i := range e {
		if hasSpace {
			r1 += e[i] + " " + sep + " "
		} else {
			r1 += e[i] + sep
		}
	}

	if hasSpace {
		r1 = r1[:len(r1)-len(sep)-2]
	} else {
		r1 = r1[:len(r1)-len(sep)]
	}

	return r1
}

func SaveFile(path, name, ext string, content []string) {

	// create diretcory if it does not exist
	confDir := path
	_, err := os.Stat(confDir)
	if os.IsNotExist(err) {
		os.MkdirAll(confDir, os.ModePerm)
	}

	// create file if it does not exist && truncate otherwise
	file, err := os.Create(confDir + "/" + name)
	if err != nil {
		ErrorHandler(GetFunction(), "File "+path+"/"+name+"."+ext+"not created!!")
	}
	defer file.Close()

	// save data
	for i := range content {
		_, err = file.WriteString(content[i])
		if err != nil {
			ErrorHandler(GetFunction(), "File "+path+"/"+name+"."+ext+"not saved!!")
		}
	}
	err = file.Sync()
	if err != nil {
		ErrorHandler(GetFunction(), "File "+path+"/"+name+"."+ext+"not synced!!")
	}
	defer file.Close()
}

func GetTypeAndBehaviour(file string) (string, string) {
	typeName := ""
	behaviour := ""
	foundType := false
	foundBehaviour := false

	// Open file
	f, err := os.Open(file)
	if err != nil {
		ErrorHandler(GetFunction(), err.Error())
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(f)

	// Read file & indentify types/behaviours
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, BEHAVIOUR_TAG) {
			behaviour = line[strings.Index(line, BEHAVIOUR_TAG)+len(BEHAVIOUR_TAG)+1:]
			foundBehaviour = true
		}
		if strings.Contains(line, TYPE_TAG) {
			typeName = strings.TrimSpace(line[strings.Index(line, TYPE_TAG)+len(TYPE_TAG)+1:])
			foundType = true
		}

		if foundType && foundBehaviour {
			break
		}
	}

	// Check wether type/behaviour information is complete or not
	if !foundType || !foundBehaviour {
		ErrorHandler(GetFunction(), "Tags '"+BEHAVIOUR_TAG+"' or '"+TYPE_TAG+"' are missing in '"+file+"''")
	}

	return typeName, behaviour
}

func CompatibleComponents(componentTypeName1, componentTypeName2 string) bool {
	return (strings.Contains(componentTypeName1, "SRH") && strings.Contains(componentTypeName2, "SRH")) ||
		   (strings.Contains(componentTypeName1, "CRH") && strings.Contains(componentTypeName2, "CRH")) ||
		   (componentTypeName1 == componentTypeName2)
}

func GetComponentTypeByNameFromRAM(componentName string) interface{} {
	return ArchitecturalComponentTypes[componentName]
}