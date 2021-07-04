package managing

import (
	"adaptive/adaptiveV6/environment/env"
	"adaptive/adaptiveV6/environment/plugins/manager"
	"adaptive/adaptiveV6/selfadaptivesystem/managed"
	"adaptive/adaptiveV6/sharedadaptive"
	"fmt"
	"github.com/streadway/amqp"
	"reflect"
	"shared"
	"time"
)

type Symptom struct {
	Name string
	Info interface{}
}

type Monitor interface {
	SetManaged(managed.Managed)
	SetMAPEK(MAPEK)
	Start()
	ProbePluginsRemote() interface{}
	ProbeManaged() int
	GetProbes() []func(impl MonitorImpl) interface{}
	GenerateSymptoms(interface{}) []Symptom
}

type MonitorImpl struct {
	Env           env.Environment
	ManagedSystem managed.Managed
	Mapek         MAPEK
	Probes        []func(MonitorImpl) interface{}
	MonitorQueue  amqp.Queue
	Ch            amqp.Channel
}

func NewMonitor(probes []func(MonitorImpl) interface{}) Monitor {
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

func (m MonitorImpl) Start() {
	for {
		// Take a time
		time.Sleep(sharedadaptive.MONITOR_TIME * time.Second)

		// Invoke the probe(s)
		plugins := m.Probes[0](m) // From Environment only TODO

		// Generate symptoms
		symptoms := m.GenerateSymptoms(plugins)
		fmt.Println(symptoms)

		// Publish symptom
		/*
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

		*/
	}
}

func (monitor MonitorImpl) GenerateSymptoms(info interface{}) []Symptom {
	r := []Symptom{}

	if reflect.TypeOf(info).String() == "[]manager.MyPlugin" {
		p := info.([]manager.MyPlugin)
		fmt.Println(p)
		if len(p) > 0 {
			r = append(r, Symptom{Name: "NewPlugin",Info:p[len(p)-1]}) // TODO
		}
	}
	return r
}

func (monitor *MonitorImpl) SetMAPEK(mapek MAPEK) {
	monitor.Mapek = mapek
}

func (monitor *MonitorImpl) SetManaged(ms managed.Managed) {

	monitor.ManagedSystem = ms
}

func (monitor MonitorImpl) GetProbes() []func(MonitorImpl) interface{} {
	return monitor.Probes
}

func (m MonitorImpl) ProbePluginsRemote() interface{} {
	r := m.Env.SensePlugins(sharedadaptive.DIR_SOURCE_PLUGINS_REMOTE)

	return r
}

func (monitor MonitorImpl) ProbeManaged() int {
	r := monitor.ManagedSystem.Sense()

	return r
}
