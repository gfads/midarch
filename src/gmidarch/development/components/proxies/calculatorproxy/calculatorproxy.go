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
	GenericProxy generic.Proxy
}

// Internal channels
var ChOut, ChIn chan messages.SAMessage

// Factory
func NewCalculatorProxy() Calculatorproxy { // TODO
	var proxy Calculatorproxy

	// Create internal channels  TODO
	//ChIn = make(chan messages.SAMessage)
	//ChOut = make(chan messages.SAMessage)

	fmt.Println(shared.GetFunction())

	// Configure proxy
    genericProxy := generic.Proxy{Host:shared.CALCULATOR_HOST,Port:shared.CALCULATOR_PORT}
    proxy = Calculatorproxy{GenericProxy:genericProxy}

	return proxy
}

// Architectural operations
func (Calculatorproxy) I_In(id string, msg *messages.SAMessage, info *interface{}) {

	// Create internal channels
	ChIn = make(chan messages.SAMessage)
	ChOut = make(chan messages.SAMessage)

	fmt.Println(shared.GetFunction(),ChIn)

	// Receive request from Client through the invocation of operations of the functional interface
	*msg = <-ChIn
}
func (Calculatorproxy) I_Out(id string, msg *messages.SAMessage, info *interface{}) {

	fmt.Println(shared.GetFunction())

	// Send message to Client
	ChOut <- *msg
}

// Functional operations
func (p Calculatorproxy) Add(_p1,_p2 int) int {

	port := p.GenericProxy.Port
	host := p.GenericProxy.Host
	_endPoint := messages.EndPoint{Host:host,Port:port}

	fmt.Println(shared.GetFunction(),_endPoint.Host)

	_params := []interface{}{_p1, _p2}

	_functionalRequest := messages.FunctionalRequest{Op: "Add", Params: _params}
	_msg := messages.Invocation{Endpoint:_endPoint,Functionalrequest:_functionalRequest} // Naming endpoint defined at architectural level

	_samMsg := messages.SAMessage{Payload: _msg}

	fmt.Println(shared.GetFunction(),ChIn)

	// Send request to I_In
	ChIn <- _samMsg

	// Receive response from I_Out
	response := <-ChOut

	aux := response.Payload.(messages.FunctionalReply).Rep.(map[string]interface{})
	fmt.Println(shared.GetFunction(),aux)

return 100  // TODO
}