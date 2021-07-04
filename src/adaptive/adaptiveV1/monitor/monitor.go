package evolutive

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"shared"
	"strings"
	"time"
)

type Monitor struct{}

func (Monitor) LoadFuncs() []func() {
		sourcePluginFiles := LoadSources()
		GenerateExecutable(sourcePluginFiles)

		fns := []func(){}
		for i := 0; i < len(sourcePluginFiles); i++ {
			p := sourcePluginFiles[i]
			plugin := LoadPlugin(p[strings.LastIndex(p,"/")+1:strings.LastIndex(p,".")])
			f,err := plugin.Lookup("Behaviour")
			if err != nil {
				shared.ErrorHandler(shared.GetFunction(),"Function not found in plugin!!")
			}
			fns = append(fns, f.(func()))
		}
		return fns
}

func LoadPlugin(pluginName string) plugin.Plugin {

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := "/Users/nsr/Dropbox/go/gmidarch-v14/src/adaptive/adaptiveV1/plugins/executable" + "/" + pluginName
	attempts := 0
	for {
		plg, err = plugin.Open(pluginFile)

		if err != nil {
			if attempts >= 3 { // TODO
				fmt.Printf("Shared:: Error on trying open plugin '%v' \n", pluginFile)
				os.Exit(0)
			} else {
				attempts++
				time.Sleep(10 * time.Millisecond) // TODO
			}
		} else {
			break
		}
	}

	return *plg
}

func LoadSources() []string {
	r := []string{}

	folders, err1 := ioutil.ReadDir("/Users/nsr/Dropbox/go/gmidarch-v14/src/adaptive/adaptiveV1/plugins/source")
	if err1 != nil {
		shared.ErrorHandler(shared.GetFunction(), err1.Error())
	}

	temp := []os.FileInfo{}

	for folder := range folders {
		temp, err1 = ioutil.ReadDir("/Users/nsr/Dropbox/go/gmidarch-v14/src/adaptive/adaptiveV1/plugins/source/" + folders[folder].Name())
		if err1 != nil {
			shared.ErrorHandler(shared.GetFunction(), err1.Error())
		}

		for file := range temp {
			fullPathName := "/Users/nsr/Dropbox/go/gmidarch-v14/src/adaptive/adaptiveV1/plugins/source/" + folders[folder].Name() + "/" + temp[file].Name()
			r = append(r, fullPathName)
		}
	}
	return r
}

func GenerateExecutable(sources [] string) {
	for i := range sources {
		plugin := sources[i]
		name := plugin[strings.LastIndex(plugin, "/")+1:]
		pOut := "/Users/nsr/Dropbox/go/gmidarch-v14/src/adaptive/adaptiveV1/plugins/executable" + "/" + name[:strings.LastIndex(name, ".")]
		pIn := plugin

		_, err := exec.Command(shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), "Something wrong in generating plugin '"+pIn+"' in "+pOut+" "+err.Error())
		}
	}
}
