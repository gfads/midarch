package managing

import (
	"adaptive/adaptiveV6/sharedadaptive"
	"encoding/json"
	"github.com/streadway/amqp"
	"shared"
)

type PlannerInfo struct { // Adaptation plan
	Source  int
	Command string
	Params  interface{}
}

type Planner interface {
	SetMAPEK(MAPEK)
	Start()
}

type PlannerImpl struct {
	Mapek         MAPEK
	Info          PlannerInfo
	AnalyserQueue amqp.Queue
	PlannerQueue  amqp.Queue
	Ch            amqp.Channel
	ChCons        <-chan amqp.Delivery
}

func NewPlanner() Planner {
	var p Planner

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
	plannerQueue, err := ch.QueueDeclare(
		"plannerQueue", false, false, false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	analyserQueue, err := ch.QueueDeclare(
		"analyserQueue", false, false, false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Create consumer to Analyser queue
	consumerChannel, err := ch.Consume(analyserQueue.Name, "", true, false,
		false, false, nil, )

	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	p = &PlannerImpl{AnalyserQueue: analyserQueue, PlannerQueue: plannerQueue, Ch: *ch, ChCons: consumerChannel}

	return p
}

func (p *PlannerImpl) SetMAPEK(mapek MAPEK) {
	p.Mapek = mapek
}

func (p PlannerImpl) Start() {

	for {
		// Receive analysed data from the Analyser (Analyser queue)
		x := <-p.ChCons

		// Deserialise received data
		analysedData := AnalyserInfo{}
		err := json.Unmarshal(x.Body, &analysedData)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// Create adaptation plan - TODO
		p.Info = PlannerInfo{Source: sharedadaptive.FROM_ENV, Command: sharedadaptive.CMD_UPDATE, Params: analysedData}

		// Serialise adaptation plan
		msgBytes, err := json.Marshal(p.Info)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// Publish adaptation plan
		err = p.Ch.Publish("", p.PlannerQueue.Name, false, false,
			amqp.Publishing{ContentType: "text/plain", Body: msgBytes,})
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
}
