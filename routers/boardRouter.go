package routers

import (
	"database/sql"
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

	// TODO: add routes for below
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

	handler.router.Handle(boardsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.GetAllBoards))).Methods("GET")
	handler.router.Handle(boardsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.CreateBoard))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetBoard))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateBoard))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteBoard))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}/details", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetCompleteBoard))).Methods("GET")

	// Because a board member is known through a role, these routes could possilby be removed or refactroed to call the role routes
	// handler.router.Handle(fmt.Sprintf("%s/{boardId}/members", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetBoardMembers))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}/members", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.AddMemberToBoard))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{boardId}/members/{memberId}", boardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.RemoveMemberFromBoard))).Methods("DELETE")

	handler.router.Handle(panelsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.GetPanels))).Methods("GET")
	handler.router.Handle(panelsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.CreatePanel))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{panelId}", panelsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetPanel))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{panelId}", panelsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdatePanel))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{panelId}", panelsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeletePanel))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{panelId}/details", panelsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetCompletePanel))).Methods("GET")

	handler.router.Handle(stacksPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.GetStacks))).Methods("GET")
	handler.router.Handle(stacksPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.CreateStack))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{stackId}", stacksPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetStack))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{stackId}", stacksPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateStack))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{stackId}", stacksPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteStack))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{stackId}/details", stacksPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetCompleteStack))).Methods("GET")

	handler.router.Handle(cardsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.GetCards))).Methods("GET")
	handler.router.Handle(cardsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.CreateCard))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{cardId}", cardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetCard))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{cardId}", cardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateCard))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{cardId}", cardsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteCard))).Methods("DELETE")

	return handler.router
}

func (handler *boardHandler) GetAllBoards(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	visibleBoards, err := handler.controller.GetViewableBoardsInOrg(ctx, tokenCustomClaims, organizationId, userId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No organization found with id %s", organizationId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get viewable boards: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(visibleBoards)
}

func (handler *boardHandler) CreateBoard(writer http.ResponseWriter, request *http.Request) {
	title := request.FormValue("title")
	isPrivate, err := strconv.ParseBool(request.FormValue("isPrivate"))
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to parse isPrivate: %s", err.Error()), http.StatusBadRequest)
		return
	}
	params := mux.Vars(request)
	orgId := params["organizationId"]
	if title == "" {
		http.Error(writer, "No Title Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	// No check if can create board?
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

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	board, err := handler.controller.GetBoardById(ctx, boardId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No board found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get board: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if board.IsPrivate {
		readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
		canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
		if !canReadBoard {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(board)
}

func (handler *boardHandler) UpdateBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]

	title := request.FormValue("title")
	ownerId := request.FormValue("ownerId")

	isPrivate, err := strconv.ParseBool(request.FormValue("isPrivate"))
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to parse isPrivate: %s", err.Error()), http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	updateBoardPerm := fmt.Sprintf("org%s:board%s:update", organizationId, boardId)
	canUpdateBoard := tokenCustomClaims.HasPermission(updateBoardPerm)
	if !canUpdateBoard {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to update board with id: %s", userId, boardId), http.StatusForbidden)
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

func (handler *boardHandler) DeleteBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read org with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	deleteBoardPerm := fmt.Sprintf("org%s:board%s:delete", organizationId, boardId)
	canDeleteBoard := tokenCustomClaims.HasPermission(deleteBoardPerm)
	if !canDeleteBoard {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to delete board with id: %s", userId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.DeleteBoardById(ctx, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to delete board with id %s: %s", boardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *boardHandler) GetCompleteBoard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	board, err := handler.controller.GetCompleteBoardById(ctx, boardId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No board found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get board: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if board.IsPrivate {
		readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
		canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
		if !canReadBoard {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(board)
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
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	addUsersPerm := fmt.Sprintf("org%s:board%s:add_members", organizationId, boardId)
	canAddUsers := tokenCustomClaims.HasPermission(addUsersPerm)
	if !canAddUsers {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to add users to org %s board with id %s", userId, organizationId, boardId), http.StatusForbidden)
		return
	}

	err := handler.controller.AddMemberToBoard(memberId, organizationId, boardId)
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
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	removeUsersPerm := fmt.Sprintf("org%s:board%s:remove_members", organizationId, boardId)
	canRemoveUsers := tokenCustomClaims.HasPermission(removeUsersPerm)
	if !canRemoveUsers {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to add users to org %s board with id %s", userId, organizationId, boardId), http.StatusForbidden)
		return
	}
	// fmt.Println("hit testing a on member: ", memberId)
	err := handler.controller.RemoveMemberFromBoard(memberId, organizationId, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users in board with id %s: %s", boardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *boardHandler) GetPanels(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s  does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
	canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
	if !canReadBoard {
		http.Error(writer, fmt.Sprintf("User with id %s  does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	panels, err := handler.controller.GetPanelsByBoardId(ctx, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get panels from board with id %s: %s", boardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(panels)
}

func (handler *boardHandler) CreatePanel(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]

	title := request.FormValue("title")
	if title == "" {
		http.Error(writer, "No Title Found", http.StatusBadRequest)
		return
	}
	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	createPanelPerm := fmt.Sprintf("org%s:board%s:create_panel", organizationId, boardId)
	canCreatePanel := tokenCustomClaims.HasPermission(createPanelPerm)
	if !canCreatePanel {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to create panel on board with id: %s", userId, boardId), http.StatusForbidden)
		return
	}

	err := handler.controller.CreatePanel(request.Context(), title, boardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create panel: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *boardHandler) GetPanel(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	panelId := params["panelId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
	canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
	if !canReadBoard {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	panel, err := handler.controller.GetPanelById(ctx, panelId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No panel with id %s found", panelId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get panel: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(panel)
}

func (handler *boardHandler) UpdatePanel(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	panelId := params["panelId"]

	title := request.FormValue("title")
	var position *int
	if request.FormValue("position") != "" {
		var err error
		tempPosition, err := strconv.Atoi(request.FormValue("position"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to parse position: %s", err.Error()), http.StatusBadRequest)
			return
		}
		position = &tempPosition
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	updatePanelPerm := fmt.Sprintf("org%s:board%s:update_panel", organizationId, boardId)
	canUpdatePanel := tokenCustomClaims.HasPermission(updatePanelPerm)
	if !canUpdatePanel {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to update panel %s on board with id: %s", userId, panelId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.UpdatePanelById(ctx, boardId, panelId, title, position)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to update panel with id %s: %s", panelId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *boardHandler) DeletePanel(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	panelId := params["panelId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	deletePanelPerm := fmt.Sprintf("org%s:board%s:delete_panel", organizationId, boardId)
	canDeletePanel := tokenCustomClaims.HasPermission(deletePanelPerm)
	if !canDeletePanel {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to delete panel %s on board with id: %s", userId, panelId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.DeletePanelById(ctx, boardId, panelId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to delete panel with id %s: %s", panelId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *boardHandler) GetCompletePanel(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	panelId := params["panelId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	board, err := handler.controller.GetBoardById(ctx, boardId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No board found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get board: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}
	if board.IsPrivate {
		readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
		canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
		if !canReadBoard {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
			return
		}
	}

	panel, err := handler.controller.GetCompletePanelById(ctx, panelId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No panel with id %s found", panelId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get panel: %s", err.Error()), http.StatusInternalServerError)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(panel)
}

func (handler *boardHandler) GetStacks(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	panelId := params["panelId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	stacks, err := handler.controller.GetStacksByPanelId(ctx, panelId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get stacks from panel with id %s: %s", panelId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(stacks)
}

func (handler *boardHandler) CreateStack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	panelId := params["panelId"]

	title := request.FormValue("title")
	if title == "" {
		http.Error(writer, "No Title Found", http.StatusBadRequest)
		return
	}
	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	createStackPerm := fmt.Sprintf("org%s:board%s:create_stack", organizationId, boardId)
	canCreateStack := tokenCustomClaims.HasPermission(createStackPerm)
	if !canCreateStack {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to create stack on board with id: %s", userId, boardId), http.StatusForbidden)
		return
	}

	err := handler.controller.CreateStack(request.Context(), title, panelId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create stack: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *boardHandler) GetStack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	stackId := params["stackId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	board, err := handler.controller.GetBoardById(ctx, boardId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No board found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get board: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if board.IsPrivate {
		readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
		canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
		if !canReadBoard {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
			return
		}
	}

	stack, err := handler.controller.GetStackById(ctx, stackId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No stack found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get stack: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(stack)
}

func (handler *boardHandler) UpdateStack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	panelId := params["panelId"]
	stackId := params["stackId"]

	title := request.FormValue("title")
	var position *int
	if request.FormValue("position") != "" {
		var err error
		tempPosition, err := strconv.Atoi(request.FormValue("position"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to parse position: %s", err.Error()), http.StatusBadRequest)
			return
		}
		position = &tempPosition
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	updateStackPerm := fmt.Sprintf("org%s:board%s:update_stack", organizationId, boardId)
	canUpdateStack := tokenCustomClaims.HasPermission(updateStackPerm)
	if !canUpdateStack {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to update stack %s in board with id: %s", userId, stackId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.UpdateStackById(ctx, panelId, stackId, title, position)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to update stack with id %s: %s", stackId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *boardHandler) DeleteStack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	panelId := params["panelId"]
	stackId := params["stackId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	deleteStackPerm := fmt.Sprintf("org%s:board%s:delete_stack", organizationId, boardId)
	canDeleteStack := tokenCustomClaims.HasPermission(deleteStackPerm)
	if !canDeleteStack {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to delete stack %s on board with id: %s", userId, panelId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.DeleteStackById(ctx, panelId, stackId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to delete stack with id %s: %s", stackId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *boardHandler) GetCompleteStack(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	stackId := params["stackId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	board, err := handler.controller.GetBoardById(ctx, boardId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No board found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get board: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if board.IsPrivate {
		readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
		canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
		if !canReadBoard {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
			return
		}
	}

	stack, err := handler.controller.GetCompleteStackById(ctx, stackId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No stack with id %s found", stackId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get stack: %s", err.Error()), http.StatusInternalServerError)
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(stack)
}

func (handler *boardHandler) GetCards(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	stackId := params["stackId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	cards, err := handler.controller.GetCardsByStackId(ctx, stackId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get cards from stack with id %s: %s", stackId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(cards)
}

func (handler *boardHandler) CreateCard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	stackId := params["stackId"]

	title := request.FormValue("title")
	if title == "" {
		http.Error(writer, "No Title Found", http.StatusBadRequest)
		return
	}
	description := request.FormValue("description")
	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	createCardPerm := fmt.Sprintf("org%s:board%s:create_card", organizationId, boardId)
	canCreateCard := tokenCustomClaims.HasPermission(createCardPerm)
	if !canCreateCard {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to create a card on board with id: %s", userId, boardId), http.StatusForbidden)
		return
	}

	err := handler.controller.CreateCard(request.Context(), title, description, stackId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create card: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *boardHandler) GetCard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	cardId := params["cardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	board, err := handler.controller.GetBoardById(ctx, boardId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No board found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get board: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	if board.IsPrivate {
		readBoardPerm := fmt.Sprintf("org%s:board%s:read", organizationId, boardId)
		canReadBoard := tokenCustomClaims.HasPermission(readBoardPerm)
		if !canReadBoard {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read board with id: %s", userId, boardId), http.StatusForbidden)
			return
		}
	}

	stack, err := handler.controller.GetCardById(ctx, cardId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No card found with id %s", boardId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get card: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(stack)
}

func (handler *boardHandler) UpdateCard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	stackId := params["stackId"]
	cardId := params["cardId"]

	title := request.FormValue("title")
	description := request.FormValue("description")
	var position *int
	if request.FormValue("position") != "" {
		var err error
		tempPosition, err := strconv.Atoi(request.FormValue("position"))
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to parse position: %s", err.Error()), http.StatusBadRequest)
			return
		}
		position = &tempPosition
	}
	newStackId := request.FormValue("stack_id")

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	updateCardPerm := fmt.Sprintf("org%s:board%s:update_card", organizationId, boardId)
	canUpdateCard := tokenCustomClaims.HasPermission(updateCardPerm)
	if !canUpdateCard {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to update card %s in board with id: %s", userId, cardId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.UpdateCardById(ctx, boardId, stackId, cardId, newStackId, title, description, position)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to update card with id %s: %s", cardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *boardHandler) DeleteCard(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	boardId := params["boardId"]
	stackId := params["stackId"]
	cardId := params["cardId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	deleteCardPerm := fmt.Sprintf("org%s:board%s:delete_card", organizationId, boardId)
	canDeleteCard := tokenCustomClaims.HasPermission(deleteCardPerm)
	if !canDeleteCard {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to delete card %s on board with id: %s", userId, cardId, boardId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.DeleteCardById(ctx, stackId, cardId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to delete card with id %s: %s", cardId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

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
