package gserver

import (
	pb "github.com/failfilin/otus_hw/internal/grpc/proto"
	"github.com/failfilin/otus_hw/internal/models"
)

func restaurantModelToProto(r *models.Restaurant) *pb.Restaurant {
	return &pb.Restaurant{
		Id:     r.Id.String(),
		Name:   r.Name,
		Logo:   r.Logo,
		Active: r.Active,
	}
}
