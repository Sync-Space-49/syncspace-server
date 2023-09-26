package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/controllers/user"
	"github.com/Sync-Space-49/syncspace-server/db"
)

type userHandler struct {
	router     *mux.Router
	controller *user.Controller
}

func registerUserRoutes(cfg *config.Config, db *db.DB) *mux.Router {
	handler := &userHandler{
		router:     mux.NewRouter(),
		controller: user.NewController(cfg, db),
	}

	handler.router.HandleFunc(fmt.Sprintf("%s/{userId}", usersPrefix), handler.GetUser).Methods("GET")
	// TODO: add more routes here

	return handler.router
}

func (handler *userHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(writer, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := handler.controller.GetUserById(userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user: %s", err.Error()), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}
