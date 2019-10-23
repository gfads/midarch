package impl

import (
	"encoding/json"
	miop2 "gmidarch/development/miop"
	"log"
)

type MarshallerImpl struct {}

func (MarshallerImpl) Marshall(msg miop2.Packet) []byte {

	r, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Marshaller:: Marshall:: %s", err)
	}

	return r
}

func (MarshallerImpl) Unmarshall(msg []byte) miop2.Packet {

	r := miop2.Packet{}
	err := json.Unmarshal(msg, &r)
	if err != nil {
		log.Fatalf("Marshaller:: Unmarshall:: %s", err)
	}
	return r
}


