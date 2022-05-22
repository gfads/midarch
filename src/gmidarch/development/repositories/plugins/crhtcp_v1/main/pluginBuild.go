package main

import (
	"fmt"
	"gmidarch/development/repositories/plugins/crhtcp_v1"
)

func GetType() interface{} {
	fmt.Println("Chamou GetType do pluginBuild.model")
	return crhtcp.CRHTCP{}
}