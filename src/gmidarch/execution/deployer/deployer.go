package deployer

import (
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/components"
	"gmidarch/execution/engine"
	"reflect"
	"shared"
)

type Deployer struct {
	MADLX madl.MADL
}

func NewEE() Deployer {
	r := new(Deployer)
	return *r
}

func (d Deployer) Start() {

	for i := range d.MADLX.Components {
		elem := d.MADLX.Components[i].Type
		graph := d.MADLX.Components[i].Graph

		// Configure Unit's Info with Element's Info (Only components)
		switch d.MADLX.Components[i].TypeName {
		case reflect.TypeOf(components.Unit{}).Name():
			tempElem := *d.MADLX.Components[i].Info[0]
			unit := elem.(components.Unit)
			unit.UnitId = d.MADLX.Components[i].ElemId
			unit.ElemOfUnit = tempElem.(madl.Element).Type
			unit.GraphOfElem = tempElem.(madl.Element).Graph
			unit.ElemOfUnitInfo = tempElem.(madl.Element).Info
			elem = unit
			go engine.Engine{}.Execute(elem, unit.ElemOfUnitInfo, graph, shared.EXECUTE_FOREVER)
		default:
			go engine.Engine{}.Execute(elem, d.MADLX.Components[i].Info, graph, shared.EXECUTE_FOREVER)
		}
	}

	for i := range d.MADLX.Connectors {
		go engine.Engine{}.Execute(d.MADLX.Connectors[i].Type, d.MADLX.Components[i].Info, d.MADLX.Connectors[i].Graph, shared.EXECUTE_FOREVER)
	}
}

func (d *Deployer) DeployApp(mee madl.MADL, mapp madl.MADL) {

	elems := []madl.Element{}
	for i := range mapp.Components {
		elems = append(elems, mapp.Components[i])
	}
	for i := range mapp.Connectors {
		elems = append(elems, mapp.Connectors[i])
	}

	idx := 0
	for i := range mee.Components {
		if mee.Components[i].TypeName == reflect.TypeOf(components.Unit{}).Name() {
			infoTemp := make([]*interface{}, 1)
			infoTemp[0] = new(interface{})
			*infoTemp[0] = elems[idx]
			mee.Components[i].Info = infoTemp
			idx++
		}
	}
	d.MADLX = mee
}

func (d *Deployer) Deploy(m madl.MADL) {

	d.MADLX = m
}
