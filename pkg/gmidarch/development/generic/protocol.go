package generic

type Protocol interface {
	StartServer(ip, port string, initialConnections int)
	StopServer()
	AvailableConnectionFromPool() (available bool, idx int)
	WaitForConnection(cliIdx int) (cl *Client)
	ConnectToServer(ip, port string)
	CloseConnection()

	ReadString() string
	WriteString(message string)
	Receive() ([]byte, error)
	Send(msgToServer []byte) error

	GetClients() (clients []*Client)
	GetClient(idx int) (client Client)
	GetClientFromAddr(addr string) (client Client)
	AddClient(client Client, idx int)
	ResetClients() // Close connections and remove all clientes from the pool
}

type Client interface {
	Address() string
	AdaptId() int
	SetAdaptId(adaptId int)

	Connection() (conn interface{})
	CloseConnection()

	Read(b []byte) (err error)
	ReadString() (message string)
	WriteString(message string)
	Receive() ([]byte, error)
	Send(msgToServer []byte) error
}
