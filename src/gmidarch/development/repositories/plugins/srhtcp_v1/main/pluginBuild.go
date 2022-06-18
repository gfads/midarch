package main

import (
	"fmt"
	"gmidarch/development/repositories/plugins/srhtcp_v1"
)

func GetType() interface{} {
	fmt.Println("Chamou GetType do pluginBuild.model")
	return &srhtcp.SRHTCP{}
}