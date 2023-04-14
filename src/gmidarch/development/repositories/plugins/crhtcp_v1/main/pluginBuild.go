package main

import (
	// "fmt"
	"github.com/gfads/midarch/src/gmidarch/development/repositories/plugins/crhtcp_v1"
)

func GetType() interface{} {
	// fmt.Println("Chamou GetType do pluginBuild.model")
	return &crhtcp.CRHTCP{}
}
