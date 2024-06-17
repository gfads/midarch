package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	sendFileProxy "github.com/gfads/midarch/examples/sendfiledistributed/sendfileProxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/frontend"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

func main() {
	// Start profiling
	// f, err := os.Create("client.prof")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	// Wait for namingserver and server to get up
	timeToRun, _ := strconv.Atoi(shared.EnvironmentVariableValueWithDefault("TIME_TO_START_CLIENT", "13"))
	lib.PrintlnDebug("Waiting", timeToRun, "seconds for naming server and server to get up")
	time.Sleep(time.Duration(timeToRun) * time.Second)

	// Example setting environment variable MIDARCH_BUSINESS_COMPONENTS_PATH on code, may be set on system environment variables too
	os.Setenv("MIDARCH_BUSINESS_COMPONENTS_PATH",
		shared.DIR_BASE+"/examples/sendfiledistributed/sendfileProxy")

	var FILE_SIZE string
	var SAMPLE_SIZE, AVERAGE_WAITING_TIME int
	if len(os.Args) >= 2 {
		FILE_SIZE = os.Args[1]
		SAMPLE_SIZE, _ = strconv.Atoi(os.Args[2])
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(os.Args[3])
	} else {
		FILE_SIZE = shared.EnvironmentVariableValueWithDefault("FILE_SIZE", "md")
		SAMPLE_SIZE, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("SAMPLE_SIZE", "100"))
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("AVERAGE_WAITING_TIME", "60"))
	}
	fmt.Println("dateTime;info;sequential;response_time") //"FILE_SIZE / SAMPLE_SIZE / AVERAGE_WAITING_TIME:", FILE_SIZE, "/", SAMPLE_SIZE, "/", AVERAGE_WAITING_TIME)

	fe := frontend.NewFrontend()

	// Configure port of SRHs/CRHs used in the configuration.
	// The order of Ip/hosts must the same as one in which
	// these elements appear in the configuration
	args := make(map[string]messages.EndPoint)
	args["crh"] = messages.EndPoint{Host: shared.CALCULATOR_HOST, Port: shared.CALCULATOR_PORT}

	// Deploy configuration
	fe.Deploy(frontend.DeployOptions{
		FileName: "SendFileDistributedClientMid.madl",
		Args:     args,
		Components: map[string]interface{}{
			"SendFileProxy": &sendFileProxy.SendFileProxy{},
		}})

	// proxy to naming service
	// endPoint := messages.EndPoint{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}
	// namingProxy := namingproxy.NewNamingproxy(endPoint)

	// aux, ok := namingProxy.Lookup("SendFile")
	// if !ok {
	// 	shared.ErrorHandler(shared.GetFunction(), "Service 'SendFile' not found in Naming Service")
	// }

	// sendFile := aux.(*sendFileProxy.SendFileProxy)

	sendFile := &sendFileProxy.SendFileProxy{}
	proxyConfig := generic.ProxyConfig{Host: shared.CALCULATOR_HOST, Port: shared.CALCULATOR_PORT}
	sendFile.Configure(proxyConfig)
	time.Sleep(2 * time.Second)

	fileBytes := getFile(FILE_SIZE)
	rand.Seed(time.Now().UnixNano())
	for x := 0; x < SAMPLE_SIZE; x++ {
		ok := false

		for !ok {
			// TODO dcruzb: getImage based on FILE_SIZE

			t1 := time.Now()

			r := sendFile.SendFile(fileBytes)
			//time.Sleep(200 * time.Millisecond)

			t2 := time.Now()

			duration := t2.Sub(t1)
			if r {
				ok = true
				// lib.PrintlnMessage(x+1, float64(duration.Nanoseconds())/1000000)
				log.Printf(";ok;%d;%f\n", x+1, float64(duration.Nanoseconds())/1000000)
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

func getFile(size string) []byte {
	var fileName string
	switch size {
	case "sm":
		fileName = shared.DIR_BASE + "/examples/sendfiledistributed/client/36x36.png"
	case "md":
		fileName = shared.DIR_BASE + "/examples/sendfiledistributed/client/2k.png"
	case "lg":
		fileName = shared.DIR_BASE + "/examples/sendfiledistributed/client/4k.jpg" // Foto de <a href="https://unsplash.com/pt-br/@francesco_ungaro?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash">Francesco Ungaro</a> na <a href="https://unsplash.com/pt-br/fotografias/cardume-de-peixes-no-corpo-de-agua-MJ1Q7hHeGlA?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash">Unsplash</a>
	}

	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Error opening file: "+fileName)
	}

	return fileBytes
}
