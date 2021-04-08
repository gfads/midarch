package components

import (
	"github.com/vmihailenco/msgpack"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"gmidarch/development/miop"
	"log"
	"os"
)

type RequestorM struct {
	Graph     graphs.ExecGraph
	Behaviour string
}

func NewRequestorM() RequestorM {

	r := new(RequestorM)
	r.Behaviour = "B = InvP.e1 -> I_In -> InvR.e2 -> TerR.e2 -> I_Out -> TerP.e1 -> B"

	return *r
}

func (e RequestorM) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	if op[2] == 'I' { // I_In
		e.I_In(msg, info)
	} else { // "I_Out"
		e.I_Out(msg, info)
	}
}

func (RequestorM) I_In(msg *messages.SAMessage, info [] *interface{}) {
	inv := msg.Payload.(messages.Invocation)
	// assembly packet
	reqHeader := miop.RequestHeader{Context: "TODO", RequestId: 13, ResponseExpected: true, Key: 131313, Operation: inv.Op}
	reqBody := miop.RequestBody{Body: inv.Args}
	miopHeader := miop.Header{Magic: "M.I.O.P.", Version: "version", MessageType: 1, Size: 131313, ByteOrder: true}
	miopBody := miop.Body{ReqHeader: reqHeader, ReqBody: reqBody, RepHeader: miop.ReplyHeader{Context: "", RequestId: 0, Status: 0}, RepBody: miop.ReplyBody{OperationResult: 0}}
	miopPacket := miop.Packet{Hdr: miopHeader, Bd: miopBody}

	// marshall packet
	pckt, err := msgpack.Marshal(miopPacket)
	//pckt, err := json.Marshal(miopPacket)
	if err != nil {
		log.Fatalf("RequestorM:: %s", err)
		os.Exit(1)
	}

	//fmt.Printf("Requestor:: %v\n", pckt)

	// store host & port in 'info'
	*info[0] = inv.Host
	*info[1] = inv.Port

	toCRH := make([]interface{}, 3, 3)
	toCRH[0] = inv.Host
	toCRH[1] = inv.Port
	toCRH[2] = pckt

	*msg = messages.SAMessage{Payload: toCRH}
}

func (RequestorM) I_Out(msg *messages.SAMessage, info [] *interface{}) {
	miopPacket := miop.Packet{}

	err := msgpack.Unmarshal(msg.Payload.([]byte), &miopPacket)
	if err != nil {
		log.Fatalf("Requestorwithmarshaller:: %s", err)
	}

	*msg = messages.SAMessage{Payload: miopPacket.Bd.RepBody.OperationResult}
}
