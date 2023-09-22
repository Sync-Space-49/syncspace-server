package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
)

func NewAPI(cfg *config.Config, db *db.DB) http.Handler {
	handler := mux.NewRouter()
	handler.PathPrefix("api/users").Handler(registerUserRoutes(cfg, db))
	handler.PathPrefix("api/organizations").Handler(registerOrganizationRoutes(cfg, db))
	handler.PathPrefix("api/boards").Handler(registerBoardRoutes(cfg, db))

	return handler
}
