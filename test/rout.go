package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	ssov1 "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := ssov1.RegisterAuthHandlerFromEndpoint(
		ctx,
		mux,
		"localhost:50551",
		opts,
	); err != nil {
		log.Fatalf("err in register endpoint %v", err)
	}

	fmt.Println("Server start in port 50550")
	log.Fatal(http.ListenAndServe(":50550", mux))

}
