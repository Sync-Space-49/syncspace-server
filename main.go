package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/internal/config"

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

	// db, err := db.New(cfg.DB.DBUser, cfg.DB.DBPass, cfg.DB.DBURI, cfg.DB.DBName)

	mainRouter := mux.NewRouter()
	// test route
	mainRouter.HandleFunc("/", func(writer http.ResponseWriter, reader *http.Request) {
		// send hello world as json
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	http.ListenAndServe(cfg.APIHost, mainRouter)
	return nil
}
