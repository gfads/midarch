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

func createStack(kind TransportProtocolFactor) error {
	return executeCommand(kind.createStackCommand())
}

func removeStack(kind TransportProtocolFactor) error {
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

func buildExperiment(transportProtocol TransportProtocolFactor, remoteOperation RemoteOperationFactor, fiboPlace int, imageSize string, sampleSize int) {
	dcModel := shared.DOCKER_COMPOSE_PATH
	log.Println("Vai ler:", dcModel)
	input, err := os.ReadFile(dcModel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	specificEnvClient := ""
	switch remoteOperation {
	case Fibonacci:
		specificEnvClient = "FIBONACCI_PLACE: \"" + strconv.Itoa(fiboPlace) + "\""
	case SendFile:
		specificEnvClient = "IMAGE_SIZE: \"" + imageSize + "\""
	}

	sampleSizeStr := strconv.Itoa(sampleSize)
	avarageWaitingTime := "200"
	experimentVersion := "1.14.0"
	baseImageName := "midarch/<remote.operation>:<version>-<app>-<tranport.protocol>"
	imageName := strings.Replace(baseImageName, "<remote.operation>", remoteOperation.toString(), 1)
	imageName = strings.Replace(imageName, "<version>", experimentVersion, 1)
	imageName = strings.Replace(imageName, "<tranport.protocol>", transportProtocol.toString(), 1)
	imageNameNamingServer := strings.Replace(imageName, "<app>", "namingserver", 1)
	imageNameServer := strings.Replace(imageName, "<app>", "server", 1)
	imageNameClient := strings.Replace(imageName, "<app>", "client", 1)

	output := bytes.Replace(input, []byte("<image.namingserver>"), []byte(imageNameNamingServer), -1)
	output = bytes.Replace(output, []byte("<image.server>"), []byte(imageNameServer), -1)
	output = bytes.Replace(output, []byte("<image.client>"), []byte(imageNameClient), -1)
	output = bytes.Replace(output, []byte("<specific.env.client>"), []byte(specificEnvClient), -1)
	output = bytes.Replace(output, []byte("<sample.size>"), []byte(sampleSizeStr), -1)
	output = bytes.Replace(output, []byte("<average.waiting.time>"), []byte(avarageWaitingTime), -1)

	outputPath := shared.DIR_EXPERIMENTS_RESULTS + "/" + remoteOperation.toString() + "-" + experimentVersion + "-" + transportProtocol.toString() + "-" + strconv.Itoa(fiboPlace) + "-" + imageSize + "-" + sampleSizeStr
	log.Println("Save to:", outputPath)
	os.Mkdir(outputPath, os.ModePerm)
	if err = os.WriteFile(outputPath+"/dc-experiment.yml", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RunFibonacciExperiment(transportProtocol TransportProtocolFactor, fiboPlace int, sampleSize int) {
	RunExperiment(transportProtocol, Fibonacci, fiboPlace, "", sampleSize)
}

func RunSendFileExperiment(transportProtocol TransportProtocolFactor, imageSize string, sampleSize int) {
	RunExperiment(transportProtocol, SendFile, 0, imageSize, sampleSize)
}

func RunExperiment(transportProtocol TransportProtocolFactor, remoteOperation RemoteOperationFactor, fiboPlace int, imageSize string, sampleSize int) {
	log.Println("Preparing experiment", transportProtocol.toString(), "(fiboPlace:", fiboPlace, "sampleSize:", sampleSize, ") experiment fiboPlace!")
	log.Println("Building experiment")
	log.Println()
	buildExperiment(transportProtocol, remoteOperation, fiboPlace, imageSize, sampleSize)
	return

	for {
		log.Println("Preparing to run", transportProtocol.toString(), "(fiboPlace:", fiboPlace, "sampleSize:", sampleSize, ") experiment fiboPlace!")
		log.Println()
		err := processExperiment(transportProtocol, fiboPlace, sampleSize)
		if err != nil {
			log.Println()
			log.Println("Error while processing experiment. Will try again in 10 seconds. Error:", err)
			log.Println()
			log.Println()
			time.Sleep(10 * time.Second)
		} else {
			log.Println()
			log.Println()
			log.Println("Finished running", transportProtocol.toString(), "(fiboPlace:", fiboPlace, "sampleSize:", sampleSize, ") experiment fiboPlace!")
			log.Println("Waiting 10 seconds to exit")
			time.Sleep(10 * time.Second)
			return
		}
	}
}

func processExperiment(kind TransportProtocolFactor, fiboPlace int, sampleSize int) []error {
	var experimentErrors []error
	stackRemoved := false

	err := createStack(kind)
	if err != nil {
		experimentErrors = append(experimentErrors, err)
		return experimentErrors
	}
	defer func() {
		if !stackRemoved {
			stackRemoved = true
			removeStack(kind)
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
		err = monitorExperiment(kind, fiboPlace, sampleSize, cli, ctx, cancel, containers["client"].ID, "client")
		if err != nil {
			experimentErrors = append(experimentErrors, err)
		}
		wg.Done()
		if !stackRemoved {
			stackRemoved = true
			removeStack(kind)
		}
	}()
	go func() {
		err = monitorExperiment(kind, fiboPlace, sampleSize, cli, ctx, cancel, containers["server"].ID, "server")
		if err != nil {
			experimentErrors = append(experimentErrors, err)
		}
		wg.Done()
		if !stackRemoved {
			stackRemoved = true
			removeStack(kind)
		}
	}()
	wg.Wait()
	return experimentErrors
}

func monitorExperiment(kind TransportProtocolFactor, fiboPlace int, sampleSize int, cli *client.Client, ctx context.Context, cancel context.CancelFunc, containerId string, containerType string) error {
	stats, err := cli.ContainerStats(ctx, containerId, true)
	if err != nil {
		//_ = fmt.Errorf("%s", err.Error())
		return err
	}
	decoder := json.NewDecoder(stats.Body)

	var containerStats ContainerStats

	filename := filepath.Join("evaluation",
		"results",
		"docker",
		"log_"+
			kind.toString()+"_"+containerType+"_"+
			strconv.Itoa(fiboPlace)+"_"+
			strconv.Itoa(sampleSize)+"_"+
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
	err = saveContainerLogsToFile(ctx, cli, containerId, containerType, kind, fiboPlace, sampleSize)
	return err
}
