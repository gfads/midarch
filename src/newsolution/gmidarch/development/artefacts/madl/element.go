package madl

import "newsolution/gmidarch/development/artefacts/graphs"

type Element struct {
	ElemId       string           // madl file
	TypeName string           // madl file
	Type     interface{}      // repository
	Behaviour      string           // repository
	Info     interface{}      // TODO
	Graph    graphs.ExecGraph //
}
