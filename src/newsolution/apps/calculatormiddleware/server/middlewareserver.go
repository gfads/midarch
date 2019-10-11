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
	frontend.FrontEnd{}.Deploy("middlewareserver.madl")

	fmt.Scanln()
}

func oldmain() {
	var chn [17] chan messages.SAMessage

	for i := 0; i < 17; i++ {
		chn[i] = make(chan messages.SAMessage)
	}

	srh := components.NewSRH()
	srh.Configure(&chn[0], &chn[1])
	t1 := connectors.NewRequestReply()
	t1.Configure(&chn[0], &chn[1], &chn[2], &chn[3])
	invoker := components.NewInvoker()
	invoker.Configure(&chn[2],&chn[3],&chn[4],&chn[5],&chn[6],&chn[7],&chn[13],&chn[14])
	t2 := connectors.NewRequestReply()
	t2.Configure(&chn[4], &chn[5], &chn[8], &chn[9])
	marshaller1 := components.NewMarshaller()
	marshaller1.Configure(&chn[8],&chn[9])
	marshaller2 := components.NewMarshaller()
	marshaller2.Configure(&chn[15],&chn[16])
	t3 := connectors.NewRequestReply()
	t3.Configure(&chn[6], &chn[7], &chn[10], &chn[11])
	t4 := connectors.NewRequestReply()
	t4.Configure(&chn[13], &chn[14], &chn[15], &chn[16])
	calculatorServer := components.Newcalculatorserver()
	calculatorServer.Configure(&chn[10],&chn[11])
	executor := components.NewExecutor()
	executor.Configure(&chn[12])

	unit1 := components.NewUnit()
	unit1.ConfigureUnit(calculatorServer, &chn[12])
	unit2 := components.NewUnit()
	unit2.ConfigureUnit(marshaller1, &chn[12])
	unit3 := components.NewUnit()
	unit3.ConfigureUnit(invoker, &chn[12])
	unit4 := components.NewUnit()
	unit4.ConfigureUnit(srh, &chn[12])
	unit5 := components.NewUnit()
	unit5.ConfigureUnit(t1, &chn[12])
	unit6 := components.NewUnit()
	unit6.ConfigureUnit(t2, &chn[12])
	unit7 := components.NewUnit()
	unit7.ConfigureUnit(t3, &chn[12])
	unit8 := components.NewUnit()
	unit8.ConfigureUnit(marshaller2, &chn[12])
	unit9 := components.NewUnit()
	unit9.ConfigureUnit(t4, &chn[12])

	//var wg sync.WaitGroup
	//wg.Add(7)
	go engine.Engine{}.Execute(unit1, unit1.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit2, unit2.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit3, unit3.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit4, unit4.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit5, unit5.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit6, unit6.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit7, unit7.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit8, unit8.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(unit9, unit9.Graph, parameters.EXECUTE_FOREVER)
	go engine.Engine{}.Execute(executor, executor.Graph, parameters.EXECUTE_FOREVER)
	//wg.Wait()

	fmt.Scanln()
}
