package grpc

import (
	"context"
	pb "micro/proto"
)

func (m Micro) Test(ctx context.Context, req *pb.Micro) (*pb.Micro, error) {
	return &pb.Micro{
		Id: req.GetId(),
	}, nil
}
