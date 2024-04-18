package main

import (
	"context"
	"log"
	"net"

	sendfile "github.com/gfads/midarch/examples/externalApps/sendFileGRPC/proto"
	sendFileImpl "github.com/gfads/midarch/examples/sendfiledistributed/sendfileImpl"
	"github.com/gfads/midarch/pkg/shared"
	"google.golang.org/grpc"
)

type SendFileServer struct {
	sendfile.UnimplementedSendFileServiceServer
}

func (sfs *SendFileServer) Upload(ctx context.Context, request *sendfile.Request) (response *sendfile.Response, err error) {
	sendFile := sendFileImpl.SendFile{}
	response = &sendfile.Response{Saved: sendFile.Save(request.Base64File)}
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
	sendfile.RegisterSendFileServiceServer(grpcServer, &SendFileServer{})

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+shared.FIBONACCI_PORT)
	if err != nil {
		log.Fatal("Error while resolving IP address: ", err)
	}
	ln, err := net.ListenTCP("tcp", addr)

	if err := grpcServer.Serve(ln); err != nil {
		log.Fatalf("gRPC: Failed to serve: %s", err)
	}
}
