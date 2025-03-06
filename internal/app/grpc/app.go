// package appgrpc
package main

import (
	"fmt"
	"log/slog"
	gRPCHudler "main/internal/transport/grpc"
	"net"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc"
	// hundler "main/internal/transport/http"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// func New(log *slog.Logger, port int) *App {

// 	// recoveryOpts := []recovery.Option{
// 	// 	// recovery.WithRecoveryHandler(func(p interface{}) error {
// 	// 	// 	log.Error("Recovered from panic", slog.Any("panic", p))
// 	// 	// 	return status.Error(codes.Internal, "Internal error")
// 	// 	// }),
// 	// }
// 	opts := []grpc.ServerOption{}
// 	gRPCServer := grpc.NewServer(opts...) // grpc.ChainStreamInterceptor(recovery.StreamServerInterceptor(...),)

// 	return &App{
// 		log:        log,
// 		gRPCServer: gRPCServer,
// 		port:       port,
// 	}
// }

func Run() error { //(a *App)

	opts := []grpc.ServerOption{}
	gRPCServer := grpc.NewServer(opts...) // grpc.ChainStreamInterceptor(recovery.StreamServerInterceptor(...),)
	const op = "appgrpc.Run"
	list, err := net.Listen("tcp", ":50053")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// a.log.Info("grpc server started", slog.String("addr", list.Addr().String()))
	fmt.Println("Start server " + list.Addr().String())

	protoc.RegisterAuthServer(gRPCServer, &gRPCHudler.ServerAPI{})
	if err := gRPCServer.Serve(list); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func main() {
	fmt.Println(Run())
}
