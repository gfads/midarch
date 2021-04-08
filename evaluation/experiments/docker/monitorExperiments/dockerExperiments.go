package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

// According to: https://docs.docker.com/engine/api/v1.41/#operation/ContainerStats
// used_memory = memory_stats.usage - memory_stats.stats.cache
// available_memory = memory_stats.limit
// Memory usage % = (used_memory / available_memory) * 100.0
// cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
// system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
// number_cpus = lenght(cpu_stats.cpu_usage.percpu_usage) or cpu_stats.online_cpus
// CPU usage % = (cpu_delta / system_cpu_delta) * number_cpus * 100.0


type ContainerStats struct {
	Read      time.Time `json:"read"`
	//Preread   time.Time `json:"preread"`
	//PidsStats struct {
	//	Current int `json:"current"`
	//} `json:"pids_stats"`
	//BlkioStats struct {
	//	IoServiceBytesRecursive []interface{} `json:"io_service_bytes_recursive"`
	//	IoServicedRecursive     []interface{} `json:"io_serviced_recursive"`
	//	IoQueueRecursive        []interface{} `json:"io_queue_recursive"`
	//	IoServiceTimeRecursive  []interface{} `json:"io_service_time_recursive"`
	//	IoWaitTimeRecursive     []interface{} `json:"io_wait_time_recursive"`
	//	IoMergedRecursive       []interface{} `json:"io_merged_recursive"`
	//	IoTimeRecursive         []interface{} `json:"io_time_recursive"`
	//	SectorsRecursive        []interface{} `json:"sectors_recursive"`
	//} `json:"blkio_stats"`
	//NumProcs     int `json:"num_procs"`
	//StorageStats struct {
	//} `json:"storage_stats"`
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			PercpuUsage       []int `json:"percpu_usage"`
			//UsageInKernelmode int64 `json:"usage_in_kernelmode"`
			//UsageInUsermode   int   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int   `json:"online_cpus"`
		//ThrottlingData struct {
		//	Periods          int `json:"periods"`
		//	ThrottledPeriods int `json:"throttled_periods"`
		//	ThrottledTime    int `json:"throttled_time"`
		//} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CPUUsage struct {
			TotalUsage        int64 `json:"total_usage"`
			//PercpuUsage       []int `json:"percpu_usage"`
			//UsageInKernelmode int64 `json:"usage_in_kernelmode"`
			//UsageInUsermode   int   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int64 `json:"system_cpu_usage"`
		//OnlineCpus     int   `json:"online_cpus"`
		//ThrottlingData struct {
		//	Periods          int `json:"periods"`
		//	ThrottledPeriods int `json:"throttled_periods"`
		//	ThrottledTime    int `json:"throttled_time"`
		//} `json:"throttling_data"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage    int `json:"usage"`
		//MaxUsage int `json:"max_usage"`
		Stats    struct {
			//ActiveAnon              int   `json:"active_anon"`
			//ActiveFile              int   `json:"active_file"`
			Cache                   int   `json:"cache"`
			//Dirty                   int   `json:"dirty"`
			//HierarchicalMemoryLimit int   `json:"hierarchical_memory_limit"`
			//HierarchicalMemswLimit  int64 `json:"hierarchical_memsw_limit"`
			//InactiveAnon            int   `json:"inactive_anon"`
			//InactiveFile            int   `json:"inactive_file"`
			//MappedFile              int   `json:"mapped_file"`
			//Pgfault                 int   `json:"pgfault"`
			//Pgmajfault              int   `json:"pgmajfault"`
			//Pgpgin                  int   `json:"pgpgin"`
			//Pgpgout                 int   `json:"pgpgout"`
			//Rss                     int   `json:"rss"`
			//RssHuge                 int   `json:"rss_huge"`
			//TotalActiveAnon         int   `json:"total_active_anon"`
			//TotalActiveFile         int   `json:"total_active_file"`
			//TotalCache              int   `json:"total_cache"`
			//TotalDirty              int   `json:"total_dirty"`
			//TotalInactiveAnon       int   `json:"total_inactive_anon"`
			//TotalInactiveFile       int   `json:"total_inactive_file"`
			//TotalMappedFile         int   `json:"total_mapped_file"`
			//TotalPgfault            int   `json:"total_pgfault"`
			//TotalPgmajfault         int   `json:"total_pgmajfault"`
			//TotalPgpgin             int   `json:"total_pgpgin"`
			//TotalPgpgout            int   `json:"total_pgpgout"`
			//TotalRss                int   `json:"total_rss"`
			//TotalRssHuge            int   `json:"total_rss_huge"`
			//TotalUnevictable        int   `json:"total_unevictable"`
			//TotalWriteback          int   `json:"total_writeback"`
			//Unevictable             int   `json:"unevictable"`
			//Writeback               int   `json:"writeback"`
		} `json:"stats"`
		Limit int `json:"limit"`
	} `json:"memory_stats"`
	Name     string `json:"name"`
	ID       string `json:"id"`
	//Networks struct {
	//	Eth0 struct {
	//		RxBytes   int `json:"rx_bytes"`
	//		RxPackets int `json:"rx_packets"`
	//		RxErrors  int `json:"rx_errors"`
	//		RxDropped int `json:"rx_dropped"`
	//		TxBytes   int `json:"tx_bytes"`
	//		TxPackets int `json:"tx_packets"`
	//		TxErrors  int `json:"tx_errors"`
	//		TxDropped int `json:"tx_dropped"`
	//	} `json:"eth0"`
	//	Eth1 struct {
	//		RxBytes   int `json:"rx_bytes"`
	//		RxPackets int `json:"rx_packets"`
	//		RxErrors  int `json:"rx_errors"`
	//		RxDropped int `json:"rx_dropped"`
	//		TxBytes   int `json:"tx_bytes"`
	//		TxPackets int `json:"tx_packets"`
	//		TxErrors  int `json:"tx_errors"`
	//		TxDropped int `json:"tx_dropped"`
	//	} `json:"eth1"`
	//} `json:"networks"`
}







//type ContainerStats struct {
//	Id       string `json:"id"`
//	Read     string `json:"read"`
//	Preread  string `json:"preread"`
//	MemoryStats MemoryStats `json:"memory_stats"`
//	CpuStats cpu `json:"cpu_stats"`
//}
//
//type MemoryStats struct {
//	Usage float64 `json:"cpu_usage"`
//	Stats
//}
//
//type cpu struct {
//	Usage cpuUsage `json:"cpu_usage"`
//}
//
//type cpuUsage struct {
//	Total float64 `json:"total_usage"`
//}


//func main() {
//	ctx := context.Background()
//	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
//	if err != nil {
//		panic(err)
//	}
//
//	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
//	if err != nil {
//		panic(err)
//	}
//
//	for _, container := range containers {
//		fmt.Println(container.ID, container.Image, container.Names, container.Ports, container.Status, "State:", container.State)
//		cli.ContainerStats(ctx, container.ID, true)
//	}
//
//}



func main() {
	kind := "udp"
	fiboPlace, sampleSize := 38, 10000

	runExperiment(kind, fiboPlace, sampleSize)
}

func runExperiment(kind string, fiboPlace int, sampleSize int) {
	err := createStack(kind)
	if err != nil {
		panic(err)
	}
	defer removeStack(kind)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Minute)
	//cli, e := client.NewEnvClient()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	container := getClientContainer(ctx, cli)
	if container == nil {
		panic("No client container detected!")
	}

	stats, err := cli.ContainerStats(ctx, container.ID, true)
	if err != nil {
		fmt.Errorf("%s", err.Error())
	}
	decoder := json.NewDecoder(stats.Body)
	var containerStats ContainerStats

	filename := filepath.Join("evaluation",
		"results",
		"docker",
		"log_"+
			kind+"_"+
			strconv.Itoa(fiboPlace)+"_"+
			strconv.Itoa(sampleSize)+"_"+
			time.Now().Format("20060102_150405")+".monitor.csv")
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("container_name;container_status;used_memory(MB);available_memory(MB);memory_usage(%);cpu_delta;system_cpu_delta;number_cpus;cpu_usage(%);total_cpu_usage;pre_total_cpu_usage\n")
OuterLoop:
	for {
		select {
		case <-ctx.Done():
			stats.Body.Close()
			fmt.Println("Stop logging")
			return
		default:
			if err = decoder.Decode(&containerStats); err == io.EOF {
				return
			} else if err != nil {
				cancel()
			}
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
				saveContainerLogsToFile(ctx, cli, container.ID, kind, fiboPlace, sampleSize)
				break OuterLoop
			}
		}
	}
}

func createStack(kind string) error {
	command := ""
	switch kind {
	case "udp": command = "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-udp.yml fibomiddleware-udp"
	}
	return executeCommand(command)
}

func removeStack(kind string) error {
	command := ""
	switch kind {
	case "udp": command = "docker stack rm fibomiddleware-udp"
	}
	return executeCommand(command)
}

func saveLogs(kind string) error {
	command := ""
	switch kind {
	case "udp": command = "docker logs fibormq_client.1.zmxmmz1530uv3287lvb0u7c17 >& evaluation/results/docker/log_E_RMQ_38_10000_$(date +'%Y%m%d_%H%M%S').txt"
	}
	return executeCommand(command)
}

func executeCommand(command string) error {
	log.Printf("Running command and waiting for it to finish...")
	out, err := exec.Command("cmd", "/C", command).Output()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	log.Printf("Logs from command execution: %s", out)
	return err
}

func getClientContainer(ctx context.Context, cli *client.Client) *types.Container {
	for i := 0; i < 20; i++ {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, containerItem := range containers {
			fmt.Println(containerItem.ID, containerItem.Image, containerItem.Names, containerItem.Ports, containerItem.Status, "State:", containerItem.State)
			if strings.Contains(containerItem.Image, "client") {
				return &containerItem
			}
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func getContainerStatus(ctx context.Context, cli *client.Client, containerID string) string {
	filter := filters.NewArgs()
	filter.Add("id", containerID)

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		//panic(err)
		return "no container"
	}

	for _, containerItem := range containers {
//			fmt.Println(containerItem.ID, containerItem.Image, containerItem.Names, containerItem.Ports, containerItem.Status, "State:", containerItem.State)
		return containerItem.State
	}

	return "no container"
}

func saveContainerLogsToFile(ctx context.Context, cli *client.Client, containerID string, kind string, fiboPlace int, sampleSize int) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		//Details: true,
		//Timestamps: true,
		Tail: "11000",
	}

	logs, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		panic(err)
	}

	filename := filepath.Join("evaluation",
		"results",
		"docker",
		"log_" +
		kind + "_" +
		strconv.Itoa(fiboPlace) + "_" +
		strconv.Itoa(sampleSize) + "_" +
		time.Now().Format("20060102_150405")  + ".results.txt")
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	header := make([]byte, 8)
	for {
		_, err = logs.Read(header)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		size := binary.BigEndian.Uint32(header[4 : 8])
		payload := make([]byte, size)
		_, err = logs.Read(payload)
		if err != nil && err != io.EOF {
			panic(err)
		}
		file.Write(payload)
	}

	//_, err = io.Copy(file, logs)
	//if err != nil && err != io.EOF {
	//	log.Fatal(err)
	//}
}