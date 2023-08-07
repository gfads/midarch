package messages

import (
	"github.com/gfads/midarch/pkg/gmidarch/development/generic"
	"net"

	"github.com/quic-go/quic-go"
)

type CRHInfo struct {
	EndPoint    EndPoint
	Protocols   map[string]generic.Protocol
	Conns       map[string]net.Conn
	QuicConns   map[string]quic.Connection
	QuicStreams map[string]quic.Stream
}
