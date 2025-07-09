package internal

import (
	"context"
	"log/slog"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	appgrpc "github.com/lunyashon/auth/internal/app/grpc"
	"github.com/lunyashon/auth/internal/config"
	database "github.com/lunyashon/auth/internal/database/psql"
	jwtsso "github.com/lunyashon/auth/internal/lib/jwt"
	"github.com/lunyashon/auth/internal/lib/rabbit"
	"github.com/lunyashon/auth/internal/services/authgo"
)

type App struct {
	GRPCServer *appgrpc.App
	Rabbit     *rabbit.RabbitService
}

func New(log *slog.Logger, yaml *config.ConfigYaml, cfg *config.ConfigEnv, db *database.StructDatabase) *App {
	rabbit := rabbit.New(log, cfg, yaml.Rabbit)
	if err := rabbit.Connect.Connect(); err != nil {
		log.Error("Failed to connect to RabbitMQ", "error", err)
	}
	if err := rabbit.Connect.Channel(); err != nil {
		log.Error("Failed to create channel", "error", err)
	}

	queueForgotToken, err := rabbit.Send.DeclareQueue(context.Background(), cfg.RABBIT_QUEUE_FORGOT_TOKEN)
	if err != nil {
		log.Error("Failed to declare queue", "error", err)
	}

	queueConfirmEmail, err := rabbit.Send.DeclareQueue(context.Background(), cfg.RABBIT_QUEUE_CONFIRM_EMAIL)
	if err != nil {
		log.Error("Failed to declare queue", "error", err)
	}

	auth := authgo.AuthData{
		Log:               log,
		DB:                db,
		Rabbit:            rabbit,
		Yaml:              yaml,
		KeysStore:         parseKeys(cfg.PRIVATE_KEY),
		QueueForgotToken:  queueForgotToken,
		QueueConfirmEmail: queueConfirmEmail,
	}
	appGRPC := appgrpc.New(log, yaml.GRPS.Port, &auth)
	return &App{
		GRPCServer: appGRPC,
		Rabbit:     rabbit,
	}
}

func parseKeys(privateKey string) *jwtsso.KeysStore {
	privateKey = strings.ReplaceAll(privateKey, "\\n", "\n")

	privateKeyParse, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		panic(err)
	}

	return &jwtsso.KeysStore{
		PrivateKey: privateKeyParse,
		PublicKey:  &privateKeyParse.PublicKey,
	}
}

func (a *App) Shutdown(log *slog.Logger) {
	if err := a.Rabbit.Connect.CloseConnection(); err != nil {
		log.Error("Failed to close RabbitMQ connection", "error", err)
	}
	if err := a.Rabbit.Connect.CloseChannel(); err != nil {
		log.Error("Failed to close RabbitMQ channel", "error", err)
	}
	log.Info("RabbitMQ connection closed")
}
