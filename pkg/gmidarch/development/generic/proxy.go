package generic

type Proxy interface {
	Configure(config ProxyConfig)
}

type ProxyConfig struct {
	Host      string
	Port      string
	ProxyName string
}
