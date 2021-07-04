package managing

import (
	"adaptive/adaptiveV4/environment/env"
	"adaptive/adaptiveV4/selfadaptivesystem/managed"
	"adaptive/adaptiveV4/sharedadaptive"
	"time"
)

type MonitorInfo struct {
	EnvInfo     env.EnvironmentInfo
	ManagedInfo managed.ManagedInfo
}

type Monitor interface {
	SetManaged(managed.Managed)
	SetMAPEK(MAPEK)
	Start()
	ProbeSourceRemote() MonitorInfo
	ProbeExecutableLocal() MonitorInfo
	ProbeManaged() MonitorInfo
	GetProbes() []func(impl MonitorImpl) MonitorInfo
}

type MonitorImpl struct {
	Env           env.Environment
	ManagedSystem managed.Managed
	Mapek         MAPEK
	Probes        []func(MonitorImpl) MonitorInfo
}

func NewMonitor(probes []func(MonitorImpl) MonitorInfo) Monitor {
	var r Monitor

	r = &MonitorImpl{Probes: probes, Env: env.NewEnvironment()}

	return r
}

func (monitor *MonitorImpl) SetMAPEK(mapek MAPEK) {
	monitor.Mapek = mapek
}

func (monitor *MonitorImpl) SetManaged(ms managed.Managed) {
	monitor.ManagedSystem = ms
}

func (monitor MonitorImpl) GetProbes() []func(MonitorImpl) MonitorInfo {
	return monitor.Probes
}

func (monitor MonitorImpl) ProbeExecutableLocal() MonitorInfo {
	r := MonitorInfo{EnvInfo: monitor.Env.Sense(sharedadaptive.EXECUTABLE, sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL)}
	return r
}

func (monitor MonitorImpl) ProbeSourceRemote() MonitorInfo {
	r := MonitorInfo{EnvInfo: monitor.Env.Sense(sharedadaptive.SOURCE, sharedadaptive.DIR_SOURCE_PLUGINS_REMOTE)}
	return r
}

func (monitor MonitorImpl) ProbeManaged() MonitorInfo {
	r := MonitorInfo{ManagedInfo: monitor.ManagedSystem.Sense()}
	return r
}

func (monitor MonitorImpl) Start() {
	for {
		// Take a time
		time.Sleep(sharedadaptive.MONITOR_TIME * time.Second)

		// Invoke the probe(s)
		p1Result := monitor.Probes[0](monitor) // From Environment TODO
		envInfo := p1Result.EnvInfo

		//p2Result := monitor.Probes[1]()  // TODO
		//managedInfo := p2.Result.ManagedInfo  // TODO

		// Configure Monitored information
		//monitorInfo := MonitorInfo{ManagedInfo: managedInfo, EnvInfo: envInfo}
		monitorInfo := MonitorInfo{EnvInfo: envInfo} // TODO

		// Send monitored information
		if len(envInfo.ExecutablePlugins) > 0 { // Local Analyser TODO
			NewAnalyser().ToLAnalyser(monitorInfo)
		} else if len(envInfo.SourcePlugins) > 0 { // Remote Analyser  TODO
			NewAnalyser().ToRAnalyser(monitorInfo)
		}
	}
}
