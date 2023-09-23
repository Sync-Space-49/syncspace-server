package routers

import (
	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
)

func registerUserRoutes(cfg *config.Config, db *db.DB) *mux.Router {
	router := mux.NewRouter()

	return router
}
