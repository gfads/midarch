package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"os"
	"strconv"
	"strings"
)

type HttpProxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      string
}

var i_PreInvRHP = make(chan messages.SAMessage) // Todo: Make global variables local
var i_PosTerRHP = make(chan messages.SAMessage) // Todo: Make global variables local

func NewHttpProxy() HttpProxy {

	r := new(HttpProxy)
	r.Behaviour = "B = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> B"

	return *r
}

func (e HttpProxy) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { //"I_Out"
		e.I_Out(msg, info)
	}
}

func (e HttpProxy) Fibo(_p1 int) int {
	_args := []interface{}{"GET","p="+strconv.Itoa(_p1)}

	//request := messages.HttpRequest{
	//	Method:          "GET",
	//	Route:           "/Fibo",
	//	QueryParameters: "p=" + strconv.Itoa(_p1),
	//	Protocol:        "HTTP/1.1",
	//}
	//var payload []interface{}
	//payload = append(payload, e.Host)
	//payload = append(payload, e.Port)
	//payload = append(payload, request.Marshal())

	_reqMsg := messages.SAMessage{ messages.Invocation{Host: e.Host, Port: e.Port, Op: "/Fibo", Args: _args}}

	i_PreInvRHP  <- _reqMsg
	_repMsg := <-i_PosTerRHP

	//_reply := _repMsg.Payload.([]byte)
	//response := messages.HttpResponse{}
	//response.Unmarshal(string(_reply))
	response := _repMsg.Payload.(messages.HttpResponse)
	if !strings.HasPrefix(response.Status, "2") {
		fmt.Printf("HttpProxy:: Response Status: %v\n", response.Status)
		os.Exit(1)
	}
	//fmt.Printf("HttpProxy:: Response Body: %v\n", response.Body)

	_r,err := strconv.Atoi(response.Body)
	if err != nil {
		_r = -1
		fmt.Printf("HttpProxy:: %v\n", err)
		os.Exit(1)
	}

	return _r
}

func (HttpProxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	*msg = <- i_PreInvRHP
}

func (HttpProxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	i_PosTerRHP <- *msg
}