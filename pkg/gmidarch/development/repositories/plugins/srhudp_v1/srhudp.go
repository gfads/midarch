package srhudp

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/gfads/midarch/pkg/gmidarch/development/components/middleware"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
)

// @Type: SRHUDP
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHUDP struct{}

func (s SRHUDP) availableConnectionFromPool(clientsPtr *[]*messages.Client) (bool, int) {
	clients := *clientsPtr

	//if len(clients) > 1 {
	//	clients = clients[0:0]
	//}
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(clients) < 1 { //shared.MAX_NUMBER_OF_CONNECTIONS { // UDP don't open different connections
		//fmt.Println("Nenhum cliente")
		client := messages.Client{
			Ip:            "",
			Connection:    nil,
			UDPConnection: nil,
		}
		*clientsPtr = append(clients, &client)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(*clientsPtr) - 1
	}

	//fmt.Println("Chegou Aqui for idx, client := range clients")
	for idx, client := range clients {
		if idx >= 1 { //shared.MAX_NUMBER_OF_CONNECTIONS
			break
		}
		if client == nil {
			//fmt.Println("Zerou client")
			client := messages.Client{
				Ip:            "",
				Connection:    nil,
				UDPConnection: nil,
			}
			clients[idx] = &client
			return true, idx
		}
		if client.Connection != nil {
			//fmt.Println("Zerou Connection")
			client.Ip = ""
			client.Connection.Close()
			client.Connection = nil
			return true, idx
		}
		if client.UDPConnection == nil {
			return true, idx
		}
	}

	return false, -1
}

func (s SRHUDP) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHUDP Version 1 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	// srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

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

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients)
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//log.Println("------------------------------>", shared.GetFunction(), "end", "SRHUDP Version 1 adapted - No connection available")
		return
	}

	//go func() {
	client := srhInfo.Clients[availableConenctionIndex]
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex)

	// Accept connections
	//conn, err := net.ListenUDP("udp4", servAddr) //, err := srhInfo.Ln.Accept()
	//if err != nil {
	//	shared.ErrorHandler(shared.GetFunction(), err.Error())
	//}

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
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)

	// Update info
	*info = srhInfo

	// Start goroutine
	go s.handler(info, availableConenctionIndex)
	//}()
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version 1 adapted")
	return
}

func (s SRHUDP) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHUDP Version 1 adapted")
	//fmt.Println(shared.GetFunction(), "HERE")
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
			//fmt.Println("SRHUDP Version 1 adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", tempMsgReceived)
			msg.ToAddr = tempMsgReceived.ToAddress
		}
	default:
		{
			*reset = true
			return
		}
	}

	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version 1 adapted")
	return
}

func (s SRHUDP) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHUDP Version 1 adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	//fmt.Println("msg.ToAddr", msg.ToAddr, "srhInfo.Clients", srhInfo.Clients)
	client := srhInfo.GetClientFromAddr(msg.ToAddr, srhInfo.Clients)
	conn := client.UDPConnection //srhInfo.CurrentConn
	if conn == nil {
		*reset = true
		return
	}
	//fmt.Println("UDP conn:", conn)
	msgTemp := msg.Payload.([]byte)
	addr := strings.Split(msg.ToAddr, ":")
	ip := net.ParseIP(addr[0])
	port, _ := strconv.Atoi(addr[1])

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))

	//miop := middleware.Jsonmarshaller{}.Unmarshall(msgTemp)
	//return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
	//fmt.Println("SRHUDP_v1 Client:", client.Ip, client.Connection, client.UDPConnection, "Connection is:", conn, "msg.ToAddr is:", msg.ToAddr, "msgTemp is:", miop.Bd.ReqHeader.Operation, miop.Bd.RepBody)
	_, err := conn.WriteTo(size, &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	//json := middleware.Jsonmarshaller{}
	//unmarshalledMsg := json.Unmarshall(msgTemp)
	//fmt.Println("<<<<<<<<<<<<  <<<<<<<<<<  <<<<<<<<<  SRHUDP Version 1 adapted => Msg: ", unmarshalledMsg.Bd.RepBody.OperationResult)
	// send message
	_, err = conn.WriteTo(msgTemp, &net.UDPAddr{IP: ip, Port: port})
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version 1 adapted")
	return
}

func (s *SRHUDP) handler(info *interface{}, connectionIndex int) {
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHUDP Version 1 adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	conn := srhInfo.Clients[connectionIndex].UDPConnection //CurrentConn
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		//fmt.Println("----------------------------------------->", shared.GetFunction(), "FOR", "SRHUDP Version 1 adapted")
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		_, addr, err := conn.ReadFromUDP(size)
		srhInfo.Clients[connectionIndex].Ip = addr.String()
		//fmt.Println("----------------------------------------->", shared.GetFunction(), "Read efetuado", "SRHUDP Version 1 adapted")
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				srhInfo.Clients[connectionIndex] = nil
				//fmt.Println("Não Vai matar o app EOF")
				break
			} else if err != nil && err != io.EOF {
				//fmt.Println("Vai matar o app, erro mas não EOF")
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		// receive message
		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = conn.Read(msgTemp)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				srhInfo.Clients[connectionIndex] = nil
				//fmt.Println("Não Vai matar o app EOF")
				break
			} else if err != nil && err != io.EOF {
				//fmt.Println("Vai matar o app, erro mas não EOF")
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		if changeProtocol, miopPacket := s.isAdapt(msgTemp); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				//fmt.Println("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
				break
			}
		}
		//fmt.Println("SRHUDP Version 1 adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message from", addr)
		rcvMessage := messages.ReceivedMessages{Msg: msgTemp, Conn: nil, ToAddress: addr.String()}
		if !*executeForever {
			break
		}
		srhInfo.RcvedMessages <- rcvMessage
		//fmt.Println("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHUDP Version 1 adapted")
	}
	//fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHUDP Version 1 adapted")
}

func (s SRHUDP) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop := middleware.Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}
