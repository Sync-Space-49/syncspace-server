package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
	"github.com/Sync-Space-49/syncspace-server/routes"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := run(); err != nil {
		log.Err(err).Msg("failed to run server")
	}
}

func run() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	db, err := db.New(cfg.DB.DBUser, cfg.DB.DBPass, cfg.DB.DBURI, cfg.DB.DBName)

	server := &http.Server{
		Addr:    cfg.APIHost,
		Handler: routes.NewAPI(cfg, db),
	}
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed while running server: %w", err)
	}

	return nil
}
