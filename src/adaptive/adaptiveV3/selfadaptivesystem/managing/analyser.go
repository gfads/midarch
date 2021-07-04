package managing

import (
	"adaptive/adaptiveV3/environment/plugins/manager"
	"adaptive/adaptiveV3/sharedadaptive"
	"encoding/json"
	"net"
	"shared"
)

type Analyser struct{}

var CnAnalyser net.Conn

func (Analyser) StartLocal(fromMonitor chan InfoToAnalyser, toPlanner chan InfoToAnalyser) {
	for {
		info := <-fromMonitor
		switch info.Source {
		case sharedadaptive.FROM_ENV:
			if len(info.FromLocalEnv) != 0 { // New plugin available
				toPlanner <- info
			}
		case sharedadaptive.FROM_MANAGED:
			// TODO
		}
	}
}

func (Analyser) StartRemote(fromMonitor chan InfoToAnalyser, toPlanner chan InfoToAnalyser) {
	for {
		info := fromRemoteMonitor()
		switch info.Source {
		case sharedadaptive.FROM_ENV:
			if len(info.FromLocalEnv) != 0 { // New plugin available
				toPlanner <- info
			}
		case sharedadaptive.FROM_MANAGED:
			// TODO
		}
	}
}

func fromRemoteMonitor() InfoToAnalyser {
	r := InfoToAnalyser{}

	// Resolve server address
	addr := "localhost:1313"
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	//  Create connection (if not create yet)
	if CnAnalyser == nil {
		CnAnalyser, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		//defer conn.Close()
	}

	d := json.NewDecoder(CnAnalyser)
	err = d.Decode(&r)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	if len(r.FromRemoteEnv) > 0 {
		for i := range r.FromRemoteEnv {
			manager.MyPlugin{}.SaveSource(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, i, r.FromRemoteEnv[i])
			manager.MyPlugin{}.GenerateExecutable(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL, i)
		}
	}

	executableLocalPlugins := manager.MyPlugin{}.LoadExecutables(sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL)
	r.FromLocalEnv = executableLocalPlugins

	return r
}
