package sendFileProxy

import (
	"os"

	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
)

// @Type: SendFileProxy
// @Behaviour: Behaviour = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> Behaviour
type SendFileProxy struct {
	Config generic.ProxyConfig
}

// Internal channels
var ChOut, ChIn chan messages.SAMessage

// Factory
func NewSendFileProxy(proxyConfig generic.ProxyConfig) SendFileProxy { // TODO
	var proxy SendFileProxy

	// Create internal channels  TODO
	//ChIn = make(chan messages.SAMessage)
	//ChOut = make(chan messages.SAMessage)

	//fmt.Println(shared.GetFunction())

	// Configure proxy
	//genericProxy := generic.Proxy{Host: shared.CALCULATOR_HOST, Port: shared.CALCULATOR_PORT} // TODO dcruzb: Host and ports should not be constants
	proxy = SendFileProxy{Config: proxyConfig}

	return proxy
}

func (p *SendFileProxy) Configure(config generic.ProxyConfig) {
	p.Config = config
}

// Architectural operations
func (SendFileProxy) I_In(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {

	// Create internal channels
	ChIn = make(chan messages.SAMessage)
	ChOut = make(chan messages.SAMessage)

	//fmt.Println(shared.GetFunction(),ChIn)

	// Receive request from Client through the invocation of operations of the functional interface
	*msg = <-ChIn
}

func (SendFileProxy) I_Out(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {

	//fmt.Println(shared.GetFunction())

	// Send message to Client
	ChOut <- *msg
}

// Functional operations
func (p SendFileProxy) SendFile(file *os.File) string {
	port := p.Config.Port
	host := p.Config.Host
	_endPoint := messages.EndPoint{Host: host, Port: port}

	//fmt.Println(shared.GetFunction(), _endPoint.Host)

	_params := []interface{}{file.getBytes()}

	// _functionalRequest := messages.FunctionalRequest{Op: "F", Params: _params}
	_functionalRequest := messages.FunctionalRequest{Op: "I", Params: _params}              // TODO dcruzb: Test to get base64 image
	_msg := messages.Invocation{Endpoint: _endPoint, Functionalrequest: _functionalRequest} // Naming endpoint defined at architectural level

	_samMsg := messages.SAMessage{Payload: _msg}

	//fmt.Println(shared.GetFunction(), ChIn)

	// Send request to I_In
	ChIn <- _samMsg

	var response messages.SAMessage
	// Receive response from I_Out
	response = <-ChOut

	var result string
	// Try again if there is no valid response
	if response.Payload.(messages.FunctionalReply).Rep == nil {
		//// Send request to I_In
		//ChIn <- _samMsg
		//
		//// Receive response from I_Out
		//response = <-ChOut
		result = "0"
	} else {
		// result = response.Payload.(messages.FunctionalReply).Rep.(float64)
		result = response.Payload.(messages.FunctionalReply).Rep.(string)
	}
	//fmt.Println(shared.GetFunction(), result)

	return result //int(result)
}
