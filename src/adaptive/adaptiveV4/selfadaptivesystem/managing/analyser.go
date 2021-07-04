package managing

import (
	"adaptive/adaptiveV4/environment/env"
	"adaptive/adaptiveV4/environment/plugins/manager"
	"adaptive/adaptiveV4/sharedadaptive"
	"encoding/json"
	"net"
	"shared"
)

var Ln1, Ln2 net.Listener
var Cn1, Cn2 net.Conn

type Analyser interface {
	Start()
	SetMAPEK(MAPEK)
	ToLAnalyser(MonitorInfo)
	ToRAnalyser(MonitorInfo)
}

type AnalyserInfo struct {
	Info MonitorInfo // TODO
}

type AnalyserImpl struct {
	Mapek MAPEK
	Info  AnalyserInfo
}

func NewAnalyser() Analyser {
	var a Analyser

	a = &AnalyserImpl{}

	return a
}

func (a *AnalyserImpl) SetMAPEK(mapek MAPEK) {
	a.Mapek = mapek
}

func (a AnalyserImpl) Start() {

	for {
		// Create a listener if not created yet
		if Ln1 == nil {
			// create listener
			servAddr, err := net.ResolveTCPAddr("tcp", "localhost:1313")
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
			Ln1, err = net.ListenTCP("tcp", servAddr)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}

			// Accept connections
			Cn1, err = Ln1.Accept()
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}

		encoder := json.NewDecoder(Cn1)

		// Receive monitored data from remote monitor
		infoSources := map[string]string{}
		err := encoder.Decode(&infoSources)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		envInfo := env.EnvironmentInfo{SourcePlugins: infoSources} // Source plugins only
		info := MonitorInfo{EnvInfo: envInfo}
		if len(info.EnvInfo.SourcePlugins) > 0 {
			for i := range info.EnvInfo.SourcePlugins {
				manager.MyPlugin{}.SaveSource(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, i, info.EnvInfo.SourcePlugins[i])
				manager.MyPlugin{}.GenerateExecutable(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL, i)
			}
		}

		executableLocalPlugins := manager.MyPlugin{}.LoadExecutables(sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL)
		info.EnvInfo.ExecutablePlugins = executableLocalPlugins

		a.Info = AnalyserInfo{Info: info}

		// Send to planner
		p := a.Mapek.P
		if len(a.Info.Info.EnvInfo.ExecutablePlugins) != 0 || len(a.Info.Info.EnvInfo.SourcePlugins) != 0 {
			// Send to Planner
			p.ToPlanner(a.Info)
		}
	}
}

func (a *AnalyserImpl) ToLAnalyser(info MonitorInfo) {

	// Receive from Local Monitor
	a.Info = AnalyserInfo{Info: info}
	if len(info.EnvInfo.ExecutablePlugins) != 0 || len(info.EnvInfo.SourcePlugins) != 0 {
		// Send to Planner
		a.Mapek.P.ToPlanner(a.Info)
	}
}

func (a *AnalyserImpl) ToRAnalyser(info MonitorInfo) {

	// Resolve server address
	addr := "localhost:1313"
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	//  Create connection (if not created yet)
	if Cn2 == nil {
		Cn2, err = net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		//defer conn.Close()
	}

	d := json.NewEncoder(Cn2)
	err = d.Encode(info.EnvInfo.SourcePlugins)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
}
