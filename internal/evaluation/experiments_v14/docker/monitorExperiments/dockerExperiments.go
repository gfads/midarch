package main

import (
	docker2 "github.com/gfads/midarch/internal/evaluation/experiments_v14/docker/monitorExperiments/docker"
)

func main() {
	sampleSize := 10000
	var fiboPlaces []int = []int{2, 11, 38}              //
	var imageSizes []string = []string{"sm", "md", "lg"} //
	var transportProtocolFactors []docker2.TransportProtocolFactor = []docker2.TransportProtocolFactor{
		docker2.Udp, docker2.Tcp, docker2.Tls, docker2.Http, docker2.Https, docker2.Http2, docker2.Rpc, docker2.Quic, docker2.TcpTls, docker2.RpcHttp, docker2.TcpHttp, docker2.TlsHttp2, docker2.RpcQuic, docker2.QuicHttp2}
	var adaptationIntervals []int = []int{120, 300}

	for _, fiboPlace := range fiboPlaces {
		for _, transportProtocolFactor := range transportProtocolFactors {
			if transportProtocolFactor.IsEvolutive() {
				for _, adaptationInterval := range adaptationIntervals {
					docker2.RunFibonacciExperiment(transportProtocolFactor, adaptationInterval, fiboPlace, sampleSize)
				}
			} else {
				docker2.RunFibonacciExperiment(transportProtocolFactor, -1, fiboPlace, sampleSize)
			}
		}
	}

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
