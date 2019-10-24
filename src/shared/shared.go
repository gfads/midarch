package shared

import (
	"errors"
	"fmt"
	"gmidarch/development/messages"
	"io/ioutil"
	"log"
	"os"
	"plugin"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ValidActions = map[string]bool{
	INVP: true,
	TERP: true,
	INVR: true,
	TERR: true}

var SetOfPorts = map[string]int{
	"NAMING_PORT":     NAMING_PORT,
	"CALCULATOR_PORT": CALCULATOR_PORT,
	"FIBONACCI_PORT":  FIBONACCI_PORT,
	"QUEUEING_PORT":   QUEUEING_PORT}

var IS_EVOLUTIVE = false
var IS_CORRECTIVE = false
var IS_PROACTIVE = false

var MONITOR_TIME time.Duration   // seconds
var INJECTION_TIME time.Duration // seconds
var REQUEST_TIME time.Duration   // milliseconds
var STRATEGY = 0                 // 1 - no change 2 - change once 3 - change same plugin 4 - alternate plugins
var SAMPLE_SIZE = 0
var NAMING_HOST = ""
var QUEUEING_HOST = ""

// MAPE-K Types
type MonitoredCorrectiveData string   // used in channel Monitor -> Analyser (Corrective)

type MonitoredEvolutiveData [] string // used in channel Monitor -> Analyser (Evolutive)

type MonitoredProactiveData [] string // used in channel Monitor -> Analyser (Proactive)

type EvolutiveAnalysisResult struct {
	NeedAdaptation         bool
	MonitoredEvolutiveData MonitoredEvolutiveData
}

type UnitCommand struct {
	Cmd    string
	Params plugin.Plugin
	Type   interface{}
}

type AdaptationPlan struct {
	Operations [] string
	Params     map[string][]string
}

type Request struct {
	Op   string
	Args []interface{}
}

type Invocation struct {
	Host string
	Port int
	Req  Request
}

type QueueingInvocation struct {
	Op   string
	Args []interface{}
}

type QueueingTermination struct {
	R interface{}
}

type ParMapActions struct {
	F1 func(*chan messages.SAMessage, *messages.SAMessage)     // External action
	F2 func(any interface{}, name string, args ...interface{}) // Internal action
	P1 *messages.SAMessage
	P2 *chan messages.SAMessage
	P3 interface{}
	P4 string
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
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

func LoadParameters(args []string) {
	for i := range args {
		variable := strings.Split(args[i], "=")
		switch strings.TrimSpace(variable[0]) {
		case "SAMPLE_SIZE":
			SAMPLE_SIZE, _ = strconv.Atoi(variable[1])
		case "REQUEST_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			REQUEST_TIME = temp2
		case "INJECTION_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			INJECTION_TIME = temp2
		case "MONITOR_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			MONITOR_TIME = temp2
		case "STRATEGY":
			STRATEGY, _ = strconv.Atoi(variable[1])
		case "NAMING_HOST":
			NAMING_HOST = variable[1]
		case "QUEUEING_HOST":
			QUEUEING_HOST = variable[1]
		default:
			fmt.Println("Shared:: Parameter '" + variable[0] + "' does not exist")
			os.Exit(0)
		}
	}
}

func ShowExecutionParameters(s bool) {
	if s {
		fmt.Println("******************************************")
		fmt.Println("Sample size                : " + strconv.Itoa(SAMPLE_SIZE))
		fmt.Println("Direrctory of base code    : " + BASE_DIR)
		fmt.Println("Directory of plugins       : " + DIR_PLUGINS)
		fmt.Println("Directory of CSP specs     : " + DIR_CSP)
		fmt.Println("Directory of Configurations: " + DIR_MADL)
		fmt.Println("Directory of Go compiler   : " + DIR_GO)
		fmt.Println("Directory of FDR           : " + DIR_FDR)
		fmt.Println("------------------------------------------")
		fmt.Println("Naming Host     : " + NAMING_HOST)
		fmt.Println("Naming Port     : " + strconv.Itoa(NAMING_PORT))
		fmt.Println("Calculator Port : " + strconv.Itoa(CALCULATOR_PORT))
		fmt.Println("Fibonacci Port  : " + strconv.Itoa(FIBONACCI_PORT))
		fmt.Println("Queueing Port   : " + strconv.Itoa(QUEUEING_PORT))
		fmt.Println("------------------------------------------")
		fmt.Println("Plugin Base Name: " + PLUGIN_BASE_NAME)
		fmt.Println("Max Graph Size  : " + strconv.Itoa(GRAPH_SIZE))
		fmt.Println("------------------------------------------")
		fmt.Println("Adaptability  ")
		fmt.Println("Corrective        : " + strconv.FormatBool(IS_CORRECTIVE))
		fmt.Println("Evolutive         : " + strconv.FormatBool(IS_EVOLUTIVE))
		fmt.Println("Proactive         : " + strconv.FormatBool(IS_PROACTIVE))
		fmt.Println("Monitor Time (s)  : " + (MONITOR_TIME * time.Second).String())
		fmt.Println("Injection Time (s): " + (INJECTION_TIME * time.Second).String())
		fmt.Println("Request Time (ms) : " + REQUEST_TIME.String())
		fmt.Println("Strategy (0-NOT DEFINED 1-No change 2-Change once 3-change same plugin 4-alternate plugins): " + strconv.Itoa(STRATEGY))
		fmt.Println("******************************************")
	}
}

func Log(args ...string) {
	if strings.Contains(args[1], "Proxy") || strings.Contains(args[1], "XXX") {
		fmt.Println(args[0] + ":" + args[1] + ":" + args[2] + ":" + args[3])
	}
}

func Invoke(any interface{}, name string, msg *messages.SAMessage, info [] *interface{}) {
	inputs := make([]reflect.Value, 2)
	inputs[0] = reflect.ValueOf(msg)
	inputs[1] = reflect.ValueOf(info)

	//fmt.Printf("Shared:: %v %v %v %v\n",reflect.TypeOf(any),name, msg, info)

	reflect.ValueOf(any).MethodByName(name).Call(inputs)

	inputs = nil
	return
}

func LoadPlugins() map[string]time.Time {
	listOfPlugins := make(map[string]time.Time)

	pluginsDir := DIR_PLUGINS
	OSDir, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		fmt.Printf("Shared:: Folder '%v' is unreadeable\n", pluginsDir)
		os.Exit(0)
	}
	for i := range OSDir {
		fileName := OSDir[i].Name()
		pluginFile := pluginsDir + "/" + fileName
		info, err := os.Stat(pluginFile)
		if err != nil {
			fmt.Printf("Shared:: Plugin '%v' not readeable\n", pluginFile)
			os.Exit(0)
		}
		listOfPlugins[fileName] = info.ModTime()
	}
	return listOfPlugins
}

func CheckForNewPlugins(listOfOldPlugins map[string]time.Time, listOfNewPlugins map[string]time.Time) [] string {
	var newPlugins [] string

	// check for new plugins
	for key := range listOfNewPlugins {
		val1, _ := listOfNewPlugins[key]
		val2, ok2 := listOfOldPlugins[key]
		if ok2 {
			if val1.After(val2) { // newer version of an old plugin is available
				newPlugins = append(newPlugins, key)
			}
		} else {
			newPlugins = append(newPlugins, key) // a new plugin is available
		}
	}
	return newPlugins
}

func LoadPlugin(pluginName string) (plugin.Plugin) {

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := DIR_PLUGINS + "/" + pluginName
	attempts := 0
	for {
		plg, err = plugin.Open(pluginFile)

		if err != nil {
			if attempts >= 1000 { // TODO
				fmt.Printf("Shared:: Error on trying open plugin '%v' \n", pluginFile)
				os.Exit(0)
			} else {
				attempts++
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			break
		}
	}

	// look for a exportable function/variable within the plugin
	//fx, err := lib.Lookup(symbolName)
	//if err != nil {
	//	fmt.Printf( "Shared:: Function '%v' not found in plugin '%v'\n",symbolName,pluginName)
	//	os.Exit(0)
	//}
	//return fx

	return *plg
}

func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:: %s", msg, err)
		os.Exit(1)
	}
}

func SkipLine(line string) bool {

	if line == "" || strings.TrimSpace(line)[:2] == ADL_COMMENT {
		return true
	} else {
		return false
	}
}

func IsAdaptationType(line string) bool {
	r := false

	line = strings.TrimSpace(strings.ToUpper(line))
	if line == CORRECTIVE || line == EVOLUTIVE || line == PROACTIVE || line == EMPTY_LINE {
		r = true
	}
	return r
}

func MyTokenize(s string) [] string {
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

func CheckFileName(fileName string) error {
	r := *new(error)
	r = nil

	len := len(fileName)

	if len <= 5 {
		r = errors.New("File Name Invalid")
	} else {
		if fileName[len-5:] != MADL_EXTENSION {
			r = errors.New("Invalid extension of '" + fileName + "'")
		}
	}

	return r
}

func SaveFile(path, name, ext string, content []string) {

	// create diretcory if it does not exist
	confDir := path
	_, err := os.Stat(confDir);
	if os.IsNotExist(err) {
		os.MkdirAll(confDir, os.ModePerm);
	}

	// create file if it does not exist && truncate otherwise
	file, err := os.Create(confDir + "/" + name + ext)
	if err != nil {
		fmt.Println("Shared:: File " + path + "/" + name + "." + ext + "not created!!")
		os.Exit(0)
	}
	defer file.Close()

	// save data
	for i := range content {
		_, err = file.WriteString(content[i])
		if err != nil {
			fmt.Println("Shared:: File " + path + "/" + name + "." + ext + "not saved!!")
			os.Exit(0)
		}
	}
	err = file.Sync()
	if err != nil {
		fmt.Println("Shared:: File " + path + "/" + name + "." + ext + "not synced!!")
		os.Exit(0)
	}
	defer file.Close()
}

// ******************* PARAMETERS

//const BASE_DIR  = "/go/midarch-go"  // docker
const BASE_DIR = "/Users/nsr/Dropbox/go/midarch-go-v13"
const DIR_PLUGINS = BASE_DIR + "/src/gmidarch/execution/repositories/plugins"
const DIR_PLUGINS_SOURCE = BASE_DIR + "/src/gmidarch/development/repositories/plugins"
const DIR_CSP = BASE_DIR + "/src/apps/artefacts/csps"
const DIR_MADL = BASE_DIR + "/src/apps/artefacts/madls"
const DIR_GO = "/usr/local/go/bin"
const DIR_FDR = "/Volumes/Macintosh HD/Applications/FDR4-2.app/Contents/MacOS"
const DIR_CSPARSER = BASE_DIR + "/src/verificationtools/cspdot/csparser"
const DIR_DOT = BASE_DIR + "/src/gmidarch/development/repositories/dot"
const MADL_EXTENSION = ".madl"
const CSP_EXTENSION = ".csp"
const DOT_EXTENSION = ".dot"
const RUNTIME_BEHAVIOUR = "RUNTIME"

const DEADLOCK_PROPERTY = "assert " + CORINGA + " :[deadlock free]"
const CORINGA = "XXX"

const ADL_COMMENT = "//"

// Ports
const NAMING_PORT = 4040
const CALCULATOR_PORT = 2020
const FIBONACCI_PORT = 2030
const QUEUEING_PORT = 2040

const CHAN_BUFFER_SIZE = 1
const PLUGIN_BASE_NAME = "receiver"
const GRAPH_SIZE = 30

const PREFIX_INTERNAL_ACTION = "I_"
const INVP = "InvP"
const TERP = "TerP"
const INVR = "InvR"
const TERR = "TerR"
const EVOLUTIVE = "EVOLUTIVE"
const CORRECTIVE = "REACTIVE"
const PROACTIVE = "PROACTIVE"
const EMPTY_LINE = "NONE"

const EXECUTE_FOREVER = true

const REPLACE_COMPONENT = "REPLACE_COMPONENT"
const FDR_COMMAND = "refines"
