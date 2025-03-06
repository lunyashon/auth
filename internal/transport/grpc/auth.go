package gRPCHudler

import (
	"context"
	"main/internal/services/grpc/authgo"

	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	protoc.UnimplementedAuthServer
	Auth authgo.Auth
}

func Register(gRPCServer *grpc.Server, auth authgo.Auth) {
	// В случае изменения pbf файла делает защиту на то, чтобы наш код был универсальным
	protoc.RegisterAuthServer(gRPCServer, &ServerAPI{})
}

// Реализуем метод авторизации пользователя {Login}
func (s *ServerAPI) Login(
	ctx context.Context,
	data *protoc.LoginRequest,
) (*protoc.LoginResponse, error) {
	return nil, nil
}

// Реализуем метод регистрации пользователя {Register}
func (s *ServerAPI) Register(
	ctx context.Context,
	data *protoc.RegisterRequest,
) (*protoc.RegisterResponse, error) {

	id, err := authgo.AuthData.RegisterUser(authgo.AuthData{}, ctx, data)
	if err != nil {
		return nil, err
	}
	return &protoc.RegisterResponse{
		UserId: id,
	}, nil

}

// реализуем метод регистрацию токена {Token}
func (s *ServerAPI) Token(
	ctx context.Context,
	data *protoc.TokenRequest,
) (*protoc.TokenResponse, error) {

	res, err := authgo.AuthData.RegisterToken(authgo.AuthData{}, ctx, data)
	if err != nil {
		return nil, err
	}

	return &protoc.TokenResponse{
		Result: res,
	}, nil
}
