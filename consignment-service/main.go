package main

import (
	micro "github.com/micro/go-micro"
	"context"
	pb "github.com/myshippy/consignment-service/proto/consignment"
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
	repo *Repository
}

func (s *ConsignmentService) Create(cxt context.Context, c *pb.Consignment, r *pb.Response) error{
	s.repo.Create(c)
	r.Result=true
	r.Consignment=c
	return nil
}

func (s *ConsignmentService) GetAll(cxt context.Context, e *pb.EmptyRequest, r *pb.Response) error{
	r.Result=true
	r.Consignments=s.repo.GetAll()
	return nil
}

func main(){
	s := micro.NewService(
		micro.Name("cz.go.microservices.consignment"),
		micro.Version("latest"),
	)

	s.Init()

	err := pb.RegisterConsignmentServiceHandler(s.Server(), &ConsignmentService{repo:&Repository{}})
	if err != nil {
		log.Fatalf("fail to RegisterConsignmentServiceHandler: %v", err)
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}