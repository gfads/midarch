package managing

import (
	"adaptive/adaptiveV2/sharedadaptive"
	"shared"
)

type Planner struct{}

type AdaptationPlan struct{
	Source int
	Command string
	Params interface{}
}

func (Planner) Start(fromAnalyser chan InfoToAnalyser, toExecutor chan AdaptationPlan) {
	for {
		info := <-fromAnalyser
		switch info.Source {
		case sharedadaptive.FROM_ENV:
			if len(info.FromEnv) > 0{
				plan := AdaptationPlan{Source: sharedadaptive.FROM_ENV,Command: sharedadaptive.CMD_UPDATE,Params:info.FromEnv}
				toExecutor <- plan
			}
		case sharedadaptive.FROM_MANAGED:
			shared.ErrorHandler(shared.GetFunction(),"TODO")
		}
	}
}
