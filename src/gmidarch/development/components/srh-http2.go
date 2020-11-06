package components

import (
	"apps/http2server/impl"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net/http"
	"shared"
)

type SRHHttp2 struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

//var http2 net.Http
var first = true

var c1Http2 = make(chan []byte)
var c2Http2 = make(chan []byte)

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

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (e SRHHttp2) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	//host := "localhost" // TODO

	log.Println("I_Receive.Begin")

	if first { // listener was not created yet
		first = false
		http.HandleFunc("/", makeHandler(impl.Handler))
		log.Fatal(http.ListenAndServeTLS(":"+port, shared.CRT_PATH, shared.KEY_PATH, nil))
	}

	log.Println("I_Receive.End")

	//switch stateHttp2 {
	//case 0:
	//	go acceptAndReadHttp2(currentConnectionHttp2, c1Http2)
	//	stateHttp2 = 1
	//case 1:
	//	go acceptAndReadHttp2(currentConnectionHttp2, c1Http2)
	//	stateHttp2 = 2
	//case 2:
	//	go acceptAndReadHttp2(currentConnectionHttp2, c1Http2)
	//}
	//
	////go acceptAndReadHttp2(currentConnectionHttp2, c1Http2, done)
	////go readHttp2(currentConnectionHttp2, c2Http2, done)
	//
	//select {
	//case msgTemp := <-c1Http2:
	//	*msg = messages.SAMessage{Payload: msgTemp}
	//case msgTemp := <-c2Http2:
	//	*msg = messages.SAMessage{Payload: msgTemp}
	//}
	//
	//currentConnectionHttp2 = nextConnectionHttp2()
}

func (e SRHHttp2) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]byte)

	//// send message
	//_, err := ConnsSRHHttp2[currentConnectionHttp2].Write(msgTemp)
	//if err != nil {
	//	fmt.Printf("SRHHttp2:: %v\n", err)
	//	os.Exit(1)
	//}
	//
	//ConnsSRHHttp2[currentConnectionHttp2].Close()

	log.Println(msgTemp)
}