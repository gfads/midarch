package namingproxy

import (
	"gmidarch/development/components/proxies/calculatorproxy"
	"gmidarch/development/generic"
	"gmidarch/development/messages"
	"reflect"
)

//@Type: Namingproxy
//@Behaviour: Behaviour = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> Behaviour

var ProxiesRepo = map[string]interface{}{ // Update for each new proxy
	reflect.TypeOf(Namingproxy{}).Name():                     Namingproxy{},
	reflect.TypeOf(calculatorproxy.Calculatorproxy{}).Name(): calculatorproxy.Calculatorproxy{}}

// Internal channels
var ChOut, ChIn chan messages.SAMessage

type Namingproxy struct {
	GenericProxy generic.Proxy
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
	genericProxy := generic.Proxy{Host: endPoint.Host, Port: endPoint.Port, ProxyName: reflect.TypeOf(Namingproxy{}).Name()}
	proxy = Namingproxy{GenericProxy: genericProxy}

	return proxy
}

// Functional operations
func (Namingproxy) Register(_p1 string, _p2 interface{}) bool {

	aux := reflect.ValueOf(_p2).FieldByName("GenericProxy")
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

func (p Namingproxy) Lookup(_p1 string) (interface{}, bool) {
	_params := []interface{}{_p1}

	_functionalRequest := messages.FunctionalRequest{Op: "Lookup", Params: _params}
	_msg := messages.Invocation{Endpoint:messages.EndPoint{},Functionalrequest:_functionalRequest} // Naming endpoint defined at architectural level
	_samMsg := messages.SAMessage{Payload: _msg}

	// Send request to I_In
	ChIn <- _samMsg

	// Receive response from I_Out
	response := <-ChOut

	aux := response.Payload.(messages.FunctionalReply).Rep.(map[string]interface{})
	//host := aux["host"].(string)
	//port:= aux["port"].(string)
	proxy := ProxiesRepo[aux["proxy"].(string)]

	return proxy, bool(true) // TODO
}

func (p Namingproxy) List() [] interface{} {
	return *new([]interface{})
}
