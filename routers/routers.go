package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
)

func NewAPI(cfg *config.Config, db *db.DB) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("api/users").Handler(registerUserRoutes(cfg, db))
	router.PathPrefix("api/organizations").Handler(registerOrganizationRoutes(cfg, db))
	router.PathPrefix("api/boards").Handler(registerBoardRoutes(cfg, db))

	// send hello world as json in temp route
	router.HandleFunc("/", func(writer http.ResponseWriter, reader *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	return router
}
