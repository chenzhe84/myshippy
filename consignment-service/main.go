package main

import (
	"context"
	micro "github.com/micro/go-micro"
	pb "github.com/myshippy/consignment-service/proto/consignment"
	vesselPb "github.com/myshippy/vessel-service/proto/vessel"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	"log"
	//"net"
	"fmt"
)

type Repository struct {
	consignments []*pb.Consignment
}

func (r *Repository) Create(c *pb.Consignment) {
	r.consignments = append(r.consignments, c)
}

func (r *Repository) GetAll() []*pb.Consignment {
	return r.consignments
}

type ConsignmentService struct {
	repo         Repository
	vesselClient vesselPb.VesselService
}

func (s *ConsignmentService) Create(cxt context.Context, c *pb.Consignment, r *pb.Response) error {
	rep, err := s.vesselClient.GetAvailable(context.Background(), &vesselPb.Request{
		Capacity:  int32(len(c.Containers)),
		MaxWeight: c.Weight,
	})
	if err != nil {
		return err
	}

	log.Printf("Found vessel: %s \n", rep.Vessel.Name)
	c.VesselId = rep.Vessel.Id
	s.repo.Create(c)

	r.Result = true
	r.Consignment = c
	return nil
}

func (s *ConsignmentService) GetAll(cxt context.Context, e *pb.EmptyRequest, r *pb.Response) error {
	r.Result = true
	r.Consignments = s.repo.GetAll()
	return nil
}

func main() {
	s := micro.NewService(
		micro.Name("cz.go.microservices.consignment"),
		micro.Version("latest"),
	)

	s.Init()

	handler := &ConsignmentService{
		repo:         Repository{},
		vesselClient: vesselPb.NewVesselService("cz.go.microservices.vessel", s.Client()),
	}
	err := pb.RegisterConsignmentServiceHandler(s.Server(), handler)
	if err != nil {
		log.Fatalf("fail to RegisterConsignmentServiceHandler: %v", err)
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}
