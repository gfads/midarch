package ee

import (
	madl2 "gmidarch/development/artefacts/madl"
	components2 "gmidarch/development/components"
	"gmidarch/execution/engine"
	"shared/shared"
)

type EE struct {
	MADLX madl2.MADL
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
			unit := elem.(components2.Unit)
			unit.UnitId = e.MADLX.Components[i].ElemId
			unit.ElemOfUnit= tempElem.(madl2.Element).Type
			unit.GraphOfElem = tempElem.(madl2.Element).Graph
			elem = unit
		}
		go engine.Engine{}.Execute(elem, graph, shared.EXECUTE_FOREVER)
	}

	for i := range e.MADLX.Connectors {
		go engine.Engine{}.Execute(e.MADLX.Connectors[i].Type, e.MADLX.Connectors[i].Graph, shared.EXECUTE_FOREVER)
	}
}

func (ee *EE) DeployApp(mee madl2.MADL, mapp madl2.MADL) {

	elems := []madl2.Element{}
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

func (ee *EE) Deploy(m madl2.MADL) {

	ee.MADLX = m
}
