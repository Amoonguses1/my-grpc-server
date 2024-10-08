package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/bank"
	"github.com/amoonguses1/grpc-proto-study/protogen/go/hello"
	resl "github.com/amoonguses1/grpc-proto-study/protogen/go/resiliency"
	"github.com/amoonguses1/my-grpc-server/internal/port"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcAdaptor struct {
	helloService      port.HelloServicePort
	bankService       port.BankServicePort
	resiliencyService port.ResiliencyServicePort
	grpcPort          int
	server            *grpc.Server
	hello.HelloServiceServer
	bank.BankServiceServer
	resl.ResiliencyServiceServer
	resl.ResiliencyWithMetadataServiceServer
}

func NewGrpcAdaptor(helloService port.HelloServicePort, bankService port.BankServicePort, resiliencyService port.ResiliencyServicePort, grpcPort int) *GrpcAdaptor {
	return &GrpcAdaptor{
		helloService:      helloService,
		bankService:       bankService,
		resiliencyService: resiliencyService,
		grpcPort:          grpcPort,
	}
}

func (a *GrpcAdaptor) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v\n", a.grpcPort, err)
	}
	log.Printf("Server listening on port %d\n", a.grpcPort)

	creds, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")

	if err != nil {
		log.Fatalln("Cannot create server credentials :", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
	// grpc.ChainStreamInterceptor(
	// 	interceptor.LogStreamServerInterceptor(),
	// 	interceptor.BasicStreamServerInterceptor(),
	// ),
	// grpc.ChainUnaryInterceptor(
	// 	interceptor.LogUnaryServerInterceptor(),
	// 	interceptor.BasicUnaryServerInterceptor(),
	// ),
	)
	a.server = grpcServer

	hello.RegisterHelloServiceServer(grpcServer, a)
	bank.RegisterBankServiceServer(grpcServer, a)
	resl.RegisterResiliencyServiceServer(grpcServer, a)
	resl.RegisterResiliencyWithMetadataServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve on port %d: %v\n", a.grpcPort, err)
	}
}

func (a *GrpcAdaptor) Stop() {
	a.server.Stop()
}
