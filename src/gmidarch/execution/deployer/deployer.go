package deployer

import (
	"gmidarch/development/artefacts/madl"
	"gmidarch/development/components/adaptive"
	"gmidarch/development/components/component"
	"gmidarch/development/connectors"
	"gmidarch/development/messages"
	"gmidarch/execution/core"
	"reflect"
	"shared"
)

type DeployParameters struct {
	Args map[string]messages.EndPoint
}

type Deployer interface {
	//Deploy(madl.MADL)
	Start()
}

type DeployerImpl struct {
	Madl madl.MADL
}

func NewDeployer(m madl.MADL) Deployer {
	return DeployerImpl{Madl: m}
}

func NewEEDeployer(m madl.MADL, mee madl.MADL) Deployer {
	components := []component.Component{}
	for i := range m.Components {
		components = append(components, m.Components[i])
	}
	connectors := []connectors.Connector{}
	for i := range m.Connectors {
		connectors = append(connectors, m.Connectors[i])
	}

	idx := 0
	for i := range mee.Components {
		if mee.Components[i].TypeName == reflect.TypeOf(adaptive.Unit{}).Name() {
			infoTemp := make([]*interface{}, 1)
			infoTemp[0] = new(interface{})
			*infoTemp[0] = components[idx]
			mee.Components[i].Info = infoTemp
			idx++
		}
	}

	return NewDeployer(mee)
}

func (d DeployerImpl) Start() {
	for i := range d.Madl.Components {
		engine := engine.NewEngine()
		go engine.Execute(&d.Madl.Components[i], shared.EXECUTE_FOREVER)
	}
}

/*
func (d *DeployerImpl) Deploy(mee madl.MADL, mapp madl.MADL) {

	comps := []component.Component{}
	for i := range mapp.Components {
		comps = append(comps, mapp.Components[i])
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
*/

/*
func (d *Deployer) Deploy(m madl.MADL) {

	d.MADLX = m
}
*/
