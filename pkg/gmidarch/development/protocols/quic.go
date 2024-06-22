package protocols

import (
	"bufio"
	"context"
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
	"github.com/quic-go/quic-go"
)

type QUICClient struct {
	connection quic.Connection
	stream     quic.Stream
	Ip         string
	adaptId    int
}

func (cl *QUICClient) Address() string {
	return cl.Ip
}

func (cl *QUICClient) AdaptId() int {
	return cl.adaptId
}

func (cl *QUICClient) SetAdaptId(adaptId int) {
	cl.adaptId = adaptId
}

func (cl *QUICClient) Connection() interface{} {
	return cl.stream
}

func (cl *QUICClient) CloseConnection() {
	cl.Ip = ""
	if cl.stream != nil {
		err := cl.stream.Close() // TODO dcruzb: There is need to close the connection too?
		if err != nil {
			lib.PrintlnError(err)
		}
	}
}

func (cl *QUICClient) ReadString() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.stream).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket QUIC. Details:", err)
	}

	return message
}

func (cl *QUICClient) WriteString(message string) {
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

	_, err := cl.stream.Write([]byte(message + "\n"))
	if err != nil {
		lib.PrintlnError("Error while writing message to socket QUIC. Details:", err)
		os.Exit(1)
	}
}

func (cl *QUICClient) Read(b []byte) (n int, err error) {
	n, err = cl.stream.Read(b)
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
	return n, nil
}

func (cl *QUICClient) Receive() (fullMessage []byte, err error) {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
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
		}
		// lib.PrintlnInfo("Received(read):for3")
	}
}

func (cl *QUICClient) Send(msg []byte) error {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msg)))
	_, err := cl.stream.Write(sizeOfMsgSize)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}

	// send message
	_, err = cl.stream.Write(msg)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}
	return nil
}

type QUIC struct {
	// Server attributes
	ip                 string
	port               string
	listener           *quic.Listener
	initialConnections int
	clients            []*generic.Client
	// Client attributes
	serverConnection quic.Connection
	stream           quic.Stream
}

func (st *QUIC) StartServer(ip, port string, initialConnections int) {
	servAddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
	}
	quicConfig := quic.Config{KeepAlivePeriod: 60 * time.Second}
	ln, err := quic.ListenAddr(servAddr.String(), lib.GetServerTLSConfig("h2"), &quicConfig)
	if err != nil {
		lib.PrintlnError("Error while starting QUIC server. Details: ", err)
	}
	st.listener = ln
	st.initialConnections = initialConnections
	st.clients = make([]*generic.Client, st.initialConnections)
}

func (st *QUIC) StopServer() {
	st.ResetClients()
	err := st.listener.Close()
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *QUIC) AvailableConnectionFromPool() (available bool, idx int) {
	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(st.clients) < st.initialConnections {
		client := &QUICClient{}
		// *clientsPtr = append(clients, &client)
		st.AddClient(client, -1)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>	>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(st.clients) - 1
	}

	for idx, client := range st.clients {
		if client == nil {
			st.AddClient(&QUICClient{}, idx)
			return true, idx
		}
	}

	return false, -1
}

func (st *QUIC) ConnectToServer(ip, port string) {
	addr := ip + ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	// lib.PrintlnDebug("Resolved addr", tcpAddr)
	//localTcpAddr := c.getLocalTcpAddr()

	for {
		st.serverConnection, err = quic.DialAddr(context.Background(), tcpAddr.String(), lib.GetClientTLSConfig("h2"), nil)
		// lib.PrintlnInfo("Dialed", st.serverConnection)
		if err != nil {
			lib.PrintlnError("Dial error", st.serverConnection, err)
			time.Sleep(200 * time.Millisecond)
			//shared.ErrorHandler(shared.GetFunction(), err.Error())
		} else {
			st.stream, err = st.serverConnection.OpenStreamSync(context.Background())
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
			break
		}
	}

	// lib.PrintlnDebug("Connected", st.serverConnection)
	if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
		//lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
		shared.LocalAddr = st.serverConnection.LocalAddr().String()
		// lib.PrintlnDebug("Got local addr", st.serverConnection)
	}
}

func (st *QUIC) WaitForConnection(cliIdx int) (cl *generic.Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	// lib.PrintlnInfo("Before accept")
	conn, err := st.listener.Accept(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "use of closed network connection") {
			return nil
		}
		shared.ErrorHandler(shared.GetFunction(), "Error while waiting for connection: "+err.Error())
	}
	// lib.PrintlnInfo("After accept (cliIdx", cliIdx, ")")
	if len(st.clients) > cliIdx {
		(*st.clients[cliIdx]).(*QUICClient).connection = conn
		(*st.clients[cliIdx]).(*QUICClient).Ip = conn.RemoteAddr().String()
		(*st.clients[cliIdx]).(*QUICClient).stream, err = conn.AcceptStream(context.Background())
		//(*st.clients[cliIdx]).(*QUICClient).Stream, err := tempConn.OpenStreamSync(context.Background())
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		return st.clients[cliIdx]
	} else {
		return nil
	}
}

func (st *QUIC) CloseConnection() {
	err := st.serverConnection.CloseWithError(0, "Closing connection")
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (st *QUIC) ReadString() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(st.stream).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket QUIC. Details:", err)
	}

	return message
}

func (st *QUIC) WriteString(message string) {
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

	_, err := st.stream.Write([]byte(message + "\n"))
	if err != nil {
		lib.PrintlnError("Error while writing message to socket QUIC. Details:", err)
		os.Exit(1)
	}
}

func (st *QUIC) Receive() ([]byte, error) {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "QUIC Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// receive reply's size
	_, err := st.stream.Read(sizeOfMsgSize)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "QUIC read size")
	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = st.stream.Read(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "QUIC read message")
	return msgFromServer, nil
}

func (st *QUIC) Send(msgToServer []byte) error {
	// lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := st.stream.Write(sizeOfMsgSize)
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
		_, err = st.stream.Write(fragment)
		if err != nil {
			//fmt.Println("Erro no envio do sizeOfMsgSize(", sizeOfMsgSize, ") Connection:", reflect.TypeOf(crhInfo.Conns[addr]).Elem().Name())
			//shared.ErrorHandler(shared.GetFunction(), err.Error())
			lib.PrintlnError("Error while writing fragment to server, error:", err)
			return err
		}

		fragmentedMessage = fragmentedMessage[fragmentSize:]
		if len(fragmentedMessage) > 0 {
			// time.Sleep(1 * time.Millisecond) // uncomment when using with debug enabled to avoid message loss
		} else {
			break
		}
	}

	return nil
}

func (st *QUIC) GetClients() (client []*generic.Client) {
	return st.clients
}

func (st *QUIC) GetClient(idx int) (client generic.Client) {
	return *st.clients[idx]
}

func (st *QUIC) AddClient(client generic.Client, idx int) {
	if idx < 0 {
		st.clients = append(st.clients, &client)
	} else if idx < st.initialConnections {
		st.clients[idx] = &client
	}
}

func (st *QUIC) GetClientFromAddr(addr string) (client generic.Client) {
	for _, client := range st.clients {
		if (*client).Address() == addr {
			return *client
		}
	}

	log.Println("IP without client from the ip:", addr)
	return nil
}

func (st *QUIC) ResetClients() {
	// log.Println("QUIC.ResetClients")
	for _, client := range st.clients {
		if client != nil {
			(*client).CloseConnection()
		}
	}
	st.clients = st.clients[:0]
	// log.Println("QUIC.ResetClients clients length:", len(st.clients))
}

// func Remove(slice []*QUICClient, idx int) []*QUICClient {
// 	var newSlice []*QUICClient

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
