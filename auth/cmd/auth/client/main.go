package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "github.com/mig-elgt/chatws/auth/proto/auth"
)

func main() {
	address := flag.String("h", ":8080", "address for the service")
	username := flag.String("user", "foo", "username")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not get connection to address %s: %v", *address, err)
	}

	client := pb.NewAuthServiceClient(conn)

	resp, err := client.Login(context.Background(), &pb.LoginRequest{
		Username: *username,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
