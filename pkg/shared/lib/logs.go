package lib

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/gfads/midarch/pkg/shared"
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

func GetURIParameters(uri string) (parameters map[string]interface{}) {
	decodedPathParam, err := url.PathUnescape(uri[8:])
	if err != nil {
		PrintlnInfo("Error decoding path parameter:", err)
		return
	}
	return map[string]interface{}{"param": decodedPathParam}

	// parameters = make(map[string]interface{})
	// paramRegex, err := regexp.Compile("([?][\\w|=|%|(|)|+|-|.|:]+)|([&][\\w|=|%|(|)|+|-|.|:]+)")
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// params := paramRegex.FindAllString(uri, -1)
	// for _, param := range params {
	// 	param := param[1:]
	// 	keyValue := strings.Split(param, "=")
	// 	parameters[keyValue[0]] = keyValue[1]
	// }
	// return parameters
}

func GetServerTLSConfig(proto string) *tls.Config {
	if shared.CRT_PATH == "" {
		log.Fatal("TLS:: Error:: Environment variable 'CRT_PATH' not configured\n")
	}

	if shared.KEY_PATH == "" {
		log.Fatal("TLS:: Error:: Environment variable 'KEY_PATH' not configured\n")
	}

	cert, err := tls.LoadX509KeyPair(shared.CRT_PATH, shared.KEY_PATH)
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{proto},
	}
	return tlsConfig
}

func GetClientTLSConfig(proto string) *tls.Config {
	if shared.CA_PATH == "" {
		log.Fatal("CRHSsl:: Error:: Environment variable 'CA_PATH' not configured\n")
	}
	trustCert, err := ioutil.ReadFile(shared.CA_PATH)
	if err != nil {
		fmt.Println("Error loading trust certificate. ", err)
	}
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(trustCert) {
		fmt.Println("Error installing trust certificate.")
	}

	// connect to server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            certs,
		NextProtos:         []string{proto},
	}
	return tlsConfig
}
