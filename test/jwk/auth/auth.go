package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type JWKSClient struct {
	conn      *grpc.ClientConn
	ssoClient sso.AuthClient
	key       *rsa.PublicKey
	mu        sync.RWMutex
	expiry    time.Time
}

func NewJWKSClient(grpcAddr string) (*JWKSClient, error) {
	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &JWKSClient{
		conn:      conn,
		ssoClient: sso.NewAuthClient(conn),
	}, nil
}

func (c *JWKSClient) Close() error {
	return c.conn.Close()
}

func (c *JWKSClient) LoadKeys(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.ssoClient.GetJWKS(
		ctx,
		&sso.Empty{},
	)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, jwk := range resp.Keys {
		if jwk.Kty != "RSA" {
			return status.Error(codes.DataLoss, "kty is not RSA method")
		}

		nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to expanent public keys: %v", err)
		}

		eBytes, err := base64.RawStdEncoding.DecodeString(jwk.E)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to modulus public keys: %v", err)
		}

		publicKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}

		c.key = publicKey
	}

	c.expiry = time.Now().Add(1 * time.Hour)

	return nil
}

func (c *JWKSClient) CheckTokenTime() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Now().After(c.expiry) {
		return status.Errorf(codes.DeadlineExceeded, "JWK keys expetied")
	}

	return nil
}

func (c *JWKSClient) VerifyToken(ctx context.Context, tokenString string) (*jwt.Token, error) {

	parser := new(jwt.Parser)

	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse token: %v", err)
	}

	if err := c.CheckTokenTime(); err != nil {
		if loadErr := c.LoadKeys(ctx); loadErr != nil {
			return nil, status.Error(codes.Internal, "failed to load keys")
		}

		if err := c.CheckTokenTime(); err != nil {
			return nil, status.Error(codes.Internal, "key not found after refresh")
		}
	}

	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, status.Error(codes.InvalidArgument, "unexpected signing method")
		}
		return c.key, nil
	})
}
