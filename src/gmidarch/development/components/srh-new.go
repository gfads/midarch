package components

/*
import (
	"encoding/binary"
	"fmt"
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
	//	Conn      net.Conn
	//	Ln        net.Listener
}

var ConnSRH net.Conn
var LnSRH net.Listener

func NewSRH() SRH {

	r := new(SRH)

	//r.Host = "localhost"           // TODO
	//r.Port = shared.FIBONACCI_PORT // TODO
	//r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B [] I_Accept -> B"
	r.Behaviour = "B = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B [] I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"
	//r.Conns = make(map[string]net.Conn, shared.NUM_MAX_CONNECTIONS)
	//r.Lns = make(map[string]net.Listener, shared.NUM_MAX_CONNECTIONS)
	//r.Conn = *new(net.Conn)
	//r.Ln = *new(net.Listener)

	return *r
}

func (e SRH) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}) {
	switch op[2] {

	case 'R':  // I_Receive
		elem.(SRH).I_Receive(msg, info, elemInfo)
	case 'S':  // "I_Send"
		elem.(SRH).I_Send(msg, info, elemInfo)
	case 'A':
		elem.(SRH).I_Accept(msg, info, elemInfo)
	}
}

func (e SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) {

	fmt.Printf("SRH:: I_Receive\n")
	if ConnSRH != nil {
		fmt.Printf("SRH:: I_Receive:: (%v,%v)\n", ConnSRH.LocalAddr(), ConnSRH.RemoteAddr())

		// receive size & message
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		_, err := ConnSRH.Read(size)
		if err == io.EOF {
			os.Exit(0)
		} else if err != nil && err != io.EOF {
			log.Fatalf("SRH::: %s", err)
		}

		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = ConnSRH.Read(msgTemp)
		if err != nil {
			log.Fatalf("SRH:: %s", err)
		}
		*msg = messages.SAMessage{Payload: msgTemp}
	} else {
		*msg = messages.SAMessage{Payload: *msg}
	}
}

func (e SRH) I_Accept(msg *messages.SAMessage, info [] *interface{}, elemInfo [] *interface{}) { // TODO Host & Port

	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "localhost"

	var err error
	if LnSRH == nil { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		LnSRH, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
	}

	// accept connections
	ConnSRH, err = LnSRH.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	fmt.Printf("SRH:: Accept:: Connection Accepted:: %v %v\n", ConnSRH.LocalAddr(), ConnSRH.RemoteAddr())
	//}

	*msg = messages.SAMessage{Payload: *msg} // TODO
}

func (e SRH) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	fmt.Printf("SRH:: I_Send:: (%v,%v)\n", ConnSRH.LocalAddr(), ConnSRH.RemoteAddr())

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
*/