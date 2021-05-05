package factories

import (
	"gmidarch/development/components"
	"shared"
)

func LocateNaming() components.Namingproxy {
	p := components.Namingproxy{Host: shared.NAMING_HOST,Port: shared.NAMING_PORT}
	return p
}

func GetHttpProxy(url string, port string) components.HttpProxy {
	p := components.HttpProxy{Host: url,Port: port}
	return p
}

func GetHttp2Proxy(url string, port string) components.Http2Proxy {
	p := components.Http2Proxy{Host: url,Port: port}
	return p
}

func GetFibonacciProxy(url string, port string) components.Fibonacciproxy {
	p := components.NewFibonacciproxy()
	p.Host = url
	p.Port = port
	return p
}

func FactoryQueueing() components.Notificationengineproxy {
	p := components.Notificationengineproxy{Host: shared.QUEUEING_HOST, Port: shared.QUEUEING_PORT}
	return p
}
