package docker

type Kind int

const (
	Udp Kind = iota
	Tcp
	Tls
	Quic
	Rpc
	Http
	Https
	Http2
	E_Rpc
	E_Grpc
	E_Rmq
	UdpTcp
)

func (kind Kind) toString() string{
	switch kind {
	case Udp: return "UDP"
	case Tcp: return "TCP"
	case Tls: return "TLS"
	case Quic: return "QUIC"
	case Rpc: return "RPC"
	case Http: return "HTTP"
	case Https: return "HTTPS"
	case Http2: return "HTTP2"
	case E_Rpc: return "E_RPC"
	case E_Grpc: return "E_GRPC"
	case E_Rmq: return "E_RMQ"
	case UdpTcp: return "UdpTcp"
	}
	panic("Kind conversion to string using unlisted kind")
}

func (kind Kind) createStackCommand() string {
	switch kind {
	case Udp:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-newfibomiddleware-udp.yml newfibomiddleware-udp"
	case Tcp:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-newfibomiddleware-tcp.yml newfibomiddleware-tcp"
	case Tls:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-newfibomiddleware-tls.yml newfibomiddleware-tls"
	case Quic:  return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-quic.yml fibomiddleware-quic"
	case Rpc:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-rpc.yml fibomiddleware-rpc"
	case Http:  return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-http.yml fibomiddleware-http"
	case Https: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-https.yml fibomiddleware-https"
	case Http2: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-http2.yml fibomiddleware-http2"
	case E_Rpc: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fiborpc.yml fiborpc"
	case E_Grpc: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibogrpc.yml fibogrpc"
	case E_Rmq: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibormq.yml fibormq"
	case UdpTcp: return "docker stack deploy -c ./evaluation/experiments/docker/dc-newfibomiddleware-udptcp.yml newfibomiddleware-udptcp"
	}
	panic("Kind create stack command using unlisted kind")
}

func (kind Kind) removeStackCommand() string {
	switch kind {
	case Udp: return "docker stack rm newfibomiddleware-udp"
	case Tcp: return "docker stack rm newfibomiddleware-tcp"
	case Tls: return "docker stack rm newfibomiddleware-tls"
	case Quic: return "docker stack rm fibomiddleware-quic"
	case Rpc: return "docker stack rm fibomiddleware-rpc"
	case Http: return "docker stack rm fibomiddleware-http"
	case Https: return "docker stack rm fibomiddleware-https"
	case Http2: return "docker stack rm fibomiddleware-http2"
	case E_Rpc: return "docker stack rm fiborpc"
	case E_Grpc: return "docker stack rm fibogrpc"
	case E_Rmq: return "docker stack rm fibormq"
	case UdpTcp: return "docker stack rm newfibomiddleware-udptcp"
	}
	panic("Kind remove stack command using unlisted kind")
}
