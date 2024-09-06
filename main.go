package main

import (
	"context"
	l "log"

	"matchamking/src/config"
	"matchamking/src/core"
	"matchamking/src/logger"
	"matchamking/src/server"
	"matchamking/src/server/command"
	"matchamking/src/storage"
)

func main() {
	ctx := context.Background()
	log, err := logger.NewLogger(logger.PRODUCTION)
	if err != nil {
		l.Fatalf("cannot create logger: %s", err.Error())
	}
	cfg, err := config.NewMSConfig()
	if err != nil {
		log.Fatal("cannot load configuration with error: %s", err.Error())
	}
	if err = storage.InitStorage(ctx, cfg.Storage); err != nil {
		log.Fatal("cannot create storage with error: %s", err.Error())
	}
	core.InitMatchmaker(cfg.MatchmakerConfig)
	manager := command.NewManager(log)
	s := server.NewServer(cfg.ServerConfig, manager, log)
	if err = s.Run(ctx); err != nil {
		log.Fatal("%s", err.Error())
	}
}
