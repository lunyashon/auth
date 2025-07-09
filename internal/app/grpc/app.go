package appgrpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/lunyashon/auth/internal/services/authgo"
	gRPCHudler "github.com/lunyashon/auth/internal/transport/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	// hundler "main/internal/transport/http"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, auth *authgo.AuthData) *App {

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) error {
			log.Error("Recovered from panic", slog.Any("panic", p))
			return status.Error(codes.Internal, "Internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainStreamInterceptor(
		recovery.StreamServerInterceptor(recoveryOpts...),
	))

	gRPCHudler.Register(gRPCServer, auth)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error { //(a *App)

	reflection.Register(a.gRPCServer)
	const op = "appgrpc.Run"
	list, err := net.Listen("tcp", ":50551")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", list.Addr().String()))

	if err := a.gRPCServer.Serve(list); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
