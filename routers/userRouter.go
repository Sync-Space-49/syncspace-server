package routers

import (
	// "encoding/json"
	"fmt"
	// "io"
	"net/http"
	// "os"
	// "path/filepath"
	// "strconv"
	// "time"

	"github.com/gorilla/mux"
	// "golang.org/x/crypto/bcrypt"

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
	// params := mux.Vars(request)
	// username := params["username"]
	// email := params["email"]
	// password := params["password"]
	// if email == "" || password == "" || username == "" {
	// 	http.Error(writer, "Invalid email or password", http.StatusBadRequest)
	// 	return
	// }

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	http.Error(writer, "Failed to hash password", http.StatusInternalServerError)
	// 	return
	// }

	// err = request.ParseMultipartForm(10 << 20) // 10 MB limit on pfp size
	// if err != nil {
	// 	http.Error(writer, "Unable to parse form (is image too large?)", http.StatusBadRequest)
	// 	return
	// }
	// TODO: abstract this into function to be used in other places
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
	// TODO: Add user to database
	// TODO: Add user to their organization
	// writer.Header().Set("Content-Type", "application/json")
	// writer.WriteHeader(http.StatusCreated)
	// json.NewEncoder(writer).Encode(user)
}

func (handler *userHandler) SignInUser(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// credential := params["credential"]
	// password := params["password"]
	// if credential == "" || password == "" {
	// 	http.Error(writer, "Invalid email or password", http.StatusBadRequest)
	// 	return
	// }

	// TODO: GetUserByCredentials (email or username)
	// Send specific error if user not found/error fr thrown
	// user, err := handler.controller.GetUserByCredentials(credential)
	// if err != nil {
	// 	http.Error(writer, "Failed to get user", http.StatusInternalServerError)
	// 	return
	// }

	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	// 	http.Error(writer, "Incorrect password", http.StatusUnauthorized)
	// 	return
	// }

	// writer.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(writer).Encode(user)
}

func (handler *userHandler) UpdateUser(writer http.ResponseWriter, request *http.Request) {
	// TODO: verify this is signed in user (possilby with middleware)
	// Need some logic checking what here the user is updating
	// params := mux.Vars(request)
	// userId, err := strconv.Atoi(params["userId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid User ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// email := params["email"]
	// username := params["username"]
	// password := params["password"]
	// TODO: see if user is updating profile picture (possilby handle this in abstracted function)
	// TODO: update user in database
	// TODO: send back 204
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
