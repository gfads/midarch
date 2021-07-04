package factories

import (
	"gmidarch/development/components/proxies/namingproxy"
	"gmidarch/development/messages"
	"shared"
)

func LocateNaming() namingproxy.Namingproxy {
	chin := make(chan messages.SAMessage)
	chout := make(chan messages.SAMessage)

	p := namingproxy.Namingproxy{Host: shared.NAMING_HOST,Port: shared.NAMING_PORT,ChIn:chin,ChOut:chout}

	return p
}

//func FactoryQueueing() components.Notificationengineproxy {
//	p := components.Notificationengineproxy{Host: sharedadaptive.QUEUEING_HOST, Port: sharedadaptive.QUEUEING_PORT}
//	return p
//}
