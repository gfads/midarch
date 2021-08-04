package namingproxy

import (
	"gmidarch/development/components/proxies/calculatorproxy"
	"gmidarch/development/generic"
	"gmidarch/development/messages"
	"reflect"
)

//@Type: Namingproxy
//@Behaviour: Behaviour = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> Behaviour
var ProxiesRepo = map[string]generic.Proxy { // Update for each new proxy
	reflect.TypeOf(Namingproxy{}).Name():                     &Namingproxy{},
	reflect.TypeOf(calculatorproxy.Calculatorproxy{}).Name(): &calculatorproxy.Calculatorproxy{}}

// Internal channels
var ChOut, ChIn chan messages.SAMessage

type Namingproxy struct {
	Config generic.ProxyConfig
}

func (p *Namingproxy) Configure(config generic.ProxyConfig) {
	p.Config = config
}

// Architectural operations
func (Namingproxy) I_In(id string, msg *messages.SAMessage, info *interface{}) {
	// Receive request from Request
	*msg = <-ChIn
}

func (Namingproxy) I_Out(id string, msg *messages.SAMessage, info *interface{}) {
	ChOut <- *msg
}

func NewNamingproxy(endPoint messages.EndPoint) Namingproxy {
	var proxy Namingproxy

	// Create internal channels
	ChIn = make(chan messages.SAMessage)
	ChOut = make(chan messages.SAMessage)

	// Configure proxy
	proxyConfig := generic.ProxyConfig{Host: endPoint.Host, Port: endPoint.Port, ProxyName: reflect.TypeOf(Namingproxy{}).Name()}
	proxy = Namingproxy{Config: proxyConfig}

	return proxy
}

// Functional operations
func (Namingproxy) Register(_p1 string, _p2 interface{}) bool {

	aux := reflect.ValueOf(_p2).FieldByName("Config")
	port := aux.FieldByName("Port").String()
	host := aux.FieldByName("Host").String()
	aor := messages.AOR{Host: host, Port: port, Id: 123456, ProxyName: reflect.TypeOf(_p2).Name()} // TODO
	_params := []interface{}{_p1, aor}

	_functionalRequest := messages.FunctionalRequest{Op: "Register", Params: _params}
	_msg := messages.Invocation{Endpoint:messages.EndPoint{},Functionalrequest:_functionalRequest} // Naming endpoint defined at architectural level

	_samMsg := messages.SAMessage{Payload: _msg}

	// Send request to I_In
	ChIn <- _samMsg

	// Receive response from I_Out
	response := <-ChOut

	return response.Payload.(messages.FunctionalReply).Rep.(bool)
}

func (p Namingproxy) Lookup(_p1 string) (generic.Proxy, bool) {
	_params := []interface{}{_p1}

	_functionalRequest := messages.FunctionalRequest{Op: "Lookup", Params: _params}
	_msg := messages.Invocation{Endpoint:messages.EndPoint{},Functionalrequest:_functionalRequest} // Naming endpoint defined at architectural level
	_samMsg := messages.SAMessage{Payload: _msg}

	// Send request to I_In
	ChIn <- _samMsg

	// Receive response from I_Out
	response := <-ChOut

	aor := response.Payload.(messages.FunctionalReply).Rep.(map[string]interface{})
	host := aor["host"].(string)
	port:= aor["port"].(string)
	proxy := ProxiesRepo[aor["proxy"].(string)]
	proxyConfig := generic.ProxyConfig{Host: host, Port: port} // TODO dcruzb: Host and ports should not be constants, change to aor["host"].(string) and aor["port"].(string), uncomment previous lines
	proxy.Configure(proxyConfig)

	return proxy, bool(true) // TODO
}

func (p Namingproxy) List() [] interface{} {
	return *new([]interface{})
}
