package manager

import (
	"adaptive/adaptiveV2/sharedadaptive"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"shared"
	"time"
)

type MyPlugin struct{}

func (MyPlugin) GenerateSource(pluginName string) []string {
	pluginCode := []string{}

	pluginCode = append(pluginCode, "package main \n")
	pluginCode = append(pluginCode, "import \"fmt\" \n")
	pluginCode = append(pluginCode, "func Behaviour(){\n")
	pluginCode = append(pluginCode, "	fmt.Println(\"Behaviour (Plugin '"+pluginName+"')\") \n")
	pluginCode = append(pluginCode, "} \n")

	return pluginCode
}

func (MyPlugin) SaveCode(pluginName string, code []string) {

	dir := sharedadaptive.DIR_SOURCE_PLUGINS + "/" + pluginName

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	file := dir + "/" + pluginName + ".go"

	f, err := os.Create(file)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	for i := range code {
		f.WriteString(code[i])
	}

	defer f.Close()
}

func (MyPlugin) InitialiseRepository() {

	// Source files
	if _, err := os.Stat(sharedadaptive.DIR_SOURCE_PLUGINS); os.IsNotExist(err) {
		err := os.Mkdir(sharedadaptive.DIR_SOURCE_PLUGINS, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	} else {
		err := os.RemoveAll(sharedadaptive.DIR_SOURCE_PLUGINS)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		err = os.Mkdir(sharedadaptive.DIR_SOURCE_PLUGINS, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	// Executable files
	if _, err := os.Stat(sharedadaptive.DIR_EXECUTABLE_PLUGINS); os.IsNotExist(err) {
		err := os.Mkdir(sharedadaptive.DIR_EXECUTABLE_PLUGINS, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	} else {
		err := os.RemoveAll(sharedadaptive.DIR_EXECUTABLE_PLUGINS)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		err = os.Mkdir(sharedadaptive.DIR_EXECUTABLE_PLUGINS, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
}

func (MyPlugin) LoadPlugin(pluginName string) plugin.Plugin {

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := sharedadaptive.DIR_EXECUTABLE_PLUGINS + "/" + pluginName + "/" + pluginName
	attempts := 0
	for {
		plg, err = plugin.Open(pluginFile)

		if err != nil {
			if attempts >= 3 { // TODO
				shared.ErrorHandler(shared.GetFunction(), "Error on trying open plugin: "+pluginFile)
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

func (MyPlugin) LoadSources() []string {
	r := []string{}

	folders, err1 := ioutil.ReadDir(sharedadaptive.DIR_SOURCE_PLUGINS)
	if err1 != nil {
		shared.ErrorHandler(shared.GetFunction(), err1.Error())
	}

	temp := []os.FileInfo{}

	for folder := range folders {
		temp, err1 = ioutil.ReadDir(sharedadaptive.DIR_SOURCE_PLUGINS + "/" + folders[folder].Name())
		if err1 != nil {
			shared.ErrorHandler(shared.GetFunction(), err1.Error())
		}

		for file := range temp {
			fullPathName := sharedadaptive.DIR_SOURCE_PLUGINS + "/" + folders[folder].Name() + "/" + temp[file].Name()
			r = append(r, fullPathName)
		}
	}
	return r
}

func (p MyPlugin) LoadPlugins() map[string]plugin.Plugin {
	r := map[string]plugin.Plugin{}

	folders, err1 := ioutil.ReadDir(sharedadaptive.DIR_EXECUTABLE_PLUGINS)
	if err1 != nil {
		shared.ErrorHandler(shared.GetFunction(), err1.Error())
	}

	pluginNames := []os.FileInfo{}

	for folder := range folders {
		pluginNames, err1 = ioutil.ReadDir(sharedadaptive.DIR_EXECUTABLE_PLUGINS + "/" + folders[folder].Name())
		if err1 != nil {
			shared.ErrorHandler(shared.GetFunction(), err1.Error())
		}

		for file := range pluginNames {
			r[pluginNames[file].Name()] = p.LoadPlugin(pluginNames[file].Name())
		}
	}
	return r
}

func (MyPlugin) GenerateExecutable(pluginName string) {

	pOut := sharedadaptive.DIR_EXECUTABLE_PLUGINS + "/" + pluginName + "/" + pluginName
	pIn := sharedadaptive.DIR_SOURCE_PLUGINS +"/"+pluginName+"/"+ pluginName+".go"

	_, err := exec.Command(shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Something wrong in generating plugin '"+pIn+"' in "+pOut+" "+err.Error())
	}
}
