package gserver

import (
	"log"
	"net"

	pb "github.com/failfilin/otus_hw/internal/grpc/proto"
	"google.golang.org/grpc"
)

func Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRestaurantServiceServer(
		grpcServer,
		NewRestaurantServer(),
	)

	log.Println("gRPC server started on", addr)
	return grpcServer.Serve(lis)
}
