package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gfads/midarch/pkg/gmidarch/development/repositories/architectural"
	evolutive "github.com/gfads/midarch/pkg/injector"
	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/pluginUtils"
)

func main() {
	server()
}

func client() {
	architectural.LoadArchitecturalRepository(map[string]interface{}{})
	componentName := "CRHTCP"
	cmdType := shared.GetComponentTypeByNameFromRAM(componentName)
	typeof := reflect.ValueOf(cmdType).Elem().Type()
	name := typeof.Name()
	fmt.Println("name:", name)
}

func server() {
	pluginName := "srhtcp_v2"
	evolutive.GeneratePlugin(pluginName, "srhtcp", pluginName)
	plg := pluginUtils.LoadPlugin(pluginName + ".so")
	fmt.Println("Executor.I_Process::plugin loaded:", pluginName)
	log.Println("Executor.I_Process::Will lookup Gettype:", pluginName)
	getType, _ := plg.Lookup("GetType")
	elemType := getType.(func() interface{})()
	cmdType := elemType

	name := reflect.ValueOf(cmdType).Elem().Type().Name()
	fmt.Println("name:", name)
	name = reflect.TypeOf(cmdType).Name()
	fmt.Println("name:", name)
}
