package main

import (
	"context"
	l "log"

	"matchamking/src/config"
	"matchamking/src/core"
	"matchamking/src/logger"
	"matchamking/src/server"
	"matchamking/src/server/command"
)

func main() {
	log, err := logger.NewLogger(logger.PRODUCTION)
	if err != nil {
		l.Fatalf("cannot create logger: %s", err.Error())
	}
	cfg, err := config.NewMSConfig()
	if err != nil {
		log.Fatal("cannot load configuration with error: %s", err.Error())
	}
	manager := command.NewManager(log)
	core.InitMatchmaker(cfg.MatchmakerConfig)
	s := server.NewServer(cfg.ServerConfig, manager, log)
	if err = s.Run(context.Background()); err != nil {
		log.Fatal("%s", err.Error())
	}
}
