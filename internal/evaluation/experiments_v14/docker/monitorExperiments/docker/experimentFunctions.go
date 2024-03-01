package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

func createStack(composePath string) error {
	return executeCommand("docker stack deploy -c " + composePath + " midarch")
}

// deprecated
func createStackFromKind(kind TransportProtocolFactor) error {
	return executeCommand(kind.createStackCommand())
}

func removeStack() error {
	return executeCommand("docker stack rm midarch")
}

func removeStackFromKind(kind TransportProtocolFactor) error {
	return executeCommand(kind.removeStackCommand())
}

func executeCommand(command string) error {
	log.Printf("Running command and waiting for it to finish...")
	//out, err := exec.Command("cmd", "/C", command).Output() // Windows
	log.Printf("Command: %s", command)
	commands := strings.Split(command, " ")
	out, err := exec.Command(commands[0], commands[1:]...).Output() // Linux
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	log.Printf("Logs from command execution: %s", out)
	return err
}

func buildExperiment(experimentDescription string, transportProtocol TransportProtocolFactor, adaptationInterval int, remoteOperation RemoteOperationFactor, sampleSize int, fiboPlace int, imageSize string) (outputPath string) {
	specificEnvClient := ""
	baseAppGoFilePath := "$GMIDARCHDIR/examples"
	switch remoteOperation {
	case Fibonacci:
		specificEnvClient = "FIBONACCI_PLACE: \"" + strconv.Itoa(fiboPlace) + "\""
		baseAppGoFilePath += "/fibonaccidistributed"
	case SendFile:
		specificEnvClient = "FILE_SIZE: \"" + imageSize + "\""
		baseAppGoFilePath += "/sendfiledistributed"
	}
	sampleSizeStr := strconv.Itoa(sampleSize)
	avarageWaitingTime := "200"
	baseImageName := "midarch/" + strings.ToLower(experimentDescription) + "-<app>"
	imageNameNamingServer := strings.Replace(baseImageName, "<app>", "naming", 1)
	imageNameServer := strings.Replace(baseImageName, "<app>", "server", 1)
	imageNameClient := strings.Replace(baseImageName, "<app>", "client", 1)

	outputPath = shared.DIR_EXPERIMENTS_RESULTS + "/" + strings.ReplaceAll(experimentDescription, ":", "-")
	if remoteOperation == Fibonacci {
		outputPath += "-" + strconv.Itoa(fiboPlace)
	} else {
		outputPath += "-" + imageSize
	}
	log.Println("Building in:", outputPath)
	os.Mkdir(outputPath, os.ModePerm)

	generateServerFile(transportProtocol, remoteOperation)
	generateDockerfile("naming", "$GMIDARCHDIR/naming", baseAppGoFilePath+"/naming/naming.go", outputPath)
	generateDockerfile("server", "$GMIDARCHDIR/server", baseAppGoFilePath+"/server/server.go", outputPath)
	generateDockerfile("client", "$GMIDARCHDIR/client", baseAppGoFilePath+"/client/client.go", outputPath)
	generateComposeFile(imageNameNamingServer, imageNameServer, imageNameClient, adaptationInterval, specificEnvClient, sampleSizeStr, avarageWaitingTime, outputPath)

	buildImages(shared.DIR_BASE, outputPath, remoteOperation, imageNameNamingServer, imageNameServer, imageNameClient, transportProtocol)

	publishImages(imageNameNamingServer, imageNameServer, imageNameClient)

	return outputPath
}

func publishImages(imageNameNamingServer string, imageNameServer string, imageNameClient string) {
	log.Println("Publishing images")
	log.Println()
	err := executeCommand("docker push " + imageNameNamingServer)
	if err != nil {
		log.Println("Error while publishing image", imageNameNamingServer, "Error:", err)
	}
	err = executeCommand("docker push " + imageNameServer)
	if err != nil {
		log.Println("Error while publishing image", imageNameServer, "Error:", err)
	}
	err = executeCommand("docker push " + imageNameClient)
	if err != nil {
		log.Println("Error while publishing image", imageNameClient, "Error:", err)
	}
}

func buildImages(contextPath string, outputPath string, remoteOperation RemoteOperationFactor, imageNameNamingServer string, imageNameServer string, imageNameClient string, transportProtocolFactor TransportProtocolFactor) {
	log.Println("Building images")
	log.Println()
	var protocolFactor string
	if transportProtocolFactor.IsEvolutive() {
		firstProtocolFactor, _ := transportProtocolFactor.getEvolutiveProtocols()
		protocolFactor = strings.ToUpper(firstProtocolFactor.toString())
	} else {
		protocolFactor = strings.ToUpper(transportProtocolFactor.toString())
	}
	err := shared.GenerateFromModel(shared.DIR_EXPERIMENTS_MODELS+"/naming.model.madl", shared.DIR_MADL+"/naming.madl", map[string]string{"<protocol>": protocolFactor})
	if err != nil {
		log.Println("Error while building image", imageNameNamingServer, "Error:", err)
	}
	err = executeCommand("docker build -t " + imageNameNamingServer + " " + contextPath + " -f " + outputPath + "/Dockerfile.naming")
	if err != nil {
		log.Println("Error while building image", imageNameNamingServer, "Error:", err)
	}

	if remoteOperation == Fibonacci {
		err = shared.GenerateFromModel(shared.DIR_EXPERIMENTS_MODELS+"/FibonacciDistributedServerMid.model.madl", shared.DIR_MADL+"/FibonacciDistributedServerMid.madl", map[string]string{"<protocol>": protocolFactor})
	} else {
		err = shared.GenerateFromModel(shared.DIR_EXPERIMENTS_MODELS+"/SendFileDistributedServerMid.model.madl", shared.DIR_MADL+"/SendFileDistributedServerMid.madl", map[string]string{"<protocol>": protocolFactor})
	}
	if err != nil {
		log.Println("Error while building image", imageNameServer, "Error:", err)
	}
	err = executeCommand("docker build -t " + imageNameServer + " " + contextPath + " -f " + outputPath + "/Dockerfile.server")
	if err != nil {
		log.Println("Error while building image", imageNameServer, "Error:", err)
	}

	if remoteOperation == Fibonacci {
		err = shared.GenerateFromModel(shared.DIR_EXPERIMENTS_MODELS+"/FibonacciDistributedClientMid.model.madl", shared.DIR_MADL+"/FibonacciDistributedClientMid.madl", map[string]string{"<protocol>": protocolFactor})
	} else {
		err = shared.GenerateFromModel(shared.DIR_EXPERIMENTS_MODELS+"/SendFileDistributedClientMid.model.madl", shared.DIR_MADL+"/SendFileDistributedClientMid.madl", map[string]string{"<protocol>": protocolFactor})
	}
	if err != nil {
		log.Println("Error while building image", imageNameClient, "Error:", err)
	}
	err = executeCommand("docker build -t " + imageNameClient + " " + contextPath + " -f " + outputPath + "/Dockerfile.client")
	if err != nil {
		log.Println("Error while building image", imageNameClient, "Error:", err)
	}
}

func generateServerFile(transportProtocol TransportProtocolFactor, remoteOperation RemoteOperationFactor) {
	serverFileModel := shared.DIR_EXPERIMENTS_MODELS
	serverFilePath := shared.DIR_BASE + "/examples"
	if remoteOperation == Fibonacci {
		serverFileModel += "/FibonacciServer.model.go"
		serverFilePath += "/fibonaccidistributed/server/server.go"
	} else {
		serverFileModel += "/SendFileServer.model.go"
		serverFilePath += "/sendfiledistributed/server/server.go"
	}
	log.Println("Vai ler:", serverFileModel)
	input, err := os.ReadFile(serverFileModel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var output []byte
	if transportProtocol.IsEvolutive() {
		output = bytes.Replace(input, []byte("<evolutive.import>"), []byte("evolutive \"github.com/gfads/midarch/pkg/injector\""), -1)
		output = bytes.Replace(output, []byte("<interval.between.injections>"), []byte("intervalBetweenInjections, _ := strconv.Atoi(shared.EnvironmentVariableValueWithDefault(\"INJECTION_INTERVAL\", \"120\"))"), -1)
		evolutiveInjectionCode := []byte("evolutive.EvolutiveInjector{}.StartEvolutiveProtocolInjection(\"<protocol1>\", \"<protocol2>\", time.Duration(intervalBetweenInjections)*time.Second)")
		protocol2, protocol1 := transportProtocol.getEvolutiveProtocols()
		evolutiveInjectionCode = bytes.Replace(evolutiveInjectionCode, []byte("<protocol1>"), []byte(protocol1.getComponentName()), -1)
		evolutiveInjectionCode = bytes.Replace(evolutiveInjectionCode, []byte("<protocol2>"), []byte(protocol2.getComponentName()), -1)
		output = bytes.Replace(output, []byte("<evolutive.injection>"), evolutiveInjectionCode, -1)
	} else {
		output = bytes.Replace(input, []byte("<evolutive.import>"), []byte(""), -1)
		output = bytes.Replace(output, []byte("<interval.between.injections>"), []byte(""), -1)
		output = bytes.Replace(output, []byte("<evolutive.injection>"), []byte(""), -1)
	}

	log.Println("Vai escrever:", serverFilePath)
	if err = os.WriteFile(serverFilePath, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateDockerfile(app string, appFilePath string, goFilePath string, outputPath string) {
	dockerfileModel := shared.DIR_EXPERIMENTS_MODELS + "/Dockerfile.model"
	log.Println("Vai ler:", dockerfileModel)
	input, err := os.ReadFile(dockerfileModel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	output := bytes.Replace(input, []byte("<gofile.path>"), []byte(goFilePath), -1)
	// output = bytes.Replace(output, []byte("<app.path>"), []byte(appFilePath), -1)

	if err = os.WriteFile(outputPath+"/Dockerfile."+app, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateComposeFile(imageNameNamingServer string, imageNameServer string, imageNameClient string, adaptationInterval int, specificEnvClient string, sampleSizeStr string, avarageWaitingTime string, outputPath string) {
	dcModel := shared.DIR_EXPERIMENTS_MODELS + "/dc-experiment.model.yml"
	log.Println("Vai ler:", dcModel)
	input, err := os.ReadFile(dcModel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	output := bytes.Replace(input, []byte("<image.naming>"), []byte(imageNameNamingServer), -1)
	output = bytes.Replace(output, []byte("<image.server>"), []byte(imageNameServer), -1)
	output = bytes.Replace(output, []byte("<image.client>"), []byte(imageNameClient), -1)
	output = bytes.Replace(output, []byte("<adaptation.interval>"), []byte(strconv.Itoa(adaptationInterval)), -1)
	output = bytes.Replace(output, []byte("<specific.env.client>"), []byte(specificEnvClient), -1)
	output = bytes.Replace(output, []byte("<sample.size>"), []byte(sampleSizeStr), -1)
	output = bytes.Replace(output, []byte("<average.waiting.time>"), []byte(avarageWaitingTime), -1)

	if err = os.WriteFile(outputPath+"/dc-experiment.yml", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunFibonacciExperiment(transportProtocol TransportProtocolFactor, adaptationInterval int, fiboPlace int, sampleSize int) {
	RunExperiment(transportProtocol, adaptationInterval, Fibonacci, sampleSize, fiboPlace, "")
}

func RunSendFileExperiment(transportProtocol TransportProtocolFactor, adaptationInterval int, imageSize string, sampleSize int) {
	RunExperiment(transportProtocol, adaptationInterval, SendFile, sampleSize, 0, imageSize)
}

func RunExperiment(transportProtocolFactor TransportProtocolFactor, adaptationInterval int, remoteOperation RemoteOperationFactor, sampleSize int, fiboPlace int, imageSize string) {
	experimentVersion := "1.14.0"
	experimentLevel := ""
	if remoteOperation == Fibonacci {
		experimentLevel = strconv.Itoa(fiboPlace)
	} else {
		experimentLevel = imageSize
	}
	var adaptationIntervalFactor string
	if transportProtocolFactor.IsEvolutive() {
		adaptationIntervalFactor = "on" + strconv.Itoa(adaptationInterval) + "s"
	} else {
		adaptationIntervalFactor = "off"
	}

	baseExperimentDescription := "<remote.operation>:<version>-<transport.protocol>-<adaptation.interval>-<sample.size>"
	baseExperimentDescription = strings.Replace(baseExperimentDescription, "<remote.operation>", remoteOperation.toString(), 1)
	baseExperimentDescription = strings.Replace(baseExperimentDescription, "<version>", experimentVersion, 1)
	baseExperimentDescription = strings.Replace(baseExperimentDescription, "<transport.protocol>", transportProtocolFactor.toString(), 1)
	baseExperimentDescription = strings.Replace(baseExperimentDescription, "<adaptation.interval>", adaptationIntervalFactor, 1)
	baseExperimentDescription = strings.Replace(baseExperimentDescription, "<sample.size>", strconv.Itoa(sampleSize), 1)
	experimentDescription := baseExperimentDescription + "-" + experimentLevel

	log.Println("Preparing experiment", experimentDescription)
	log.Println("Building experiment")
	log.Println()
	ouputPath := buildExperiment(baseExperimentDescription, transportProtocolFactor, adaptationInterval, remoteOperation, sampleSize, fiboPlace, imageSize)
	//return

	for {
		log.Println("Preparing to run", experimentDescription)
		log.Println()
		err := processExperiment(experimentDescription, ouputPath, transportProtocolFactor, fiboPlace, sampleSize)
		if err != nil {
			log.Println()
			log.Println("Error while processing experiment. Will try again in 10 seconds. Error:", err)
			log.Println()
			log.Println()
			time.Sleep(10 * time.Second)
		} else {
			log.Println()
			log.Println()
			log.Println("Finished running", experimentDescription)
			log.Println("Waiting 10 seconds to exit")
			time.Sleep(10 * time.Second)
			return
		}
	}
}

func processExperiment(experimentDescription string, outputPath string, kind TransportProtocolFactor, fiboPlace int, sampleSize int) []error {
	var experimentErrors []error
	stackRemoved := false

	err := createStack(outputPath + "/dc-experiment.yml")
	if err != nil {
		experimentErrors = append(experimentErrors, err)
		return experimentErrors
	}
	defer func() {
		if !stackRemoved {
			stackRemoved = true
			removeStack()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		experimentErrors = append(experimentErrors, err)
		return experimentErrors
	}

	containers := getClientContainer(ctx, cli)
	if containers == nil || len(containers) <= 1 {
		experimentErrors = append(experimentErrors, errors.New("not all containers detected"))
		return experimentErrors
	}

	lib.PrintlnDebug("containers['server'].ID:", containers["server"].ID)
	lib.PrintlnDebug("containers['server'].Names:", containers["server"].Names)
	lib.PrintlnDebug("containers['server'].Image:", containers["server"].Image)
	lib.PrintlnDebug("containers['client'].ID:", containers["client"].ID)
	lib.PrintlnDebug("containers['client'].Names:", containers["client"].Names)
	lib.PrintlnDebug("containers['client'].Image:", containers["client"].Image)

	wg := sync.WaitGroup{}
	wg.Add(2) // Wait for any of the experiments to finish
	go func() {
		err = monitorExperiment(experimentDescription, outputPath, kind, fiboPlace, sampleSize, cli, ctx, cancel, containers["client"].ID, "client")
		if err != nil {
			experimentErrors = append(experimentErrors, err)
		}
		wg.Done()
		if !stackRemoved {
			stackRemoved = true
			removeStack()
		}
	}()
	go func() {
		err = monitorExperiment(experimentDescription, outputPath, kind, fiboPlace, sampleSize, cli, ctx, cancel, containers["server"].ID, "server")
		if err != nil {
			experimentErrors = append(experimentErrors, err)
		}
		wg.Done()
		if !stackRemoved {
			stackRemoved = true
			removeStack()
		}
	}()
	wg.Wait()
	return experimentErrors
}

func monitorExperiment(experimentDescription string, outputPath string, kind TransportProtocolFactor, fiboPlace int, sampleSize int, cli *client.Client, ctx context.Context, cancel context.CancelFunc, containerId string, containerType string) error {
	stats, err := cli.ContainerStats(ctx, containerId, true)
	if err != nil {
		//_ = fmt.Errorf("%s", err.Error())
		return err
	}
	decoder := json.NewDecoder(stats.Body)

	var containerStats ContainerStats

	filename := filepath.Join(outputPath,
		"log_"+
			strings.Replace(experimentDescription, ":", "-", 1)+"-"+
			containerType+"_"+
			time.Now().Format("20060102_150405")+".monitor.csv")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("dateTime;container_name;container_status;used_memory(MB);available_memory(MB);memory_usage(%);cpu_delta;system_cpu_delta;number_cpus;cpu_usage(%);total_cpu_usage;pre_total_cpu_usage\n")
	//start := time.Now()
OuterLoop:
	for {
		select {
		case <-ctx.Done():
			stats.Body.Close()
			fmt.Println("Stop logging")
			return errors.New("timeout while processing experiment")
		default:
			if err = decoder.Decode(&containerStats); err == io.EOF {
				return errors.New("eof while decoding container stats")
			} else if err != nil {
				cancel()
			}
			t := time.Now()
			//duration := t.Sub(start)
			//fmt.Println(containerStats.CPUStats.CPUUsage.TotalUsage)
			usedMemory := containerStats.MemoryStats.Usage - containerStats.MemoryStats.Stats.Cache                  // used_memory = memory_stats.usage - memory_stats.stats.cache
			availableMemory := containerStats.MemoryStats.Limit                                                      // available_memory = memory_stats.limit
			percMemoryUsage := (float64(usedMemory) / float64(availableMemory)) * 100                                // Memory usage % = (used_memory / available_memory) * 100.0
			cpuDelta := containerStats.CPUStats.CPUUsage.TotalUsage - containerStats.PrecpuStats.CPUUsage.TotalUsage // cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
			systemCpuDelta := containerStats.CPUStats.SystemCPUUsage - containerStats.PrecpuStats.SystemCPUUsage     // system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
			numberCpus := containerStats.CPUStats.OnlineCpus                                                         // number_cpus = length(cpu_stats.cpu_usage.percpu_usage) or cpu_stats.online_cpus
			cpuUsage := float64(0)
			if systemCpuDelta > 0 {
				cpuUsage = (float64(cpuDelta) / float64(systemCpuDelta)) * float64(numberCpus) * 100 // CPU usage % = (cpu_delta / system_cpu_delta) * number_cpus * 100.0
			}
			containerStatus := getContainerStatus(ctx, cli, containerId)

			//fmt.Println(
			file.WriteString(
				t.Format("2006/01/02 15:04:05.999") + ";" +
					//duration.String() + ";" +
					//fmt.Sprintf("%f", duration.Seconds()) + ";" +
					containerStats.Name + ";" +
					containerStatus + ";" +
					fmt.Sprintf("%f", float64(usedMemory)/1024/1024) + ";" +
					fmt.Sprintf("%f", float64(availableMemory)/1024/1024) + ";" +
					fmt.Sprintf("%f", percMemoryUsage) + ";" +
					strconv.FormatInt(cpuDelta, 10) + ";" +
					strconv.FormatInt(systemCpuDelta, 10) + ";" +
					strconv.Itoa(numberCpus) + ";" +
					fmt.Sprintf("%f", cpuUsage) + ";" +

					//containerStats.CPUStats.SystemCPUUsage + ";" +
					//containerStats.PrecpuStats.SystemCPUUsage,

					strconv.FormatInt(containerStats.CPUStats.CPUUsage.TotalUsage, 10) + ";" +
					strconv.FormatInt(containerStats.PrecpuStats.CPUUsage.TotalUsage, 10) + "\n",
			)

			if containerStatus == "no container" {
				break OuterLoop
			}
		}
	}
	err = saveContainerLogsToFile(experimentDescription, outputPath, ctx, cli, containerId, containerType, kind, fiboPlace, sampleSize)
	return err
}
