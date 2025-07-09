package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lunyashon/auth/test/jwk/auth"
	ssov1 "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UpdateAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func Validate(jwksClient *auth.JWKSClient) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "failed to parse authorization token"},
			)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "failed to `bearer token`"},
			)
			return
		}

		tokenString := tokenParts[1]

		conn, err := grpc.NewClient(
			"localhost:50551",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
		}

		sso := ssov1.NewAuthClient(conn)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		resp, err := sso.ValidateToken(
			ctx,
			&ssov1.ValidateRequest{
				AccessToken: tokenString,
				Service:     "notes",
			})

		fmt.Println(resp, err)

		c.Next()
	}
}

func UpdateAccessToken(jwksClient *auth.JWKSClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "failed to parse authorization token"},
			)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "failed to `bearer token`"},
			)
			return
		}

		var data *UpdateAccessTokenRequest

		body, err := c.GetRawData()
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
		}

		json.Unmarshal(body, &data)

		conn, err := grpc.NewClient(
			"localhost:50551",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": err.Error()},
			)
		}

		sso := ssov1.NewAuthClient(conn)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		resp, err := sso.UpdateAccessToken(
			ctx,
			&ssov1.AccessTokenRequest{
				RefreshToken: data.RefreshToken,
			},
		)

		fmt.Println(resp, err)

		c.Next()
	}
}
