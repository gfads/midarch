package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gfads/midarch/examples/sendfiledistributed/middleware"
	sendFileProxy "github.com/gfads/midarch/examples/sendfiledistributed/sendfileProxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/proxies/namingproxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	<evolutive.import>
	"github.com/gfads/midarch/pkg/shared"
)

func main() {
	// Wait for namingserver to get up
	timeToRun, _ := strconv.Atoi(shared.EnvironmentVariableValueWithDefault("TIME_TO_START_SERVER", "8"))
	fmt.Println("Waiting", timeToRun, "seconds for naming server to get up")
	time.Sleep(time.Duration(timeToRun) * time.Second)

	// Example setting environment variable MIDARCH_BUSINESS_COMPONENTS_PATH on code, may be set on system environment variables too
	os.Setenv("MIDARCH_BUSINESS_COMPONENTS_PATH",
		shared.DIR_BASE+"/examples/sendfiledistributed/middleware")

	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	args["srh"] = messages.EndPoint{Host: "0.0.0.0", Port: shared.CALCULATOR_PORT}

	// Deploy configuration
	fe.Deploy(frontend.DeployOptions{FileName: "SendFileDistributedServerMid.madl", Args: args, Components: map[string]interface{}{
		"SendFileInvoker": &middleware.SendFileInvoker{},
	}})

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	// Create proxy to Fibonacci
	sendFileProxy := sendFileProxy.NewSendFileProxy(generic.ProxyConfig{
		Host: shared.CALCULATOR_HOST,
		Port: shared.CALCULATOR_PORT,
	})

	// Register Fibonacci in Lookup
	ok := namingProxy.Register("SendFile", sendFileProxy)

	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "'SendFile' already registered in the Naming Server")
	}

	fmt.Printf("SendFile server is running at Port: %v \n", shared.CALCULATOR_PORT)

	<interval.between.injections>
	<evolutive.injection>
	// intervalBetweenInjections, _ := strconv.Atoi(shared.EnvironmentVariableValueWithDefault("INJECTION_INTERVAL", "20"))
	// evolutive.EvolutiveInjector{}.StartEvolutiveProtocolInjection("srhtcp", "srhrpc", time.Duration(intervalBetweenInjections)*time.Second)

	//fmt.Scanln()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
