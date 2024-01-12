package middleware

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

// @Type: SRHUDP
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHUDP struct{}

func (s SRHUDP) availableConnectionFromPool(clientsPtr *[]*messages.Client) (bool, int) {
	clients := *clientsPtr
	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(clients) < 1 { //shared.MAX_NUMBER_OF_CONNECTIONS { // UDP don't open different connections
		client := messages.Client{
			Ip:            "",
			Connection:    nil,
			UDPConnection: nil,
		}
		*clientsPtr = append(clients, &client)
		lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		lib.PrintlnDebug("len(clients) < 1")
		return true, len(*clientsPtr) - 1
	}

	for idx, client := range clients {
		if idx >= 1 { //shared.MAX_NUMBER_OF_CONNECTIONS
			lib.PrintlnDebug("idx >= 1")
			break
		}
		if client == nil {
			lib.PrintlnDebug("client == nil")
			client := messages.Client{
				Ip:            "",
				Connection:    nil,
				UDPConnection: nil,
			}
			clients[idx] = &client
			return true, idx
		}
		if client.Connection != nil {
			lib.PrintlnDebug("Zerou Connection")
			client.Ip = ""
			client.Connection.Close()
			client.Connection = nil
			return true, idx
		}
		if client.UDPConnection == nil {
			lib.PrintlnDebug("UDPConnection == nil")
			return true, idx
		}
	}

	return false, -1
}

func (s SRHUDP) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// srhInfo.Counter++
	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//lib.PrintlnDebug("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

	var servAddr *net.UDPAddr
	var err error
	// check if a listen has been already created
	//if srhInfo.UDPConnection == nil { // no listen created
	//	servAddr, err = net.ResolveUDPAddr("udp4", srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port)
	//	if err != nil {
	//		shared.ErrorHandler(shared.GetFunction(), err.Error())
	//	}
	//	//srhInfo.Ln, err = net.ListenUDP("udp4", servAddr)
	//	//if err != nil {
	//	//	shared.ErrorHandler(shared.GetFunction(), err.Error())
	//	//}
	//	srhInfo.UDPConnection, err = net.ListenUDP("udp4", servAddr) //, err := srhInfo.Ln.Accept()
	//	if err != nil {
	//		shared.ErrorHandler(shared.GetFunction(), err.Error())
	//	}
	//}

	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients)
	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		time.Sleep(1 * time.Millisecond)
		//log.Println("------------------------------>", shared.GetFunction(), "end", "SRHUDP Version Not adapted - No connection available")
		return
	}

	//go func() {
	client := srhInfo.Clients[availableConenctionIndex]
	lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex)
	lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients", srhInfo.Clients, srhInfo.Clients[availableConenctionIndex])

	// Accept connections
	//conn, err := net.ListenUDP("udp4", servAddr) //, err := srhInfo.Ln.Accept()
	//if err != nil {
	//	shared.ErrorHandler(shared.GetFunction(), err.Error())
	//}
	if client == nil {
		*reset = true
		return
	}

	if client.UDPConnection == nil { // no listen created
		servAddr, err = net.ResolveUDPAddr("udp4", srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		//srhInfo.Ln, err = net.ListenUDP("udp4", servAddr)
		//if err != nil {
		//	shared.ErrorHandler(shared.GetFunction(), err.Error())
		//}
		client.UDPConnection, err = net.ListenUDP("udp4", servAddr) //, err := srhInfo.Ln.Accept()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	srhInfo.Conns = append(srhInfo.Conns, client.UDPConnection)
	//srhInfo.CurrentConn = conn

	client.Ip = "" //conn.RemoteAddr().String() UDP dont start with RemoteAddr
	//client.UDPConnection = srhInfo.UDPConnection
	lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)

	// Update info
	*info = srhInfo

	// Start goroutine
	go s.handler(info, availableConenctionIndex)
	//}()
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
	return
}

func (s SRHUDP) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version Not adapted")
	//lib.PrintlnDebug(shared.GetFunction(), "HERE")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)

	select {
	case tempMsgReceived := <-srhInfo.RcvedMessages:
		{
			// Receive message from handlers
			//srhInfo.CurrentConn = tempMsgReceived.Conn

			// Update info
			*info = srhInfo
			msg.Payload = tempMsgReceived.Msg
			lib.PrintlnDebug("SRHUDP Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", tempMsgReceived)
			//if isNewConnection, miopPacket := s.isNewConnection(tempMsgReceived.Msg); isNewConnection { // TODO dcruzb: move to I_Receive
			//	fmt.Println("SRHUDP Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", miopPacket)
			//	*reset = true
			//	return
			//}
			msg.ToAddr = tempMsgReceived.ToAddress
		}
	default:
		{
			*reset = true
			return
		}
	}

	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
	return
}

func (s SRHUDP) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	lib.PrintlnDebug("msg.ToAddr", msg.ToAddr, "srhInfo.Clients", srhInfo.Clients)
	client := srhInfo.GetClientFromAddr(msg.ToAddr, srhInfo.Clients)
	conn := client.UDPConnection //srhInfo.CurrentConn
	if conn == nil {
		fmt.Println("SRHUDP.send reset when conn == nil")
		*reset = true
		return
	}
	lib.PrintlnDebug("UDP conn:", conn)
	msgTemp := msg.Payload.([]byte)
	addr := strings.Split(msg.ToAddr, ":")
	ip := net.ParseIP(addr[0])
	port, _ := strconv.Atoi(addr[1])

	s.send(conn, &net.UDPAddr{IP: ip, Port: port}, msgTemp)

	// update info
	*info = srhInfo
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
	return
}

func (s SRHUDP) send(conn *net.UDPConn, addr *net.UDPAddr, msgTemp []byte) {
	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := conn.WriteTo(size, addr)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	//json := Jsonmarshaller{}
	//unmarshalledMsg := json.Unmarshall(msgTemp)
	//lib.PrintlnDebug("<<<<<<<<<<<<  <<<<<<<<<<  <<<<<<<<<  SRHUDP Version Not adapted => Msg: ", unmarshalledMsg.Bd.RepBody.OperationResult)
	// send message
	_, err = conn.WriteTo(msgTemp, addr)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
}

func (s *SRHUDP) handler(info *interface{}, connectionIndex int) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHUDP Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	conn := srhInfo.Clients[connectionIndex].UDPConnection //CurrentConn
	executeForever := srhInfo.ExecuteForever
	for {
		if !*executeForever {
			break
		}
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR", "SRHUDP Version Not adapted")
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		//err := conn.SetReadDeadline(time.Now().Add(10*time.Second))
		//if err != nil {
		//	lib.PrintlnError(shared.GetFunction(), err.Error())
		//}
		_, addr, err := conn.ReadFromUDP(size)
		if len(srhInfo.Clients) == 0 || srhInfo.Clients[connectionIndex] == nil {
			lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Handler without client")
			break
		}
		// lib.PrintlnInfo("Received:", size)
		srhInfo.Clients[connectionIndex].Ip = addr.String()
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Read efetuado", "SRHUDP Version Not adapted")
		if err != nil {
			if err == io.EOF ||
				strings.Contains(err.Error(), "use of closed network connection") ||
				strings.Contains(err.Error(), "i/o timeout") {
				srhInfo.Clients[connectionIndex] = nil
				lib.PrintlnDebug("Não Vai matar o app EOF", err)
				break
			} else if err != nil && err != io.EOF {
				lib.PrintlnDebug("Vai matar o app, erro mas não EOF. Error:", err)
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		// receive message
		// msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		//err = conn.SetReadDeadline(time.Now().Add(10*time.Second))
		//if err != nil {
		//	lib.PrintlnError(shared.GetFunction(), err.Error())
		//}
		// _, err = conn.Read(msgTemp)
		// lib.PrintlnInfo("Received:", msgTemp)
		msgTemp, err := s.read(binary.LittleEndian.Uint32(size), conn)
		if err != nil {
			if err == io.EOF ||
				strings.Contains(err.Error(), "use of closed network connection") ||
				strings.Contains(err.Error(), "i/o timeout") {
				srhInfo.Clients[connectionIndex] = nil
				lib.PrintlnDebug("Não Vai matar o app EOF")
				break
			} else if err != nil && err != io.EOF {
				lib.PrintlnDebug("Vai matar o app, erro mas não EOF")
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		if changeProtocol, miopPacket := s.isAdapt(msgTemp); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
				break
			}
		}
		if isNewConnection, miopPacket := s.isNewConnection(msgTemp); isNewConnection { // TODO dcruzb: move to I_Receive
			lib.PrintlnDebug("Is New Connection")
			miopPacket := miop.CreateReqPacket("Connect", []interface{}{miopPacket.Bd.ReqBody.Body[0], "Ok"}, miopPacket.Bd.ReqBody.Body[0].(int)) // idx is the Connection ID
			msgPayload := Jsonmarshaller{}.Marshall(miopPacket)

			lib.PrintlnDebug("Before send")
			s.send(conn, addr, msgPayload)
			lib.PrintlnDebug("After send")
			//if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
			//	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
			//	break
			//}
			continue
		}
		lib.PrintlnDebug("SRHUDP Version Not adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")
		rcvMessage := messages.ReceivedMessages{Msg: msgTemp, Conn: nil, ToAddress: addr.String()}
		if !*executeForever {
			break
		}
		srhInfo.RcvedMessages <- rcvMessage
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHUDP Version Not adapted")
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version Not adapted")
}

func (s SRHUDP) read(size uint32, conn *net.UDPConn) (fullMessage []byte, err error) {
	// msgTemp := make([]byte, size)
	//err = conn.SetReadDeadline(time.Now().Add(10*time.Second))
	//if err != nil {
	//	lib.PrintlnError(shared.GetFunction(), err.Error())
	//}
	// _, err = conn.Read(msgTemp)
	// fullMessage = make([]byte, size)
	// lib.PrintlnInfo("Received(read):size", size)
	const maxBufferSize = shared.MAX_UDP_PACKET_SIZE
	for {
		bufferSize := int(size) - len(fullMessage)
		if bufferSize > maxBufferSize {
			bufferSize = maxBufferSize
		}
		buffer := make([]byte, bufferSize, bufferSize)
		// lib.PrintlnInfo("Received(read-ini):size", size, "len(fullMessage)", len(fullMessage), "bufferSize", bufferSize, "remaining", int(size)-len(fullMessage))

		// lib.PrintlnInfo("Received(read):for1")
		n, _, err := conn.ReadFromUDP(buffer)
		// lib.PrintlnInfo("Received(read):", buffer)

		if err != nil {
			lib.PrintlnError("Error while reading message. Error:", err)
			return nil, err
		}

		fullMessage = append(fullMessage, buffer[:n]...)
		// lib.PrintlnInfo("Received(read):for2")
		// lib.PrintlnInfo("Received(read-end):size", size, "len(fullMessage)", len(fullMessage), "bufferSize", bufferSize)
		// Check if the message is complete (you need a way to determine this based on your protocol)
		if len(fullMessage) >= int(size) {
			return fullMessage, nil
		}
		// lib.PrintlnInfo("Received(read):for3")
	}
}

func (s SRHUDP) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (s SRHUDP) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}
