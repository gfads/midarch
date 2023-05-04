package factories

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/components/proxies/namingproxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
)

func LocateNaming() namingproxy.Namingproxy {
	chin := make(chan messages.SAMessage)
	chout := make(chan messages.SAMessage)

	p := namingproxy.Namingproxy{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT, ChIn: chin, ChOut: chout}

	return p
}

//func FactoryQueueing() components.Notificationengineproxy {
//	p := components.Notificationengineproxy{Host: sharedadaptive.QUEUEING_HOST, Port: sharedadaptive.QUEUEING_PORT}
//	return p
//}
