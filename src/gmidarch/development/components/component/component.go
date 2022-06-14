package component

import (
	"gmidarch/development/artefacts/graphs/dot"
	"gmidarch/development/connectors"
	"gmidarch/development/messages"
	"shared"
)

type Component struct {
	Id        		string
	TypeName  		string
	Type      		interface{}
	Behaviour 		string
	Buffer    		messages.SAMessage
	Graph     		dot.DOTGraph
	Info      		interface{}
	ExecuteForever	*bool // TODO dcruzb: move to the Type attribute (needs to modify the current struct to add the attribute dynamically) this will make possible the start of new srh while the old one is still executing
	Executing		*bool // TODO dcruzb: move to the Type attribute (needs to modify the current struct to add the attribute dynamically) this will make possible the start of new srh while the old one is still executing
}

func (Component) InvR(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}, reset *bool) {
	connector.Run(msg, shared.INVR, invoker)
}
func (Component) InvP(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}, reset *bool) {
	connector.Run(msg, shared.INVP, invoker)
}
func (Component) TerP(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}, reset *bool) {
	connector.Run(msg, shared.TERP, invoker)
}
func (Component) TerR(msg *messages.SAMessage, connector connectors.Connector, invoker string, info *interface{}, reset *bool) {
	connector.Run(msg, shared.TERR, invoker)
}
