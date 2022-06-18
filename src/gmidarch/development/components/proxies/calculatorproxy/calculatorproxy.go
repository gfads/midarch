package calculatorproxy

import (
	"fmt"
	"gmidarch/development/generic"
	"gmidarch/development/messages"
	"shared"
)

//@Type: Calculatorproxy
//@Behaviour: Behaviour = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> Behaviour
type Calculatorproxy struct{
	Config generic.ProxyConfig
}

// Internal channels
var ChOut, ChIn chan messages.SAMessage

// Factory
func NewCalculatorProxy(proxyConfig generic.ProxyConfig) Calculatorproxy { // TODO
	var proxy Calculatorproxy

	// Create internal channels  TODO
	//ChIn = make(chan messages.SAMessage)
	//ChOut = make(chan messages.SAMessage)

	fmt.Println(shared.GetFunction())

	// Configure proxy
    //genericProxy := generic.Proxy{Host: shared.CALCULATOR_HOST, Port: shared.CALCULATOR_PORT} // TODO dcruzb: Host and ports should not be constants
    proxy = Calculatorproxy{Config: proxyConfig}

	return proxy
}

func (p *Calculatorproxy) Configure(config generic.ProxyConfig) {
	p.Config = config
}

// Architectural operations
func (Calculatorproxy) I_In(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {

	// Create internal channels
	ChIn = make(chan messages.SAMessage)
	ChOut = make(chan messages.SAMessage)

	fmt.Println(shared.GetFunction(),ChIn)

	// Receive request from Client through the invocation of operations of the functional interface
	*msg = <-ChIn
}

func (Calculatorproxy) I_Out(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {

	fmt.Println(shared.GetFunction())

	// Send message to Client
	ChOut <- *msg
}

// Functional operations
func (p Calculatorproxy) Add(_p1,_p2 int) int {

	port := p.Config.Port
	host := p.Config.Host
	_endPoint := messages.EndPoint{Host:host, Port:port}

	fmt.Println(shared.GetFunction(), _endPoint.Host)

	_params := []interface{}{_p1, _p2}

	_functionalRequest := messages.FunctionalRequest{Op: "Add", Params: _params}
	_msg := messages.Invocation{Endpoint:_endPoint, Functionalrequest:_functionalRequest} // Naming endpoint defined at architectural level

	_samMsg := messages.SAMessage{Payload: _msg}

	fmt.Println(shared.GetFunction(), ChIn)

	// Send request to I_In
	ChIn <- _samMsg

	var response messages.SAMessage
	// Receive response from I_Out
	response = <-ChOut

	var result float64
	// Try again if there is no valid response
	if response.Payload.(messages.FunctionalReply).Rep == nil {
		//// Send request to I_In
		//ChIn <- _samMsg
		//
		//// Receive response from I_Out
		//response = <-ChOut
		result = 0
	} else {
		result = response.Payload.(messages.FunctionalReply).Rep.(float64)
	}
	fmt.Println(shared.GetFunction(), result)

	return int(result)
}