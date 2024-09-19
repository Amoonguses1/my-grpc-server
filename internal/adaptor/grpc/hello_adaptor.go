package grpc

import (
	"context"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/hello"
)

func (a *GrpcAdaptor) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	greet := a.helloService.GenerateHello(req.Name)

	return &hello.HelloResponse{
		Greet: greet,
	}, nil
}
