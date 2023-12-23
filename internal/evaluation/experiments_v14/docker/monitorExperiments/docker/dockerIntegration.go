package docker

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/gfads/midarch/pkg/shared/lib"
)

func getClientContainer(ctx context.Context, cli *client.Client) map[string]types.Container {
	filteredContainers := make(map[string]types.Container)
	clientOk := false
	serverOk := false
	for i := 0; i < 20; i++ {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}
		lib.PrintlnDebug("containers:", len(containers))
		for _, containerItem := range containers {
			//fmt.Println(containerItem.ID, containerItem.Image, containerItem.Names, containerItem.Ports, containerItem.Status, "State:", containerItem.State)
			if strings.Contains(containerItem.Image, "client") {
				lib.PrintlnDebug("Found client:", containerItem.Image)
				filteredContainers["client"] = containerItem
				clientOk = true
			} else if !strings.Contains(containerItem.Image, "naming") && strings.Contains(containerItem.Image, "server") {
				lib.PrintlnDebug("Found server:", containerItem.Image)
				filteredContainers["server"] = containerItem
				serverOk = true
			}
		}
		if clientOk && serverOk {
			return filteredContainers
		}
		time.Sleep(1 * time.Second)
	}
	lib.PrintlnDebug("filteredContainers:", len(filteredContainers))
	return filteredContainers
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

func saveContainerLogsToFile(ctx context.Context, cli *client.Client, containerID string, containerType string, kind TransportProtocolFactor, fiboPlace int, sampleSize int) error {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		//Details: true,
		//Timestamps: true,
		Tail: "11000",
	}

	logs, err := cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return err
	}

	filename := filepath.Join("evaluation",
		"results",
		"docker",
		"log_"+
			kind.toString()+"_"+containerType+"_"+
			strconv.Itoa(fiboPlace)+"_"+
			strconv.Itoa(sampleSize)+"_"+
			time.Now().Format("20060102_150405")+".results.txt")
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	header := make([]byte, 8)
	lines := 0
	for {
		_, err = logs.Read(header)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		size := binary.BigEndian.Uint32(header[4:8])
		payload := make([]byte, size)
		_, err = logs.Read(payload)
		if err != nil && err != io.EOF {
			return err
		}
		file.Write(payload)
		lines++
	}

	if containerType == "client" && lines < sampleSize {
		return errors.New("saveContainerLogsToFile: less line logs than sample size")
	}

	return nil
}
