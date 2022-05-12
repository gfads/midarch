package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"plugin"
	"reflect"
	"shared"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	pluginName := "test_v1"
	plugin := LoadPlugin(pluginName, 1)

	fmt.Println("main::will lookup GetType:", pluginName)
	getType, _ := plugin.Lookup("GetType")
	fmt.Println("main::loaded GetType")
	elemType := getType.(func() interface{})()
	fmt.Println("main::elem Type", elemType)
	t := reflect.TypeOf(elemType)
	fmt.Println("t:", t)
	o := reflect.New(t)
	fmt.Println("o:", o)
	e := o.Elem()
	fmt.Println("e:", e)
	f := e.Interface()
	fmt.Println("f:", f)
	reflect.ValueOf(f).MethodByName("Print").Call([]reflect.Value{}) //[]) //[])

	//fmt.Println("Will set value for Name:", reflect.ValueOf(e).FieldByName("Name"))
	//reflect.ValueOf(e).FieldByName("Name").SetString("Test Name")
	//fmt.Println("Will set value for Info:", f)
	//reflect.ValueOf(f).FieldByName("Info").SetString("Test Info")
	//
	//reflect.ValueOf(f).MethodByName("Print").Call([]reflect.Value{})   //[]) //[])

	//s, ok := f.(Test)
	//fmt.Println("s:", s)

	//if !ok {
	//	println("bad type")
	//}
	//s.name = "Nome"
	//s.info = "Info"
	//s.Print()

	//log.Println("Will load the second one")
	//log.Println("")
	//log.Println("")
	//log.Println("")
	//
	//input, err := ioutil.ReadFile("pluginsSrc/pluginBuild.model")
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	////pluginSourcePath := shared.DIR_PLUGINS_SOURCE + "/" + pluginName
	//pluginSourcePath := "adaptive/pluginTest/pluginsSrc" + "/" + pluginName
	//output := bytes.Replace(input, []byte("<pluginName>"), []byte(pluginSourcePath), -1)
	//output = bytes.Replace(output, []byte("<pluginType>"), []byte("test.Test{}"), -1)
	//output = bytes.Replace(output, []byte("Chamou GetType"), []byte("Segundo pluginBuild"), -1)
	//
	//if err = ioutil.WriteFile("pluginsSrc/pluginBuild.go", output, 0666); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//
	//log.Println("")
	//log.Println("")
	//log.Println("")

	plugin = LoadPlugin(pluginName, 2)

	fmt.Println("main::will lookup GetType:", pluginName)
	getType, _ = plugin.Lookup("GetType")
	fmt.Println("main::loaded GetType")
	elemType = getType.(func() interface{})()
	fmt.Println("main::elem Type", elemType)
	t = reflect.TypeOf(elemType)
	fmt.Println("t:", t)
	o = reflect.New(t)
	fmt.Println("o:", o)
	e = o.Elem()
	fmt.Println("e:", e)
	f = e.Interface()
	fmt.Println("f:", f)
	reflect.ValueOf(f).MethodByName("Print").Call([]reflect.Value{}) //[]) //[])
}

func LoadPlugin(pluginName string, retry int) plugin.Plugin {
	pluginSource := "pluginBuild"
	GeneratePlugin(pluginSource, pluginName, retry)

	var plg *plugin.Plugin
	var err error

	// Open and load plugin
	pluginFile := "./plugins/" + pluginName + ".so"
	attempts := 0
	for {
		fmt.Println("pluginUtils.LoadPlugin::will open plugin:", pluginFile)
		plg, err = plugin.Open(pluginFile)
		fmt.Println("pluginUtils.LoadPlugin::opened plugin:", pluginFile)
		if err != nil {
			fmt.Println("pluginUtils.LoadPlugin::error while opening plugin:", err)
			if attempts >= 5 { // TODO
				fmt.Printf("Shared:: Error on trying open plugin '%v' \n", pluginFile)
				os.Exit(0)
			} else {
				attempts++
				time.Sleep(40 * time.Second) // TODO
			}
		} else {
			break
		}
	}

	return *plg
}

func GeneratePlugin(source, pluginName string, retry int) {
	input, err := ioutil.ReadFile("pluginsSrc/pluginBuild.model")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//pluginSourcePath := shared.DIR_PLUGINS_SOURCE + "/" + pluginName
	pluginSourcePath := "adaptive/pluginTest/pluginsSrc" + "/" + pluginName
	output := bytes.Replace(input, []byte("<pluginName>"), []byte(pluginSourcePath), -1)
	output = bytes.Replace(output, []byte("<pluginType>"), []byte("test.Test{}"), -1)

	if retry == 2 {
		output = bytes.Replace(output, []byte("Chamou GetType"), []byte("Segundo pluginBuild"), -1)
	}

	if err = ioutil.WriteFile("pluginsSrc/pluginBuild.go", output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//evolutive.GeneratePlugin(source, plugin)
	generatePlugin(source, pluginName)
}

func generatePlugin(source, plugin string) {
	//pOut := shared.DIR_PLUGINS + "/" + plugin
	pOut := "./plugins" + "/" + plugin + ".so"
	//pIn := shared.DIR_PLUGINS_SOURCE + "/" + source + "/" + source + ".go"
	pIn := "./pluginBuild.go"

	fmt.Println("injector::evolutive.generatePlugin::will build plugin:", source)
	//fmt.Println("command:", shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn)
	fmt.Println("command:", "go", "build", "-buildmode=plugin", "-o", pOut, pIn)
	//_, err := exec.Command(shared.DIR_GO+"/go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	_, err := exec.Command("go", "build", "-buildmode=plugin", "-o", pOut, pIn).CombinedOutput()
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Something wrong in generating plugin '"+pIn+"' in \n '"+pOut+"': "+err.Error()+"\n")
	}
	fmt.Println("injector::evolutive.generatePlugin::plugin built:", source)
}
