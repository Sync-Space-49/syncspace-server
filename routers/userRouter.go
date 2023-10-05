package routers

import (
	"fmt"
	"net/http"
	"strconv"

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

	handler.router.HandleFunc(fmt.Sprintf("%s/{userId}", usersPrefix), handler.GetUser).Methods("GET")
	handler.router.HandleFunc(fmt.Sprintf("%s/signup", usersPrefix), handler.SignUpUser).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/signin", usersPrefix), handler.SignInUser).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{userId}", usersPrefix), handler.UpdateUser).Methods("PUT")
	handler.router.HandleFunc(fmt.Sprintf("%s/{userId}", usersPrefix), handler.DeleteUser).Methods("DELETE")

	return handler.router
}

func (handler *userHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// userId, err := strconv.Atoi(params["userId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid User ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }

	// user, err := handler.controller.GetUserById(userId)
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Failed to get user: %s", err.Error()), http.StatusInternalServerError)
	// 	return
	// }

	// writer.Header().Set("Content-Type", "application/json")
	// writer.WriteHeader(http.StatusOK)
	// json.NewEncoder(writer).Encode(user)
}

func (handler *userHandler) SignUpUser(writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")
	if email == "" || password == "" || username == "" {
		http.Error(writer, "Missing email, username, or password", http.StatusBadRequest)
		return
	}

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

	ctx := request.Context()
	err := handler.controller.CreateUser(ctx, username, email, password, nil)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create user: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// TODO: Add user to their organization
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *userHandler) SignInUser(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	credential := request.FormValue("credential")
	password := request.FormValue("password")
	if credential == "" || password == "" {
		http.Error(writer, "Invalid username, email or password", http.StatusBadRequest)
		return
	}

	user, err := handler.controller.GetUserByCredentials(ctx, credential, password)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user: %s", err.Error()), http.StatusUnauthorized)
		return
	}

	token, err := auth.CreateLoginToken(*user)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create token: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	writer.WriteHeader(http.StatusOK)
}

func (handler *userHandler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	// TODO: verify this is signed in user (possilby with middleware)
	// Need some logic checking what here the user is updating
	params := mux.Vars(request)
	userId, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid User ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	email := request.FormValue("email")
	username := request.FormValue("username")
	password := request.FormValue("password")
	// TODO: see if user is updating profile picture (possilby handle this in abstracted function)
	var pfpUrl *string = nil

	token := request.Header.Get("Authorization")
	if token == "" {
		http.Error(writer, "Missing Authorization Header", http.StatusUnauthorized)
		return
	}
	claims, err := auth.AuthenticateLoginToken(token)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to Authenticate Token: %s", err.Error()), http.StatusUnauthorized)
		return
	}
	if claims.Id != userId {
		http.Error(writer, "Unauthorized to Update This User", http.StatusUnauthorized)
	}

	ctx := request.Context()
	toUpdateUser, err := handler.controller.GetUserById(ctx, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to Find User: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if email == "" {
		email = toUpdateUser.Email
	}
	if username == "" {
		username = toUpdateUser.Username
	}
	if password == "" {
		password = toUpdateUser.HashedPassword
	} else {
		password, err = handler.controller.HashPassword(password)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to Hash Password: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
	err = handler.controller.UpdateUser(ctx, userId, email, username, password, pfpUrl)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to Update User: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *userHandler) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// userId, err := strconv.Atoi(params["userId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid User ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify this is signed in user (possilby with middleware)
	// TODO: Delete user from database
	// TODO: send back 204
}
