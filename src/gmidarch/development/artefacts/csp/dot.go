package csp

import (
	"strings"
	"strconv"
	"newsolution/shared/parameters"
	"errors"
	"gmidarch/development/artefacts/graphs"
)

type DOT struct {
	SourceDotFile DOTFile
	Dotgraph      graphs.GraphDot
}

func (DOT) Create(csp CSP) (map[string]DOT, error) {
	r1 := map[string]DOT{}
	r2 := *new(error)

	// Dots
	dotFilesMid, err := DOTFile{}.LoadDotFiles(csp)
	if err != nil {
		r2 := errors.New("DOT" + err.Error())
		return r1, r2
	}
	for i := range dotFilesMid {
		dot, err := create(dotFilesMid[i])
		if err != nil {
			r2 = errors.New("DOT: " + err.Error())
			return r1, r2
		}
		r1[i] = dot
	}

	return r1, r2
}

func create(dotFile DOTFile) (DOT, error) {
	r1 := DOT{}
	r2 := *new(error)

	// Configure source dot
	r1.SourceDotFile = dotFile

	// Configure digraph
	dotGraph := graphs.NewGraphDot(parameters.GRAPH_SIZE)

	for l := range dotFile.Content {
		line := dotFile.Content[l]
		if strings.Contains(line, "->") {
			from, _ := strconv.Atoi(strings.TrimSpace(line[:strings.Index(line, "->")]))
			to, _ := strconv.Atoi(strings.TrimSpace(line[strings.Index(line, "->")+2 : strings.Index(line, "[")]))
			label := line[strings.Index(line, "=")+2 : strings.LastIndex(line, "]")-2]
			dotGraph.AddEdge(from, to, label)
		}
	}
	r1.Dotgraph = *dotGraph
	return r1, r2
}
