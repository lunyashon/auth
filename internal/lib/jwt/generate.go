package jwtsso

import (
	"crypto/rsa"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(
	id int64,
	services []string,
	nameSSO string,
	accessTokenTTL, refreshTokenTTL time.Duration,
	privateKey *rsa.PrivateKey,
	ip, device string,
) (*TokenPair, error) {

	services = append(services, "sso.auth", "sso.auth.change.password")

	accessSigned, err := GenerateAccessToken(
		id,
		services,
		nameSSO,
		accessTokenTTL,
		privateKey,
		ip,
		device,
	)
	if err != nil {
		return nil, err
	}

	refreshSigned, err := GenerateRefreshToken(
		id,
		services,
		nameSSO,
		refreshTokenTTL,
		privateKey,
		ip,
		device,
	)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessSigned,
		RefreshToken: refreshSigned,
	}, nil
}

func GenerateRefreshToken(
	id int64,
	services []string,
	nameSSO string,
	refreshTokenTTL time.Duration,
	privateKey *rsa.PrivateKey,
	ip, device string,
) (string, error) {

	refreshClaims := UserClaims{
		TokenType: "refresh",
		IP:        jwt.ClaimStrings{ip},
		Device:    jwt.ClaimStrings{device},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(int(id)),
			Issuer:    nameSSO,
			Audience:  services,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshSigned, err := refreshToken.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return refreshSigned, nil
}

func GenerateAccessToken(
	id int64,
	services []string,
	nameSSO string,
	accessTokenTTL time.Duration,
	privateKey *rsa.PrivateKey,
	ip, device string,
) (string, error) {

	accessClaims := UserClaims{
		TokenType: "access",
		IP:        jwt.ClaimStrings{ip},
		Device:    jwt.ClaimStrings{device},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    nameSSO,
			Audience:  services,
			Subject:   strconv.Itoa(int(id)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessSigned, err := accessToken.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return accessSigned, nil
}
