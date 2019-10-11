package impl

import (
	"newsolution/gmidarch/development/miop"
	"encoding/json"
	"log"
)

type MarshallerImpl struct {}

func (MarshallerImpl) Marshall(msg miop.Packet) []byte {

	r, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Marshaller:: Marshall:: %s", err)
	}

	return r
}

func (MarshallerImpl) Unmarshall(msg []byte) miop.Packet {

	r := miop.Packet{}
	err := json.Unmarshal(msg, &r)
	if err != nil {
		log.Fatalf("Marshaller:: Unmarshall:: %s", err)
	}
	return r
}


