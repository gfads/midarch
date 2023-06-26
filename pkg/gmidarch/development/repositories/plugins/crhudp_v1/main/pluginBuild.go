package main

import (
	// "fmt"
	"github.com/gfads/midarch/pkg/gmidarch/development/repositories/plugins/crhudp_v1"
)

func GetType() interface{} {
	// fmt.Println("Chamou GetType do pluginBuild.model")
	return &crhudp.CRHUDP{}
}