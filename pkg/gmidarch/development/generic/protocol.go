package generic

type Protocol interface {
	StartServer(ip, port string)
	StopServer()
	WaitForConnection(cliIdx int) // (cl *Client)
	ConnectToServer(ip, port string)
	CloseConnection()
	Send()
	Receive()
}
