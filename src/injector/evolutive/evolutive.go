package evolutive

import (
	"fmt"
	"os"
	"os/exec"
	"shared"
	"strings"
	"time"
)

type EvolutiveInjector struct{}

func (EvolutiveInjector) Start(elem string) {

	// Remove old plugins
	outputLS, err := exec.Command("/bin/ls", shared.DIR_PLUGINS).CombinedOutput()
	if err != nil {
		fmt.Printf("Injector:: Something wrong in dir '%v'", shared.DIR_PLUGINS)
		os.Exit(0)
	}
	oldPlugins := strings.Split(string(outputLS), "\n")
	for plugin := range oldPlugins {
		exec.Command("/bin/rm", shared.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
	}

	// Strategies for replacing

	//go noChange()
	//go changeOnce(elem)
	//go changeSamePluginSeveralTimes(elem)
	go alternatePlugins(elem)
}

func noChange() {}

func changeOnce(elem string) {
	generatePlugin(elem, elem+"_v1")
}

func changeSamePluginSeveralTimes(elem string) {

	for {
		generatePlugin(elem, elem+"_v1")
		time.Sleep(shared.INJECTION_TIME * time.Second)
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
			//time.Sleep(shared.INJECTION_TIME * time.Second)
		}
		time.Sleep(1 * time.Second)
	}
}

func generatePlugin(source, plugin string) {

	pOut := shared.DIR_PLUGINS + "/" + plugin
	pIn := shared.DIR_PLUGINS_SOURCE + "/" + source + "/" + source + ".go"

	_, err := exec.Command(shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	if err != nil {
		fmt.Printf("Injector:: Something wrong in generating plugin '%v' in \n '%v' \n", pIn, pOut)
		os.Exit(0)
	}
}
