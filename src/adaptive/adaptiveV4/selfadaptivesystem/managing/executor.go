package managing

import (
	"adaptive/adaptiveV4/selfadaptivesystem/managed"
	"adaptive/adaptiveV4/sharedadaptive"
	"plugin"
	"shared"
)

type Executor interface {
	ToExecutor(PlannerInfo)
	SetManaged(managed.Managed)
	SetMAPEK(MAPEK)
}

type ExecutorImpl struct {
	Mapek   MAPEK
	Managed managed.Managed
}

func NewExecutor() Executor {
	var e Executor

	e = &ExecutorImpl{}

	return e
}

func (e *ExecutorImpl) SetManaged(ms managed.Managed) {
	e.Managed = ms
}

func (e *ExecutorImpl) SetMAPEK(mapek MAPEK) {
	e.Mapek = mapek
}

func (e ExecutorImpl) ToExecutor(info PlannerInfo) {
	var s plugin.Symbol

	switch info.Command {
	case sharedadaptive.CMD_UPDATE:
		if info.Source == sharedadaptive.FROM_ENV {
			analyserInfo := info.Params.(AnalyserInfo)
			monitorInfo := analyserInfo.Info
			envInfo := monitorInfo.EnvInfo

			managed := e.Managed
			executablePlugins := envInfo.ExecutablePlugins

			if len(executablePlugins) > 0 {
				// Take the most recente plugin
				last := ""
				for i := range executablePlugins {
					if i > last {
						last = i
						p := executablePlugins[i]
						s, _ = p.Lookup("Behaviour")
					}
				}
				managed.Adapt(s.(func()))
			}
		} else {
			shared.ErrorHandler(shared.GetFunction(), "TODO")
		}
	default:
		shared.ErrorHandler(shared.GetFunction(), "TODO")
	}
}
