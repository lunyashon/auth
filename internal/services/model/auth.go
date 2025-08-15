package model

import jwtsso "github.com/lunyashon/auth/internal/lib/jwt"

type UserAuth struct {
	Services []string
	UID      int64
	Tokens   jwtsso.TokenPair
}
