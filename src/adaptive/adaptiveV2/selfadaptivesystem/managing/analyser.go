package managing

import (
	"adaptive/adaptiveV2/sharedadaptive"
)

type Analyser struct{}

func (Analyser) Start(fromMonitor chan InfoToAnalyser, toPlanner chan InfoToAnalyser) {
	for {
		info := <-fromMonitor
		switch info.Source {
		case sharedadaptive.FROM_ENV:
			if len(info.FromEnv) != 0 { // New plugin available
				toPlanner <- info
			}
		case sharedadaptive.FROM_MANAGED:
			// TODO
		}
	}
}

