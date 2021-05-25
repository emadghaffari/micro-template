package controller

import (
	"context"
	"micro/proto/pb"
)

func (m Micro) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "Hello Master!",
	}, nil
}
