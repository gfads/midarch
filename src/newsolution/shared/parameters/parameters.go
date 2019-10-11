package parameters

import "time"

// Dirs
//const BASE_DIR  = "/go/midarch-go"  // docker
const BASE_DIR = "/Users/nsr/Dropbox/go/midarch-go-v13"
const DIR_PLUGINS = BASE_DIR + "/src/plugins"
const DIR_CSP = BASE_DIR + "/src/newsolution/apps/csps"
const DIR_MADL = BASE_DIR + "/src/newsolution/apps/madls"
const DIR_GO = "/usr/local/go/bin"
const DIR_FDR = "/Volumes/Macintosh HD/Applications/FDR4-2.app/Contents/MacOS"
const DIR_CSPARSER = BASE_DIR + "/src/verificationtools/cspdot/csparser"
const DIR_DOT = BASE_DIR + "/src/newsolution/gmidarch/development/repositories/dot"
const COMPONENTS_PATH = "components"
const CONNECTORS_PATH = "connectors"
const NAMINGCLIENTPROXY_PATH = "namingclientproxy"
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

//Java parameters
const JAVA_COMMAND = "java"
const JAR_COMMAND = "-jar"

var SetOfPorts = map[string]int{
	"NAMING_PORT":     NAMING_PORT,
	"CALCULATOR_PORT": CALCULATOR_PORT,
	"FIBONACCI_PORT":  FIBONACCI_PORT,
	"QUEUEING_PORT":   QUEUEING_PORT}

const NO_CHANGE = 0
const REACTIVE_CHANGE = 1
const EVOLUTIVE_CHANGE = 2
const PROACTIVE_CHANGE = 3

const CHAN_BUFFER_SIZE = 1
const QUEUE_SIZE = 100

//const PLUGIN_BASE_NAME  = "calculatorinvoker"
//const PLUGIN_BASE_NAME = "fibonacciinvoker"
const PLUGIN_BASE_NAME = "receiver"
const GRAPH_SIZE = 30

const MAX_NUMBER_OF_ACTIVE_CONSUMERS = 10

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

const PREFIX_ACTION = "->"

//const CHOICE = "[]"
const PREFIX_INTERNAL_ACTION = "I_"
const INVP = "InvP"
const TERP = "TerP"
const INVR = "InvR"
const TERR = "TerR"
const EVOLUTIVE = "EVOLUTIVE"
const CORRECTIVE = "REACTIVE"
const PROACTIVE = "PROACTIVE"
const EMPTY_LINE = "NONE"
const DEFAULT_INFO = ""
const PREFIX_ADL_EXECUTION_ENVIRONMENT = "ExecutionEnvironment"

const MADL_EE = 0
const MADL_MID = 1

const EXECUTE_FOREVER = true

const REPLACE_COMPONENT = "REPLACE_COMPONENT"
const STOP              = "STOP"
const FDR_COMMAND       = "refines"