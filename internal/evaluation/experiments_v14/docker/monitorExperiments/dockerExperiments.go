package main

import (
	docker2 "github.com/gfads/midarch/internal/evaluation/experiments_v14/docker/monitorExperiments/docker"
)

func main() {
	sampleSize := 100
	//var fiboPlaces []int = []int{}           //{2, 11, 38}
	var imageSizes []string = []string{"sm"} //, "md", "lg"}
	var transportProtocolFactors []docker2.TransportProtocolFactor = []docker2.TransportProtocolFactor{docker2.Tcp, docker2.Tls}

	//	for _, transportProtocolFactor := range transportProtocolFactors {
	//		for _, fiboPlace := range fiboPlaces {
	//			docker2.RunFibonacciExperiment(transportProtocolFactor, fiboPlace, sampleSize)
	//		}
	//	}

	for _, transportProtocolFactor := range transportProtocolFactors {
		for _, imageSize := range imageSizes {
			docker2.RunSendFileExperiment(transportProtocolFactor, imageSize, sampleSize)
		}
	}

}
