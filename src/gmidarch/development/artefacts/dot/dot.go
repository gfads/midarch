package dot

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"shared"
	"strconv"
	"strings"
)

type DOT struct{}

func (DOT) Read(file string) DOTGraph {

	// Check DOT file name
	err := checkfilename(file)
	if err != nil {
		fmt.Println("DOT:: " + err.Error())
		os.Exit(0)
	}

	fullPathFileName := shared.DIR_DOT + "/" + file

	// Read DOT file
	fileContent := []string{}
	fileTemp, err := os.Open(fullPathFileName)
	if err != nil {
		fmt.Println("DOT:: "+err.Error())
		os.Exit(0)
	}

	defer fileTemp.Close()

	scanner := bufio.NewScanner(fileTemp)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// Configure DOT digraph
	dotGraph := NewDOTGraph(shared.GRAPH_SIZE)

	for l := range fileContent {
		line := fileContent[l]
		if strings.Contains(line, "->") {
			from, _ := strconv.Atoi(strings.TrimSpace(line[:strings.Index(line, "->")]))
			to, _ := strconv.Atoi(strings.TrimSpace(line[strings.Index(line, "->")+2 : strings.Index(line, "[")]))
			label := line[strings.Index(line, "=")+2 : strings.LastIndex(line, "]")-2]
			dotGraph.AddEdge(from, to, label)
		}
	}
	return *dotGraph
}

func checkfilename(file string) error {
	r1 := *new(error)

	l := len(file)

	if l <= len(shared.DOT_EXTENSION) {
		r1 = errors.New("File Name Invalid")
	} else {
		if file[l-len(shared.DOT_EXTENSION):] != shared.DOT_EXTENSION {
			r1 = errors.New("Invalid extension of '" + file + "'")
		} else {
			r1 = nil
		}
	}
	return r1
}
