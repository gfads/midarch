package protocols

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"
)

type UDPClient struct {
	connection *net.UDPConn
	Ip         string
	adaptId    int
}

func (cl *UDPClient) Address() string {
	return cl.Ip
}

func (cl *UDPClient) AdaptId() int {
	return cl.adaptId
}

func (cl *UDPClient) SetAdaptId(adaptId int) {
	cl.adaptId = adaptId
}

func (cl *UDPClient) Connection() interface{} {
	return cl.connection
}

func (cl *UDPClient) CloseConnection() {
	cl.Ip = ""
	if cl.connection != nil {
		err := cl.connection.Close()
		if err != nil {
			lib.PrintlnError(err)
		}
	}
}

func (cl *UDPClient) ReadString() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.connection).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket UDP. Details:", err)
	}

	return message
}

func (cl *UDPClient) WriteString(message string) {
	// envia resposta

	// Vários tipos diferentes de se escrever utilizando Writer, todos funcionam
	//_, err := fmt.Fprintf(conn, msgToServer+"\n")
	//_, err := conn.Write([]byte( msgToServer + "\n"))
	/*reader := bufio.NewWriter(conn)
	_, err := reader.WriteString( msgToServer + "\n")
	reader.Flush()*/
	/*reader := bufio.NewWriter(conn)
	_, err := io.WriteString(reader, msgToServer + "\n")
	reader.Flush()*/
	//_, err := io.WriteString(conn, msgToServer+"\n")

	ipAddr := &net.UDPAddr{IP: net.ParseIP(cl.Ip)}
	_, err := cl.connection.WriteTo([]byte(message+"\n"), ipAddr)
	if err != nil {
		lib.PrintlnError("Error while writing message to socket UDP. Details:", err)
		os.Exit(1)
	}
}

func (cl *UDPClient) Read(b []byte) (n int, err error) {
	n, addr, err := cl.connection.ReadFromUDP(b)
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
			cl = nil
			lib.PrintlnError("EOF Error: Will not kill app")
			return n, err
		} else if err != nil && err != io.EOF {
			lib.PrintlnError("Error, not EOF, will kill the app")
			shared.ErrorHandler(shared.GetFunction(), err.Error())
			return n, err
		}
	}
	cl.Ip = addr.String()
	return n, nil
}

func (cl *UDPClient) Receive() (fullMessage []byte, err error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
	// receive reply's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	cl.Read(size)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	// TODO dcruzb: validate if size is smaller than shared.NUM_MAX_MESSAGE_BYTES
	// receive reply
	//msg = make([]byte, binary.LittleEndian.Uint32(size), binary.LittleEndian.Uint32(size))
	//err = cl.Read(msg)
	//if err != nil {
	//	lib.PrintlnError(shared.GetFunction(), err)
	//	//shared.ErrorHandler(shared.GetFunction(), err.Error())
	//	return nil, err
	//}

	addr, err := net.ResolveUDPAddr("udp", cl.Ip)
	if err != nil {
		return nil, err
	}
	msgSize := binary.LittleEndian.Uint32(size)
	const maxBufferSize = shared.MAX_PACKET_SIZE
	for {
		bufferSize := int(msgSize) - len(fullMessage)
		if bufferSize > maxBufferSize {
			bufferSize = maxBufferSize
		}
		buffer := make([]byte, bufferSize, bufferSize)
		// lib.PrintlnInfo("Received(read-ini):size", size, "len(fullMessage)", len(fullMessage), "bufferSize", bufferSize, "remaining", int(size)-len(fullMessage))

		// lib.PrintlnInfo("Received(read):for1")
		n, err := cl.Read(buffer)
		// lib.PrintlnInfo("Received(read):", buffer)
		if err != nil {
			lib.PrintlnError("Error while reading message. Error:", err)
			return nil, err
		}

		fullMessage = append(fullMessage, buffer[:n]...)
		// lib.PrintlnInfo("Received(read):for2")
		// lib.PrintlnInfo("Received(read-end):size", size, "len(fullMessage)", len(fullMessage), "bufferSize", bufferSize)
		// Check if the message is complete (you need a way to determine this based on your protocol)
		if len(fullMessage) >= int(msgSize) {
			return fullMessage, nil
		} else {
			_, err = cl.connection.WriteTo([]byte("ack"), addr)
			if err != nil {
				return nil, err
			}
		}
		// lib.PrintlnInfo("Received(read):for3")
	}
}

func (cl *UDPClient) Send(msg []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
	addr, err := net.ResolveUDPAddr("udp", cl.Ip)
	if err != nil {
		return err
	}
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msg)))
	_, err = cl.connection.WriteTo(sizeOfMsgSize, addr)
	if err != nil {
		return err
	}
	// send message
	_, err = cl.connection.WriteTo(msg, addr)
	if err != nil {
		return err
	}
	return nil
}

type UDP struct {
	// Server attributes
	ip                 string
	port               string
	connection         *net.UDPConn
	initialConnections int
	clients            []*generic.Client
	// Client attributes
	serverConnection *net.UDPConn
}

func (st *UDP) StartServer(ip, port string, initialConnections int) {
	st.ip = ip
	st.port = port
	servAddr, err := net.ResolveUDPAddr("udp4", ip+":"+port)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
	}
	connection, err := net.ListenUDP("udp4", servAddr) //, err := srhInfo.Ln.Accept()

	if err != nil {
		lib.PrintlnError("Error while starting UDP server. Details: ", err)
	}
	st.connection = connection
	st.initialConnections = initialConnections
	st.clients = make([]*generic.Client, st.initialConnections)
}

func (st *UDP) StopServer() {
	st.ResetClients()
	err := st.connection.Close()
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *UDP) AvailableConnectionFromPool() (available bool, idx int) {
	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(st.clients) < st.initialConnections {
		client := &UDPClient{}
		// *clientsPtr = append(clients, &client)
		st.AddClient(client, -1)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>	>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(st.clients) - 1
	}

	for idx, client := range st.clients {
		if client == nil {
			st.AddClient(&UDPClient{}, idx)
			return true, idx
		}
	}

	return false, -1
}

func (st *UDP) ConnectToServer(ip, port string) {
	addr := ip + ":" + port
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	lib.PrintlnDebug("Resolved addr", udpAddr)
	//localUdpAddr := c.getLocalUdpAddr()

	for {
		st.serverConnection, err = net.DialUDP("udp", nil, udpAddr)
		// lib.PrintlnInfo("Dialed", st.serverConnection)
		if err != nil {
			lib.PrintlnError("Dial error", st.serverConnection, err)
			time.Sleep(200 * time.Millisecond)
			//shared.ErrorHandler(shared.GetFunction(), err.Error())
		} else {
			break
		}
	}

	if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
		//lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
		shared.LocalAddr = st.serverConnection.LocalAddr().String()
		lib.PrintlnDebug("Got local addr", st.serverConnection)
	}
}

func (st *UDP) WaitForConnection(cliIdx int) (cl *generic.Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	// lib.PrintlnInfo("Before accept")
	// conn, err := st.listener.Accept()
	// if err != nil {
	// 	if strings.Contains(err.Error(), "use of closed network connection") {
	// 		return nil
	// 	}
	// 	shared.ErrorHandler(shared.GetFunction(), "Error while waiting for connection: "+err.Error())
	// }
	// lib.PrintlnInfo("After accept (cliIdx", cliIdx, ")")
	if len(st.clients) > cliIdx {
		(*st.clients[cliIdx]).(*UDPClient).connection = st.connection
		(*st.clients[cliIdx]).(*UDPClient).Ip = "" //conn.RemoteAddr().String()

		return st.clients[cliIdx]
	} else {
		return nil
	}
}

func (st *UDP) CloseConnection() {
	err := st.serverConnection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (st *UDP) ReadString() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(st.serverConnection).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket UDP. Details:", err)
	}

	return message
}

func (st *UDP) WriteString(message string) {
	// envia resposta

	// Vários tipos diferentes de se escrever utilizando Writer, todos funcionam
	//_, err := fmt.Fprintf(conn, msgToServer+"\n")
	//_, err := conn.Write([]byte( msgToServer + "\n"))
	/*reader := bufio.NewWriter(conn)
	_, err := reader.WriteString( msgToServer + "\n")getLocalTcpAddr
	reader.Flush()*/
	/*reader := bufio.NewWriter(conn)
	_, err := io.WriteString(reader, msgToServer + "\n")
	reader.Flush()*/
	//_, err := io.WriteString(conn, msgToServer+"\n")

	_, err := st.serverConnection.Write([]byte(message + "\n"))
	if err != nil {
		lib.PrintlnError("Error while writing message to socket UDP. Details:", err)
		os.Exit(1)
	}
}

func (st *UDP) Receive() ([]byte, error) {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "UDP Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// receive reply's size
	_, err := st.serverConnection.Read(sizeOfMsgSize)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "UDP read size")
	// receive reply
	// TODO dcruzb: validate if size is smaller than shared.NUM_MAX_MESSAGE_BYTES
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), binary.LittleEndian.Uint32(sizeOfMsgSize))
	_, err = st.serverConnection.Read(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "UDP read message")
	return msgFromServer, nil
}

func (st *UDP) Send(msgToServer []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHUDP Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := st.serverConnection.Write(sizeOfMsgSize)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}

	// send message
	const maxPacketSize = shared.MAX_PACKET_SIZE
	// send message
	fragmentedMessage := msgToServer
	for {
		fragmentSize := len(fragmentedMessage)
		if fragmentSize > maxPacketSize {
			fragmentSize = maxPacketSize
		}
		fragment := fragmentedMessage[:fragmentSize]
		// lib.PrintlnInfo("Send: fragment:", fragment)
		// lib.PrintlnInfo("Send(read-ini):size", len(msgToServer), "len(fragmentedMessage)-remaining:", len(fragmentedMessage), "maxPacketSize", maxPacketSize)
		_, err = st.serverConnection.Write(fragment)
		if err != nil {
			//fmt.Println("Erro no envio do sizeOfMsgSize(", sizeOfMsgSize, ") Connection:", reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name())
			//shared.ErrorHandler(shared.GetFunction(), err.Error())
			lib.PrintlnError("Error while writing fragment to server, error:", err)
			return err
		}

		fragmentedMessage = fragmentedMessage[fragmentSize:]
		if len(fragmentedMessage) > 0 {
			ackBuffer := make([]byte, 3)
			_, err := st.serverConnection.Read(ackBuffer)
			if err != nil || (strings.TrimSpace(string(ackBuffer)) != "ack" && strings.TrimSpace(string(ackBuffer)) != "ok") {
				lib.PrintlnError("Error while reading message. ackBuffer: '"+strings.TrimSpace(string(ackBuffer))+"'. Error:", err)
				if err != nil {
					return err
				} else {
					// if the message is not an error, and not the expected info, then reads again to get one more byte, should be the message size
					sizeaux := make([]byte, 30)
					_, err := st.serverConnection.Read(sizeaux)
					if err != nil {
						return err
					}
					msgSize := ackBuffer
					msgSize = append(msgSize, sizeaux[0])
					lib.PrintlnError("msgSize: '" + strings.TrimSpace(string(msgSize)))
				}
			}
		} else {
			break
		}
	}

	return nil
}

func (st *UDP) GetClients() (client []*generic.Client) {
	return st.clients
}

func (st *UDP) GetClient(idx int) (client generic.Client) {
	return *st.clients[idx]
}

func (st *UDP) AddClient(client generic.Client, idx int) {
	if idx < 0 {
		st.clients = append(st.clients, &client)
	} else if idx < st.initialConnections {
		st.clients[idx] = &client
	}
}

func (st *UDP) GetClientFromAddr(addr string) (client generic.Client) {
	for _, client := range st.clients {
		if (*client).Address() == addr {
			return *client
		}
	}

	log.Println("IP without client from the ip:", addr)
	return nil
}

func (st *UDP) ResetClients() {
	// log.Println("UDP.ResetClients")
	for _, client := range st.clients {
		if client != nil {
			(*client).CloseConnection()
		}
	}
	st.clients = st.clients[:0]
	// log.Println("UDP.ResetClients clients length:", len(st.clients))
}

// func Remove(slice []*UDPClient, idx int) []*UDPClient {
// 	var newSlice []*UDPClient

// 	if len(slice) == idx+1 {
// 		newSlice = append(slice[:idx])
// 	} else {
// 		newSlice = append(slice[:idx], slice[idx+1:]...)
// 	}

// 	return newSlice
// }

// for len(srhInfo.Protocol.GetClients()) > 0 {
// 	tmpClient := srhInfo.Protocol.GetClients()[0]
// 	lib.PrintlnInfo("Will initialize:", tmpClient)
// 	srhInfo.Clients = messages.Remove(srhInfo.Clients, len(srhInfo.Clients)-1)
// 	tmpClient.Initialize()
// 	lib.PrintlnInfo("Initialized")
// }
