package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Locke23/go-grpc-test/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUsers(client)
	// AddUserVerbose(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Lucas",
		Email: "l@l.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)

}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Lucas",
		Email: "l@l.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status: ", stream.Status)
	}

}

func AddUsers(client pb.UserServiceClient) {

	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "lucas",
			Email: "lucas@lu.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "lucas 2",
			Email: "lucas2@lu.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "lucas 3",
			Email: "lucas3@lu.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "lucas 4",
			Email: "lucas4@lu.com",
		},
		&pb.User{
			Id:    "l5",
			Name:  "lucas 5",
			Email: "lucas5@lu.com",
		},
		&pb.User{
			Id:    "l6",
			Name:  "lucas 6",
			Email: "lucas6@lu.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving response: %v", err)
	}

	fmt.Println(res)

}

func AddUserStreamBoth(client pb.UserServiceClient) {

	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "lucas",
			Email: "lucas@lu.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "lucas 2",
			Email: "lucas2@lu.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "lucas 3",
			Email: "lucas3@lu.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "lucas 4",
			Email: "lucas4@lu.com",
		},
		&pb.User{
			Id:    "l5",
			Name:  "lucas 5",
			Email: "lucas5@lu.com",
		},
		&pb.User{
			Id:    "l6",
			Name:  "lucas 6",
			Email: "lucas6@lu.com",
		},
	}

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data from server: %v", err)
				break
			}
			fmt.Printf("Receiving user: %v, com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}
