package dot

import (
	"bufio"
	"github.com/gfads/midarch/src/shared"
	"os"
	"strconv"
	"strings"
)

type DOTCreator interface {
	Create(string) DOTGraph
}

type DOTLoaderImpl struct{}

func NewMADLLoader() DOTCreator {
	var mp DOTCreator

	mp = DOTLoaderImpl{}

	return mp
}

func (DOTLoaderImpl) Create(fileName string) DOTGraph {

	// Check DOT file name
	shared.CheckFileName(fileName, shared.DOT_EXTENSION)

	fullPathFileName := shared.DIR_DOT + "/" + fileName

	// Read DOT file
	fileContent := []string{}
	fileTemp, err := os.Open(fullPathFileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	defer fileTemp.Close()

	scanner := bufio.NewScanner(fileTemp)
	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	// Create DOT graph
	dotGraph := NewDOTGraph(shared.MAXIMUM_GRAPH_SIZE)
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
