package managing

import (
	"adaptive/adaptiveV5/environment/plugins/manager"
	"adaptive/adaptiveV5/selfadaptivesystem/managed"
	"adaptive/adaptiveV5/sharedadaptive"
	"encoding/json"
	"github.com/streadway/amqp"
	"plugin"
	"shared"
)

type Executor interface {
	SetManaged(managed.Managed)
	SetMAPEK(MAPEK)
	Start()
}

type ExecutorImpl struct {
	Mapek        MAPEK
	Managed      managed.Managed
	PlannerQueue amqp.Queue
	Ch           amqp.Channel
	ChCons       <-chan amqp.Delivery
}

func NewExecutor() Executor {
	var r Executor

	// Connect to messaging server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Declare queue
	plannerQueue, err := ch.QueueDeclare(
		"plannerQueue", false, false, false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	// Create consumer to Planner queue
	consumerChannel, err := ch.Consume(plannerQueue.Name, "", true, false,
		false, false, nil, )
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	r = &ExecutorImpl{PlannerQueue: plannerQueue, Ch: *ch, ChCons: consumerChannel}

	return r
}

func (e ExecutorImpl) Start() {
	for {
		// Receive adaptation plan (Planner queue)
		x := <-e.ChCons

		// Deserialise adaptation plan
		plannerData := PlannerInfo{}
		err := json.Unmarshal(x.Body, &plannerData)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		// Check actions of the plan
		switch plannerData.Command { // Single action TODO
		case sharedadaptive.CMD_UPDATE:
			if plannerData.Source == sharedadaptive.FROM_ENV {
				// Recover source plugins
				x0 := plannerData.Params.(map[string]interface{})
				x1 := x0["Info"].(map[string]interface{})
				x2 := x1["EnvInfo"].(map[string]interface{})
				sourcePlugins := x2["Plugins"].(map[string]interface{})

				// Take the most recent plugin
				last := ""
				for i := range sourcePlugins {
					if i > last {
						last = i
					}
				}

				// Generate executable plugin
				manager.MyPlugin{}.SaveSource(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, last, sourcePlugins[last].(string))
				manager.MyPlugin{}.GenerateExecutable(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL, sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL, last)

				// Load executable plugins
				var sym plugin.Symbol
				executablePlugins := manager.MyPlugin{}.LoadExecutables(sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL)

				// Take the most recent executable plugin
				if len(executablePlugins) > 0 {
					// Take the most recente plugin
					last := ""
					for i := range executablePlugins {
						if i > last {
							last = i
							p := executablePlugins[i]
							sym, _ = p.Lookup("Behaviour")
						}
					}

					// Adapt managed system
					e.Managed.Adapt(sym.(func()))
				}
			}
		default:
			shared.ErrorHandler(shared.GetFunction(), "TODO")
		}
	}
}

func (e *ExecutorImpl) SetManaged(ms managed.Managed) {
	e.Managed = ms
}

func (e *ExecutorImpl) SetMAPEK(mapek MAPEK) {
	e.Mapek = mapek
}
