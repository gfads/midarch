package pluginUtils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"strings"
	"time"

	"github.com/gfads/midarch/pkg/shared"
)

func LoadPlugins() map[string]time.Time {
	listOfPlugins := make(map[string]time.Time)

	pluginsDir := shared.DIR_PLUGINS
	OSDir, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		fmt.Printf("github.com/gfads/midarch/src/shared:: Folder '%v' is unreadeable\n", pluginsDir)
		os.Exit(0)
	}
	for i := range OSDir {
		fileName := OSDir[i].Name()
		pluginFile := pluginsDir + "/" + fileName
		info, err := os.Stat(pluginFile)
		if err != nil {
			fmt.Printf("github.com/gfads/midarch/src/shared:: Plugin '%v' not readeable\n", pluginFile)
			os.Exit(0)
		}
		listOfPlugins[fileName] = info.ModTime()
	}
	return listOfPlugins
}

func CheckForNewPlugins(listOfOldPlugins map[string]time.Time, listOfNewPlugins map[string]time.Time) []string {
	var newPlugins []string

	// check for new plugins
	for key := range listOfNewPlugins {
		val1, _ := listOfNewPlugins[key]
		val2, ok2 := listOfOldPlugins[key]
		if ok2 {
			if val1.After(val2) { // newer version of an old plugin is available
				newPlugins = append(newPlugins, key)
			}
		} else {
			newPlugins = append(newPlugins, key) // a new plugin is available
		}
	}
	return newPlugins
}

func LoadPlugin(pluginName string) plugin.Plugin {

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := shared.DIR_PLUGINS + "/" + pluginName
	attempts := 0
	for {
		//fmt.Println("pluginUtils.LoadPlugin::will open plugin:", pluginFile)
		plg, err = plugin.Open(pluginFile)
		//fmt.Println("pluginUtils.LoadPlugin::opened plugin:", pluginFile)
		if err != nil {
			fmt.Println("pluginUtils.LoadPlugin::error while opening plugin:", err)
			if attempts >= shared.ATTEMPTS_TO_OPEN_A_PLUGIN { // TODO
				fmt.Printf("github.com/gfads/midarch/src/shared:: Error on trying open plugin '%v' \n", pluginFile)
				os.Exit(0)
			} else {
				attempts++
				time.Sleep(shared.MONITOR_TIME) // TODO
			}
		} else {
			break
		}
	}

	// look for a exportable function/variable within the plugin
	//fx, err := lib.Lookup(symbolName)
	//if err != nil {
	//	fmt.Printf( "github.com/gfads/midarch/src/shared:: Function '%v' not found in plugin '%v'\n",symbolName,pluginName)
	//	os.Exit(0)
	//}
	//return fx

	return *plg
}

func GeneratePlugin(pluginName, versionedPluginName string) {
	// log.Println("Vai ler:", shared.DIR_PLUGINS_SOURCE+"/pluginBuild.model")
	input, err := ioutil.ReadFile(shared.DIR_PLUGINS_SOURCE + "/pluginBuild.model")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pluginType, _, err := shared.GetTypeAndBehaviour(shared.DIR_PLUGINS_SOURCE + "/middleware/" + pluginName + ".go")
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	pluginSourcePath := shared.DIR_PLUGINS_IMPORT + "/middleware" //+ versionedPluginName
	// log.Println("Type and Behaviour:", pluginType)

	// log.Println("pluginSourcePath:", pluginSourcePath)
	//pluginSourcePath := "adaptive/pluginTest/pluginsSrc" + "/" + pluginName
	output := bytes.Replace(input, []byte("<pluginName>"), []byte(pluginSourcePath), -1)
	output = bytes.Replace(output, []byte("<pluginType>"), []byte("middleware."+pluginType+"{}"), -1)

	// log.Println(shared.DIR_PLUGINS_SOURCE + "/" + versionedPluginName + "/main/pluginBuild.go")
	os.Mkdir(shared.DIR_PLUGINS_SOURCE+"/middleware/main", os.ModePerm)
	if err = ioutil.WriteFile(shared.DIR_PLUGINS_SOURCE+"/middleware/main/pluginBuild.go", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	buildPlugin(versionedPluginName)
}

func buildPlugin(plugin string) {
	pOut := shared.DIR_PLUGINS + "/" + plugin + ".so"
	//pIn := shared.DIR_PLUGINS_SOURCE + "/" + source + "/" + source + ".go"
	// pIn := shared.DIR_PLUGINS_SOURCE + "/" + plugin + "/main/pluginBuild.go"
	pIn := shared.DIR_PLUGINS_SOURCE + "/middleware/main/pluginBuild.go"

	//fmt.Println("injector::evolutive.buildPlugin::will build plugin:", source)
	//fmt.Println("command:", shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn)
	_, err := exec.Command("go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput() // If running outside containers add to args: "-gcflags", "all=-N -l",
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Something wrong in generating plugin '"+pIn+"' in \n '"+pOut+"': "+err.Error()+"\n")
	}
	//fmt.Println("injector::evolutive.buildPlugin::plugin built:", source)
}

func RemoveOldPlugins() {
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
