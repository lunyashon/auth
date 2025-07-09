package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	corsgin "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lunyashon/auth/test/jwk/auth"
	middleware "github.com/lunyashon/auth/test/jwk/check"
	sso "github.com/lunyashon/protobuf/auth/gen/go/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// jwk, err := auth.NewJWKSClient("localhost:50551")
	// if err != nil {
	// 	panic(err)
	// }
	// defer jwk.Close()

	// ctx := context.Background()
	// if err := jwk.LoadKeys(ctx); err != nil {
	// 	panic(err)
	// }

	var wg sync.WaitGroup

	wg.Add(2)

	// go func() {
	// 	ticker := time.NewTicker(59 * time.Minute)
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// 			defer cancel()

	// 			if err := jwk.LoadKeys(ctx); err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 	}
	// }()

	go ginHundler(&auth.JWKSClient{}, &wg)
	go gRPCHundler(&wg)

	wg.Wait()

}

func ginHundler(jwks *auth.JWKSClient, wg *sync.WaitGroup) {

	defer wg.Done()

	router := gin.Default()

	router.Use(corsgin.New(corsgin.Config{
		AllowOrigins:     []string{"https://24microservice.ru", "https://sso,24microservice.ru"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/public", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Public enpoint"})
	})

	protected := router.Group("/private")
	protected.Use(middleware.Validate(jwks))
	{
		protected.POST("/profile", func(c *gin.Context) {
			// c.JSON(200, gin.H{
			// 	"user_id": c.MustGet("userID").(float64),
			// })
		})
	}
	protected.Use(middleware.UpdateAccessToken(jwks))
	{
		protected.POST("/update.access.token", func(c *gin.Context) {
			// c.JSON(200, gin.H{
			// 	"user_id": c.MustGet("userID").(float64),
			// })
		})
	}

	if err := router.Run(":50555"); err != nil {
		panic("failed to start server on port")
	}
}

func gRPCHundler(wg *sync.WaitGroup) {

	defer wg.Done()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := sso.RegisterAuthHandlerFromEndpoint(
		context.Background(),
		mux,
		"localhost:50551",
		opts,
	); err != nil {
		log.Fatalf("err in register endpoint %v", err)
	}

	fmt.Println("Server start in port 50550")
	log.Fatal(http.ListenAndServe(":50550", mux))
}
