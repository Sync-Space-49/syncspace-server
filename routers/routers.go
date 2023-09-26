package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
)

var (
	usersPrefix  = "api/users"
	orgsPrefix   = "api/organizations"
	boardsPrefix = "api/boards"
)

func NewAPI(cfg *config.Config, db *db.DB) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix(usersPrefix).Handler(registerUserRoutes(cfg, db))
	router.PathPrefix(orgsPrefix).Handler(registerOrganizationRoutes(cfg, db))
	router.PathPrefix(boardsPrefix).Handler(registerBoardRoutes(cfg, db))

	// send hello world as json in temp route
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	return router
}
