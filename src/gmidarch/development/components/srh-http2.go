package components

import (
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"net/http"
	"shared"
)

type SRHHttp2 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var first = true

var c1Http2 = make(chan messages.HttpMessage)
var c2Http2 = make(chan messages.HttpMessage)

func NewSRHHttp2() SRHHttp2 {

	r := new(SRHHttp2)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRHHttp2) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHHttp2).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHHttp2).I_Send(msg, info, elemInfo)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	//log.Println("Before c1Http")
	c1Http2 <- messages.HttpMessage{w, r}
	// Awaiting for message processing to return
	<- c2Http2
	//response := <- c2Http2
	//log.Println("Message:", response)
}

func (e SRHHttp2) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	//host := "localhost" // TODO

	//log.Println("I_Receive.Begin")

	if first { // listener was not created yet
		first = false
		http.HandleFunc("/", handler) //makeHandler(impl.Handler))
		go http.ListenAndServeTLS(":"+port, shared.CRT_PATH, shared.KEY_PATH, nil)
	}
	//log.Println("Before receive c1Http")
	httpMessage := <- c1Http2

	msg.Payload = httpMessage
	//log.Println("HttpMessage:", httpMessage)
	//log.Println("I_Receive.End")
}

func (e SRHHttp2) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	//log.Println("I_Send.Begin")
	httpMessage := msg.Payload.(messages.HttpMessage)

	// Report that the message was sent
	c2Http2 <- httpMessage

	//log.Println("I_Send.End")
}