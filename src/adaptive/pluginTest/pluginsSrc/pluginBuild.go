package main

import (
	"adaptive/pluginTest/pluginsSrc/test_v1"
	"fmt"
)

func GetType() interface{} {
	fmt.Println("Segundo pluginBuild")
	return test.Test{}
}