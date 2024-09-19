package main

import (
	"log"

	mygrpc "github.com/amoonguses1/my-grpc-server/internal/adaptor/grpc"
	app "github.com/amoonguses1/my-grpc-server/internal/application"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	hs := &app.HelloService{}
	grpcAdaptor := mygrpc.NewGrpcAdaptor(hs, 9090)
	grpcAdaptor.Run()
}
