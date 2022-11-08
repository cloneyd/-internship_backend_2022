package main

import (
	"avito-balance-service/config"
	"avito-balance-service/internal/server"
	"avito-balance-service/internal/storage/postgres"
	"context"
	"log"
)

func main() {
	log.Println("starting balance service...")
	cfgFile, err := config.LoadConfig("docker-config")
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	psqlDB, err := postgres.NewDB(context.Background(), cfg)
	if err != nil {
		log.Fatalf("error creating database connection: %v", err)
	}

	s := server.NewServer(cfg, psqlDB)
	if err = s.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
