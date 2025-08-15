package authgo

import (
	"context"
	"encoding/base64"
	"math/big"

	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
)

func (a AuthData) GetJWK(ctx context.Context) *sso.JWK {

	publicKey := a.KeysStore.PublicKey

	jwk := &sso.JWK{
		Kid: "",
		Kty: "RSA",
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   base64.RawStdEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}
	return jwk
}
