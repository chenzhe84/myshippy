package main

import (
	pb "../vessel-service/proto/vessel"
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	"log"
)

type Repository struct {
	vessels []*pb.Vessel
}

func (s *Repository) GetAvailable(req *pb.Request) (*pb.Vessel, error) {
	for _, v := range s.vessels {
		if v.Capacity >= req.Capacity && v.MaxWeight >= req.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("can not find available")
}

type VesselService struct {
	repo Repository
}

func (s *VesselService) GetAvailable(cxt context.Context, req *pb.Request, res *pb.Response) error {
	v, err := s.repo.GetAvailable(req)
	if err != nil {
		return err
	}
	res.Vessel = v
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("cz.go.microservices.vessel"),
	)

	service.Init()

	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	handler := &VesselService{repo: Repository{vessels: vessels}}

	if err := pb.RegisterVesselServiceHandler(service.Server(), handler); err != nil {
		log.Fatalf("fail to RegisterVesselServiceHandler: %v", err)
	}

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
