package main

import (
	"fmt"
	"gmidarch/execution/frontend"
	"injector/evolutive"
)

func main() {

	/*
	// profiling
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
*/
	fe := frontend.FrontEnd{}
	fe.Deploy("senderreceiver.madl")

	// Start evolutive injector
	inj := evolutive.EvolutiveInjector{}
	inj.Start("receiver")

	fmt.Scanln()
}
