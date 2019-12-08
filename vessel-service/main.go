package main

import (
	"context"
	"errors"
	"log"
	pb "shippy/vessel-service/proto/vessel"

	"github.com/micro/go-micro"
)

type Repository interface {
	FindAvailable(spec *pb.Specification) (*pb.Vessel, error)
}


type VesselRepository struct {
	vessels []*pb.Vessel
}

func ( vrepo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	//选择一条合适的货轮
	for _, item := range vrepo.vessels {
		if item.Capacity >= spec.Capacity && item.MaxWeight >= spec.MaxWeight {
			return item, nil
		}
	}
	return nil, errors.New("can not find available vessel")
}

type service struct {
	repo Repository
}

func (s *service) FindAvailable(ctx context.Context, sp *pb.Specification, resp *pb.Response) error {
	v, err := s.repo.FindAvailable(sp)
	if err != nil {
		return err
	}
	resp.Vessel = v
	return nil
}

func main()  {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	vesselRepo := &VesselRepository{vessels:vessels}
	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
		)
	server.Init()

	//注册服务
	pb.RegisterVesselServiceHandler(server.Server(), &service{vesselRepo})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

