package manager

import (
	"adaptive/adaptiveV3/sharedadaptive"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"shared"
	"strings"
	"time"
)

type MyPlugin struct{}

func (MyPlugin) GenerateSource(pluginName string) string {

	pluginCode := "package main \n"
	pluginCode = pluginCode +  "import \"fmt\" \n"
	pluginCode = pluginCode +  "func Behaviour(){\n"
	pluginCode = pluginCode +  "	fmt.Println(\"Behaviour (Plugin '"+pluginName+"')\") \n"
	pluginCode = pluginCode + "} \n"

	return pluginCode
}

func (MyPlugin) SaveSource(d string, pluginName string, code string) {

	dir := d + "/" + pluginName

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

	f.WriteString(code)

	defer f.Close()
}

func (MyPlugin) InitialiseRepository() {

	// Source files
	removeFiles(sharedadaptive.DIR_SOURCE_PLUGINS_LOCAL)
	removeFiles(sharedadaptive.DIR_SOURCE_PLUGINS_REMOTE)
	removeFiles(sharedadaptive.DIR_EXECUTABLE_PLUGINS_LOCAL)
	removeFiles(sharedadaptive.DIR_EXECUTABLE_PLUGINS_REMOTE)
}

func (MyPlugin) LoadExecutable(d string, pluginName string) plugin.Plugin {

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := d + "/" + pluginName + "/" + pluginName
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

func (MyPlugin) LoadSource(d string, pluginName string) string {

	// Open and load plugin
	pluginFile := d + "/" + pluginName + "/" + pluginName + ".go"
	data, err := ioutil.ReadFile(pluginFile)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}
	return string(data)
}

func (p MyPlugin) LoadSources(d string) map[string]string {
	r := map[string]string{}

	folders, err := ioutil.ReadDir(d)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	temp := []os.FileInfo{}

	for folder := range folders {
		temp, err = ioutil.ReadDir(d + "/" + folders[folder].Name())
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		for file := range temp {
			key := temp[file].Name()[:strings.LastIndex(temp[file].Name(),".")]
			r[key] = p.LoadSource(d, key)
		}
	}
	return r
}

func (p MyPlugin) LoadExecutables(d string) map[string]plugin.Plugin {
	r := map[string]plugin.Plugin{}

	folders, err := ioutil.ReadDir(d)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	pluginNames := []os.FileInfo{}

	for folder := range folders {
		pluginNames, err = ioutil.ReadDir(d + "/" + folders[folder].Name())
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}

		for file := range pluginNames {
			key := pluginNames[file].Name()
			r[key] = p.LoadExecutable(d, key)
		}
	}
	return r
}

func (MyPlugin) GenerateExecutable(dirSource, dirExecutable, pluginName string) {

	pIn := dirSource + "/" + pluginName + "/" + pluginName + ".go"
	pOut := dirExecutable + "/" + pluginName + "/" + pluginName

	_, err := exec.Command(shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Something wrong in generating plugin '"+pIn+"' in "+pOut+" "+err.Error())
	}
}

func removeFiles(d string) {

	if _, err := os.Stat(d); os.IsNotExist(err) {
		err := os.Mkdir(d, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	} else {
		err := os.RemoveAll(d)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
		err = os.Mkdir(d, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}
}
