package Handlers

type HandlerNotify struct {
	Host   string
	Port int
}

// Channel used to communicate
var HandlerChan = make(chan interface{})

func (HN HandlerNotify) Start() {
	/*
	var conn net.Conn
	var err error
	var ln net.Listener

	// Create server to wait for notifications from 'Notification Consumer'
	addr := netshared.ResolveHostIp() + ":" + strings.TrimSpace(strconv.Itoa(HN.Port))
	ln, err = net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		myError := error.MyError{Source: "HandlerNotify", Message: "Unable to listen on port " + strconv.Itoa(HN.Port)}
		myError.ERROR()
	}

	if ln != nil {
		conn, err = ln.Accept()
		if err != nil {
			fmt.Println(err)
			myError := error.MyError{Source: "HandlerNotify", Message: "Unable to accept connections at port " + strconv.Itoa(HN.Port)}
			myError.ERROR()
		}
	}

	// Loop to receive data
	for {
		jsonDecoder := json.NewDecoder(conn)
		msgMOM := messages.MessageMOM{}
		err = jsonDecoder.Decode(&msgMOM)

		if err != nil {
			fmt.Println(err)
			myError := error.MyError{Source: "HandlerNotify", Message: "Unable to read data"}
			myError.ERROR()
		}
		HandlerChan <- msgMOM.PayLoad
	}
	return
*/
}

func (HN HandlerNotify) StartHandler() {
	go HN.Start()
}

// Return to 'Subscriber' the notifications received from 'Notification Consumer'
func (HandlerNotify) GetResult() interface{} {

	return <-HandlerChan
}
