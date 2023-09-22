package routes

import (
	"net/http"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
)

func NewAPI(cfg *config.Config, db *db.DB) http.Handler {
	return nil
}
