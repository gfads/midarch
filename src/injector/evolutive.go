package evolutive

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"shared"
	"strings"
	"time"
)

type EvolutiveInjector struct{}

func (EvolutiveInjector) Start(firstElem, secondElem string, interval time.Duration) {
	// Replacing strategies
	//go noChange()
	//go changeOnce(firstElem, interval)
	//go changeSamePluginSeveralTimes(elem)
	go alternatePlugins(firstElem, secondElem, interval)
}

func noChange() {}

func changeOnce(elem string, interval time.Duration) {
	//removeOldPlugins()
	time.Sleep(interval)
	elemNew := elem + "_v1"
	GeneratePlugin(elemNew, elem, elemNew)
}

func changeSamePluginSeveralTimes(elem string) {

	for {
		removeOldPlugins()
		generatePlugin(elem, elem+"_v1")
		time.Sleep(shared.INJECTION_TIME * time.Second)
	}
}

func alternatePlugins(firstElem, secondElem string, interval time.Duration) {
	//removeOldPlugins()

	elemNew := ""
	elemOld := ""

	currentPlugin := 1
	for {
		fmt.Printf("Evolutive:: Next plugin '%v' will be generated in %v !! \n", elemNew, interval)
		time.Sleep(interval)

		switch currentPlugin {
		case 1: // Plugin 01
			currentPlugin = 2
			elemOld = firstElem + "_v2"
			elemNew = firstElem + "_v2"
			GeneratePlugin(elemOld, firstElem, elemNew)
		case 2: // Plugin 02
			currentPlugin = 1
			elemOld = secondElem + "_v1"
			elemNew = secondElem + "_v1"
			GeneratePlugin(elemOld, secondElem, elemNew)
		}
	}
}

func GeneratePlugin(source, pluginName, versionedPluginName string) {
	//log.Println("Vai ler:", shared.DIR_PLUGINS_SOURCE+"/pluginBuild.model")
	input, err := ioutil.ReadFile(shared.DIR_PLUGINS_SOURCE + "/pluginBuild.model")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pluginType, _ := shared.GetTypeAndBehaviour(shared.DIR_PLUGINS_SOURCE + "/" + versionedPluginName + "/" + pluginName + ".go")
	pluginSourcePath := shared.DIR_PLUGINS_IMPORT + "/" + versionedPluginName

	//log.Println("pluginSourcePath:", pluginSourcePath)
	//pluginSourcePath := "adaptive/pluginTest/pluginsSrc" + "/" + pluginName
	output := bytes.Replace(input, []byte("<pluginName>"), []byte(pluginSourcePath), -1)
	output = bytes.Replace(output, []byte("<pluginType>"), []byte(pluginName+"."+pluginType+"{}"), -1)

	//log.Println("Vai gravar:", shared.DIR_PLUGINS_SOURCE+"/pluginBuild.go")
	os.Mkdir(shared.DIR_PLUGINS_SOURCE+"/"+versionedPluginName+"/main", os.ModePerm)
	if err = ioutil.WriteFile(shared.DIR_PLUGINS_SOURCE+"/"+versionedPluginName+"/main/pluginBuild.go", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generatePlugin(source, versionedPluginName)
}

func generatePlugin(source, plugin string) {
	pOut := shared.DIR_PLUGINS + "/" + plugin + ".so"
	//pIn := shared.DIR_PLUGINS_SOURCE + "/" + source + "/" + source + ".go"
	pIn := shared.DIR_PLUGINS_SOURCE + "/" + plugin + "/main/pluginBuild.go"

	//fmt.Println("injector::evolutive.generatePlugin::will build plugin:", source)
	//fmt.Println("command:", shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn)
	_, err := exec.Command(shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput() //"-gcflags", "all=-N -l",
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Something wrong in generating plugin '"+pIn+"' in \n '"+pOut+"': "+err.Error()+"\n")
	}
	//fmt.Println("injector::evolutive.generatePlugin::plugin built:", source)
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
