package messages

import "net"

type CRHInfo struct {
	EndPoint EndPoint
	Conns    map[string]net.Conn
}
