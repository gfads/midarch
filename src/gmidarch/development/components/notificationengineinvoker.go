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

func (e Notificationengineinvoker) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{},r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg)
	} else { // I_Out
		e.I_Out(msg)
	}
}

func (Notificationengineinvoker) I_In(msg *messages.SAMessage) {

	// unmarshall
	payload := msg.Payload.([]byte)

	miopPacket := miop.Packet{}
	err := msgpack.Unmarshal(payload, &miopPacket)
	if err != nil {
		log.Fatalf("Notificationegnineinvoker:: %s\n", err)
		os.Exit(1)
	}

	var inv shared.Request
	op := miopPacket.Bd.ReqHeader.Operation
	switch op {
	case "Publish":
		_packetMOM := miopPacket.Bd.ReqBody.Body[0].(map[string]interface{})
		_header := _packetMOM["Header"].(map[string]interface{})
		_msg := _packetMOM["Payload"].(string)  // TODO
		_destination := _header["Destination"]
		argsTemp := make([]interface{}, 2, 2)
		argsTemp[0] = _destination
		argsTemp[1] = _msg
		inv = shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	case "Subscribe":
		_p0 := miopPacket.Bd.ReqBody.Body[0].(string)
		_p1 := miopPacket.Bd.ReqBody.Body[1].(string)
		_p2 := miopPacket.Bd.ReqBody.Body[2].(string)
		argsTemp := make([]interface{}, 3, 3)
		argsTemp[0] = _p0
		argsTemp[1] = _p1
		argsTemp[2] = _p2
		inv = shared.Request{Op: miopPacket.Bd.ReqHeader.Operation, Args: argsTemp}
	case "Unsubscribe":
	default:
		fmt.Printf("Notificationegnineinvoker:: Operation '%v' not implemented by Naming Service\n", op)
		os.Exit(1)
	}
	*msg = messages.SAMessage{Payload: inv}
}

func (Notificationengineinvoker) I_Out(msg *messages.SAMessage) { // TODO

	// assembly packet
	repHeader := miop.ReplyHeader{Context: "TODO", RequestId: 13, Status: 131313}
	repBody := miop.ReplyBody{OperationResult: *msg}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 2, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{RepHeader: repHeader, RepBody: repBody}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// configure message
	r, err := msgpack.Marshal(miopPacket)
	if err != nil {
		log.Fatalf("Notificationegnineinvoker:: %s", err)
		os.Exit(1)
	}

	toSRH := make([]interface{}, 1, 1)
	toSRH[0] = r

	*msg = messages.SAMessage{Payload: toSRH}
}
