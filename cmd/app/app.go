package main

import (
	internal "github.com/lunyashon/auth/internal/app"
	"github.com/lunyashon/auth/internal/config"
	database "github.com/lunyashon/auth/internal/database/psql"
	loger "github.com/lunyashon/auth/internal/lib/log"
)

func main() {
	config.Load()

	cfg := config.LoadEnv()
	yaml := config.LoadYaml()

	log := loger.ExecLog(yaml.Env, yaml.PathToLog)

	// initialization database struct (Singlton)
	db, err := database.GetInstance(*cfg, log)

	if err != nil {
		panic(err)
	}
	defer db.Base.Close()

	application := internal.New(log, yaml, cfg, db)
	defer application.Shutdown(log)
	if err := application.GRPCServer.Run(); err != nil {
		panic(err)
	}
}
