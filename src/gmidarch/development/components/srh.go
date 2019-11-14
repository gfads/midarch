package components

import (
	"encoding/binary"
	"github.com/vmihailenco/msgpack"
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
	//Host      string
	//Port      string
	Conns map[string]net.Conn
	Lns   map[string]net.Listener
}

func NewSRH() SRH {

	r := new(SRH)

	//r.Host = "localhost"           // TODO
	//r.Port = shared.FIBONACCI_PORT // TODO
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"
	r.Conns = make(map[string]net.Conn, shared.NUM_MAX_CONNECTIONS)
	r.Lns = make(map[string]net.Listener, shared.NUM_MAX_CONNECTIONS)

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
	if _, ok := e.Lns[key]; !ok { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", key)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		e.Lns[key], err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}

		// accept connections
		e.Conns[key], err = e.Lns[key].Accept()
		if err != nil {
			log.Fatalf("SRH:: %s", err)
		}
	}

	// configure conn to be used
	conn := e.Conns[key]

	// receive size & message
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	_, err := conn.Read(size)
	if err == io.EOF {
		os.Exit(0)
	} else if err != nil && err != io.EOF {
		log.Fatalf("SRH::: %s", err)
	}

	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
	_, err = conn.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgTemp} // TODO
}

func (e SRH) I_Send(msg *messages.SAMessage, info [] *interface{}, elemInfo []*interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// configure conn to be used
	tempPort := *elemInfo[0]
	port := tempPort.(string)
	host := "localhost"
	key := host + ":" + port
	conn := e.Conns[key]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := conn.Write(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}

// Used in *_*EncDec
var encSRH msgpack.Encoder
var decSRH msgpack.Decoder
