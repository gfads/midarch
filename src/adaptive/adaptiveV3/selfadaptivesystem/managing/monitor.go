package managing

import (
	"adaptive/adaptiveV3/sharedadaptive"
	"encoding/json"
	"net"
	"plugin"
	"shared"
	"time"
)

type Monitor struct{}

var Ln net.Listener
var CnMonitor net.Conn

func (Monitor) StartLocal(fromEnv chan map[string]plugin.Plugin, fromManaged chan int, toAnalyser chan InfoToAnalyser) {
	for {
		select {
		case plugins := <-fromEnv:
			toAnalyser <- InfoToAnalyser{Source: sharedadaptive.FROM_ENV, FromLocalEnv: plugins}
		case n := <-fromManaged:
			toAnalyser <- InfoToAnalyser{Source: sharedadaptive.FROM_MANAGED, FromManaged: n}
		}
		time.Sleep(sharedadaptive.MONITOR_TIME * time.Second)
	}
}

func (Monitor) StartRemote(fromEnv chan map[string]string) {
	for {
		select {
		case plugins := <-fromEnv:
			toAnalyser(InfoToAnalyser{Source: sharedadaptive.FROM_ENV, FromRemoteEnv: plugins})
		}
		time.Sleep(sharedadaptive.MONITOR_TIME * time.Second)
	}
}

func toAnalyser(info InfoToAnalyser) {

	// Create a listner if not created yet
	if Ln == nil {
		// create listener
		servAddr, err := net.ResolveTCPAddr("tcp", "localhost:1313")
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		Ln, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// Accept connections
		CnMonitor, err = Ln.Accept()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	encoder := json.NewEncoder(CnMonitor)

	// send data
	err := encoder.Encode(info)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
}
