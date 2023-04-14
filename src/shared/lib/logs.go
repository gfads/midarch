package lib

import (
	"github.com/gfads/midarch/src/shared"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var SHOW_MESSAGES = []DebugLevel{} //ERROR, INFO, MESSAGE} //, DEBUG}

type DebugLevel int

const (
	ERROR   DebugLevel = 0
	INFO    DebugLevel = 1
	MESSAGE DebugLevel = 2
	DEBUG   DebugLevel = 3
)

func (d DebugLevel) ToInt() int {
	return [...]int{0, 1, 2, 3}[d]
}

func ConfigureDebugLevel() {
	debugLevel := shared.EnvironmentVariableValueWithDefault("DEBUG_LEVEL", "ERROR, INFO, MESSAGE")
	if strings.Contains(debugLevel, "ERROR") {
		SHOW_MESSAGES = append(SHOW_MESSAGES, ERROR)
	}
	if strings.Contains(debugLevel, "INFO") {
		SHOW_MESSAGES = append(SHOW_MESSAGES, INFO)
	}
	if strings.Contains(debugLevel, "MESSAGE") {
		SHOW_MESSAGES = append(SHOW_MESSAGES, MESSAGE)
	}
	if strings.Contains(debugLevel, "DEBUG") {
		SHOW_MESSAGES = append(SHOW_MESSAGES, DEBUG)
	}
}

func FunctionName() string {
	pc, _, _, _ := runtime.Caller(1)

	name := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	return name[len(name)-1]
}

func Println(messageLevel DebugLevel, message ...interface{}) {
	if len(SHOW_MESSAGES) > 0 {
		if inArrayDL(messageLevel, SHOW_MESSAGES) {
			_, file, line, ok := runtime.Caller(2)
			if !ok {
				file = "???"
				line = 0
			}

			switch messageLevel {
			case INFO:
				var logs []interface{}
				logs = append(logs, "- \"", file+":"+strconv.Itoa(line), "- INFO -")
				logs = append(logs, message...)
				logs = append(logs, "\"")
				log.Println(logs...)
			case DEBUG:
				var logs []interface{}
				logs = append(logs, file+":"+strconv.Itoa(line), "- DEBUG -", FunctionName())
				logs = append(logs, message...)
				log.Println(logs...)
			case MESSAGE:
				log.Println(message...)
			case ERROR:
				log.Println("- \"", file+":"+strconv.Itoa(line), "***** ERROR *****", message, "\"")
			}
		}
	}
}

func PrintlnInfo(message ...interface{}) {
	Println(INFO, message...)
}

func PrintlnDebug(message ...interface{}) {
	Println(DEBUG, message...)
}

func PrintlnMessage(message ...interface{}) {
	Println(MESSAGE, message...)
}

func PrintlnError(message ...interface{}) {
	Println(ERROR, message...)
}

func FailOnError(err error, msg string) {
	if err != nil {
		Println(ERROR, msg, ":", err)
		os.Exit(1)
	}
}

func InArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func inArrayDL(a DebugLevel, list []DebugLevel) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
