package middleware

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	evolutive "github.com/gfads/midarch/pkg/injector"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
	"io/ioutil"
	"log"
	"net"
	"reflect"
	"time"
)

// @Type: CRHTLS
// @Behaviour: Behaviour = InvP.e1 -> I_Process -> TerP.eNot -> Behaviour
type CRHTLS struct{}

func (c CRHTLS) getLocalTcpAddr() *net.TCPAddr {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
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

func (c CRHTLS) I_Process(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
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
	if _, ok := crhInfo.Conns[addr]; !ok || reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name() != "Conn" { // no connection open yet
		//fmt.Println("Entrou", crhInfo.Conns[addr])
		lib.PrintlnDebug("Entrou", crhInfo.Conns[addr])
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		//log.Println("Resolveu", crhInfo.Conns[addr])
		//localTcpAddr := c.getLocalTcpAddr()

		for {
			crhInfo.Conns[addr], err = tls.Dial("tcp4", tcpAddr.String(), getClientTLSConfig())
			//log.Println("Discou", crhInfo.Conns[addr])
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
	if changeProtocol, miopPacket := c.isAdapt(msgFromServer); changeProtocol {
		lib.PrintlnDebug("Adapting, miopPacket.Bd.ReqBody.Body:", miopPacket.Bd.ReqBody.Body)

		shared.AdaptId = miopPacket.Bd.ReqBody.Body[1].(int)

		miopPacket := miop.CreateReqPacket("ChangeProtocol", []interface{}{miopPacket.Bd.ReqBody.Body[0], shared.AdaptId, "Ok"}, shared.AdaptId) // idx is the Connection ID
		msgPayload := Jsonmarshaller{}.Marshall(miopPacket)
		c.send(sizeOfMsgSize, msgPayload, conn)

		if miopPacket.Bd.ReqBody.Body[0] == "udp" {
			lib.PrintlnInfo("Adapting => UDP")
			evolutive.GeneratePlugin("crhudp_v1", "crhudp", "crhudp_v1")
		} else if miopPacket.Bd.ReqBody.Body[0] == "tcp" {
			lib.PrintlnInfo("Adapting => TCP")
			evolutive.GeneratePlugin("crhtcp_v1", "crhtcp", "crhtcp_v1")
		} else {
			msgFromServer, _ = c.read(conn, sizeOfMsgSize)
			//fmt.Println("=================> ############### ============> ########### TCP: Leu o read")
		}
	}

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func (c CRHTLS) send(sizeOfMsgSize []byte, msgToServer []byte, conn net.Conn) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
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

func (c CRHTLS) read(conn net.Conn, size []byte) ([]byte, error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
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

func (c CRHTLS) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func getClientTLSConfig() *tls.Config {
	if shared.CA_PATH == "" {
		log.Fatal("CRHSsl:: Error:: Environment variable 'CA_PATH' not configured\n")
	}
	trustCert, err := ioutil.ReadFile(shared.CA_PATH)
	if err != nil {
		fmt.Println("Error loading trust certificate. ", err)
	}
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(trustCert) {
		fmt.Println("Error installing trust certificate.")
	}

	// connect to server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		RootCAs:            certs,
		NextProtos:         []string{"h2"},
	}
	return tlsConfig
}