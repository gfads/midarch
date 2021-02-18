package components

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"gmidarch/development/artefacts/graphs"
	"gmidarch/development/messages"
	"io/ioutil"
	"log"
	"os"
	"shared"
)

type CRHSsl struct {
	Behaviour string
	Graph     graphs.ExecGraph
	Conns     map[string]*tls.Conn
}

func NewCRHSsl() CRHSsl {
	r := new(CRHSsl)
	r.Behaviour = "B = InvP.e1 -> I_Process -> TerP.e1 -> B"
	r.Conns = make(map[string]*tls.Conn, shared.NUM_MAX_CONNECTIONS)

	return *r
}

func (CRHSsl) Selector(elem interface{}, elemInfo [] *interface{}, op string, msg *messages.SAMessage, info []*interface{}, r *bool) {
	elem.(CRHSsl).I_Process(msg, info)
}

func (c CRHSsl) I_Process(msg *messages.SAMessage, info [] *interface{}) {

	// check message
	payload := msg.Payload.([]interface{})
	host := "localhost" //"127.0.0.1"                // host TODO
	port := payload[1].(string)        // port
	msgToServer := payload[2].([]byte)

	addr := host + ":" + port
	var err error
	if _, ok := c.Conns[addr]; !ok { // no connection open yet
		//tcpAddr, err := net.ResolveTCPAddr("tcp", key)
		//if err != nil {
		//	log.Fatalf("CRHSsl:: %s", err)
		//}

		c.Conns[addr], err = tls.Dial("tcp4", addr, getClientTLSConfig())
		if err != nil {
			fmt.Printf("CRHSsl:: %v\n", err)
			os.Exit(1)
		}
	}

	// connect to server
	conn := c.Conns[addr]

	// send message's size
	size := make([]byte, shared.SIZE_OF_MESSAGE_SIZE, shared.SIZE_OF_MESSAGE_SIZE)
	binary.LittleEndian.PutUint32(size, uint32(len(msgToServer)))
	_, err = conn.Write(size)
	if err != nil {
		fmt.Printf("CRHSsl:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHSsl:: %v \n\n",size)

	// send message
	//fmt.Printf("CRHSsl:: Message to server:: %v %v >> %v << \n\n",msgToServer, len(msgToServer), binary.LittleEndian.Uint32(size))
	_, err = conn.Write(msgToServer)
	if err != nil {
		fmt.Printf("CRHSsl:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHSsl:: Message sent to Server [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	// receive reply's size
	_, err = conn.Read(size)
	if err != nil {
		fmt.Printf("CRHSsl:: %v\n", err)
		os.Exit(1)
	}

	// receive reply
	msgFromServer := make([]byte, binary.LittleEndian.Uint32(size), shared.NUM_MAX_MESSAGE_BYTES)
	_, err = conn.Read(msgFromServer)
	if err != nil {
		fmt.Printf("CRHSsl:: %v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("CRHSsl:: Message received from Server:: [%v,%v] \n",conn.LocalAddr(),conn.RemoteAddr())

	*msg = messages.SAMessage{Payload: msgFromServer}
}

func getClientTLSConfig() *tls.Config {
	if shared.CA_PATH == "" {
		log.Fatal("CRHSsl:: Error:: Environment variable 'CA_PATH' not configured\n")
	}
	trustCert, err := ioutil.ReadFile(shared.CA_PATH)
	if err != nil {
		fmt.Println("Error loading trust certificate. ",err)
	}
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(trustCert) {
		fmt.Println("Error installing trust certificate.")
	}

	// connect to server
	tlsConfig := &tls.Config{
		//InsecureSkipVerify: true,
		RootCAs: certs,
		NextProtos:         []string{"exemplo"},
	}
	return tlsConfig
}