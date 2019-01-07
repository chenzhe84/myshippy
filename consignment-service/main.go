package main

import (
	"context"
	pb "github.com/myshippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const port = ":50050"

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
	repo *Repository
}

func (s *ConsignmentService) Create(cxt context.Context, c *pb.Consignment) (*pb.Response, error) {
	s.repo.Create(c)
	return &pb.Response{
		Result:      true,
		Consignment: c,
	}, nil
}

func (s *ConsignmentService) GetAll(cxt context.Context, r *pb.EmptyRequest) (*pb.Response, error) {
	cs := s.repo.GetAll()
	return &pb.Response{
		Result:       true,
		Consignments: cs,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("fail to listen %s: %v", port, err)
	}

	s := grpc.NewServer()
	pb.RegisterConsignmentServiceServer(s, &ConsignmentService{repo: &Repository{}})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}
}
