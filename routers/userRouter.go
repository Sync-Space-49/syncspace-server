package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/Sync-Space-49/syncspace-server/aws"
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

	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUser))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateUser))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteUser))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/organizations", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserOrganizations))).Methods("GET")
	return handler.router
}

func (handler *userHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]
	if userId == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}

	user, err := handler.controller.GetUserById(userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(user)
}

func (handler *userHandler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]
	email := request.FormValue("email")
	username := request.FormValue("username")
	password := request.FormValue("password")

	if userId == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}

	err := request.ParseMultipartForm(10 << 20) // 10 MB limit on pfp size
	if err != nil {
		http.Error(writer, "Unable to parse form (is image too large?)", http.StatusBadRequest)
		return
	}
	var pfpUrl *string
	file, header, err := request.FormFile("profile_picture")
	if err == nil {
		fileExtension := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%s-pfp%s", userId, fileExtension)
		pfpUrl, err = aws.UploadPfp(file, filename)
		if err != nil {
			http.Error(writer, "Unable to upload file", http.StatusInternalServerError)
			return
		}
		defer file.Close()
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenUserId := token.RegisteredClaims.Subject
	if tokenUserId != userId {
		http.Error(writer, "Unauthorized to Update This User", http.StatusUnauthorized)
		return
	}

	err = handler.controller.UpdateUserById(userId, email, username, password, pfpUrl)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to Update User: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *userHandler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]
	if userId == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}
	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenUserId := token.RegisteredClaims.Subject
	if tokenUserId != userId {
		http.Error(writer, "Unauthorized to Delete This User", http.StatusUnauthorized)
		return
	}
	err := handler.controller.DeleteUserById(userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to Delete User: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *userHandler) GetUserOrganizations(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]
	if userId == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to get organizations for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	organizations, err := handler.controller.GetUserOrganizationsById(ctx, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get organizations for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(organizations)
}
