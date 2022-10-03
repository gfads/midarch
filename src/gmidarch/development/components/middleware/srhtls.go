package middleware

import (
	"crypto/tls"
	"encoding/binary"
	"gmidarch/development/messages"
	"gmidarch/development/messages/miop"
	"io"
	"log"
	"net"
	"shared"
	"shared/lib"
	"strings"
	"time"
)

//@Type: SRHTLS
//@Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHTLS struct{}

func (s SRHTLS) availableConnectionFromPool(clientsPtr *[]*messages.Client, ip string) (bool, int) {
	clients := *clientsPtr

	if ip != "" {
		for idx, client := range clients {
			if client.Ip == ip {
				return true, idx
			}
			if client.UDPConnection != nil {
				client.UDPConnection = nil
				return true, idx
			}
		}
	}

	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(clients) < 2 { //shared.MAX_NUMBER_OF_CONNECTIONS { TODO: dcruzb go back the env var
		client := messages.Client{
			Ip:            "",
			Connection:    nil,
			UDPConnection: nil,
		}
		*clientsPtr = append(clients, &client)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(*clientsPtr) - 1
	}

	for idx, client := range clients {
		if client == nil {
			client := messages.Client{
				Ip:            "",
				Connection:    nil,
				UDPConnection: nil,
			}
			clients[idx] = &client
			return true, idx
		}
		if client.UDPConnection != nil {
			client.UDPConnection.Close()
			client.UDPConnection = nil
			return true, idx
		}
	}

	return false, -1
}

func (s SRHTLS) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version 2 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

	// check if a listen has been already created
	if srhInfo.Ln == nil { // no listen created
		servAddr, err := net.ResolveTCPAddr("tcp", srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		srhInfo.Ln, err = tls.Listen("tcp4", servAddr.String(), getServerTLSConfig()) //net.ListenTCP("tcp", servAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients, "")
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHTLS Version 2 adapted - No connection available")
		time.Sleep(1 * time.Millisecond)
		return
	}

	go func() {
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex)

		// Accept connections
		conn, err := srhInfo.Ln.Accept()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		srhInfo.Conns = append(srhInfo.Conns, conn)
		//srhInfo.CurrentConn = conn

		lib.PrintlnDebug("SRHTLS Version 2 adapted >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Buscou nova conexão, ip:", conn.RemoteAddr().String())
		//connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients, conn.RemoteAddr().String())
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(
		//srhInfo.Clients))
		//if !connectionAvailable {
		//	lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHTLS Version 2 adapted - No connection available")
		//	return
		//}
		if len(srhInfo.Clients) <= availableConenctionIndex {
			lib.PrintlnDebug("SRHTLS Got len(srhInfo.Clients) <= availableConenctionIndex")
			*reset = true
			return
		}
		client := srhInfo.Clients[availableConenctionIndex]
		client.Ip = conn.RemoteAddr().String()
		client.Connection = conn
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)

		// Update info
		*info = srhInfo

		// Start goroutine
		go s.handler(info, availableConenctionIndex)
	}()
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
	return
}

func (s SRHTLS) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version Not adapted")
	//lib.PrintlnDebug(shared.GetFunction(), "HERE")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)

	select {
	case tempMsgReceived := <-srhInfo.RcvedMessages:
		{
			// Receive message from handlers
			//srhInfo.CurrentConn = tempMsgReceived.Chn

			// Update info
			*info = srhInfo
			msg.Payload = tempMsgReceived.Msg
			lib.PrintlnDebug("SRHTLS Version 2 adapted: tempMsgReceived", tempMsgReceived)
			lib.PrintlnDebug("SRHTLS Version 2 adapted: tempMsgReceived.Chn", tempMsgReceived.Chn)
			if tempMsgReceived.Chn == nil {
				*reset = true
				return
			}
			if isNewConnection, miopPacket := s.isNewConnection(tempMsgReceived.Msg); isNewConnection { // TODO dcruzb: move to I_Receive
				lib.PrintlnDebug("SRHTLS Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", miopPacket)
				*reset = true
				return
			}
			msg.ToAddr = tempMsgReceived.ToAddress //Chn.RemoteAddr().String()
		}
	default:
		{
			*reset = true
			return
		}
	}

	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
	return
}

func (s SRHTLS) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	lib.PrintlnDebug("msg.ToAddr", msg.ToAddr, "srhInfo.Clients", srhInfo.Clients)
	client := srhInfo.GetClientFromAddr(msg.ToAddr, srhInfo.Clients)
	conn := client.Connection //srhInfo.CurrentConn
	if conn == nil {
		*reset = true
		return
	}
	lib.PrintlnDebug("SRHTLS Version 2 adapted   >>>>> TCP => msg.ToAddr:", msg.ToAddr, "TCP conn:", conn, "AdaptId:", client.AdaptId)
	msgTemp := msg.Payload.([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := conn.Write(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	//json := Jsonmarshaller{}
	//unmarshalledMsg := json.Unmarshall(msgTemp)
	//log.Println("<<<<<<<<<<<<  <<<<<<<<<<  <<<<<<<<<  SRHTLS Version Not adapted => Msg: ", unmarshalledMsg.Bd.RepBody.OperationResult)
	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
	return
}

func (s SRHTLS) handler(info *interface{}, connectionIndex int) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHTLS Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	conn := srhInfo.Clients[connectionIndex].Connection //CurrentConn
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR", "SRHTLS Version Not adapted")
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		_, err := conn.Read(size)
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Read efetuado", "SRHTLS Version Not adapted")
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				srhInfo.Clients[connectionIndex] = nil
				lib.PrintlnError("Não Vai matar o app EOF")
				break
			} else if err != nil && err != io.EOF {
				lib.PrintlnError("Vai matar o app, erro mas não EOF")
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		// receive message
		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = conn.Read(msgTemp)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				srhInfo.Clients[connectionIndex] = nil
				lib.PrintlnError("Não Vai matar o app EOF")
				break
			} else if err != nil && err != io.EOF {
				lib.PrintlnError("Vai matar o app, erro mas não EOF")
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		if changeProtocol, miopPacket := s.isAdapt(msgTemp); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHTLS Version Not adapted")
				break
			}
		}
		if isNewConnection, _ := s.isNewConnection(msgTemp); isNewConnection { // TODO dcruzb: move to I_Receive
			//newConnection = true
			lib.PrintlnDebug("TCP Is New Connection")
			//miopPacket := miop.CreateReqPacket("Connect", []interface{}{miopPacket.Bd.ReqBody.Body[0], "Ok"}, miopPacket.Bd.ReqBody.Body[0].(int)) // idx is the Connection ID
			//msgPayload := Jsonmarshaller{}.Marshall(miopPacket)

			lib.PrintlnDebug("TCP Before send")
			//s.send(conn, addr, msgPayload)
			lib.PrintlnDebug("TCP After send")
			//if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
			//	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
			//	break
			//}
			continue
		}

		rcvMessage := messages.ReceivedMessages{Msg: msgTemp, Chn: conn, ToAddress: srhInfo.Clients[connectionIndex].Ip}
		lib.PrintlnDebug("SRHTLS Version 2 adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")
		if !*executeForever {
			break
		}
		srhInfo.RcvedMessages <- rcvMessage
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHTLS Version Not adapted")
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHTLS Version Not adapted")
}

func (s SRHTLS) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (s SRHTLS) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTLS Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}

func getServerTLSConfig() *tls.Config {
	if shared.CRT_PATH == "" {
		log.Fatal("SRHSsl:: Error:: Environment variable 'CRT_PATH' not configured\n")
	}

	if shared.KEY_PATH == "" {
		log.Fatal("SRHSsl:: Error:: Environment variable 'KEY_PATH' not configured\n")
	}

	cert, err := tls.LoadX509KeyPair(shared.CRT_PATH, shared.KEY_PATH)
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2"},
	}
	return tlsConfig
}
