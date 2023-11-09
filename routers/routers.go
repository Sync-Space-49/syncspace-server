package routers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
)

const (
	usersPrefix         = "/api/users"
	organizationsPrefix = "/api/organizations"
	boardsPrefix        = "/api/organizations/{organizationId}/boards"
	panelsPrefix        = "/api/organizations/{organizationId}/boards/{boardId}/panels"
	stacksPrefix        = "/api/organizations/{organizationId}/boards/{boardId}/panels/{panelId}/stacks"
	cardsPrefix         = "/api/organizations/{organizationId}/boards/{boardId}/panels/{panelId}/stacks/{stackId}/cards"
)

func NewAPI(cfg *config.Config, db *db.DB) http.Handler {
	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	router := mux.NewRouter()
	router.PathPrefix(usersPrefix).Handler(registerUserRoutes(router, cfg, db))
	router.PathPrefix(organizationsPrefix).Handler(registerOrganizationRoutes(router, cfg, db))

	// send hello world as json in temp route
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(map[string]string{"message": "Hello World!"})
	})

	return corsWrapper.Handler(router)
}
