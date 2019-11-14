package madl

import (
	"gmidarch/development/artefacts/graphs"
)

type Element struct {
	ElemId    string            // from madl file
	TypeName  string            // from madl file
	Type      interface{}       // from repository
	Behaviour string            // from repository
	Info      []*interface{}    // particular to each element
	Graph     graphs.ExecGraph  // from dot file
	Params    [] interface{}    // from MADL - particular to each element, e.g., Connector OneToN (N)
}
