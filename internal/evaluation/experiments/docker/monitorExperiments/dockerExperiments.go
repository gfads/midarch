package main

import (
	docker2 "github.com/gfads/midarch/internal/evaluation/experiments/docker/monitorExperiments/docker"
)

func main() {
	fiboPlace, sampleSize := 38, 10000

	//docker.RunExperiment(docker.Udp, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.Tcp, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.Tls, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.Quic, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.Rpc, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.Http, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.Https, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.Http2, fiboPlace, sampleSize)
	docker2.RunExperiment(docker2.UdpTcp, fiboPlace, sampleSize)
	docker2.RunExperiment(docker2.Tls, fiboPlace, sampleSize)
	//
	//docker.RunExperiment(docker.E_Rpc, fiboPlace, sampleSize)
	docker2.RunExperiment(docker2.E_Grpc, fiboPlace, sampleSize)
	//docker.RunExperiment(docker.E_Rmq, fiboPlace, sampleSize)

}
