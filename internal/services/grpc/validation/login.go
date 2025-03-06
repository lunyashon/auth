package validation

import (
	protoc "github.com/lunyashon/protoc/auth/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (v Validate) Login(data *protoc.LoginRequest) error {
	if data.Login == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if data.Password == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	return nil
}
