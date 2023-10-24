package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/auth"
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

	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUser))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateUser))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteUser))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/organizations", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserOrganizations))).Methods("POST")
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

	// TODO: see if user is updating profile picture (possilby handle this in abstracted function)
	// err := request.ParseMultipartForm(10 << 20) // 10 MB limit on pfp size
	// if err != nil {
	// 	http.Error(writer, "Unable to parse form (is image too large?)", http.StatusBadRequest)
	// 	return
	// }
	// // TODO: abstract this into function to be used in other places
	// file, header, err := request.FormFile("profile_picture")
	// if err != nil {
	// 	http.Error(writer, "Unable to retrieve profile picture", http.StatusBadRequest)
	// 	return
	// }
	// defer file.Close()

	// fileExtension := filepath.Ext(header.Filename)
	// filename := fmt.Sprintf("%s-%s%s", username, time.Now().Format(time.RFC3339), fileExtension)

	// destinationFile, err := os.Create(filename)
	// if err != nil {
	// 	http.Error(writer, "Unable to create or open destination file", http.StatusInternalServerError)
	// 	return
	// }
	// defer destinationFile.Close()
	// _, err = io.Copy(destinationFile, file)
	// if err != nil {
	// 	http.Error(writer, "Unable to copy file", http.StatusInternalServerError)
	// 	return
	// }
	// TODO: Add pfp to bucket
	// user currently cannot upload pfp
	pfpUrl := ""

	if userId == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenUserId := token.RegisteredClaims.Subject
	if tokenUserId != userId {
		http.Error(writer, "Unauthorized to Update This User", http.StatusUnauthorized)
		return
	}

	err := handler.controller.UpdateUserById(userId, email, username, password, pfpUrl)
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
