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

type Client struct {
	connection net.Conn
	Ip         string
}

func (cl *Client) Address() string {
	return cl.Ip
}

func (cl *Client) Connection() interface{} {
	return cl.connection
}

func (cl *Client) CloseConnection() {
	err := cl.connection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (cl *Client) ReadString() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.connection).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket TCP. Details:", err)
	}

	return message
}

func (cl *Client) WriteString(message string) {
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

	_, err := cl.connection.Write([]byte(message + "\n"))
	if err != nil {
		lib.PrintlnError("Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}
}

func (cl *Client) Read(b []byte) (err error) {
	_, err = cl.connection.Read(b)
	if err != nil {
		if err == io.EOF || strings.Contains(err.Error(), "use of closed network connection") {
			cl = nil
			lib.PrintlnError("EOF Error: Will not kill app")
			return err
		} else if err != nil && err != io.EOF {
			lib.PrintlnError("Error, not EOF, will kill the app")
			shared.ErrorHandler(shared.GetFunction(), err.Error())
			return err
		}
	}
	return nil
}

func (cl *Client) Receive() (msg []byte, err error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	// receive reply's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	cl.Read(size)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	// receive reply
	msg = make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	err = cl.Read(msg)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	return msg, nil
}

func (cl *Client) Send(msg []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msg)))
	_, err := cl.connection.Write(sizeOfMsgSize)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}

	// send message
	_, err = cl.connection.Write(msg)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}
	return nil
}

type TCP struct {
	ip                 string
	port               string
	listener           net.Listener
	serverConnection   net.Conn
	initialConnections int
	clients            []*generic.Client
}

func (st *TCP) StartServer(ip, port string, initialConnections int) {
	servAddr, err := net.ResolveTCPAddr("tcp", ip+":"+port)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err.Error())
	}
	ln, err := net.ListenTCP("tcp", servAddr)
	if err != nil {
		lib.PrintlnError("Error while starting TCP server. Details: ", err)
	}
	st.listener = ln
	st.initialConnections = initialConnections
	st.clients = make([]*generic.Client, st.initialConnections)
}

func (st *TCP) StopServer() {
	err := st.listener.Close()
	if err != nil {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *TCP) AvailableConnectionFromPool() (available bool, idx int) {
	//lib.PrintlnDebug(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Total clients", len(clients))
	if len(st.clients) < st.initialConnections {
		client := &Client{}
		// *clientsPtr = append(clients, &client)
		st.AddClient(client, -1)
		//log.Println(">>>>>>>>>>>>>>>>>>>>>>>>	>>>>>>>>>>>>>> Total Clients", len(*clientsPtr))
		return true, len(st.clients) - 1
	}

	for idx, client := range st.clients {
		if client == nil {
			st.AddClient(&Client{}, idx)
			return true, idx
		}
	}

	return false, -1
}

func (st *TCP) ConnectToServer(ip, port string) {
	addr := ip + ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	lib.PrintlnDebug("Resolved addr", tcpAddr)
	//localTcpAddr := c.getLocalTcpAddr()

	for {
		st.serverConnection, err = net.DialTCP("tcp", nil, tcpAddr)
		lib.PrintlnDebug("Dialed", st.serverConnection)
		if err != nil {
			lib.PrintlnError("Dial error", st.serverConnection, err)
			time.Sleep(200 * time.Millisecond)
			//shared.ErrorHandler(shared.GetFunction(), err.Error())
		} else {
			break
		}
	}
	lib.PrintlnDebug("Connected", st.serverConnection)
	if addr != shared.NAMING_HOST+":"+shared.NAMING_PORT && shared.LocalAddr == "" {
		//lib.PrintlnDebug("crhInfo.Conns[addr].LocalAddr().String()", crhInfo.Conns[addr].LocalAddr().String())
		shared.LocalAddr = st.serverConnection.LocalAddr().String()
		lib.PrintlnDebug("Got local addr", st.serverConnection)
	}
}

func (st *TCP) WaitForConnection(cliIdx int) (cl *generic.Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	conn, err := st.listener.Accept()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Error while waiting for connection: "+err.Error())
	}

	(*st.clients[cliIdx]).(*Client).connection = conn
	(*st.clients[cliIdx]).(*Client).Ip = conn.RemoteAddr().String()

	return st.clients[cliIdx]
}

func (st *TCP) CloseConnection() {
	err := st.serverConnection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (st *TCP) ReadString() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(st.serverConnection).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket TCP. Details:", err)
	}

	return message
}

func (st *TCP) WriteString(message string) {
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

	_, err := st.serverConnection.Write([]byte(message + "\n"))
	if err != nil {
		lib.PrintlnError("Error while writing message to socket TCP. Details:", err)
		os.Exit(1)
	}
}

func (st *TCP) Receive() ([]byte, error) {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	// receive reply's size
	_, err := st.serverConnection.Read(sizeOfMsgSize)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(sizeOfMsgSize), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = st.serverConnection.Read(msgFromServer)
	if err != nil {
		lib.PrintlnError(shared.GetFunction(), err)
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return nil, err
	}
	return msgFromServer, nil
}

func (st *TCP) Send(msgToServer []byte) error {
	lib.PrintlnDebug("----------------------------------------->", shared.GetFunction(), "CRHTCP Version Not adapted")
	sizeOfMsgSize := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE) // TODO dcruzb: create attribute to avoid doing this everytime
	binary.LittleEndian.PutUint32(sizeOfMsgSize, uint32(len(msgToServer)))
	_, err := st.serverConnection.Write(sizeOfMsgSize)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}

	// send message
	_, err = st.serverConnection.Write(msgToServer)
	if err != nil {
		//shared.ErrorHandler(shared.GetFunction(), err.Error())
		return err
	}
	return nil
}

func (st *TCP) AddClient(client generic.Client, idx int) {
	if idx < 0 {
		st.clients = append(st.clients, &client)
	} else if idx < st.initialConnections {
		st.clients[idx] = &client
	}
}

func (st *TCP) GetClient(idx int) (client generic.Client) {
	return *st.clients[idx]
}

func (st *TCP) GetClientFromAddr(addr string) (client generic.Client) {
	for _, client := range st.clients {
		if (*client).Address() == addr {
			return *client
		}
	}

	log.Println("IP without client from the ip:", addr)
	return nil
}
