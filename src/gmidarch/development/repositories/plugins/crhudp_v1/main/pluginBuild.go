package main

import (
	// "fmt"
	"gmidarch/development/repositories/plugins/crhudp_v1"
)

func GetType() interface{} {
	// fmt.Println("Chamou GetType do pluginBuild.model")
	return &crhudp.CRHUDP{}
}