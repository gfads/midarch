package main

import (
	// "fmt"
	"github.com/gfads/midarch/pkg/gmidarch/development/repositories/plugins/middleware"
)

func GetType() interface{} {
	// fmt.Println("Chamou GetType do pluginBuild.model")
	return &middleware.SRHHTTP2{}
}