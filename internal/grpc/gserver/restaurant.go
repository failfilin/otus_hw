package gserver

import (
	"context"

	pb "github.com/failfilin/otus_hw/internal/grpc/proto"
	"github.com/failfilin/otus_hw/internal/models"
	"github.com/failfilin/otus_hw/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type RestaurantServer struct {
	pb.UnimplementedRestaurantServiceServer
}

func NewRestaurantServer() *RestaurantServer {
	return &RestaurantServer{}
}

func (s *RestaurantServer) GetRestaurants(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.RestaurantList, error) {

	repository.RestSlice.Mu.Lock()
	defer repository.RestSlice.Mu.Unlock()

	resp := &pb.RestaurantList{
		Restaurants: make([]*pb.Restaurant, 0, len(repository.RestSlice.Slice)),
	}

	for _, r := range repository.RestSlice.Slice {
		resp.Restaurants = append(resp.Restaurants, restaurantModelToProto(&r))
	}

	return resp, nil
}

func (s *RestaurantServer) GetRestaurantById(
	ctx context.Context,
	req *pb.RestaurantIdRequest,
) (*pb.Restaurant, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	rest, ok := repository.RestSlice.GetByID(id)
	if !ok {
		return nil, status.Error(codes.NotFound, "restaurant not found")
	}

	return restaurantModelToProto(rest), nil
}

func (s *RestaurantServer) CreateRestaurant(
	ctx context.Context,
	req *pb.CreateRestaurantRequest,
) (*pb.IdReply, error) {

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	rest := models.Restaurant{
		Id:     uuid.New(),
		Name:   req.Name,
		Logo:   req.Logo,
		Active: false,
	}

	repository.RestSlice.Append(rest)

	if err := repository.SaveToFile(&repository.RestSlice); err != nil {
		return nil, status.Error(codes.Internal, "failed to save restaurant")
	}

	return &pb.IdReply{
		Id: rest.Id.String(),
	}, nil
}

func (s *RestaurantServer) UpdateRestaurant(
	ctx context.Context,
	req *pb.UpdateRestaurantRequest,
) (*pb.MessageReply, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	rest, ok := repository.RestSlice.GetByID(id)
	if !ok {
		return nil, status.Error(codes.NotFound, "restaurant not found")
	}

	if req.Restaurant == nil {
		return nil, status.Error(codes.InvalidArgument, "restaurant payload is required")
	}

	r := req.Restaurant
	rest.Update(
		&r.Name,
		&r.Logo,
		nil,
		&r.Active,
	)

	if err := repository.SaveToFile(&repository.RestSlice); err != nil {
		return nil, status.Error(codes.Internal, "failed to save restaurant")
	}

	return &pb.MessageReply{
		Message: "restaurant updated successfully",
	}, nil
}

func (s *RestaurantServer) DeleteRestaurant(
	ctx context.Context,
	req *pb.RestaurantIdRequest,
) (*pb.MessageReply, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	if !repository.RestSlice.RemoveByID(id) {
		return nil, status.Error(codes.NotFound, "restaurant not found")
	}

	if err := repository.SaveToFile(&repository.RestSlice); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete restaurant")
	}

	return &pb.MessageReply{
		Message: "restaurant deleted successfully",
	}, nil
}
