package authgo

import (
	"context"
	"crypto/rand"
	"math/big"

	validate "github.com/lunyashon/auth/internal/services/validation"
	protoc "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create token in database
// Return resul (true or false with error)
func (a AuthData) RegisterToken(ctx context.Context, data *protoc.TokenRequest) (string, error) {

	token := generateToken()

	if err := validate.RegisterToken(ctx, data, token, a.DB); err != nil {
		return "", err
	}

	if err := a.DB.Token.CreateToken(ctx, data, token); err != nil {
		a.Log.ErrorContext(ctx, err.Error())
		return "", status.Error(codes.Internal, "database error")
	}

	return token, nil
}

func generateToken() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 32)
	for i := range token {
		t, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		token[i] = chars[t.Int64()]
	}
	return string(token)
}
