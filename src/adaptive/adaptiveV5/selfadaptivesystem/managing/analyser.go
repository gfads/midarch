package managing

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"shared"
)

type EventConditionRule struct{}

type Analyser interface {
	Start()
	SetMAPEK(MAPEK)
}

type AnalyserInfo struct {
	Info MonitorInfo // TODO
}

type AnalyserImpl struct {
	Mapek         MAPEK
	Info          AnalyserInfo
	MonitorQueue  amqp.Queue
	AnalyserQueue amqp.Queue
	Ch            amqp.Channel
	ChCons        <-chan amqp.Delivery
}

func NewAnalyser() Analyser {
	var a Analyser

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

	// Declare queues
	monitorQueue, err := ch.QueueDeclare(
		"monitorQueue", false, false, false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	analyserQueue, err := ch.QueueDeclare(
		"analyserQueue", false, false, false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Create a consumer queue to Monitor queue
	consumerChannel, err := ch.Consume(monitorQueue.Name, "", true, false,
		false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	a = &AnalyserImpl{MonitorQueue: monitorQueue, AnalyserQueue: analyserQueue, ChCons: consumerChannel, Ch: *ch,}

	return a
}

func (a *AnalyserImpl) SetMAPEK(mapek MAPEK) {
	a.Mapek = mapek
}

func (a AnalyserImpl) Start() {

	for {
		// Receive monitored data from the Monitor (Monitor queue)
		x := <-a.ChCons

		// Deserialize received data
		monitoredData := MonitorInfo{}
		err := json.Unmarshal(x.Body, &monitoredData)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// Make a simple analysis, i.e., a new plugin is available - TODO
		if len(monitoredData.EnvInfo.Plugins) > 0 {

			// Configure analysed data
			a.Info = AnalyserInfo{Info: monitoredData}

			// Serialise analysed data
			msgBytes, err := json.Marshal(a.Info)
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}

			// Publish serialized data
			err = a.Ch.Publish("", a.AnalyserQueue.Name, false, false,
				amqp.Publishing{ContentType: "text/plain", Body: msgBytes,})
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), err.Error())
			}
		}
	}
}

