package component

import (
	"gmidarch/development/artefacts/graphs/dot"
	"gmidarch/development/connectors"
	"gmidarch/development/messages"
	"shared"
)

type Component struct {
	Id        string
	TypeName  string
	Type      interface{}
	Behaviour string
	Buffer    messages.SAMessage
	Graph     dot.DOTGraph
	Info      interface{}
}

func (Component) InvR(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}) {
	connector.Run(msg, shared.INVR, invoker)
}
func (Component) InvP(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}) {
	connector.Run(msg, shared.INVP, invoker)
}
func (Component) TerP(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}) {
	connector.Run(msg, shared.TERP, invoker)
}
func (Component) TerR(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}) {
	connector.Run(msg, shared.TERR, invoker)
}
