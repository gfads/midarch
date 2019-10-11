package interfaces

import "newsolution/development/miop"

type Marshaller interface {
	Marshall(packet miop.Packet) []byte
	Unmarshall ([]byte) miop.Packet
}
