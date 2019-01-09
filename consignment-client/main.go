package main

import (
	"context"
	"encoding/json"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	pb "github.com/myshippy/consignment-service/proto/consignment"
	//"google.golang.org/grpc"
	"io/ioutil"
	"log"
)

const (
	defaultFilename = "consignment.json"
)

func parseFile(f string) (*pb.Consignment, error) {
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	var c *pb.Consignment
	err = json.Unmarshal(bytes, &c)
	return c, err
}

func main() {
	cmd.Init()
	s := pb.NewConsignmentService("cz.go.microservices.consignment", microclient.DefaultClient)

	c, err := parseFile(defaultFilename)
	if err != nil {
		log.Fatalf("fail to parseFile: %v", err)
	}

	r, err := s.Create(context.TODO(), c)
	if err != nil {
		log.Fatalf("fail to Create: %v", err)
	}
	log.Printf("Created: %t", r.Result)

	r, err = s.GetAll(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("fail to GetAll: %v", err)
	}
	for _, c := range r.Consignments {
		log.Println(c)
	}
}
