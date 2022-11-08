package server

import (
	"avito-balance-service/config"
	"avito-balance-service/internal/storage/postgres"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	cfg        *config.Config
	httpServer *http.Server
	db         *postgres.DB
}

func NewServer(cfg *config.Config, db *postgres.DB) *Server {
	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         cfg.Server.Port,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
		},
		db: db,
	}
}

func (s *Server) Run() error {
	s.MapHandlers()

	go func() {
		log.Printf("Server is listening on port:%s\n", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error ListenAndServe: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), s.cfg.Server.CtxTimeout)
	defer shutdown()

	log.Println("Server exited properly")
	return s.httpServer.Shutdown(ctx)
}
