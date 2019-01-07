package main

import (
	"context"
	"encoding/json"
	pb "github.com/myshippy/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
)

const (
	address         = "localhost:50050"
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
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewConsignmentServiceClient(conn)

	c, err := parseFile(defaultFilename)
	if err != nil {
		log.Fatalf("fail to parseFile: %v", err)
	}

	r, err := client.Create(context.Background(), c)
	if err != nil {
		log.Fatalf("fail to Create: %v", err)
	}
	log.Printf("Created: %t", r.Result)

	r, err = client.GetAll(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("fail to GetAll: %v", err)
	}
	for _, c := range r.Consignments {
		log.Println(c)
	}
}
