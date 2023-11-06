package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/controllers/board"
	"github.com/Sync-Space-49/syncspace-server/db"
)

type boardHandler struct {
	router     *mux.Router
	controller *board.Controller
}

func registerBoardRoutes(parentRouter *mux.Router, cfg *config.Config, db *db.DB) *mux.Router {
	handler := &boardHandler{
		router:     parentRouter.NewRoute().Subrouter(),
		controller: board.NewController(cfg, db),
	}
	// Grab board with id of boardId from organization with id of organizationId
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}", organizationsPrefix, boardsPrefix), handler.GetBoard).Methods("GET")

	// Edit/replace info about a list based on params
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}", organizationsPrefix, boardsPrefix), handler.UpdateList).Methods("PUT")

	// // Delete a list from a board (including it's cards)
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}", organizationsPrefix, boardsPrefix), handler.DeleteList).Methods("DELETE")

	// // Edit/replace info about a card based on params
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{CardId}", organizationsPrefix, boardsPrefix), handler.UpdateCard).Methods("PUT")

	// // Grab all cards assigned to a user for a specific board
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{BoardMemberId}/assigned", organizationsPrefix, boardsPrefix), handler.GetAssignedCards).Methods("POST")

	// // Grab all cards assigned to a user for all boards
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardMemberId}/assigned", organizationsPrefix, boardsPrefix), handler.GetAllAssignedCards).Methods("POST")

	// // Assign a card to a user
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{ListId}/{CardId}/{BoardMemberId}", organizationsPrefix, boardsPrefix), handler.AssignCardToUser).Methods("POST")

	// // Unassign a card from a user
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{ListId}/{CardId}/{BoardMemberId}", organizationsPrefix, boardsPrefix), handler.UnassignCardFromUser).Methods("DELETE")

	// // add a tag to a board to use on tags
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/tag", organizationsPrefix, boardsPrefix), handler.CreateTag).Methods("POST")

	// // delete a tag from a board
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/tag/{TagID}", organizationsPrefix, boardsPrefix), handler.DeleteTag).Methods("DELETE")

	// // add a tag to a card
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardID}/{ListID}/{CardID}/{TagID}", organizationsPrefix, boardsPrefix), handler.AddTagToCard).Methods("POST")

	// // delete a tag from a card
	// handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardID}/{ListID}/{CardID}/{TagID}", organizationsPrefix, boardsPrefix), handler.AddTagToCard).Methods("DELETE")

	// TODO: Update 'list' methods to 'stack' and 'panel'

	handler.router.Handle("/api/organizations/{organizationId}/boards", auth.EnsureValidToken()(http.HandlerFunc(handler.CreateBoard))).Methods("POST")
	handler.router.Handle("/api/organizations/{organizationId}/boards/{boardId}", auth.EnsureValidToken()(http.HandlerFunc(handler.GetBoard))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateBoard))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteBoard))).Methods("DELETE")

	// handler.router.Handle(fmt.Sprintf("%s/{boardId}/members", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetBoardMembers))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}/members", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.AddMemberToBoard))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}/members/{memberId}", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.RemoveMemberFromBoard))).Methods("DELETE")

	return handler.router
}

func (handler *boardHandler) CreateBoard(writer http.ResponseWriter, request *http.Request) {
	title := request.FormValue("title")
	isPrivate, err := strconv.ParseBool(request.FormValue("isPrivate"))
	params := mux.Vars(request)
	orgId := params["organizationId"]
	if title == "" {
		http.Error(writer, "No Title Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	ctx := request.Context()
	board, err := handler.controller.CreateBoard(ctx, userId, title, isPrivate, orgId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create board: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = handler.controller.InitializeBoard(userId, board.Id.String(), orgId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to initialize board: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *boardHandler) GetBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	// The methods below will never run as they give a 404

	// if organizationId == "" {
	// 	http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
	// 	return
	// }
	// if boardId == "" {
	// 	http.Error(writer, "No Board ID Found", http.StatusBadRequest)
	// 	return
	// }
	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	canReadBoard, err := auth.HasPermission(userId, readBoardPerm)

	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to read organization with id: %s", organizationId), http.StatusForbidden)
		return
	}
	if !canReadBoard {
		http.Error(writer, fmt.Sprintf("User does not have permission to read board with id: %s", boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	org, err := handler.controller.GetBoardById(ctx, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get board: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(org)
}

func (handler *boardHandler) DeleteBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	// fmt.Printf("1 userId: %s", userId)

	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	deleteBoardPerm := fmt.Sprintf("org%s:board%s:delete", organizationId, boardId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	canDeleteBoard, err := auth.HasPermission(userId, deleteBoardPerm)

	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to read org with id: %s", organizationId), http.StatusForbidden)
		return
	}
	if !canDeleteBoard {
		http.Error(writer, fmt.Sprintf("User does not have permission to delete board with id: %s", boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err = handler.controller.DeleteBoardById(ctx, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to delete board with id %s: %s", boardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

}

func (handler *boardHandler) UpdateBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]

	title := request.FormValue("title")
	ownerId := request.FormValue("ownerId")
	// fmt.Printf("ownerId: %s", ownerId)

	isPrivateStr := request.FormValue("isPrivate")
	isPrivate, err := strconv.ParseBool(isPrivateStr)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to parse isPrivate: %s", err.Error()), http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	// fmt.Printf("1 userId: %s", userId)

	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	updateBoardPerm := fmt.Sprintf("org%s:board%s:update", organizationId, boardId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	canUpdateBoard, err := auth.HasPermission(userId, updateBoardPerm)

	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to read organization with id: %s", organizationId), http.StatusForbidden)
		return
	}
	if !canUpdateBoard {
		http.Error(writer, fmt.Sprintf("User does not have permission to update board with id: %s", boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err = handler.controller.UpdateBoardById(ctx, organizationId, boardId, title, isPrivate, ownerId, userId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to update board with id %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

}

func (handler *boardHandler) AddMemberToBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	boardId := params["boardId"]
	if boardId == "" {
		http.Error(writer, "No Board ID Found", http.StatusBadRequest)
		return
	}
	memberId := request.FormValue("user_id")
	if memberId == "" {
		http.Error(writer, "No Member ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	addUsersPerm := fmt.Sprintf("org%s:board%s:add_members", organizationId, boardId)
	canAddUsers, err := auth.HasPermission(userId, addUsersPerm)

	// fmt.Println("addUsersPerm: ", addUsersPerm)
	// fmt.Println("canAddcanAddUsers: ", canAddUsers)

	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canAddUsers {
		http.Error(writer, fmt.Sprintf("User does not have permission to add users to board with id: %s", organizationId), http.StatusForbidden)
		return
	}

	err = handler.controller.AddMemberToBoard(memberId, organizationId, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users in board with id %s: %s", boardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *boardHandler) RemoveMemberFromBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	boardId := params["boardId"]
	if boardId == "" {
		http.Error(writer, "No Board ID Found", http.StatusBadRequest)
		return
	}
	memberId := params["memberId"]
	if memberId == "" {
		http.Error(writer, "No Member ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	removeUsersPerm := fmt.Sprintf("org%s:board%s:remove_members", organizationId, boardId)
	canRemoveUsers, err := auth.HasPermission(userId, removeUsersPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canRemoveUsers {
		http.Error(writer, fmt.Sprintf("User does not have permission to remove users from board with id: %s", organizationId), http.StatusForbidden)
		return
	}
	// fmt.Println("hit testing a on member: ", memberId)
	err = handler.controller.RemoveMemberFromBoard(memberId, organizationId, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users in board with id %s: %s", boardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

// func (handler *boardHandler) UpdateList(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) DeleteList(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) UpdateCard(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) GetAssignedCards(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) GetAllAssignedCards(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) AssignCardToUser(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) UnassignCardFromUser(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) CreateTag(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) DeleteTag(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) AddTagToCard(writer http.ResponseWriter, request *http.Request) {

// }

// func (handler *boardHandler) RemoveTagFromCard(writer http.ResponseWriter, request *http.Request) {

// }
