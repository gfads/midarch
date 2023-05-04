package deployer

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/artefacts/madl"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/adaptive"
	"github.com/gfads/midarch/pkg/gmidarch/development/components/component"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/execution/core"
	"github.com/gfads/midarch/pkg/shared"
	"reflect"
	"strings"
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
	components := []*component.Component{}
	for i := range m.Components {
		components = append(components, &m.Components[i])
	}
	//connectors := []connectors.Connector{}
	//for i := range m.Connectors {
	//	connectors = append(connectors, m.Connectors[i])
	//}

	idx := 0
	for i := range mee.Components {
		if mee.Components[i].TypeName == reflect.TypeOf(adaptive.Unit{}).Name() {
			infoTemp := make([]*interface{}, 1)
			infoTemp[0] = new(interface{})
			*infoTemp[0] = components[idx]

			//fmt.Println("NewEEDeployer::Unit.Type", components[idx].Type)
			//fmt.Println("NewEEDeployer::Unit.Type is", reflect.TypeOf(components[idx].Type))
			//fmt.Println("NewEEDeployer::Unit.Graph", components[idx].Graph)
			//fmt.Println("NewEEDeployer::Unit.Info", components[idx].Info)

			mee.Components[i].Info = infoTemp
			idx++
		}
	}

	return NewDeployer(mee)
}

func (d DeployerImpl) Start() {
	for i := range d.Madl.Components {
		engine := engine.NewEngine()
		if d.Madl.Components[i].TypeName == reflect.TypeOf(adaptive.Unit{}).Name() { // TODO dcrzub: Transfer this to NewEEDeployer
			unit := &*(d.Madl.Components[i].Type).(*adaptive.Unit)
			//reflect.ValueOf(unit).MethodByName("PrintId").Call([]reflect.Value{})
			unit.UnitId = d.Madl.Components[i].Id
			//reflect.ValueOf(unit).MethodByName("PrintId").Call([]reflect.Value{})
			//reflect.ValueOf(&(d.Madl.Components[i].Type)).MethodByName("PrintId").Call([]reflect.Value{})
			//element := (*d.Madl.Components[i].Info.([]*interface{})[0]).(component.Component)
			//unit.ElemOfUnit = element.Type
			//unit.GraphOfElem = element.Graph
			//unit.ElemOfUnitInfo = element.Info
			//fmt.Println("elem is", reflect.TypeOf(element))
			//fmt.Println("elem kind is", reflect.TypeOf(element).Kind())
			////fmt.Println("elem kind is", reflect.TypeOf(element).Elem())
			////fmt.Println("elem.Elem() kind is", reflect.TypeOf(element).Elem().Kind())
			//
			//fmt.Println("elem value", reflect.ValueOf(element))
			////fmt.Println("elem value.Elem()", reflect.ValueOf(element).Elem())
			//fmt.Println("elem methods", reflect.ValueOf(element).Type().NumMethod())

			//fmt.Println("Deployer.Start::element.Type is", reflect.TypeOf(element.Type))
			//fmt.Println("Deployer.Start::element.Type kind is", reflect.TypeOf(element.Type).Kind())
			//
			//fmt.Println("Deployer.Start::Unit.ElemOfUnit:", unit.ElemOfUnit)
			//fmt.Println("Deployer.Start::Unit.GraphOfElem:", unit.GraphOfElem)
			//
			//fmt.Println("Deployer.Start::(*d.Madl.Components[i].Info.([]*interface{})[0]).(component.Component).Info:", (*d.Madl.Components[i].Info.([]*interface{})[0]).(component.Component).Info)
			//fmt.Println("Deployer.Start::element.Info:", element.Info)
			//fmt.Println("Deployer.Start::Unit.ElemOfUnitInfo:", unit.ElemOfUnitInfo)
		}
		d.Madl.Components[i].ExecuteForever = &shared.ExecuteForever
		//fmt.Println("Setará executeforever:", d.Madl.Components[i].TypeName)
		if strings.Contains(d.Madl.Components[i].TypeName, "SRH") {
			//fmt.Println("Setou executeforever")
			//log.Println("Setou executeforever")
			infoTemp := d.Madl.Components[i].Info
			srhInfo := infoTemp.(*messages.SRHInfo)
			srhInfo.ExecuteForever = d.Madl.Components[i].ExecuteForever
		}
		go engine.Execute(&d.Madl.Components[i], d.Madl.Components[i].ExecuteForever) // EXECUTE_FOREVER)
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
