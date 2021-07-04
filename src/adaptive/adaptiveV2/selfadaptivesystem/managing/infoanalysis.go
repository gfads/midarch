package managing

import "plugin"

type InfoToAnalyser struct {
	Source      int
	FromEnv     map[string]plugin.Plugin
	FromManaged int
}
