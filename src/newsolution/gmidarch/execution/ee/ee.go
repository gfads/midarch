package ee

import (
	"newsolution/gmidarch/development/artefacts/madl"
	"newsolution/gmidarch/development/components"
	"newsolution/gmidarch/execution/engine"
	"newsolution/shared/parameters"
)

type EE struct {
	MADLX madl.MADL
}

func NewEE() EE {
	r := new(EE)
	return *r
}

func (e EE) Start() {

	for i := range e.MADLX.Components {
		elem := e.MADLX.Components[i].Type
		graph := e.MADLX.Components[i].Graph

		// Configure Unit's Info with Element's Info (Only components)
		if e.MADLX.Components[i].TypeName == "Unit" { // TODO - Generalise for any component having 'Info'
			tempElem := *e.MADLX.Components[i].Info[0]
			unit := elem.(components.Unit)
			unit.UnitId = e.MADLX.Components[i].ElemId
			unit.ElemOfUnit= tempElem.(madl.Element).Type
			unit.GraphOfElem = tempElem.(madl.Element).Graph
			elem = unit
		}
		go engine.Engine{}.Execute(elem, graph, parameters.EXECUTE_FOREVER)
	}

	for i := range e.MADLX.Connectors {
		go engine.Engine{}.Execute(e.MADLX.Connectors[i].Type, e.MADLX.Connectors[i].Graph, parameters.EXECUTE_FOREVER)
	}
}

func (ee *EE) DeployApp(mee madl.MADL, mapp madl.MADL) {

	elems := []madl.Element{}
	for i := range mapp.Components {
		elems = append(elems, mapp.Components[i])
	}
	for i := range mapp.Connectors {
		elems = append(elems, mapp.Connectors[i])
	}

	idx := 0
	for i := range mee.Components {
		if mee.Components[i].TypeName == "Unit" { // TODO
			infoTemp := make([]*interface{}, 1)
			infoTemp[0] = new(interface{})
			*infoTemp[0] = elems[idx]
			mee.Components[i].Info = infoTemp
			idx++
		}
	}

	ee.MADLX = mee
}

func (ee *EE) Deploy(m madl.MADL) {

	ee.MADLX = m
}
