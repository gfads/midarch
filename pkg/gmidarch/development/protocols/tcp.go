package protocols

import (
	"bufio"
	"net"
	"os"

	"github.com/gfads/midarch/pkg/shared/lib"
)

type Client struct {
	connection net.Conn
}

type TCP struct {
	ip                 string
	port               string
	listener           net.Listener
	serverConnection   net.Conn
	initialConnections int
	clients            []Client
}

func (st *TCP) StartServer(ip, port string, useJson bool, initialConnections int) {
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		lib.PrintlnError("Error while starting TCP server. Details: ", err)
	}
	st.listener = ln
	st.initialConnections = initialConnections
	st.clients = make([]Client, st.initialConnections)
}

func (st *TCP) StopServer() {
	err := st.listener.Close()
	if err != nil {
		lib.PrintlnError("Error while stoping server. Details:", err)
	}
}

func (st *TCP) ConnectToServer(ip, port string) {
	// connect to server
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		lib.PrintlnError(err)
	}

	st.serverConnection = conn
}

func (st *TCP) WaitForConnection(cliIdx int) (cl *Client) { // TODO if cliIdx >= inicitalConnections => need to append to the slice
	// aceita conexões na porta
	conn, err := st.listener.Accept()
	if err != nil {
		lib.PrintlnError("Error while waiting for connection", err)
	}

	cl = &st.clients[cliIdx]

	cl.connection = conn

	return cl
}

func (st *TCP) CloseConnection() {
	err := st.serverConnection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (cl *Client) CloseConnection() {
	err := cl.connection.Close()
	if err != nil {
		lib.PrintlnError(err)
	}
}

func (st *TCP) Read() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(st.serverConnection).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket TCP. Details:", err)
	}

	return message
}

func (st *TCP) Write(message string) {
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

func (cl *Client) Read() (message string) {
	var err error
	// recebe solicitações do cliente
	message, err = bufio.NewReader(cl.connection).ReadString('\n')
	if err != nil {
		lib.PrintlnError("Error while reading message from socket TCP. Details:", err)
	}

	return message
}

func (cl *Client) Write(message string) {
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
