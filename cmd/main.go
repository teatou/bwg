package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/teatou/bwg/internal/config"
	"github.com/teatou/bwg/internal/http-server/handlers/add"
	"github.com/teatou/bwg/internal/http-server/handlers/fetch"
	"github.com/teatou/bwg/internal/storage/postgresql"
	"github.com/teatou/bwg/pkg/mylogger"
)

const (
	configEnv = "CONFIG"
)

func main() {
	conf, ok := os.LookupEnv(configEnv)
	if !ok {
		panic("no config env")
	}

	cfg, err := config.LoadConfig(conf)
	if err != nil {
		panic("uploading config error")
	}

	logger, err := mylogger.NewZapLogger(cfg.Logger.Level)
	if err != nil {
		panic("making mylogger error")
	}
	defer logger.Sync()

	storage, err := postgresql.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DbName)
	if err != nil {
		panic("failed to init storage")
	}

	r := chi.NewRouter()

	r.Post("/add_ticker", add.New(storage, logger))
	r.Get("/fetch", fetch.New(storage, logger))

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", cfg.Server.Port),
		Handler: r,
	}

	logger.Infof("starting server on port: %d", cfg.Server.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.Errorf("failed to start server")
	}

	logger.Errorf("server stopped")
}
