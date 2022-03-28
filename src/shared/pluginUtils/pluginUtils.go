package pluginUtils

import (
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"shared"
	"time"
)

func LoadPlugins() map[string]time.Time {
	listOfPlugins := make(map[string]time.Time)

	pluginsDir := shared.DIR_PLUGINS
	OSDir, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		fmt.Printf("Shared:: Folder '%v' is unreadeable\n", pluginsDir)
		os.Exit(0)
	}
	for i := range OSDir {
		fileName := OSDir[i].Name()
		pluginFile := pluginsDir + "/" + fileName
		info, err := os.Stat(pluginFile)
		if err != nil {
			fmt.Printf("Shared:: Plugin '%v' not readeable\n", pluginFile)
			os.Exit(0)
		}
		listOfPlugins[fileName] = info.ModTime()
	}
	return listOfPlugins
}

func CheckForNewPlugins(listOfOldPlugins map[string]time.Time, listOfNewPlugins map[string]time.Time) [] string {
	var newPlugins [] string

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

func LoadPlugin(pluginName string) (plugin.Plugin) {

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := shared.DIR_PLUGINS + "/" + pluginName
	attempts := 0
	for {
		fmt.Println("pluginUtils.LoadPlugin::will open plugin:", pluginFile)
		plg, err = plugin.Open(pluginFile)
		fmt.Println("pluginUtils.LoadPlugin::opened plugin:", pluginFile)
		if err != nil {
			fmt.Println("pluginUtils.LoadPlugin::error while opening plugin:", err)
			if attempts >= shared.ATTEMPTS_TO_OPEN_A_PLUGIN { // TODO
				fmt.Printf("Shared:: Error on trying open plugin '%v' \n", pluginFile)
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
	//	fmt.Printf( "Shared:: Function '%v' not found in plugin '%v'\n",symbolName,pluginName)
	//	os.Exit(0)
	//}
	//return fx

	return *plg
}