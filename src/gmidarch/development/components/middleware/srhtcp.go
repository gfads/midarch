package middleware

import (
	"encoding/binary"
	"fmt"
	"gmidarch/development/messages"
	"io"
	"log"
	"net"
	"shared"
)

//@Type: SRHTCP
//@Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHTCP struct{}

func (s SRHTCP) availableConnectionFromPool(clientsPtr *[]*messages.Client) (bool, int) {
	clients := *clientsPtr
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(clients) < shared.MAX_NUMBER_OF_CONNECTIONS {
		client := messages.Client{
			Ip:         "",
			Connection: nil,
		}
		*clientsPtr = append(clients, &client)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(*clientsPtr) -1
	}

	for idx, client := range clients {
		if client == nil {
			client := messages.Client{
				Ip:         "",
				Connection: nil,
			}
			clients[idx] = &client
			return true, idx
		}
	}

	return false, -1
}

func (s SRHTCP) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version Not adapted")
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
		srhInfo.Ln, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients)
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//log.Println("------------------------------>", shared.GetFunction(), "end", "SRHTCP Version Not adapted - No connection available")
		return
	}

	go func() {
		client := srhInfo.Clients[availableConenctionIndex]
		log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex)

		// Accept connections
		conn, err := srhInfo.Ln.Accept()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		srhInfo.Conns = append(srhInfo.Conns, conn)
		//srhInfo.CurrentConn = conn

		client.Ip = conn.RemoteAddr().String()
		client.Connection = conn
		log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)


		// Update info
		*info = srhInfo

		// Start goroutine
		go handler(info, availableConenctionIndex)
	}()
	fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version Not adapted")
	return
}

func (s SRHTCP) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version Not adapted")
	//fmt.Println(shared.GetFunction(), "HERE")
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
			msg.ToAddr = tempMsgReceived.Chn.RemoteAddr().String()
		}
	default:
		{
			*reset = true
			return
		}
	}

	fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version Not adapted")
	return
}

func (s SRHTCP) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	log.Println("msg.ToAddr", msg.ToAddr, "srhInfo.Clients", srhInfo.Clients)
	conn := srhInfo.GetClientFromAddr(msg.ToAddr, srhInfo.Clients).Connection //srhInfo.CurrentConn
	log.Println("conn:", conn)
	msgTemp := msg.Payload.([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgTemp)))
	_, err := conn.Write(size)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	json := Jsonmarshaller{}
	unmarshalledMsg := json.Unmarshall(msgTemp)
	log.Println("<<<<<<<<<<<<  <<<<<<<<<<  <<<<<<<<<  SRHTCP Version Not adapted => Msg: ", unmarshalledMsg.Bd.RepBody.OperationResult)
	// send message
	_, err = conn.Write(msgTemp)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// update info
	*info = srhInfo
	fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version Not adapted")
	return
}

func handler(info *interface{}, connectionIndex int) {
	fmt.Println("----------------------------------------->", shared.GetFunction(), "SRHTCP Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	conn := srhInfo.Clients[connectionIndex].Connection //CurrentConn

	for {
		fmt.Println("----------------------------------------->", shared.GetFunction(), "FOR", "SRHTCP Version Not adapted")
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		_, err := conn.Read(size)
		fmt.Println("----------------------------------------->", shared.GetFunction(), "Read efetuado", "SRHTCP Version Not adapted")
		if err == io.EOF {
			srhInfo.Clients[connectionIndex] = nil
			fmt.Println("Não Vai matar o app EOF")
			break
		} else if err != nil && err != io.EOF {
			fmt.Println("Vai matar o app, erro mas não EOF")
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// receive message
		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = conn.Read(msgTemp)
		if err == io.EOF {
			srhInfo.Clients[connectionIndex] = nil
			fmt.Println("Não Vai matar o app EOF")
			break
		} else if err != nil && err != io.EOF {
			fmt.Println("Vai matar o app, erro mas não EOF")
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		rcvMessage := messages.ReceivedMessages{Msg: msgTemp, Chn: conn}

		srhInfo.RcvedMessages <- rcvMessage
		fmt.Println("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHTCP Version Not adapted")
	}
	fmt.Println("----------------------------------------->", shared.GetFunction(), "end", "SRHTCP Version Not adapted")
}
