package components

import (
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Http2Proxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      string
}

var i_PreInvRH2P = make(chan messages.SAMessage) // Todo: Make global variables local
var i_PosTerRH2P = make(chan messages.SAMessage) // Todo: Make global variables local

func NewHttp2Proxy() Http2Proxy {

	r := new(Http2Proxy)
	r.Behaviour = "B = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> B"

	return *r
}

func (e Http2Proxy) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { //"I_Out"
		e.I_Out(msg, info)
	}
}

func (e Http2Proxy) Fibo(_p1 int) int {
	_args := []interface{}{"GET","place="+strconv.Itoa(_p1)}

	_reqMsg := messages.SAMessage{ messages.Invocation{Host: e.Host, Port: e.Port, Op: "/api/fibo", Args: _args}}

	i_PreInvRH2P  <- _reqMsg
	_repMsg := <- i_PosTerRH2P

	response := _repMsg.Payload.(*http.Response)
	if !strings.HasPrefix(response.Status, "2") {
		fmt.Printf("Http2Proxy:: Response Status: %v\n", response.Status)
		os.Exit(1)
	}
	//fmt.Printf("Http2Proxy:: Response Body: %v\n", response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Http2Proxy:: %v\n", err)
		os.Exit(1)
	}
	_r,err := strconv.Atoi(string(body))
	if err != nil {
		_r = -1
		fmt.Printf("Http2Proxy:: %v\n", err)
		os.Exit(1)
	}

	return _r
}

func (Http2Proxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	*msg = <- i_PreInvRH2P
}

func (Http2Proxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	i_PosTerRH2P <- *msg
}