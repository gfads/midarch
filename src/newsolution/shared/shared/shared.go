package shared

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"newsolution/gmidarch/development/messages"
	"newsolution/shared/parameters"
	"os"
	"plugin"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//type Invocation struct {
//	Method  reflect.Value
//	InArgs  []reflect.Value
//	OutArgs [] reflect.Value
//}

// MAPE-K Types
type MonitoredCorrectiveData string   // used in channel Monitor -> Analyser (Corrective)
type MonitoredEvolutiveData [] string // used in channel Monitor -> Analyser (Evolutive)
type MonitoredProactiveData [] string // used in channel Monitor -> Analyser (Proactive)

type AnalysisResult struct {
	Result   interface{}
	Analysis int
}

type AdaptationPlan struct {
	Plan string
}

var ValidActions = map[string]bool{
	parameters.INVP: true,
	parameters.TERP: true,
	parameters.INVR: true,
	parameters.TERR: true}

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
	if action[0:2] == parameters.PREFIX_INTERNAL_ACTION {
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

		if action == parameters.INVP || action == parameters.TERP || action == parameters.INVR || action == parameters.TERR {
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
			parameters.SAMPLE_SIZE, _ = strconv.Atoi(variable[1])
		case "REQUEST_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.REQUEST_TIME = temp2
		case "INJECTION_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.INJECTION_TIME = temp2
		case "MONITOR_TIME":
			temp1, _ := strconv.Atoi(variable[1])
			temp2 := time.Duration(temp1)
			parameters.MONITOR_TIME = temp2
		case "STRATEGY":
			parameters.STRATEGY, _ = strconv.Atoi(variable[1])
		case "NAMING_HOST":
			parameters.NAMING_HOST = variable[1]
		case "QUEUEING_HOST":
			parameters.QUEUEING_HOST = variable[1]
		default:
			fmt.Println("Shared:: Parameter '" + variable[0] + "' does not exist")
			os.Exit(0)
		}
	}
}

func ShowExecutionParameters(s bool) {
	if s {
		fmt.Println("******************************************")
		fmt.Println("Sample size                : " + strconv.Itoa(parameters.SAMPLE_SIZE))
		fmt.Println("Direrctory of base code    : " + parameters.BASE_DIR)
		fmt.Println("Directory of plugins       : " + parameters.DIR_PLUGINS)
		fmt.Println("Directory of CSP specs     : " + parameters.DIR_CSP)
		fmt.Println("Directory of Configurations: " + parameters.DIR_MADL)
		fmt.Println("Directory of Go compiler   : " + parameters.DIR_GO)
		fmt.Println("Directory of FDR           : " + parameters.DIR_FDR)
		fmt.Println("------------------------------------------")
		fmt.Println("Naming Host     : " + parameters.NAMING_HOST)
		fmt.Println("Naming Port     : " + strconv.Itoa(parameters.NAMING_PORT))
		fmt.Println("Calculator Port : " + strconv.Itoa(parameters.CALCULATOR_PORT))
		fmt.Println("Fibonacci Port  : " + strconv.Itoa(parameters.FIBONACCI_PORT))
		fmt.Println("Queueing Port   : " + strconv.Itoa(parameters.QUEUEING_PORT))
		fmt.Println("------------------------------------------")
		fmt.Println("Plugin Base Name: " + parameters.PLUGIN_BASE_NAME)
		fmt.Println("Max Graph Size  : " + strconv.Itoa(parameters.GRAPH_SIZE))
		fmt.Println("------------------------------------------")
		fmt.Println("Adaptability  ")
		fmt.Println("Corrective        : " + strconv.FormatBool(parameters.IS_CORRECTIVE))
		fmt.Println("Evolutive         : " + strconv.FormatBool(parameters.IS_EVOLUTIVE))
		fmt.Println("Proactive         : " + strconv.FormatBool(parameters.IS_PROACTIVE))
		fmt.Println("Monitor Time (s)  : " + (parameters.MONITOR_TIME * time.Second).String())
		fmt.Println("Injection Time (s): " + (parameters.INJECTION_TIME * time.Second).String())
		fmt.Println("Request Time (ms) : " + parameters.REQUEST_TIME.String())
		fmt.Println("Strategy (0-NOT DEFINED 1-No change 2-Change once 3-change same plugin 4-alternate plugins): " + strconv.Itoa(parameters.STRATEGY))
		fmt.Println("******************************************")
	}
}

func Log(args ...string) {
	if strings.Contains(args[1], "Proxy") || strings.Contains(args[1], "XXX") {
		fmt.Println(args[0] + ":" + args[1] + ":" + args[2] + ":" + args[3])
	}
}

/*
func Invoke(any interface{}, name string, args ... interface{}) {
	inputs := make([]reflect.Value, len(args))

	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	fmt.Println(any)
	fmt.Println(name)
	fmt.Println(len(inputs))

	reflect.ValueOf(any).MethodByName(name).Call(inputs)

	os.Exit(0)
	inputs = nil
	return
}
*/

func InvokeOld1(any interface{}, name string, args [] *interface{}) {
	inputs := make([]reflect.Value, len(args))

	for i, _ := range args {
		inputs[i] = reflect.ValueOf(*args[i])
	}

	fmt.Printf("Shared:: %v %v %v\n", reflect.TypeOf(any), name, inputs)

	reflect.ValueOf(any).MethodByName(name).Call(inputs)

	inputs = nil
	return
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

func InvokeNew(any interface{}, name string, args [] reflect.Value) {
	inputs := make([]reflect.Value, len(args))

	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	fmt.Printf("Shared:: %v %v %v\n", reflect.TypeOf(any), name, inputs)

	reflect.ValueOf(any).MethodByName(name).Call(inputs)

	inputs = nil
	return
}

func InvokeOld(any interface{}, name string, args [] reflect.Value) {

	//fmt.Printf("Shared:: %v %v %vn",reflect.TypeOf(any),name, args)
	reflect.ValueOf(any).MethodByName(name).Call(args)

	return
}

func LoadPlugins(confName string) map[string]time.Time {
	listOfPlugins := make(map[string]time.Time)

	pluginsDir := parameters.DIR_PLUGINS + "/" + confName
	OSDir, err := ioutil.ReadDir(pluginsDir)
	CheckError(err, "Shared:: Folder '"+pluginsDir+"' not read")
	for i := range OSDir {
		fileName := OSDir[i].Name()
		if strings.Contains(fileName, "_plugin") {
			pluginFile := pluginsDir + "/" + fileName
			info, err := os.Stat(pluginFile)
			CheckError(err, "Shared:: Plugin '"+pluginFile+"'not read")
			listOfPlugins[fileName] = info.ModTime()
		}
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

func LoadPlugin(confName string, pluginName string, symbolName string) (plugin.Symbol) {

	var lib *plugin.Plugin
	var err error

	pluginFile := parameters.DIR_PLUGINS + "/" + confName + "/" + pluginName
	attempts := 0
	for {
		lib, err = plugin.Open(pluginFile)

		if err != nil {
			if attempts >= 3 {
				CheckError(err, "Shared:: Error on trying open plugin "+pluginFile+" "+strconv.Itoa(attempts)+" times")
			} else {
				attempts++
				time.Sleep(10 * time.Millisecond)
			}
		} else {
			break
		}
	}

	fx, err := lib.Lookup(symbolName)
	CheckError(err, "Shared:: Message: Function "+symbolName+" not found in plugin")

	return fx
}

func CheckError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:: %s", msg, err)
		os.Exit(1)
	}
}

func SkipLine(line string) bool {

	if line == "" || strings.TrimSpace(line)[:2] == parameters.ADL_COMMENT {
		return true
	} else {
		return false
	}
}

func IsAdaptationType(line string) bool {
	r := false

	line = strings.TrimSpace(strings.ToUpper(line))
	if line == parameters.CORRECTIVE || line == parameters.EVOLUTIVE || line == parameters.PROACTIVE || line == parameters.EMPTY_LINE {
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

type Request struct {
	Op   string
	Args []interface{}
}

type Invocation struct {
	Host string
	Port int
	Req  Request
}

func CheckFileName(fileName string) error {
	r := *new(error)
	r = nil

	len := len(fileName)

	if len <= 5 {
		r = errors.New("File Name Invalid")
	} else {
		if fileName[len-5:] != parameters.MADL_EXTENSION {
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
