package manager

import (
	"adaptive/adaptiveV6/sharedadaptive"
	"io/ioutil"
	"os"
	"os/exec"
	"plugin"
	"shared"
	"strings"
	"time"
)

type MyPlugin struct {
	Type string
	Name string
	Code string
	Tme  time.Time
}

func (MyPlugin) GenerateSource(pluginName string) MyPlugin {
	r := MyPlugin{}

	pluginCode := "package main \n"
	pluginCode = pluginCode + "import \"fmt\" \n"
	pluginCode = pluginCode + "func Behaviour(){\n"
	pluginCode = pluginCode + "	fmt.Println(\"Behaviour (Plugin '" + pluginName + "')\") \n"
	pluginCode = pluginCode + "} \n"

	r.Name = pluginName
	r.Type = "TODO"
	r.Code = pluginCode

	return r
}

func (MyPlugin) SaveSource(d string, p MyPlugin) {

	dir := d + "/" + p.Name

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			shared.ErrorHandler(shared.GetFunction(), err.Error())
		}
	}

	file := dir + "/" + p.Name + ".go"

	f, err := os.Create(file)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

	_, err = f.WriteString(p.Code)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), err.Error())
	}

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

func (p MyPlugin) LoadSources(d string) []MyPlugin {
	var r [] MyPlugin

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
			n := temp[file].Name()[:strings.LastIndex(temp[file].Name(), ".")]
			c := p.LoadSource(d, n)
			t := "TODO"
			tme := temp[file].ModTime()
			p := MyPlugin{Name: n, Code: c, Type: t, Tme: tme}

			r = append(r, p)
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

func (MyPlugin) GenerateExecutable(dirSource string, dirExecutable string, p MyPlugin) {

	pIn := dirSource + "/" + p.Name + "/" + p.Name + ".go"
	pOut := dirExecutable + "/" + p.Name + "/" + p.Name

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
