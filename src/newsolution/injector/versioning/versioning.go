package versioning

import (
	"fmt"
	"newsolution/shared/parameters"
	"os"
	"os/exec"
	"strings"
	"time"
)

type VersioningInjector struct{}

func (VersioningInjector) Start(conf, elem string) {

	// remove old plugins
	outputLS, err := exec.Command("/bin/ls", parameters.DIR_PLUGINS).CombinedOutput()
	if err != nil {
		fmt.Printf("Injector:: Something wrong in dir '%v'", parameters.DIR_PLUGINS)
		os.Exit(0)
	}
	oldPlugins := strings.Split(string(outputLS), "\n")

	for plugin := range oldPlugins {
		//exec.Command("/bin/rm", "-r", parameters.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
		exec.Command("/bin/rm", parameters.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
	}

	//go noChange()
	go changeOnce(conf, elem)
	//go changeSamePlugin(plugin01, dirPlugins, source01)
	//go altenatePlugins(plugin01, plugin02, dirPlugins, source01, source02)
}

func noChange() {}

func changeOnce(conf, elem string) {
	generatePlugin("receiver","receiver_v1")
}

func changeSamePlugin(conf, elem string) {

	// configure plugin names
	plugin01 := elem + "01"
	source01 := plugin01 + ".go"

	for {
		pluginName := strings.TrimSpace(plugin01 + "_v1")
		_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", parameters.DIR_PLUGINS+"/"+pluginName, parameters.DIR_PLUGINS_SOURCE+"/"+source01).CombinedOutput()
		if err != nil {
			fmt.Printf("Injector:: Something wrong in generating plugin '%v' \n", pluginName)
			os.Exit(0)
		}
		time.Sleep(parameters.INJECTION_TIME * time.Second)
	}
}

func alternatePlugins(conf, elem string) {

	// configure plugin names
	plugin01 := elem + "01"
	source01 := plugin01 + ".go"

	plugin02 := elem + "02"
	source02 := plugin02 + ".go"

	currentPlugin := 1
	for {
		switch currentPlugin {
		case 1: // Plugin 01
			currentPlugin = 2
			pluginName := strings.TrimSpace(plugin01 + "_v1")
			_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", parameters.DIR_PLUGINS+"/"+pluginName, parameters.DIR_PLUGINS_SOURCE+"/"+source01).CombinedOutput()
			if err != nil {
				fmt.Printf("Injector:: Something is wrong in generating plugin '%v'", pluginName)
				os.Exit(0)
			}
			fmt.Println("Injector:: [PLUGIN 01 GENERATED]")
		case 2: // Plugin 02
			currentPlugin = 1
			pluginName := strings.TrimSpace(plugin02 + "_v1")
			_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", parameters.DIR_PLUGINS+"/"+pluginName, parameters.DIR_PLUGINS_SOURCE+"/"+source02).CombinedOutput()
			if err != nil {
				fmt.Printf("Injector:: Something is wrong in generating plugin '%v' in \n", pluginName)
				os.Exit(0)
			}
			fmt.Println("Injector:: [PLUGIN 02 GENERATED]")
		}
		time.Sleep(parameters.INJECTION_TIME * time.Second)
	}
}

func generatePlugin(source, plugin string) {

	pOut := parameters.DIR_PLUGINS + "/" + plugin
	pIn := parameters.DIR_PLUGINS_SOURCE + "/" + source + ".go"

	_, err := exec.Command(parameters.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	if err != nil {
		fmt.Printf("Injector:: Something wrong in generating plugin '%v' '%v' \n", source,plugin)
		os.Exit(0)
	}
}
