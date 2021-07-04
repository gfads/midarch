package managing

import (
	"adaptive/adaptiveV3/sharedadaptive"
	"plugin"
	"shared"
)

type Executor struct{}

func (Executor) Start(fromPlanner chan AdaptationPlan, toManaged chan func()) {
	var f plugin.Symbol

	for {
		plan := <-fromPlanner
		switch plan.Command {
		case sharedadaptive.CMD_UPDATE:
			if plan.Source == sharedadaptive.FROM_ENV {
				p := plan.Params.(plugin.Plugin)
				f, _ = p.Lookup("Behaviour")
				toManaged <- f.(func())
			} else {
				shared.ErrorHandler(shared.GetFunction(), "TODO")
			}
		default:
			shared.ErrorHandler(shared.GetFunction(), "TODO")
		}
	}
}
