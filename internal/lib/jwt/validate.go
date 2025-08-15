package jwtsso

import (
	"crypto/rsa"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validate acces token with publick key
// Return user ID or error
func ValidateAccessToken(
	tokenString string,
	serviceName string,
	nameSSO string,
	publicKey *rsa.PublicKey,
) (*UserClaims, error) {

	switch {
	case tokenString == "":
		return nil, status.Error(codes.InvalidArgument, "authorization token is empty")
	case publicKey == nil:
		return nil, status.Error(codes.InvalidArgument, "public key is empty")
	}

	token, err := parseToken(tokenString, publicKey)
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, status.Error(codes.Unauthenticated, "token is expired")
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		switch {
		case claims.TokenType != "access":
			return nil, status.Error(codes.InvalidArgument, "token type is not access type")
		case claims.Subject == "":
			return nil, status.Error(codes.InvalidArgument, "subject is empty")
		case claims.Issuer != nameSSO:
			return nil, status.Error(codes.InvalidArgument, "issuer is invalid")
		case !slices.Contains(claims.Audience, serviceName):
			return nil, status.Errorf(codes.PermissionDenied, "access to '%s' denied", serviceName)
		}

		return claims, nil
	}

	return nil, status.Error(codes.Internal, "failed to valid token")
}

// Validate refresh token
// Return error or nil
func ValidateRefreshToken(
	tokenString string,
	nameSSO string,
	publicKey *rsa.PublicKey,
) (*UserClaims, error) {

	switch {
	case tokenString == "":
		return nil, status.Error(codes.InvalidArgument, "token is empty")
	case publicKey == nil:
		return nil, status.Error(codes.InvalidArgument, "public key is empty")
	}

	token, err := parseToken(tokenString, publicKey)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		switch {
		case claims.TokenType != "refresh":
			return nil, status.Error(codes.InvalidArgument, "token type is not refresh type")
		case claims.Issuer != nameSSO:
			return nil, status.Error(codes.InvalidArgument, "issuer is invalid")
		}

		return claims, nil
	}

	return nil, status.Error(codes.Internal, "failed to validate refresh token")
}

// Parsing token use a public key
// Return struct token or error
func parseToken(
	tokenString string,
	publicKey *rsa.PublicKey,
) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithLeeway(5*time.Second),
	)
}
