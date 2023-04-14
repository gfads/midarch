package main

import (
	"context"
	fibonacci "github.com/gfads/midarch/evaluation/experiments/fiboApps/fibo_gRPC/proto"
	"github.com/gfads/midarch/src/apps/businesses/fibonacciImpl"
	"github.com/gfads/midarch/src/shared"
	"google.golang.org/grpc"
	"log"
	"net"
)

type FibonacciServer struct{}

func (f *FibonacciServer) Fibo(ctx context.Context, request *fibonacci.Request) (response *fibonacci.Response, err error) {
	fibo := fibonacciImpl.Fibonacci{}
	response = &fibonacci.Response{Number: int64(fibo.F(int(request.Place)))}
	return response, nil
}

func main() {
	//rand.Seed(time.Now().UnixNano())
	//for i := 0; i < 100; i++ {
	//	var rd = rand.NormFloat64()
	//	//println(rd)
	//	//println(float64(rd))
	//	//println(rd * 1000)
	//	//println(float64(int(rd * 1000)))
	//	//println(float64(int(rd * 1000))/1000)
	//	fmt.Println(rd)
	//	//fmt.Println(math.Round(rd))
	//	fmt.Println(math.Round(rd * 1000))
	//}
	//
	//return

	grpcServer := grpc.NewServer()
	fibonacci.RegisterFibonacciServiceServer(grpcServer, &FibonacciServer{})

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatal("Error while resolving IP address: ", err)
	}
	ln, err := net.ListenTCP("tcp", addr)

	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("gRPC: Failed to serve: %s", err)
	}
}
