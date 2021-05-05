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

type RPCProxy struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      string
}

var i_PreInvRRPCP = make(chan messages.SAMessage) // Todo: Make global variables local
var i_PosTerRRPCP = make(chan messages.SAMessage) // Todo: Make global variables local

func NewRPCProxy() RPCProxy {

	r := new(RPCProxy)
	r.Behaviour = "B = I_In -> InvR.e1 -> TerR.e1 -> I_Out -> B"

	return *r
}

func (e RPCProxy) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { //"I_Out"
		e.I_Out(msg, info)
	}
}

func (e RPCProxy) Fibo(_p1 int) int {
	_args := []interface{}{_p1}

	_reqMsg := messages.SAMessage{ messages.Invocation{Host: e.Host, Port: e.Port, Op: "/api/fibo", Args: _args}}

	i_PreInvRRPCP <- _reqMsg
	_repMsg := <-i_PosTerRRPCP

	response := _repMsg.Payload.(*http.Response)
	if !strings.HasPrefix(response.Status, "2") {
		fmt.Printf("RPCProxy:: Response Status: %v\n", response.Status)
		os.Exit(1)
	}
	//fmt.Printf("RPCProxy:: Response Body: %v\n", response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("RPCProxy:: %v\n", err)
		os.Exit(1)
	}
	_r,err := strconv.Atoi(string(body))
	if err != nil {
		_r = -1
		fmt.Printf("RPCProxy:: %v\n", err)
		os.Exit(1)
	}

	return _r
}

func (RPCProxy) I_In(msg *messages.SAMessage, info [] *interface{}) {
	*msg = <-i_PreInvRRPCP
}

func (RPCProxy) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	i_PosTerRRPCP <- *msg
}