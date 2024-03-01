package main

import (
	docker2 "github.com/gfads/midarch/internal/evaluation/experiments_v14/docker/monitorExperiments/docker"
)

func main() {
	sampleSize := 10000
	//var fiboPlaces []int = []int{2, 11, 38}
	var imageSizes []string = []string{"md", "lg"} //"sm"
	var transportProtocolFactors []docker2.TransportProtocolFactor = []docker2.TransportProtocolFactor{
		docker2.QuicHttp2} //docker2.UdpTcp, docker2.TcpTls}
	var adaptationIntervals []int = []int{30, 120, 300}
	//docker2.Udp, docker2.Tcp, docker2.Http, docker2.Https, docker2.Http2, docker2.Rpc, docker2.Quic}

	//for _, fiboPlace := range fiboPlaces {
	//		for _, transportProtocolFactor := range transportProtocolFactors {
	//			docker2.RunFibonacciExperiment(transportProtocolFactor, fiboPlace, sampleSize)
	//		}
	//	}

	for _, imageSize := range imageSizes {
		for _, transportProtocolFactor := range transportProtocolFactors {
			if transportProtocolFactor.IsEvolutive() {
				for _, adaptationInterval := range adaptationIntervals {
					docker2.RunSendFileExperiment(transportProtocolFactor, adaptationInterval, imageSize, sampleSize)
				}
			} else {
				docker2.RunSendFileExperiment(transportProtocolFactor, -1, imageSize, sampleSize)
			}
		}
	}
}
