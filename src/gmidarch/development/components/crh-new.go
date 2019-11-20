package components

/*
import (
	"encoding/binary"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net"
	"os"
	"shared"
)

type CRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	//Conn     net.Conn
}

var ConnCRH net.Conn

func NewCRH() CRH {

	r := new(CRH)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	//r.Conn = *new(net.Conn)

	return *r
}

func (CRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	elem.(CRH).I_Process(msg, info)
}

func (c CRH) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	payload := msg.Payload.([]interface{})
	//host := payload[0].(string)        // host
	host := "localhost"                // host TODO
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte) // message

	key := host + ":" + port
	//if _, ok := c.Conns[key]; !ok { // no connection open yet
		servAddr := key // TODO
		tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		if err != nil {
			fmt.Printf("Client:: %v\n", err)
			os.Exit(0)
		}

		ConnCRH, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Printf("Client:: %v\n", err)
			os.Exit(0)
		}
	//}

	fmt.Printf("CRH:: I_Process:: (%v,%v)\n", ConnCRH.LocalAddr(), ConnCRH.RemoteAddr())

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = ConnCRH.Write(size)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// send message
	_, err = ConnCRH.Write(msgToServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	//fmt.Printf("CRH:: HERE\n")

	// receive reply's size
	_, err = ConnCRH.Read(size)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = ConnCRH.Read(msgFromServer)
	if err != nil {
		log.Fatalf("CRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}
*/