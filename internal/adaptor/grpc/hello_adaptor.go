package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/amoonguses1/grpc-proto-study/protogen/go/hello"
)

func (a *GrpcAdaptor) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	greet := a.helloService.GenerateHello(req.Name)

	return &hello.HelloResponse{
		Greet: greet,
	}, nil
}

func (a *GrpcAdaptor) SayManyHellos(req *hello.HelloRequest, stream hello.HelloService_SayManyHellosServer) error {
	for i := 0; i < 10; i++ {
		greet := a.helloService.GenerateHello(req.Name)

		res := fmt.Sprintf("[%d] %s", i, greet)
		stream.Send(
			&hello.HelloResponse{
				Greet: res,
			},
		)

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
