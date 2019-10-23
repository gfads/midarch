package evolutive

import (
	"fmt"
	"newsolution/shared/parameters"
	"os"
	"os/exec"
	"strings"
	"time"
)

type EvolutiveInjector struct{}

func (EvolutiveInjector) Start(elem string) {

	// remove old plugins
	outputLS, err := exec.Command("/bin/ls", parameters.DIR_PLUGINS).CombinedOutput()
	if err != nil {
		fmt.Printf("Injector:: Something wrong in dir '%v'", parameters.DIR_PLUGINS)
		os.Exit(0)
	}
	oldPlugins := strings.Split(string(outputLS), "\n")
	for plugin := range oldPlugins {
		exec.Command("/bin/rm", parameters.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
	}

	// Strategies for replacing

	//go noChange()
	//go changeOnce(elem)
	//go changeSamePluginSeveralTimes(elem)
	go alternatePlugins(elem)
}

func noChange() {}

func changeOnce(elem string) {
	generatePlugin(elem,elem+"_v1")
}

func changeSamePluginSeveralTimes(elem string) {

	for {
		generatePlugin(elem, elem+"_v1")
		time.Sleep(parameters.INJECTION_TIME * time.Second)
	}
}

func alternatePlugins(elem string) {

	currentPlugin := 1
	for {
		switch currentPlugin {
		case 1: // Plugin 01
			currentPlugin = 2
			generatePlugin(elem+"_v1", elem+"_v1")
		case 2: // Plugin 02
			currentPlugin = 1
			generatePlugin(elem+"_v2", elem+"_v2")
			time.Sleep(parameters.INJECTION_TIME * time.Second)
		}
	}
}

func generatePlugin(source, plugin string) {

	pOut := parameters.DIR_PLUGINS + "/" + plugin
	pIn := parameters.DIR_PLUGINS_SOURCE + "/" + source + "/" + source + ".go"

	_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	if err != nil {
		fmt.Printf("Injector:: Something wrong in generating plugin '%v' '%v' \n", source, plugin)
		os.Exit(0)
	}
}
