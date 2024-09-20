package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
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

func (a *GrpcAdaptor) SayHelloToEveryone(stream hello.HelloService_SayHelloToEveryoneServer) error {
	res := ""

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(
				&hello.HelloResponse{
					Greet: res,
				},
			)
		}
		if err != nil {
			log.Fatalln("Error while reading from client :", err)
		}

		greet := a.helloService.GenerateHello(req.Name)

		res += greet + " "
	}
}

func (a *GrpcAdaptor) SayHelloContinuous(stream hello.HelloService_SayHelloContinuousServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalln("Error while reading from client :", err)
		}

		greet := a.helloService.GenerateHello(req.Name)

		err = stream.Send(
			&hello.HelloResponse{
				Greet: greet,
			},
		)
		if err != nil {
			log.Fatalln("Error while sending response to client :", err)
		}
	}
}
