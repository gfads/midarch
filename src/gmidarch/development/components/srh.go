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
	Host      string
	Port      string
	Conns     map[string]net.Conn
	Lns       map[string]net.Listener
}

func NewSRH() SRH {

	r := new(SRH)

	r.Host = "localhost"           // TODO
	r.Port = shared.FIBONACCI_PORT // TODO
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"
	r.Conns = make(map[string]net.Conn, shared.NUM_MAX_CONNECTIONS)
	r.Lns = make(map[string]net.Listener, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (SRH) Selector(elem interface{}, op string) func(*messages.SAMessage, []*interface{}) {

	if op == "I_Receive" {
		return func(msg *messages.SAMessage, info []*interface{}) {
			elem.(SRH).I_Receive(msg, info)
		}
	} else { // "I_Send"
		return func(msg *messages.SAMessage, info []*interface{}) {
			elem.(SRH).I_Send(msg, info)
		}
	}
}

func (s SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}) { // TODO Host & Port

	// create listener if necessary
	key := s.Host + s.Port
	if _, ok := s.Lns[key]; !ok { // listener was not created yet
		servAddr, err := net.ResolveTCPAddr("tcp", s.Host+":"+s.Port)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}
		s.Lns[key], err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatalf("SRH:: %v\n", err)
		}

		// accept connections
		s.Conns[key], err = s.Lns[key].Accept()
		if err != nil {
			log.Fatalf("SRH:: %s", err)
		}
	}

	// configure conn to be used
	conn := s.Conns[key]

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

func (s SRH) I_Send(msg *messages.SAMessage, info [] *interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// configure conn to be used
	key := s.Host + s.Port
	conn := s.Conns[key]

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
