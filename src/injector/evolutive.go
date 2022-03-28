package evolutive

import (
	"bytes"
	"fmt"
	"gmidarch/development/repositories/architectural"
	"io/ioutil"
	"log"
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

	//removeOldPlugins()

	elemNew := ""
	elemOld := ""

	currentPlugin := 1
	for {
		switch currentPlugin {
		case 1: // Plugin 01
			currentPlugin = 2
			elemOld = elem + "_v1"
			elemNew = elem + "_v1"
			log.Println("Teste")
			GeneratePlugin(elemOld, elem, elemNew)
		case 2: // Plugin 02
			currentPlugin = 1
			elemOld = elem + "_v2"
			elemNew = elem + "_v2"
			GeneratePlugin(elemOld, elem, elemNew)
		}

		fmt.Printf("Evolutive:: Next plugin '%v' will be generated in %v !! \n", elemNew, interval)
		time.Sleep(interval)
	}
}

func GeneratePlugin(source, pluginName, versionedPluginName string) {
	log.Println("Vai ler:", shared.DIR_PLUGINS_SOURCE+"/pluginBuild.model")
	input, err := ioutil.ReadFile(shared.DIR_PLUGINS_SOURCE + "/pluginBuild.model")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pluginType, _ := architectural.GetTypeAndBehaviour(shared.DIR_PLUGINS_SOURCE + "/" + versionedPluginName + "/" + pluginName + ".go")
	pluginSourcePath := shared.DIR_PLUGINS_IMPORT + "/" + versionedPluginName

	log.Println("pluginSourcePath:", pluginSourcePath)
	//pluginSourcePath := "adaptive/pluginTest/pluginsSrc" + "/" + pluginName
	output := bytes.Replace(input, []byte("<pluginName>"), []byte(pluginSourcePath), -1)
	output = bytes.Replace(output, []byte("<pluginType>"), []byte(pluginName+"."+pluginType+"{}"), -1)

	log.Println("Vai gravar:", shared.DIR_PLUGINS_SOURCE+"/pluginBuild.go")
	if err = ioutil.WriteFile(shared.DIR_PLUGINS_SOURCE+"/pluginBuild.go", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generatePlugin(source, versionedPluginName)
}

func generatePlugin(source, plugin string) {

	pOut := shared.DIR_PLUGINS + "/" + plugin + ".so"
	//pIn := shared.DIR_PLUGINS_SOURCE + "/" + source + "/" + source + ".go"
	pIn := shared.DIR_PLUGINS_SOURCE + "/pluginBuild.go"

	fmt.Println("injector::evolutive.generatePlugin::will build plugin:", source)
	fmt.Println("command:", shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn)
	_, err := exec.Command(shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Something wrong in generating plugin '"+pIn+"' in \n '"+pOut+"': "+err.Error()+"\n")
	}
	fmt.Println("injector::evolutive.generatePlugin::plugin built:", source)
}

func removeOldPlugins() {
	outputLS, err := exec.Command("/bin/ls", shared.DIR_PLUGINS).CombinedOutput()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Something wrong in dir '"+shared.DIR_PLUGINS)
	}
	oldPlugins := strings.Split(string(outputLS), "\n")

	for plugin := range oldPlugins {
		if strings.TrimSpace(oldPlugins[plugin]) != "" {
			_, err = exec.Command("/bin/rm", shared.DIR_PLUGINS+"/"+strings.TrimSpace(oldPlugins[plugin])).CombinedOutput()
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(), "Something is wrong in removing the pluginsSrc at '"+shared.DIR_PLUGINS+"' "+err.Error()+"/")
			}
		}
	}
}
