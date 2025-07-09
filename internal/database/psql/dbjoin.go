package database

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lunyashon/auth/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	instance *StructDatabase
	once     sync.Once
)

// Initialization database
// Return struct main database or error
func GetInstance(cfg config.ConfigEnv, log *slog.Logger) (*StructDatabase, error) {
	var err error
	once.Do(func() {
		provider := &DatabaseProvider{cfg: cfg, log: log}

		if err = provider.Connect(); err != nil {
			return
		}
		instance = &StructDatabase{
			Cfg:       cfg,
			Validator: provider,
			Base:      provider,
			User:      provider,
			Token:     provider,
		}
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect database %v", err)
	}

	return instance, nil
}

// Connect to Postgree
// Return error
func (s *DatabaseProvider) Connect() error {

	var err error

	dbName, dbData := s.GetData()
	db, err := sqlx.Open(dbName, dbData)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(15)
	db.SetConnMaxLifetime(5 * time.Minute)
	s.db = db
	return nil
}

// Close database
// No return
func (s *DatabaseProvider) Close() {
	s.db.Close()
}

// Get data from env file
// Return type database (Postgres) and data to connect
func (s *DatabaseProvider) GetData() (string, string) {
	return s.cfg.TYPE_DB,
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			s.cfg.HOST_DB,
			s.cfg.PORT_DB,
			s.cfg.USER_DB,
			s.cfg.PASS_DB,
			s.cfg.NAME_DB,
		)
}
