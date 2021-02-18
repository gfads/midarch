package components

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

type CRHUdp struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]net.Conn
}

func NewCRHUdp() CRHUdp {
	r := new(CRHUdp)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]net.Conn, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (CRHUdp) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(CRHUdp).I_Process(msg, info)
}

func (c CRHUdp) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	payload := msg.Payload.([]interface{})
	host := "127.0.0.1"                // host TODO
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte)

	key := host + ":" + port
	var err error
	if _, ok := c.Conns[key]; !ok { // no connection open yet
		//servAddr := key // TODO
		//tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		udpAddr, err := net.ResolveUDPAddr("udp", key)
		if err != nil {
			log.Fatalf("CRHUdp:: %s", err)
		}

		c.Conns[key], err = net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Printf("CRHUdp:: %v\n", err)
			os.Exit(1)
		}
	}

	// connect to server
	conn := c.Conns[key]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = conn.Write(size)
	if err != nil {
		fmt.Printf("CRHUdp:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHUdp:: %v \n\n",size)

	// send message
	//fmt.Printf("CRHUdp:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	_, err = conn.Write(msgToServer)
	if err != nil {
		fmt.Printf("CRHUdp:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHUdp:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	// receive reply's size
	_, err = conn.Read(size)
	if err != nil {
		fmt.Printf("CRHUdp:: %v\n", err)
		os.Exit(1)
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		fmt.Printf("CRHUdp:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHUdp:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	*msg = messages.SAMessage{Payload: msgFromServer}
}