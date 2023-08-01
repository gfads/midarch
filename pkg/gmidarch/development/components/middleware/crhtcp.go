package middleware

import (
	"encoding/binary"
	"net"
	"reflect"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: CRHTCP
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.e1 -> Behaviour
type CRHTCP struct{}

func (c CRHTCP) getLocalTcpAddr() *net.TCPAddr {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	//fmt.Println("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
	lib.PrintlnDebug("github.com/gfads/midarch/src/shared.LocalAddr:", shared.LocalAddr)
	var err error = nil
	var localTCPAddr *net.TCPAddr = nil
	//shared.LocalAddr = "127.0.0.1:37521"
	if shared.LocalAddr != "" {
		localTCPAddr, err = net.ResolveTCPAddr("tcp", shared.LocalAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
	return localTCPAddr
}

func (c CRHTCP) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	infoTemp := *info
	crhInfo := infoTemp.(messages.CRHInfo)

	// check message
	//payload := msg.Payload.([]byte)
	payload := msg.Payload.(messages.RequestorInfo).MarshalledMessage
	h := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Host
	p := msg.Payload.(messages.RequestorInfo).Inv.Endpoint.Port

	host := ""
	port := ""

	if h == "" || p == "" {
		host = crhInfo.EndPoint.Host
		port = crhInfo.EndPoint.Port
	} else {
		host = h
		port = p
	}

	msgToServer := payload

	addr := host + ":" + port
	var err error
	//fmt.Println("Vai conectar", crhInfo.Conns[addr])
	lib.PrintlnDebug("Vai conectar", crhInfo.Conns[addr])
	if _, ok := crhInfo.Conns[addr]; !ok || reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name() != "TCPConn" { // no connection open yet
		//fmt.Println("Entrou", crhInfo.Conns[addr])
		lib.PrintlnDebug("Entrou", crhInfo.Conns[addr])
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		//log.Println("Resolveu", crhInfo.Conns[addr])
		//localTcpAddr := c.getLocalTcpAddr()

		for {
			crhInfo.Conns[addr], err = net.DialTCP("tcp", nil, tcpAddr)
			//log.Println("Dialed", crhInfo.Conns[addr])
			if err != nil {
				lib.PrintlnError("Erro na discagem", crhInfo.Conns[addr], err)
				time.Sleep(200 * time.Millisecond)
				//shared.ErrorHandler(shared.GetFunction(), err.Error())
			} else {
				break
			}
		}
		if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
			//fmt.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr())
			//log.Println("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
			shared.LocalAddr = crhInfo.Conns[addr].LocalAddr().String()
		}
	}
	//fmt.Println("Terminou", crhInfo.Conns[addr])
	lib.PrintlnDebug("Terminou", crhInfo.Conns[addr])

	// send message's size
	conn := crhInfo.Conns[addr]
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	err = c.send(sizeOfMsgSize, msgToServer, conn)
	if err != nil {
		lib.PrintlnError("Error trying to send message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Conns[addr].Close()
		crhInfo.Conns[addr] = nil
		delete(crhInfo.Conns, addr)
		return
	}

	msgFromServer, err := c.read(conn, sizeOfMsgSize)
	if err != nil {
		lib.PrintlnError("Error trying to read message:", err.Error())
		*msg = messages.SAMessage{Payload: nil} // TODO dcruzb: adjust message
		crhInfo.Conns[addr].Close()
		crhInfo.Conns[addr] = nil
		delete(crhInfo.Conns, addr)
		return
	}
	VerifyAdaptation(msgFromServer, sizeOfMsgSize, conn, c.send)

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHTCP) send(sizeOfMsgSize []byte, msgToServer []byte, conn net.Conn) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := conn.Write(sizeOfMsgSize)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}

	// send message
	_, err = conn.Write(msgToServer)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}
	return nil
}

func (c CRHTCP) read(conn net.Conn, size []byte) ([]byte, error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	// receive reply's size
	_, err := conn.Read(size)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	return msgFromServer, nil
}
