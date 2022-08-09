package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func createStack(kind Kind) error {
	return executeCommand(kind.createStackCommand())
}

func removeStack(kind Kind) error {
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

func RunExperiment(kind Kind, fiboPlace int, sampleSize int) {
	for {
		log.Println("Preparing to run", kind.toString(), "(fiboPlace:",fiboPlace,"sampleSize:",sampleSize,") experiment fiboPlace!")
		log.Println()
		err := processExperiment(kind, fiboPlace, sampleSize)
		if err != nil {
			log.Println()
			log.Println("Error while processing experiment. Will try again. Error:", err)
			log.Println()
			log.Println()
			time.Sleep(10 * time.Second)
		}else {
			log.Println()
			log.Println()
			time.Sleep(10 * time.Second)
			return
		}
	}
}

func processExperiment(kind Kind, fiboPlace int, sampleSize int) error {
	err := createStack(kind)
	if err != nil {
		return err
	}
	defer removeStack(kind)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	container := getClientContainer(ctx, cli)
	if container == nil {
		return errors.New("no client container detected")
	}

	stats, err := cli.ContainerStats(ctx, container.ID, true)
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
			kind.toString()+"_"+
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
			containerStatus := getContainerStatus(ctx, cli, container.ID)

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
	err = saveContainerLogsToFile(ctx, cli, container.ID, kind, fiboPlace, sampleSize)
	return err
}