package main

import (
	"flag"
	"log"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"
	"tg-bot-golang/internal/server"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()

	srv := server.NewServer(appLogger, cfg)
	appLogger.Fatal(srv.Run())
}
