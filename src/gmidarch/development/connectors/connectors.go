package connectors

import (
	"github.com/gfads/midarch/src/gmidarch/development/messages"
	"github.com/gfads/midarch/src/shared"
)

type Connector struct {
	Id                string
	TypeName          string
	DefaultLeftArity  int
	DefaultRightArity int
	LeftArity         int
	RightArity        int
	Buffer            chan messages.SAMessage
	State             []chan bool
	Behaviour         string
}

func NewConnector(id string, typeName string, behaviour string, lArity, rArity int) Connector {

	connTemp := Connector{}

	switch typeName {
	case shared.ONEWAY:
		connTemp.Id = id
		connTemp.DefaultLeftArity = 1
		connTemp.LeftArity = connTemp.DefaultLeftArity
		connTemp.DefaultRightArity = 1
		connTemp.RightArity = connTemp.DefaultRightArity
		connTemp.TypeName = typeName
		connTemp.Buffer = make(chan messages.SAMessage, connTemp.DefaultLeftArity)
		connTemp.State = []chan bool{make(chan bool, 1), make(chan bool, 1)}
		connTemp.State[0] <- true
		connTemp.Behaviour = behaviour
		return connTemp
	case shared.REQUEST_REPLY:
		connTemp.Id = id
		connTemp.DefaultLeftArity = 1
		connTemp.LeftArity = connTemp.DefaultLeftArity
		connTemp.DefaultRightArity = 1
		connTemp.RightArity = connTemp.DefaultRightArity
		connTemp.TypeName = typeName
		connTemp.Buffer = make(chan messages.SAMessage, connTemp.DefaultLeftArity)
		connTemp.State = []chan bool{make(chan bool, 1), make(chan bool, 1), make(chan bool, 1), make(chan bool, 1)}
		connTemp.State[0] <- true
		connTemp.Behaviour = behaviour
		return connTemp
	case shared.ONETON:
		connTemp.Id = id
		connTemp.DefaultLeftArity = 1
		connTemp.LeftArity = connTemp.DefaultLeftArity
		connTemp.DefaultRightArity = shared.MAX_RIGHT_ARITY
		connTemp.RightArity = rArity
		connTemp.TypeName = typeName
		connTemp.Buffer = make(chan messages.SAMessage, rArity)
		connTemp.State = []chan bool{make(chan bool, rArity), make(chan bool, rArity)}
		for i := 0; i < rArity; i++ { // move to initial state
			connTemp.State[0] <- true
		}
		connTemp.Behaviour = behaviour
		return connTemp
	case shared.NTOONE:
		connTemp.Id = id
		connTemp.DefaultLeftArity = shared.MAX_LEFT_ARITY
		connTemp.LeftArity = lArity
		connTemp.DefaultRightArity = 1
		connTemp.RightArity = connTemp.DefaultRightArity
		connTemp.TypeName = typeName
		connTemp.Buffer = make(chan messages.SAMessage, lArity)
		connTemp.State = []chan bool{make(chan bool, lArity), make(chan bool, rArity)}
		for i := 0; i < lArity; i++ {
			connTemp.State[0] <- true
		}
		connTemp.Behaviour = behaviour
		return connTemp
	case shared.NTOONEREQREP: // B = Reqrep1 [] Reqrep2 [] reqrep3....
		connTemp.Id = id
		connTemp.DefaultLeftArity = shared.MAX_LEFT_ARITY
		connTemp.LeftArity = lArity
		connTemp.DefaultRightArity = 1
		connTemp.RightArity = connTemp.DefaultRightArity
		connTemp.TypeName = typeName
		connTemp.Buffer = make(chan messages.SAMessage, lArity)
		connTemp.State = []chan bool{make(chan bool, 1), make(chan bool, 1), make(chan bool, 1), make(chan bool, 1)}
		connTemp.State[0] <- true
		connTemp.Behaviour = behaviour
		return connTemp
	case shared.ONETONREQREP: // B = Reqrep1 -> Reqrep2 -> reqrep3....
		connTemp.Id = id
		connTemp.DefaultLeftArity = 1
		connTemp.LeftArity = connTemp.DefaultLeftArity
		connTemp.DefaultRightArity = shared.MAX_RIGHT_ARITY
		connTemp.RightArity = rArity
		connTemp.TypeName = typeName
		connTemp.Buffer = make(chan messages.SAMessage, rArity)
		connTemp.State = []chan bool{make(chan bool, 1), make(chan bool, 1), make(chan bool, 1), make(chan bool, 1)}
		connTemp.State[0] <- true
		connTemp.Behaviour = behaviour
		return connTemp
	default:
		shared.ErrorHandler(shared.GetFunction(), "Connector type '"+typeName+"' does not exist!!")
	}
	return Connector{TypeName: ""} // Invalid
}

func (c *Connector) Run(msg *messages.SAMessage, act string, invoker string) {
	switch c.TypeName {
	case shared.ONEWAY:
		c.OneWayBehaviour(act, msg, invoker)
	case shared.REQUEST_REPLY:
		c.RequestReplyBehaviour(act, msg, invoker)
	case shared.ONETON:
		c.OneToNBehaviour(act, msg, invoker)
	case shared.NTOONE:
		c.NToOneBehaviour(act, msg, invoker)
	case shared.NTOONEREQREP:
		c.NToOneReqRepBehaviour(act, msg, invoker)
	default:
		shared.ErrorHandler(shared.GetFunction(), "Connector type does not exist")
	}
}

func (c *Connector) OneWayBehaviour(act string, msg *messages.SAMessage, invoker string) {
	switch act {
	case shared.INVR: // Receive from left
		<-c.State[0]
		c.Buffer <- *msg
		c.State[1] <- true
	case shared.INVP: // Send to right
		<-c.State[1]
		*msg = <-c.Buffer
		c.State[0] <- true
	}
}

func (c *Connector) RequestReplyBehaviour(act string, msg *messages.SAMessage, invoker string) {
	switch act {
	case shared.INVR: // Receive from left
		<-c.State[0]
		msg.From = invoker
		c.Buffer <- *msg
		c.State[1] <- true
	case shared.INVP: // Send to right
		<-c.State[1]
		msg.To = invoker
		*msg = <-c.Buffer
		c.State[2] <- true
	case shared.TERP: // Receive from right
		<-c.State[2]
		msg.To = msg.From
		msg.From = invoker
		c.Buffer <- *msg
		c.State[3] <- true
	case shared.TERR: // Send to left
		<-c.State[3]
		*msg = <-c.Buffer
		c.State[0] <- true
	}
}

func (c *Connector) OneToNBehaviour(act string, msg *messages.SAMessage, invoker string) {
	switch act {
	case shared.INVR:
		for i := 0; i < c.RightArity; i++ { // enable INVR
			<-c.State[0]
		}
		for i := 0; i < c.RightArity; i++ { // put msg into connector
			c.Buffer <- *msg
		}
		for i := 0; i < c.RightArity; i++ { // trigger receivers
			c.State[1] <- true
		}
	case shared.INVP: // to right
		<-c.State[1]
		*msg = <-c.Buffer
		c.State[0] <- true // move to next state
	}
}

func (c *Connector) NToOneBehaviour(act string, msg *messages.SAMessage, invoker string) {
	switch act {
	case shared.INVR: // from left
		<-c.State[0]
		c.Buffer <- *msg
		c.State[1] <- true
	case shared.INVP: // to right
		<-c.State[1]
		*msg = <-c.Buffer
		c.State[0] <- true // move to next state
	}
}

func (c *Connector) NToOneReqRepBehaviour(act string, msg *messages.SAMessage, invoker string) {

	switch act {
	case shared.INVR: // Receive from left
		<-c.State[0]
		msg.From = invoker
		c.Buffer <- *msg
		c.State[1] <- true
	case shared.INVP: // Send to right
		<-c.State[1]
		msg.To = invoker
		*msg = <-c.Buffer
		c.State[2] <- true
	case shared.TERP: // Receive from right
		<-c.State[2]
		msg.To = msg.From
		msg.From = invoker
		c.Buffer <- *msg
		c.State[3] <- true
	case shared.TERR: // Send to left
		<-c.State[3]
		*msg = <-c.Buffer
		c.State[0] <- true
	}
}
