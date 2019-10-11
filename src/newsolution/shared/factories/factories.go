package factories

import (
	"gmidarch/sharedch/shared/parameters"
	"framework/components"
	"framework/components/namingclientproxy"
)

func FactoryQueueing() components.NotificationEngineClientProxy {
	p := components.NotificationEngineClientProxy{Host: parameters.QUEUEING_HOST, Port: parameters.QUEUEING_PORT}
	return p
}

func LocateNaming() namingclientproxy.NamingClientProxy {
	p := namingclientproxy.NamingClientProxy{Host:parameters.NAMING_HOST,Port:parameters.NAMING_PORT}
	return p
}