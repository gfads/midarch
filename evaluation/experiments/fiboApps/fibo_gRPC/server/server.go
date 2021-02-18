package main

import (
	fibonacci "apps/fiboApps/fibo_gRPC/proto"
	"apps/fibomiddleware/impl"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"shared"
)

type FibonacciServer struct{}

func (f *FibonacciServer) Fibo(ctx context.Context, request *fibonacci.Request) (response *fibonacci.Response, err error) {
	fibo := impl.Fibonacci{}
	response = &fibonacci.Response{Number: int64(fibo.F(int(request.Place)))}
	return response, nil
}

func main() {
	grpcServer := grpc.NewServer()
	fibonacci.RegisterFibonacciServiceServer(grpcServer, &FibonacciServer{})

	addr, err := net.ResolveTCPAddr("tcp", "localhost:" + shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatal("Error while resolving IP address: ", err)
	}
	ln, err := net.ListenTCP("tcp", addr)

	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("gRPC: Failed to serve: %s", err)
	}
}