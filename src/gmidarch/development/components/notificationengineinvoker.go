package components

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"os"
	"shared"
)

type Notificationengineinvoker struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

func Newnotificationengineinvoker() Notificationengineinvoker {
	r := new(Notificationengineinvoker)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e Notificationengineinvoker) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg)
	} else { // I_Out
		e.I_Out(msg)
	}
}

func (Notificationengineinvoker) I_In(msg *messages.SAMessage) {
	// unmarshall
	payload := msg.Payload.([]byte)

	//fmt.Printf("NotificationEngineInvoker:: Payload %v\n",payload)

	miopPacket := miop.Packet{}

	//err := json.Unmarshal(payload, &miopPacket)
	err := msgpack.Unmarshal(payload, &miopPacket)
	if err != nil {
		fmt.Printf("Notificationegnineinvoker:: I_In:: %v %v >> %v <<\n", err, payload,len(payload))
		os.Exit(1)
	}

 	var inv shared.Request
	op := miopPacket.Bd.ReqHeader.Operation
	switch op {
	case "Publish":
		_x1 := miopPacket.Bd.ReqBody.Body[0].(map[string]interface{})
		_x2 := _x1["Header"].(map[string]interface{})
		_destination := _x2["Destination"].(string)
		_payload := _x1["Payload"]
		_header := messages.MOMHeader{Destination: _destination}
		_msgMOM := messages.MessageMOM{Header: _header, Payload: _payload}
		_args := make([]interface{}, 1, 1)
		_args[0] = _msgMOM
		inv = shared.Request{Op: op, Args: _args}
	case "Subscribe":
		_p0 := miopPacket.Bd.ReqBody.Body[0].(string)
		_p1 := miopPacket.Bd.ReqBody.Body[1].(string)
		_p2 := miopPacket.Bd.ReqBody.Body[2].(string)
		argsTemp := make([]interface{}, 3, 3)
		argsTemp[0] = _p0
		argsTemp[1] = _p1
		argsTemp[2] = _p2
		inv = shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	case "Unsubscribe": // TODO
	default:
		fmt.Printf("Notificationegnineinvoker:: Operation '%v' not implemented by Notification Engine\n", op)
		os.Exit(1)
	}
	*msg = messages.SAMessage{Payload: inv}
}

func (Notificationengineinvoker) I_Out(msg *messages.SAMessage) { // TODO

	// assembly packet
	reqHeader := miop.RequestHeader{Context: "", RequestId: 0, ResponseExpected: true, Key: 0, Operation: ""}
	reqBody := miop.RequestBody{[]interface{}{""}}
	repHeader := miop.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
	repBody := miop.ReplyBody{OperationResult: *msg}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{RepHeader: repHeader, RepBody: repBody, ReqHeader: reqHeader, ReqBody: reqBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// configure message
	r, err := msgpack.Marshal(miopPacket)
	if err != nil {
		log.Fatalf("Notificationegnineinvoker:: I_Out %s", err)
		os.Exit(1)
	}

	toSRH := make([]interface{}, 1, 1)
	toSRH[0] = r

	*msg = messages.SAMessage{Payload: toSRH}
}
