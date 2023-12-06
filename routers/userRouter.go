package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"

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

func registerUserRoutes(parentRouter *mux.Router, cfg *config.Config, db *db.DB) *mux.Router {
	handler := &userHandler{
		router:     parentRouter.NewRoute().Subrouter(),
		controller: user.NewController(cfg, db),
	}

	handler.router.Handle(usersPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.GetAllUsers))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateUser))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{userId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteUser))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/organizations", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserOrganizations))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/organizations/owned", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserOwnedOrganizations))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/boards", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserBoards))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/boards/owned", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserOwnerBoards))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/assigned", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserAssignedCards))).Methods("GET")

	handler.router.Handle(fmt.Sprintf("%s/{userId}/boards/favourite", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetUserFavouriteBoards))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/boards/favourite/{boardId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.AddUserFavouriteBoard))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{userId}/boards/favourite/{boardId}", usersPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.RemoveUserFavouriteBoard))).Methods("DELETE")

	return handler.router
}

func (handler *userHandler) GetAllUsers(writer http.ResponseWriter, request *http.Request) {
	users, err := handler.controller.GetUsers()
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(users)
}

func (handler *userHandler) GetUser(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]

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

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenUserId := token.RegisteredClaims.Subject
	if tokenUserId != userId {
		http.Error(writer, "Unauthorized to Update This User", http.StatusUnauthorized)
		return
	}

	// err := request.ParseMultipartForm(10 << 20) // 10 MB limit on pfp size
	// if err != nil {
	// 	http.Error(writer, "Unable to parse form (is image too large?)", http.StatusBadRequest)
	// 	return
	// }
	var pfpUrl *string
	pfpFile, _, err := request.FormFile("profile_picture")
	if err == nil {
		decodedPfp, _, err := image.Decode(pfpFile)
		defer pfpFile.Close()
		if err != nil {
			http.Error(writer, "Unable to decode image", http.StatusBadRequest)
			return
		}

		pfpDimensions := 512
		var rescaledPfp image.Image
		// make smallest dimension pfpDimensions
		if decodedPfp.Bounds().Dx() < decodedPfp.Bounds().Dy() {
			rescaledPfp = resize.Resize(uint(pfpDimensions), 0, decodedPfp, resize.Lanczos3)
		} else {
			rescaledPfp = resize.Resize(0, uint(pfpDimensions), decodedPfp, resize.Lanczos3)
		}
		pfpBounds := rescaledPfp.Bounds()
		pfpWidth := pfpBounds.Dx()
		pfpHeight := pfpBounds.Dy()
		x1, y1 := (pfpWidth/2)-(pfpDimensions/2), (pfpHeight/2)-(pfpDimensions/2)
		x2, y2 := x1+pfpDimensions, y1+pfpDimensions
		cropSize := image.Rect(x1, y1, x2, y2)
		// https://stackoverflow.com/questions/32544927/cropping-and-creating-thumbnails-with-go
		croppedPfp := rescaledPfp.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(cropSize)

		pfpBuffer := new(bytes.Buffer)
		err = png.Encode(pfpBuffer, croppedPfp)
		if err != nil {
			http.Error(writer, "Unable to rescale image", http.StatusBadRequest)
			return
		}

		fileExtension := ".png"
		filename := fmt.Sprintf("%s-pfp%s", userId, fileExtension)
		pfpUrl, err = aws.UploadPfp(bytes.NewReader(pfpBuffer.Bytes()), filename)
		if err != nil {
			http.Error(writer, "Unable to upload file", http.StatusInternalServerError)
			return
		}
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

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenUserId := token.RegisteredClaims.Subject
	if tokenUserId != userId {
		http.Error(writer, "Unauthorized to Delete This User", http.StatusUnauthorized)
		return
	}
	ctx := request.Context()
	err := handler.controller.DeleteUserById(ctx, userId)
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

func (handler *userHandler) GetUserOwnedOrganizations(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to get organizations for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	organizations, err := handler.controller.GetUserOwnedOrganizationsById(ctx, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get organizations for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(organizations)
}

func (handler *userHandler) GetUserBoards(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to get organizations for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	cards, err := handler.controller.GetUserBoardsById(ctx, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get organizations for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(cards)
}

func (handler *userHandler) GetUserOwnerBoards(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to get organizations for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	cards, err := handler.controller.GetUserOwnedBoardsById(ctx, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get organizations for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(cards)
}

func (handler *userHandler) GetUserAssignedCards(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to get organizations for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	cards, err := handler.controller.GetUserAssignedCardsById(ctx, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get organizations for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(cards)
}

func (handler *userHandler) GetUserFavouriteBoards(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to get favourite boards for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	boards, err := handler.controller.GetFavouriteBoards(ctx, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get favourite boards for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(boards)
}

func (handler *userHandler) AddUserFavouriteBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]
	boardId := params["boardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to add favourite board for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.AddFavouriteBoard(ctx, userId, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to add favourite board for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *userHandler) RemoveUserFavouriteBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId := params["userId"]
	boardId := params["boardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	if signedInUserId != userId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to remove favourite board for user with id %s", signedInUserId, userId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.RemoveFavouriteBoard(ctx, userId, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to remove favourite board for user with id %s: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
