package ee

import (
	"newsolution/gmidarch/development/artefacts/madl"
	"newsolution/gmidarch/execution/engine"
	"newsolution/shared/parameters"
)

type EE struct {
	MADLTemp madl.MADL
}

func NewEE(madlTemp madl.MADL) EE {
	r := new(EE)

	r.MADLTemp = madlTemp

	return *r
}

func (e EE) Start() {

	for i := range e.MADLTemp.Components{
		go engine.Engine{}.Execute(e.MADLTemp.Components[i].Type, e.MADLTemp.Components[i].Graph, parameters.EXECUTE_FOREVER)
	}

	for i := range e.MADLTemp.Connectors{
		go engine.Engine{}.Execute(e.MADLTemp.Connectors[i].Type, e.MADLTemp.Connectors[i].Graph, parameters.EXECUTE_FOREVER)
	}
}