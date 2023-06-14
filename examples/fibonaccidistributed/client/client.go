package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gfads/midarch/examples/fibonaccidistributed/fibonacciProxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/proxies/namingproxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

func main() {
	// Wait for namingserver and server to get up
	time.Sleep(13 * time.Second)

	// Example setting environment variable MIDARCH_BUSINESS_COMPONENTS_PATH on code, may be set on system environment variables too
	os.Setenv("MIDARCH_BUSINESS_COMPONENTS_PATH",
		shared.DIR_BASE+"/examples/fibonaccidistributed/fibonacciProxy")

	var n, SAMPLE_SIZE, AVERAGE_WAITING_TIME int
	if len(os.Args) >= 2 {
		n, _ = strconv.Atoi(os.Args[1])
		SAMPLE_SIZE, _ = strconv.Atoi(os.Args[2])
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(os.Args[3])
	} else {
		n, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("FIBONACCI_PLACE", "11"))
		SAMPLE_SIZE, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("SAMPLE_SIZE", "100"))
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("AVERAGE_WAITING_TIME", "60"))
	}
	fmt.Println("FIBONACCI_PLACE / SAMPLE_SIZE / AVERAGE_WAITING_TIME:", n, "/", SAMPLE_SIZE, "/", AVERAGE_WAITING_TIME)

	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}

	// Deploy configuration
	fe.Deploy(frontend.DeployOptions{
		FileName: "FibonacciDistributedClientMid.madl",
		Args:     args,
		Components: map[string]interface{}{
			"FibonacciProxy": &fibonacciProxy.FibonacciProxy{},
		}})

	// proxy to naming service
	endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	namingProxy := namingproxy.NewNamingproxy(endPoint)

	aux, ok := namingProxy.Lookup("Fibonacci")
	if !ok {
		shared.ErrorHandler(shared.GetFunction(), "Service 'Fibonacci' not found in Naming Service")
	}

	fibonacci := aux.(*fibonacciProxy.FibonacciProxy)

	rand.Seed(time.Now().UnixNano())
	for x := 0; x < SAMPLE_SIZE; x++ {
		ok := false
		for !ok {
			t1 := time.Now()
			//fmt.Println("Result:", fibonacci.F(n))
			r := fibonacci.F(n)
			//time.Sleep(200 * time.Millisecond)

			t2 := time.Now()

			duration := t2.Sub(t1)
			if r != 0 {
				ok = true
				lib.PrintlnMessage(x+1, float64(duration.Nanoseconds())/1000000)
			}

			// Normally distributed waiting time between calls with an average of 60 milliseconds and standard deviation of 20 milliseconds
			var rd = int(math.Round((rand.NormFloat64() * float64(AVERAGE_WAITING_TIME/5)) + float64(AVERAGE_WAITING_TIME)))
			if rd > 0 {
				time.Sleep(time.Duration(rd) * time.Millisecond)
			}
		}
	}

	//fmt.Scanln()
	//var wg sync.WaitGroup
	//wg.Add(1)
	//wg.Wait()
}
