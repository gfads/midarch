package messages

import "gmidarch/development/miop"

type ToCRH struct {
	Host string
	Port int
	MIOP miop.Packet
}

