package main

import (
	"fmt"
	"gmidarch/development/components/proxies/fibonacciProxy"
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/messages"
	"gmidarch/execution/frontend"
	"os"
	"shared"
	"strconv"
	"time"
)

func main() {
	// Wait for namingserver and server to get up
	time.Sleep(15 * time.Second)

	var n, SAMPLE_SIZE, AVERAGE_WAITING_TIME int
	if len(os.Args) >= 2 {
		n, _ = strconv.Atoi(os.Args[1])
		SAMPLE_SIZE, _ = strconv.Atoi(os.Args[2])
		AVERAGE_WAITING_TIME = 60
	}else{
		n, _ = strconv.Atoi(shared.EnvironmentVariableValue("FIBONACCI_PLACE"))
		SAMPLE_SIZE, _ = strconv.Atoi(shared.EnvironmentVariableValue("SAMPLE_SIZE"))
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(shared.EnvironmentVariableValue("AVERAGE_WAITING_TIME"))
	}
	fmt.Println("FIBONACCI_PLACE / SAMPLE_SIZE / AVERAGE_WAITING_TIME:", n, "/", SAMPLE_SIZE, "/", AVERAGE_WAITING_TIME)


	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: "namingserver", Port: shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy("FibonacciDistributedClientMid.madl", args)

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	aux, ok := namingProxy.Lookup("Fibonacci")
	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "Service 'Fibonacci' not found in Naming Service")
	}

	fibonacci := aux.(*fibonacciProxy.FibonacciProxy)
	for x := 0; x < SAMPLE_SIZE; x++ {
		fmt.Println("Result:", fibonacci.F(n))
		time.Sleep(200 * time.Millisecond)
	}

	//fmt.Scanln()
	//var wg sync.WaitGroup
	//wg.Add(1)
	//wg.Wait()
}
