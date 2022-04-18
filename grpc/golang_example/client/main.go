// Copyright 2022 Kien Nguyen-Tuan <kiennt2609@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	pb "github.com/ntk148v/testing/grpc/golang_example/customer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
)

const (
	address = "localhost:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {
	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

// getCustomers calls the RPC method GetCustomer of CustomerServer
func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {
	// calling the streaming API
	stream, err := client.GetCustomer(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomer(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}
}

func main() {
	// Set up a connection to the gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Creates a new CustomerClient
	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    101,
		Name:  "Kien Nguyen-Tuan",
		Email: "kiennt2609@gmail",
		Phone: "976872724",
		Address: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "18 HHT",
				City:              "Hanoi",
				State:             "Hanoi",
				Zip:               "6969",
				IsShippingAddress: false,
			},
			&pb.CustomerRequest_Address{
				Street:            "O day",
				City:              "Thanh pho ne",
				State:             "Ne thanh pho",
				Zip:               "9696",
				IsShippingAddress: true,
			},
		},
	}

	// Create a new customer
	createCustomer(client, customer)

	customer = &pb.CustomerRequest{
		Id:    102,
		Name:  "Kien Nguyen",
		Email: "mail@mail.com",
		Phone: "932245949",
		Address: []*pb.CustomerRequest_Address{
			&pb.CustomerRequest_Address{
				Street:            "Met vl",
				City:              "Haloi",
				State:             "Haloi",
				Zip:               "6999",
				IsShippingAddress: true,
			},
		},
	}

	// Create one more new customer
	createCustomer(client, customer)
	// filter with an empty Keyword
	filter := &pb.CustomerFilter{Keyword: ""}
	getCustomers(client, filter)
}
