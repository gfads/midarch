package middleware

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/messages"
	"github.com/gfads/midarch/pkg/gmidarch/development/messages/miop"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
	"github.com/quic-go/quic-go"
)

// @Type: SRHRPC
// @Behaviour: Behaviour = I_Accept -> I_Receive -> InvR.e1 -> TerR.e1 -> I_Send -> Behaviour
type SRHRPC struct {
	// Graph exec.ExecGraph
}

// var ConnsSRHRPC []quic.Connection
// var StreamsSRHRPC []quic.Stream
// var LnSRHRPC quic.Listener

// var c1Quic = make(chan []byte)
// var c2Quic = make(chan []byte)

// var currentConnectionQuic = -1
// var stateQuic = 0

func (s SRHRPC) availableConnectionFromPool(clientsPtr *[]*messages.Client, ip string) (bool, int) {
	clients := *clientsPtr

	if ip != "" {
		for idx, client := range clients {
			if client.Ip == ip {
				return true, idx
			}
			if client.UDPConnection != nil || client.Connection != nil {
				client.UDPConnection = nil
				client.Connection = nil // TODO dcruzb: verify memory leak (didn't close the connection)
				return true, idx
			}
		}
	}

	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(clients) < 2 { //shared.MAX_NUMBER_OF_CONNECTIONS { TODO dcruzb: go back the env var
		client := messages.Client{
			Ip:             "",
			Connection:     nil,
			UDPConnection:  nil,
			QUICConnection: nil,
			QUICStream:     nil,
		}
		*clientsPtr = append(clients, &client)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(*clientsPtr) - 1
	}

	for idx, client := range clients {
		if client == nil {
			client := messages.Client{
				Ip:             "",
				Connection:     nil,
				UDPConnection:  nil,
				QUICConnection: nil,
				QUICStream:     nil,
			}
			clients[idx] = &client
			return true, idx
		}
		if client.UDPConnection != nil {
			client.UDPConnection.Close()
			client.UDPConnection = nil
			return true, idx
		}
		if client.Connection != nil {
			client.Ip = ""
			client.Connection.Close()
			client.Connection = nil
			return true, idx
		}
	}

	return false, -1
}

func (s SRHRPC) I_Accept(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHRPC Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	srhInfo.Counter++
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Cons", len(srhInfo.Clients))
	//log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< Counter", srhInfo.Counter)

	// check if a listener has been already created
	if srhInfo.QUICLn == nil { // no listener created
		//servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
		//if err != nil {
		//	log.Fatalf("SRH:: %v\n", err)
		//}
		quicConfig := quic.Config{KeepAlivePeriod: 60 * time.Second}
		ln, err := quic.ListenAddr(srhInfo.EndPoint.Host+":"+srhInfo.EndPoint.Port, getServerTLSQuicConfig(), &quicConfig)
		if err != nil {
			log.Fatalf("SRHRPC:: %v\n", err)
		}
		srhInfo.QUICLn = ln
	}

	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients, "")
	//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
	if !connectionAvailable {
		//lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHRPC Version Not adapted - No connection available")
		time.Sleep(1 * time.Millisecond)
		return
	}

	go func() {
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Clients Index", availableConenctionIndex)

		// Accept connections
		conn, err := srhInfo.QUICLn.Accept(context.Background())
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		srhInfo.QUICConns = append(srhInfo.QUICConns, conn)
		//srhInfo.CurrentConn = conn

		lib.PrintlnDebug("SRHRPC Version Not adapted >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Buscou nova conexão, ip:", conn.RemoteAddr().String())
		//connectionAvailable, availableConenctionIndex := s.availableConnectionFromPool(&srhInfo.Clients, conn.RemoteAddr().String())
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total Clients out", len(srhInfo.Clients))
		//if !connectionAvailable {
		//	lib.PrintlnDebug("------------------------------>", shared.GetFunction(), "end", "SRHRPC Version Not adapted - No connection available")
		//	return
		//}
		if len(srhInfo.Clients) <= availableConenctionIndex {
			lib.PrintlnDebug("SRHRPC Got len(srhInfo.Clients) <= availableConenctionIndex")
			*reset = true
			return
		}
		client := srhInfo.Clients[availableConenctionIndex]
		client.Ip = conn.RemoteAddr().String()
		client.QUICConnection = conn
		client.QUICStream, err = client.QUICConnection.AcceptStream(context.Background())
		//stream, err := tempConn.OpenStreamSync(context.Background())
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Connected Client", client)

		// Update info
		*info = srhInfo

		// Start goroutine
		go s.handler(info, availableConenctionIndex)
	}()
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
	return
}

func (s SRHRPC) I_Receive(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	// tempPort := *elemInfo[0]
	// port := tempPort.(string)
	// host := "0.0.0.0" //"127.0.0.1" // TODO

	// if LnSRHRPC == nil { // listener was not created yet
	// 	//servAddr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	// 	//if err != nil {
	// 	//	log.Fatalf("SRH:: %v\n", err)
	// 	//}
	// 	quicConfig := quic.Config{KeepAlivePeriod: 60 * time.Second}
	// 	ln, err := quic.ListenAddr(host+":"+port, getServerTLSQuicConfig(), &quicConfig)
	// 	if err != nil {
	// 		log.Fatalf("SRHRPC:: %v\n", err)
	// 	}
	// 	LnSRHRPC = ln
	// }

	// switch stateQuic {
	// case 0:
	// 	go acceptAndReadQuic(currentConnectionQuic, c1Quic)
	// 	stateQuic = 1
	// case 1:
	// 	go readQuic(currentConnectionQuic, c1Quic)
	// 	stateQuic = 2
	// case 2:
	// 	go readQuic(currentConnectionQuic, c1Quic)
	// }

	//go acceptAndRead(currentConnectionQuic, c1Quic, done)
	//go read(currentConnectionQuic, c2Quic, done)

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)

	select {
	case tempMsgReceived := <-srhInfo.RcvedMessages:
		// *msg = messages.SAMessage{Payload: tempMsgReceived}
		// Update info
		*info = srhInfo
		msg.Payload = tempMsgReceived.Msg
		lib.PrintlnDebug("SRHRPC Version Not adapted: tempMsgReceived", tempMsgReceived)
		lib.PrintlnDebug("SRHRPC Version Not adapted: tempMsgReceived.QUICStream", tempMsgReceived.QUICStream)
		if tempMsgReceived.QUICStream == nil {
			*reset = true
			return
		}
		if isNewConnection, miopPacket := s.isNewConnection(tempMsgReceived.Msg); isNewConnection {
			lib.PrintlnDebug("SRHRPC Version Not adapted: tempMsgReceived >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", miopPacket)
			*reset = true
			return
		}
		msg.ToAddr = tempMsgReceived.ToAddress //Chn.RemoteAddr().String()
	default:
		{
			*reset = true
			return
		}
	}

	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
	return
}

// func acceptAndReadQuic(currentConnectionQuic int, c chan []byte) {

// 	// accept connections
// 	temp, err := LnSRHRPC.Accept(context.Background())
// 	if err != nil {
// 		fmt.Printf("SRHRPC:: %v\n", err)
// 		os.Exit(1)
// 	}
// 	ConnsSRHRPC = append(ConnsSRHRPC, temp) // Quic Session
// 	currentConnectionQuic++

// 	// receive size
// 	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
// 	tempConn := ConnsSRHRPC[currentConnectionQuic]
// 	stream, err := tempConn.AcceptStream(context.Background())
// 	//stream, err := tempConn.OpenStreamSync(context.Background())
// 	StreamsSRHRPC = append(StreamsSRHRPC, stream)
// 	if err != nil {
// 		fmt.Printf("SRHRPC:: %v\n", err)
// 		os.Exit(1)
// 	}
// 	_, err = stream.Read(size)
// 	if err == io.EOF {
// 		{
// 			fmt.Printf("SRHRPC:: Accept and Read\n")
// 			os.Exit(0)
// 		}
// 	} else if err != nil && err != io.EOF {
// 		fmt.Printf("SRHRPC:: %v\n", err)
// 		os.Exit(1)
// 	}
// 	stream2 := StreamsSRHRPC[currentConnectionQuic]
// 	// receive message
// 	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
// 	_, err = stream2.Read(msgTemp)
// 	if err != nil {
// 		fmt.Printf("SRHRPC:: %v\n", err)
// 		os.Exit(1)
// 	}
// 	c <- msgTemp
// }

// func readQuic(currentConnectionQuic int, c chan []byte) {
// 	// receive size
// 	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
// 	stream := StreamsSRHRPC[currentConnectionQuic]

// 	_, err := stream.Read(size)
// 	if err == io.EOF {
// 		fmt.Printf("SRHRPC:: Read\n")
// 		os.Exit(0)
// 	} else if err != nil && err != io.EOF {
// 		fmt.Printf("SRHRPC:: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// receive message
// 	msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
// 	_, err = stream.Read(msgTemp)
// 	if err != nil {
// 		fmt.Printf("SRHRPC:: %v\n", err)
// 		os.Exit(1)
// 	}

// 	c <- msgTemp

// 	return
// }

func (e SRHRPC) I_Send(id string, msg *messages.SAMessage, info *interface{}, reset *bool) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHRPC Version Not adapted")
	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	lib.PrintlnDebug("msg.ToAddr", msg.ToAddr, "srhInfo.Clients", srhInfo.Clients)
	client := srhInfo.GetClientFromAddr(msg.ToAddr, srhInfo.Clients)
	stream := client.QUICStream //srhInfo.CurrentConn
	if stream == nil {
		*reset = true
		return
	}
	lib.PrintlnDebug("SRHRPC Version Not adapted   >>>>> QUIC => msg.ToAddr:", msg.ToAddr, "QUIC stream:", stream, "AdaptId:", client.AdaptId)
	msgPayload := msg.Payload.([]byte)

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgPayload)))
	_, err := stream.Write(size)
	if err != nil {
		fmt.Printf("SRHRPC:: %v\n", err)
		os.Exit(1)
	}

	// send message
	_, err = stream.Write(msgPayload)
	if err != nil {
		fmt.Printf("SRHRPC:: %v\n", err)
		os.Exit(1)
	}

	// update info
	*info = srhInfo
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
	return
}

func (s SRHRPC) handler(info *interface{}, connectionIndex int) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "SRHRPC Version Not adapted")

	infoTemp := *info
	srhInfo := infoTemp.(*messages.SRHInfo)
	stream := srhInfo.Clients[connectionIndex].QUICStream //CurrentConn
	executeForever := srhInfo.ExecuteForever

	for {
		if !*executeForever {
			break
		}
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR", "SRHRPC Version Not adapted")
		size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
		_, err := stream.Read(size)
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Read finished", "SRHRPC Version Not adapted")
		if err != nil {
			if err == io.EOF ||
				strings.Contains(err.Error(), "use of closed network connection") ||
				strings.Contains(err.Error(), "timeout: no recent network activity") {
				srhInfo.Clients[connectionIndex] = nil
				lib.PrintlnError("EOF error - Will not kill the app")
				break
			} else if err != nil && err != io.EOF {
				lib.PrintlnError("Not EOF error - Will kill the app")
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		// receive message
		msgTemp := make([]byte, binary.LittleEndian.Uint32(size))
		_, err = stream.Read(msgTemp)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
				srhInfo.Clients[connectionIndex] = nil
				lib.PrintlnError("EOF error - Will not kill the app")
				break
			} else if err != nil && err != io.EOF {
				lib.PrintlnError("Not EOF error - Will kill the app")
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		if changeProtocol, miopPacket := s.isAdapt(msgTemp); changeProtocol {
			if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
				lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHRPC Version Not adapted")
				break
			}
		}
		if isNewConnection, _ := s.isNewConnection(msgTemp); isNewConnection { // TODO dcruzb: move to I_Receive
			//newConnection = true
			lib.PrintlnDebug("QUIC Is New Connection")
			//miopPacket := miop.CreateReqPacket("Connect", []interface{}{miopPacket.Bd.ReqBody.Body[0], "Ok"}, miopPacket.Bd.ReqBody.Body[0].(int)) // idx is the Connection ID
			//msgPayload := Jsonmarshaller{}.Marshall(miopPacket)

			lib.PrintlnDebug("QUIC Before send")
			//s.send(conn, addr, msgPayload)
			lib.PrintlnDebug("QUIC After send")
			//if miopPacket.Bd.ReqBody.Body[2] == "Ok" {
			//	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "Received Ok to Adapt", "SRHUDP Version Not adapted")
			//	break
			//}
			continue
		}

		rcvMessage := messages.ReceivedMessages{Msg: msgTemp, QUICStream: stream, ToAddress: srhInfo.Clients[connectionIndex].Ip}
		lib.PrintlnDebug("SRHRPC Version Not adapted: handler >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> received message")
		if !*executeForever {
			break
		}
		srhInfo.RcvedMessages <- rcvMessage
		lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "FOR end", "SRHRPC Version Not adapted")
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "end", "SRHRPC Version Not adapted")
}

func (s SRHRPC) isAdapt(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "ChangeProtocol", miop
}

func (s SRHRPC) isNewConnection(msgFromServer []byte) (bool, miop.MiopPacket) {
	//log.Println("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	miop := Jsonmarshaller{}.Unmarshall(msgFromServer)
	return miop.Bd.ReqHeader.Operation == "Connect", miop
}

// func getServerTLSRPCConfig() *tls.Config {
// 	if shared.CRT_PATH == "" {
// 		log.Fatal("SRHRPC:: Error:: Environment variable 'CRT_PATH' not configured\n")
// 	}

// 	if shared.KEY_PATH == "" {
// 		log.Fatal("SRHRPC:: Error:: Environment variable 'KEY_PATH' not configured\n")
// 	}

// 	cert, err := tls.LoadX509KeyPair(shared.CRT_PATH, shared.KEY_PATH)
// 	if err != nil {
// 		log.Fatal("Error loading certificate. ", err)
// 	}

// 	tlsConfig := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 		NextProtos:   []string{"MidArchQuic"}, // TODO dcruzb: Verify what NextProtos should be
// 	}
// 	return tlsConfig
// }

// func (e SRHRPC) Selector(elem interface{}, elemInfo []*interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
// 	if op[2] == 'R' { // I_Receive
// 		elem.(SRHRPC).I_Receive(msg, info, elemInfo)
// 	} else { // "I_Send"
// 		elem.(SRHRPC).I_Send(msg, info, elemInfo)
// 	}
// }
