package jwtsso

import (
	"context"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

func RefreshToken(ctx context.Context, refreshToken string) (*sso.RefreshResponse, error) {
	return nil, nil
}
