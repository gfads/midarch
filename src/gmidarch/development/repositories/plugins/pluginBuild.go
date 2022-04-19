package main

import (
	"gmidarch/development/repositories/plugins/srhtcp_v2"
	"fmt"
)

func GetType() interface{} {
	fmt.Println("Chamou GetType do pluginBuild.model")
	return srhtcp.SRHTCP{}
}