package components

import (
	"apps/fibomiddleware/impl"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net"
	"net/rpc"
)

type SRHRpc struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var firstRPC = true

var c1RPC = make(chan messages.HttpMessage) // TODO: Adjust message type
var c2RPC = make(chan messages.HttpMessage)

func NewSRHRpc() SRHRpc {

	r := new(SRHRpc)
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	return *r
}

func (e SRHRpc) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'R' { // I_Receive
		elem.(SRHRpc).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRHRpc).I_Send(msg, info, elemInfo)
	}
}

//func handler(w http.ResponseWriter, r *http.Request) {
//	//log.Println("Before c1Http")
//	c1RPC <- messages.HttpMessage{w, r}
//	// Awaiting for message processing to return
//	<-c2RPC
//	//response := <- c2RPC
//	//log.Println("Message:", response)
//}

func (e SRHRpc) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "0.0.0.0" // TODO

	//log.Println("I_Receive.Begin")

	if firstRPC { // listener was not created yet
		firstRPC = false
		//http.HandleFunc("/", handler) //makeHandler(impl.Handler))
		//go http.ListenAndServeTLS(":"+port, shared.CRT_PATH, shared.KEY_PATH, nil)

		fibonacci := new(impl.Fibonacci)

		rpc.Register(fibonacci)

		addr, err := net.ResolveTCPAddr("tcp", host + ":" + port)//shared.FIBONACCI_PORT)
		if err != nil {
			log.Fatal("Error while resolving IP address: ", err)
		}
		ln, err := net.ListenTCP("tcp", addr)
		rpc.Accept(ln)
		fmt.Println("Chegou após accept")
	}
	fmt.Println("Chegou após first loop")
	//log.Println("Before receive c1Http")
	//httpMessage := <-c1RPC

	//msg.Payload = httpMessage
	//log.Println("HttpMessage:", httpMessage)
	//log.Println("I_Receive.End")
}

func (e SRHRpc) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	//log.Println("I_Send.Begin")
	//httpMessage := msg.Payload.(messages.HttpMessage)

	// Report that the message was sent
	//c2RPC <- httpMessage

	//log.Println("I_Send.End")
}