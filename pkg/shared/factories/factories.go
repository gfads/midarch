package factories

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/components/proxies/namingproxy"
	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/shared"
)

func LocateNaming() namingproxy.Namingproxy {
	p := namingproxy.Namingproxy{Config: generic.ProxyConfig{Host: shared.NAMING_HOST, Port: shared.NAMING_PORT}}

	return p
}

//func FactoryQueueing() components.Notificationengineproxy {
//	p := components.Notificationengineproxy{Host: sharedadaptive.QUEUEING_HOST, Port: sharedadaptive.QUEUEING_PORT}
//	return p
//}
