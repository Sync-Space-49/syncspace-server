package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
)

var (
	usersPrefix         = "/api/users"
	organizationsPrefix = "/api/organizations"
	boardsPrefix        = "/api/organization/:organizationId/boards"
	// Probably add variables for prevelege ID's here
)

func NewAPI(cfg *config.Config, db *db.DB) http.Handler {
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	router := mux.NewRouter()
	router.PathPrefix(usersPrefix).Handler(registerUserRoutes(cfg, db))
	router.PathPrefix(organizationsPrefix).Handler(registerOrganizationRoutes(cfg, db))
	router.PathPrefix(boardsPrefix).Handler(registerBoardRoutes(cfg, db))

	// send hello world as json in temp route
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	return corsWrapper.Handler(router)
}
