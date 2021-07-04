package managing

import "plugin"

type InfoToAnalyser struct {
	Source        int
	FromLocalEnv  map[string]plugin.Plugin
	FromRemoteEnv map[string]string
	FromManaged   int
}
