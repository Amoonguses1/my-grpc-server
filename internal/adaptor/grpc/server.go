package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/hello"
	"github.com/amoonguses1/my-grpc-server/internal/port"
	"google.golang.org/grpc"
)

type GrpcAdaptor struct {
	helloService port.HelloServicePort
	grpcPort     int
	server       *grpc.Server
	hello.HelloServiceServer
}

func NewGrpcAdaptor(helloService port.HelloServicePort, grpcPort int) *GrpcAdaptor {
	return &GrpcAdaptor{
		helloService: helloService,
		grpcPort:     grpcPort,
	}
}

func (a *GrpcAdaptor) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v\n", a.grpcPort, err)
	}
	log.Printf("Server listening on port %d\n", a.grpcPort)

	grpcServer := grpc.NewServer()
	a.server = grpcServer

	hello.RegisterHelloServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve on port %d: %v\n", a.grpcPort, err)
	}
}

func (a *GrpcAdaptor) Stop() {
	a.server.Stop()
}
