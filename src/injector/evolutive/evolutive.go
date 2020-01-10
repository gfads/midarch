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

func (EvolutiveInjector) Start(elem string, interval time.Duration) {

	// Replacing strategies
	//go noChange()
	//go changeOnce(elem)
	//go changeSamePluginSeveralTimes(elem)
	go alternatePlugins(elem, interval)
}

func noChange() {}

func changeOnce(elem string) {
	removeOldPlugins()
	generatePlugin(elem, elem+"_v1")
}

func changeSamePluginSeveralTimes(elem string) {

	for {
		removeOldPlugins()
		generatePlugin(elem, elem+"_v1")
		time.Sleep(shared.INJECTION_TIME * time.Second)
	}
}

func alternatePlugins(elem string, interval time.Duration) {

	removeOldPlugins()

	elemNew := ""
	elemOld := ""

	currentPlugin := 1
	for {
		switch currentPlugin {
		case 1: // Plugin 01
			currentPlugin = 2
			elemOld = elem+"_v1"
			elemNew = elem+"_v2"
			generatePlugin(elemOld, elemNew)
		case 2: // Plugin 02
			currentPlugin = 1
			elemOld = elem+"_v2"
			elemNew = elem+"_v1"
			generatePlugin(elemOld, elemNew)
		}

		fmt.Printf("Evolutive:: Next plugin '%v' will generated in %v !! \n",elemNew,interval)
		time.Sleep(interval)
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

func removeOldPlugins() {
	outputLS, err := exec.Command("/bin/ls", shared.DIR_PLUGINS).CombinedOutput()
	if err != nil {
		fmt.Printf("Injector:: Something wrong in dir '%v'", shared.DIR_PLUGINS)
		os.Exit(0)
	}
	oldPlugins := strings.Split(string(outputLS), "\n")

	for plugin := range oldPlugins {
		if strings.TrimSpace(oldPlugins[plugin]) != "" {
			_, err = exec.Command("/bin/rm", shared.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
			if err != nil {
				fmt.Printf("Injector:: Something wrong in removing the plugins at '%v' %v", shared.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin]), err)
				os.Exit(0)
			}
		}
	}
}
