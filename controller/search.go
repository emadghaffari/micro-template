package controller

import (
	"micro/proto/pb"
	"time"
)

func (m *Micro) Search(req *pb.SearchRequest, stream pb.Micro_SearchServer) error {
	msgs := []*pb.SearchResponse{
		&pb.SearchResponse{
			Name: "Emad",
		},
		&pb.SearchResponse{
			Name: "Mamad",
		},
		&pb.SearchResponse{
			Name: "Reza",
		},
		&pb.SearchResponse{
			Name: "Mina",
		},
	}

	for _, msg := range msgs {
		time.Sleep(time.Millisecond * 250)
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
