package components

import (
	"encoding/binary"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"log"
	"net"
	"shared"
	"strconv"
)

type SRH struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Host      string
	Port      int
	Conns     map[string]net.Conn
	Lns       map[string] net.Listener
}

//var lnSRH net.Listener
//var connSRH net.Conn
//var firstListenerSRH bool

/*
func NewSRHOld() SRH {

	// create a new instance of Server
	r := new(SRH)

	// configure the new instance
	r.Host = "localhost"                  // TODO
	r.Port = shared.FIBONACCI_PORT        // TODO
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"

	firstListenerSRH = true

	return *r
}
*/

func NewSRH() SRH {

	// create a new instance of Server
	r := new(SRH)

	// configure the new instance
	r.Host = "localhost"                  // TODO
	r.Port = shared.FIBONACCI_PORT        // TODO
	r.Behaviour = "B = I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> B"
	r.Conns = make(map[string]net.Conn)
	r.Lns = make(map[string]net.Listener)

	return *r
}

func (s SRH) I_Receive(msg *messages.SAMessage, info [] *interface{}) { // TODO

	host := s.Host // TODO
	port := s.Port // TODO

	// create listener
	err := *new(error)
	key := host+strconv.Itoa(port)
	if _,ok := s.Lns[key]; !ok {
		s.Lns[key], err = net.Listen("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			log.Fatalf("SRH:: %v\n",err)
		}

		// accept connections
		s.Conns[key], err = s.Lns[key].Accept()
		if err != nil {
			log.Fatalf("SRH:: %s", err)
		}
	}

	// receive message's size
	size := make([]byte, 4)
	_, err = s.Conns[key].Read(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	sizeInt := binary.LittleEndian.Uint32(size)

	// receive message
	msgTemp := make([]byte, sizeInt)
	_, err = s.Conns[key].Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgTemp} // TODO
}

/*
func (s SRH) I_ReceiveOld(msg *messages.SAMessage, info [] *interface{}) { // TODO

	host := s.Host            // TODO
	port := s.Port // TODO

	// create listener
	err := *new(error)
	if firstListenerSRH {
		firstListenerSRH = false
		lnSRH, err = net.Listen("tcp", host+":"+strconv.Itoa(port))
		if err != nil {
			log.Fatalf("SRH:: %v\n",err)
		}
	}

	// accept connections
	connSRH, err = lnSRH.Accept()
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// receive message's size
	size := make([]byte, 4)
	_, err = connSRH.Read(size)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
	sizeInt := binary.LittleEndian.Uint32(size)

	// receive message
	msgTemp := make([]byte, sizeInt)
	_, err = connSRH.Read(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	*msg = messages.SAMessage{Payload: msgTemp} // TODO
}
*/

func (s SRH) I_Send(msg *messages.SAMessage, info [] *interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	key := s.Host+strconv.Itoa(s.Port)
	size := make([]byte, 4)
	l := uint32(len(msgTemp))
	binary.LittleEndian.PutUint32(size, l)
	s.Conns[key].Write(size)

	err := *new(error)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = s.Conns[key].Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}

/*
func (SRH) I_SendOld(msg *messages.SAMessage, info [] *interface{}) {
	msgTemp := msg.Payload.([]interface{})[0].([]byte)

	// send message's size
	size := make([]byte, 4)
	l := uint32(len(msgTemp))
	binary.LittleEndian.PutUint32(size, l)
	connSRH.Write(size)

	err := *new(error)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}

	// send message
	_, err = connSRH.Write(msgTemp)
	if err != nil {
		log.Fatalf("SRH:: %s", err)
	}
}
 */