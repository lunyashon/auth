package main

import (
	"context"
	"fmt"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.NewClient("localhost:50051", opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()

	client := protoc.NewAuthClient(conn)
	request := &protoc.RegisterRequest{
		Email:    "brand",
		Login:    "brand",
		Password: "brand",
		Api:      "ewew",
		Services: "dsds",
	}
	response, err := client.Register(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	fmt.Println(response.UserId)
}
