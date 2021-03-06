package docker

type Kind int

const (
	Udp Kind = iota
	Tcp
	Ssl
	Quic
	Rpc
	Http
	Https
	Http2
	E_Rpc
	E_Grpc
	E_Rmq
)

func (kind Kind) toString() string{
	switch kind {
	case Udp: return "UDP"
	case Tcp: return "TCP"
	case Ssl: return "SSL"
	case Quic: return "QUIC"
	case Rpc: return "RPC"
	case Http: return "HTTP"
	case Https: return "HTTPS"
	case Http2: return "HTTP2"
	case E_Rpc: return "E_RPC"
	case E_Grpc: return "E_GRPC"
	case E_Rmq: return "E_RMQ"
	}
	panic("Kind conversion to string using unlisted kind")
}

func (kind Kind) createStackCommand() string {
	switch kind {
	case Udp:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-udp.yml fibomiddleware-udp"
	case Tcp:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-tcp.yml fibomiddleware-tcp"
	case Ssl:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-ssl.yml fibomiddleware-ssl"
	case Quic:  return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-quic.yml fibomiddleware-quic"
	case Rpc:   return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-rpc.yml fibomiddleware-rpc"
	case Http:  return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-http.yml fibomiddleware-http"
	case Https: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-https.yml fibomiddleware-https"
	case Http2: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibomiddleware-http2.yml fibomiddleware-http2"
	case E_Rpc: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fiborpc.yml fiborpc"
	case E_Grpc: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibogrpc.yml fibogrpc"
	case E_Rmq: return "docker stack deploy -c ./evaluation/experiments/docker/dc-fibormq.yml fibormq"
	}
	panic("Kind create stack command using unlisted kind")
}

func (kind Kind) removeStackCommand() string {
	switch kind {
	case Udp: return "docker stack rm fibomiddleware-udp"
	case Tcp: return "docker stack rm fibomiddleware-tcp"
	case Ssl: return "docker stack rm fibomiddleware-ssl"
	case Quic: return "docker stack rm fibomiddleware-quic"
	case Rpc: return "docker stack rm fibomiddleware-rpc"
	case Http: return "docker stack rm fibomiddleware-http"
	case Https: return "docker stack rm fibomiddleware-https"
	case Http2: return "docker stack rm fibomiddleware-http2"
	case E_Rpc: return "docker stack rm fiborpc"
	case E_Grpc: return "docker stack rm fibogrpc"
	case E_Rmq: return "docker stack rm fibormq"
	}
	panic("Kind remove stack command using unlisted kind")
}
