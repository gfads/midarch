package generic

type Protocol interface {
	StartServer(ip, port string, initialConnections int)
	StopServer()
	WaitForConnection(cliIdx int) (cl *Client)
	ConnectToServer(ip, port string)
	CloseConnection()
	Write(message string)
	Read() string
	Receive(size []byte) ([]byte, error)
	Send(sizeOfMsgSize []byte, msgToServer []byte) error
}

type Client interface {
	Connection() (conn interface{})
	CloseConnection()
	Read() (message string)
	Write(message string)
}
