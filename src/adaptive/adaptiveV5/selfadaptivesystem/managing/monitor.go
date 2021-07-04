package managing

import (
	"adaptive/adaptiveV5/environment/env"
	"adaptive/adaptiveV5/selfadaptivesystem/managed"
	"adaptive/adaptiveV5/sharedadaptive"
	"encoding/json"
	"github.com/streadway/amqp"
	"shared"
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
	ProbeSourceRemote()
	ProbeManaged() MonitorInfo
	GetProbes() []func(impl MonitorImpl)
}

type MonitorImpl struct {
	Env           env.Environment
	ManagedSystem managed.Managed
	Mapek         MAPEK
	Probes        []func(MonitorImpl)
	MonitorQueue  amqp.Queue
	Ch            amqp.Channel
}

func NewMonitor(probes []func(MonitorImpl)) Monitor {
	var r Monitor

	// Connect to messaging server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Create channels
	ch, err := conn.Channel()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Declare queue
	monitorQueue, err := ch.QueueDeclare(
		"monitorQueue", false, false, false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Probes
	r = &MonitorImpl{Probes: probes, Env: env.NewEnvironment(), Ch: *ch, MonitorQueue: monitorQueue}

	return r
}

func (monitor MonitorImpl) Start() {
	for {
		// Take a time
		time.Sleep(sharedadaptive.MONITOR_TIME * time.Second)

		// Invoke the probe(s)
		monitor.Probes[0](monitor) // From Environment only TODO
	}
}

func (monitor *MonitorImpl) SetMAPEK(mapek MAPEK) {
	monitor.Mapek = mapek
}

func (monitor *MonitorImpl) SetManaged(ms managed.Managed) {

	monitor.ManagedSystem = ms
}

func (monitor MonitorImpl) GetProbes() []func(MonitorImpl) {
	return monitor.Probes
}

func (m MonitorImpl) ProbeSourceRemote() {
	envInfo := m.Env.SensePlugins(sharedadaptive.DIR_SOURCE_PLUGINS_REMOTE)

	// Configure message to be published
	msg := MonitorInfo{EnvInfo: envInfo}

	// Serialise message
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Publish message
	err = m.Ch.Publish("", m.MonitorQueue.Name, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: msgBytes,})
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
}

func (monitor MonitorImpl) ProbeManaged() MonitorInfo {
	managedInfo := monitor.ManagedSystem.Sense()

	r := MonitorInfo{ManagedInfo: managedInfo}
	return r
}
