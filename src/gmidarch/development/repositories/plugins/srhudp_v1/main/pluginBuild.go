package main

import (
	"fmt"
	"gmidarch/development/repositories/plugins/srhudp_v1"
)

func GetType() interface{} {
	fmt.Println("Chamou GetType do pluginBuild.model")
	return &srhudp.SRHUDP{}
}