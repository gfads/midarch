package docker

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

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

func saveContainerLogsToFile(ctx context.Context, cli *client.Client, containerID string, kind Kind, fiboPlace int, sampleSize int) error {
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
		"log_" +
			kind.toString() + "_" +
			strconv.Itoa(fiboPlace) + "_" +
			strconv.Itoa(sampleSize) + "_" +
			time.Now().Format("20060102_150405")  + ".results.txt")
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
		size := binary.BigEndian.Uint32(header[4 : 8])
		payload := make([]byte, size)
		_, err = logs.Read(payload)
		if err != nil && err != io.EOF {
			return err
		}
		file.Write(payload)
		lines++
	}

	if lines < sampleSize {
		return errors.New("saveContainerLogsToFile: less line logs than sample size")
	}

	return nil
}
