package main

import (
	"gmidarch/development/repositories/plugins/sender_v1"
	"fmt"
)

func GetType() interface{} {
	fmt.Println("Chamou GetType do pluginBuild.model")
	return sender.Sender{}
}