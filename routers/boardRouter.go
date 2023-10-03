package routers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/controllers/board"
	"github.com/Sync-Space-49/syncspace-server/db"
)

type boardHandler struct {
	router     *mux.Router
	controller *board.Controller
}

func registerBoardRoutes(cfg *config.Config, db *db.DB) *mux.Router {
	handler := &boardHandler{
		router:     mux.NewRouter(),
		controller: board.NewController(cfg, db),
	}
	// Grab board with id of boardId from organization with id of organizationId
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}", organizationsPrefix, boardsPrefix), handler.GetBoard).Methods("GET")

	// Edit/replace info about a list based on params
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}", organizationsPrefix, boardsPrefix), handler.UpdateList).Methods("PUT")

	// Delete a list from a board (including it's cards)
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}", organizationsPrefix, boardsPrefix), handler.DeleteList).Methods("DELETE")

	// Edit/replace info about a card based on params
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{CardId}", organizationsPrefix, boardsPrefix), handler.UpdateCard).Methods("PUT")

	// Grab all cards assigned to a user for a specific board
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{BoardMemberId}/assigned", organizationsPrefix, boardsPrefix), handler.GetAssignedCards).Methods("POST")

	// Grab all cards assigned to a user for all boards
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardMemberId}/assigned", organizationsPrefix, boardsPrefix), handler.GetAllAssignedCards).Methods("POST")

	// Assign a card to a user
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{ListId}/{CardId}/{BoardMemberId}", organizationsPrefix, boardsPrefix), handler.AssignCardToUser).Methods("POST")

	// Unassign a card from a user
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardId}/{ListId}/{CardId}/{BoardMemberId}", organizationsPrefix, boardsPrefix), handler.UnassignCardFromUser).Methods("DELETE")

	// add a tag to a board to use on tags
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/tag", organizationsPrefix, boardsPrefix), handler.CreateTag).Methods("POST")

	// delete a tag from a board
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/tag/{TagID}", organizationsPrefix, boardsPrefix), handler.DeleteTag).Methods("DELETE")

	// add a tag to a card
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardID}/{ListID}/{CardID}/{TagID}", organizationsPrefix, boardsPrefix), handler.AddTagToCard).Methods("POST")

	// delete a tag from a card
	handler.router.HandleFunc(fmt.Sprintf("%s/{OrganizationId}/%s/{BoardID}/{ListID}/{CardID}/{TagID}", organizationsPrefix, boardsPrefix), handler.AddTagToCard).Methods("DELETE")

	return handler.router
}

func (handler *boardHandler) GetBoard(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) UpdateList(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) DeleteList(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) UpdateCard(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) GetAssignedCards(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) GetAllAssignedCards(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) AssignCardToUser(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) UnassignCardFromUser(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) CreateTag(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) DeleteTag(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) AddTagToCard(writer http.ResponseWriter, request *http.Request) {

}

func (handler *boardHandler) RemoveTagFromCard(writer http.ResponseWriter, request *http.Request) {

}
