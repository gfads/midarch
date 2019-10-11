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
	frontend.FrontEnd{}.Deploy("middlewareclient.madl")

	fmt.Scanln()
}

func oldmain() {
	var chn [21] chan messages.SAMessage

	for i := 0; i < 21; i++ {
		chn[i] = make(chan messages.SAMessage)
	}

	// application
	clientCalculator := components.NewCalculatorclient()
	clientCalculator.Configure(&chn[0], &chn[1])

	t1 := connectors.NewRequestReply()
	t1.Configure(&chn[0], &chn[1], &chn[2], &chn[3])

	calculatorProxy := components.NewCalculatorProxy()
	calculatorProxy.Configure(&chn[2], &chn[3], &chn[4], &chn[5])

	// middleware
	t2 := connectors.NewRequestReply()
	t2.Configure(&chn[4], &chn[5], &chn[6], &chn[7])

	requestor := components.NewRequestor()
	requestor.Configure(&chn[6], &chn[7], &chn[8], &chn[9], &chn[12], &chn[13], &chn[17], &chn[18])

	t3 := connectors.NewRequestReply()
	t3.Configure(&chn[8], &chn[9], &chn[10], &chn[11])

	marshaller1 := components.NewMarshaller()
	marshaller1.Configure(&chn[10], &chn[11])

	marshaller2 := components.NewMarshaller()
	marshaller2.Configure(&chn[19], &chn[20])
	t4 := connectors.NewRequestReply()
	t4.Configure(&chn[12], &chn[13], &chn[14], &chn[15])
	t5 := connectors.NewRequestReply()
	t5.Configure(&chn[17], &chn[18], &chn[19], &chn[20])
	crh := components.NewCRH()
	crh. Configure(&chn[14], &chn[15])
	executor := components.NewExecutor()
	executor.Configure(&chn[16])

	unit1 := components.NewUnit()
	unit1.ConfigureUnit(clientCalculator, &chn[16])

	unit2 := components.NewUnit()
	unit2.ConfigureUnit(calculatorProxy, &chn[16])

	unit3 := components.NewUnit()
	unit3.ConfigureUnit(requestor, &chn[16])

	unit4 := components.NewUnit()
	unit4.ConfigureUnit(marshaller1, &chn[16])

	unit5 := components.NewUnit()
	unit5.ConfigureUnit(crh, &chn[16])

	unit6 := components.NewUnit()
	unit6.ConfigureUnit(t1, &chn[16])

	unit7 := components.NewUnit()
	unit7.ConfigureUnit(t2, &chn[16])

	unit8 := components.NewUnit()
	unit8.ConfigureUnit(t3, &chn[16])

	unit9 := components.NewUnit()
	unit9.ConfigureUnit(t4, &chn[16])

	unit10 := components.NewUnit()
	unit10.ConfigureUnit(marshaller2, &chn[16])

	unit11 := components.NewUnit()
	unit11.ConfigureUnit(t5, &chn[16])

	//var wg sync.WaitGroup
	//wg.Add(12)
	go engine.Engine{}.Execute(unit1, unit1.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit4, unit4.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit5, unit5.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit6, unit6.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit7, unit7.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit8, unit8.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit9, unit9.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit10, unit10.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit11, unit11.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor, executor.Graph, parameters.EXECUTE_FOREVER)
	//wg.Wait()

	fmt.Scanln()
}
