package jwtsso

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	TokenType string
	IP        jwt.ClaimStrings
	Device    jwt.ClaimStrings
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type KeysStore struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}
