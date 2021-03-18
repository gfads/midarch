package components

import (
	"encoding/binary"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"net"
	"os"
	"shared"
)

type CRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]net.Conn
}

func NewCRH() CRH {

	r := new(CRH)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]net.Conn, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (CRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(CRH).I_Process(msg, info)
}

func (c CRH) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	payload := msg.Payload.([]interface{})
	host := payload[0].(string)        // host
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte)

	addr := host + ":" + port
	var err error
	if _, ok := c.Conns[addr]; !ok { // no connection open yet
		//servAddr := key // TODO
		//tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		//tcpAddr, err := net.ResolveTCPAddr("tcp", key)
		//if err != nil {
		//	log.Fatalf("CRH:: %s", err)
		//}
		//fmt.Println("addr", addr)
		c.Conns[addr], err = net.Dial("tcp", addr)//nil, tcpAddr)
		if err != nil {
			fmt.Printf("CRH:: %v\n", err)
			os.Exit(1)
		}
	}

	// connect to server
	conn := c.Conns[addr]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = conn.Write(size)
	if err != nil {
		fmt.Printf("CRH:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRH:: %v \n\n",size)

	// send message
	//fmt.Printf("CRH:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	_, err = conn.Write(msgToServer)
	if err != nil {
		fmt.Printf("CRH:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRH:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	// receive reply's size
	_, err = conn.Read(size)
	if err != nil {
		fmt.Printf("CRH:: %v\n", err)
		os.Exit(1)
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		fmt.Printf("CRH:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRH:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	*msg = messages.SAMessage{Payload: msgFromServer}
}