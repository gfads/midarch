package main

import (
	"fmt"
	"gmidarch/development/framework/messages"
	"newsolution/development/components"
	"newsolution/development/connectors"
	"newsolution/execution/environment/engine"
	"newsolution/shared/parameters"
	"newsolution/execution/environment/frontend"
)

func main() {
	frontend.FrontEnd{}.Deploy("clientserver.madl")

	fmt.Scanln()
}

func oldmain() {
	chn1 := make(chan messages.SAMessage)
	chn2 := make(chan messages.SAMessage)
	chn3 := make(chan messages.SAMessage)
	chn4 := make(chan messages.SAMessage)
	chn5 := make(chan messages.SAMessage)

	client := components.NewClient()
	client.Configure(&chn1, &chn2)
	t := connectors.NewRequestReply()
	t.Configure(&chn1, &chn2, &chn3, &chn4)
	server := components.NewServer()
	server.Configure(&chn3, &chn4)
	executor := components.NewExecutor()
	executor.Configure(&chn5)

	unit1 := components.NewUnit()
	unit1.ConfigureUnit(client,&chn5)
	unit2 := components.NewUnit()
	unit2.ConfigureUnit(server,&chn5)
	unit3 := components.NewUnit()
	unit3.ConfigureUnit(t,&chn5)

	go engine.Engine{}.Execute(unit1, unit1.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor, executor.Graph, parameters.EXECUTE_FOREVER)

	fmt.Scanln()
}
