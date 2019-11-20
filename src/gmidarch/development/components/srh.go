package components

import (
	"encoding/binary"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io"
	"log"
	"net"
	"os"
	"shared"
)

type SRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
}

var ConnSRH net.Conn
var LnSRH net.Listener

func NewSRH() SRH {

	r := new(SRH)

	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"
	return *r
}

func (e SRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	if op[2] == 'R' { // I_Receive
		elem.(SRH).I_Receive(msg, info, elemInfo)
	} else { // "I_Send"
		elem.(SRH).I_Send(msg, info, elemInfo)
	}
}

func (e SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port

	// create listener if necessary
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "localhost"
	key := host + ":" + port

	//var err error
	if LnSRH == nil { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", key)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		LnSRH, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
	}

	//conn := e.Conns[key]
	c1 := make(chan []byte)
	c2 := make(chan []byte)

	go acceptread(&ConnSRH,LnSRH, c1)
	if ConnSRH != nil {
		go read(ConnSRH, c2)
	}

	select {
	case msgTemp1 := <- c1:
		*msg = messages.SAMessage{Payload: msgTemp1}
	case msgTemp2 := <- c2:
		*msg = messages.SAMessage{Payload: msgTemp2}
	}
}

func acceptread(conn *net.Conn,ln net.Listener, c chan []byte){
	// accept connections
	var err error

	*conn, err = ln.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// receive size & message
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	x := *conn
	_, err = x.Read(size)
	if err == io.EOF {
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		log.Fatalf("SRH::: %s", err)
	}

	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = x.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	c <- msgTemp
}

func read(conn net.Conn,  c chan []byte){
	// accept connections
	var err error

		// receive size & message
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		x := conn
		_, err = x.Read(size)
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil && err != io.EOF {
			log.Fatalf("SRH::: %s", err)
		}

		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = x.Read(msgTemp)
		if err != nil {
			log.Fatalf("SRH:: %s", err)
		}

	   c <- msgTemp
}

func (e SRH) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := ConnSRH.Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = ConnSRH.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}
