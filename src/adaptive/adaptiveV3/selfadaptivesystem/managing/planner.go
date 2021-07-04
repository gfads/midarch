package managing

import (
	"adaptive/adaptiveV3/sharedadaptive"
	"plugin"
	"shared"
)

type Planner struct{}

type AdaptationPlan struct {
	Source  int
	Command string
	Params  interface{}
}

func (Planner) Start(fromAnalyser chan InfoToAnalyser, toExecutor chan AdaptationPlan) {
	for {
		info := <-fromAnalyser
		switch info.Source {
		case sharedadaptive.FROM_ENV:
			if len(info.FromLocalEnv) > 0 {

				// Take the last most recent plugin
				allNewPlugins := info.FromLocalEnv
				mostRecentPlugin := plugin.Plugin{}
				mostRecent := ""
				for i := range allNewPlugins {
					if i > mostRecent {
						mostRecent = i
						mostRecentPlugin = allNewPlugins[i]
					}
				}

				// Configure plan
				plan := AdaptationPlan{Source: sharedadaptive.FROM_ENV, Command: sharedadaptive.CMD_UPDATE, Params: mostRecentPlugin}
				toExecutor <- plan
			}
		case sharedadaptive.FROM_MANAGED:
			shared.ErrorHandler(shared.GetFunction(), "TODO")
		}
	}
}
