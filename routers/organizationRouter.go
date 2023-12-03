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
	"github.com/Sync-Space-49/syncspace-server/controllers/organization"
	"github.com/Sync-Space-49/syncspace-server/controllers/user"
	"github.com/Sync-Space-49/syncspace-server/db"
)

type organizationHandler struct {
	router     *mux.Router
	controller *organization.Controller
}

func registerOrganizationRoutes(parentRouter *mux.Router, cfg *config.Config, db *db.DB) *mux.Router {
	handler := &organizationHandler{
		router:     parentRouter.NewRoute().Subrouter(),
		controller: organization.NewController(cfg, db),
	}
	handler.router.Handle(organizationsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.CreateOrganization))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetOrganization))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateOrganization))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteOrganization))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/members", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetOrganizationMembers))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/members", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.AddMemberToOrganization))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/members/{memberId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.RemoveMemberFromOrganization))).Methods("DELETE")
	handler.router.PathPrefix("{organizationId}/roles").Handler(registerRoleRoutes(handler.router, cfg, db))
	handler.router.PathPrefix("{organizationId}/boards").Handler(registerBoardRoutes(handler.router, cfg, db))
	return handler.router
}

func (handler *organizationHandler) CreateOrganization(writer http.ResponseWriter, request *http.Request) {
	title := request.FormValue("title")
	description := request.FormValue("description")
	if title == "" {
		http.Error(writer, "No Title Found", http.StatusBadRequest)
		return
	}
	// addition of aiEnabledString variable allows for us to default to 'false' if '' is passed
	aiEnabledString := request.FormValue("aiEnabled")
	if aiEnabledString == "" {
		aiEnabledString = "false"
	}
	aiEnabled, err := strconv.ParseBool(aiEnabledString)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to parse aiEnabledString: %s", err.Error()), http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	ctx := request.Context()
	org, err := handler.controller.CreateOrganization(ctx, userId, title, &description, aiEnabled)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create organization: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = handler.controller.InitializeOrganization(userId, org.Id.String())
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to initialize organization: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(org)
}

func (handler *organizationHandler) GetOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to read organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	org, err := handler.controller.GetOrganizationById(ctx, organizationId)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			http.Error(writer, fmt.Sprintf("No organization found with id %s", organizationId), http.StatusNotFound)
		} else {
			http.Error(writer, fmt.Sprintf("Failed to get organization: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(org)
}

func (handler *organizationHandler) UpdateOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]

	title := request.FormValue("title")
	description := request.FormValue("description")
	// addition of aiEnabledString variable allows for us to default to 'false' if '' is passed
	aiEnabledString := request.FormValue("ai_enabled")
	if aiEnabledString == "" {
		aiEnabledString = "false"
	}
	aiEnabled, err := strconv.ParseBool(aiEnabledString)
	print(aiEnabled)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to parse aiEnabledString: %s", err.Error()), http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	updateOrgPerm := fmt.Sprintf("org%s:update", organizationId)
	canUpdateOrg := tokenCustomClaims.HasPermission(updateOrgPerm)
	if !canUpdateOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to update organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err = handler.controller.UpdateOrganizationById(ctx, organizationId, title, description, aiEnabled)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to update organization with id %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)

}

func (handler *organizationHandler) DeleteOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	deleteOrgPerm := fmt.Sprintf("org%s:delete", organizationId)
	canDeleteOrg := tokenCustomClaims.HasPermission(deleteOrgPerm)
	if !canDeleteOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to delete organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err := handler.controller.DeleteOrganizationById(ctx, organizationId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to delete organization with id %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}

	orgRolePrefix := fmt.Sprintf("org%s:", organizationId)
	orgRoles, err := auth.GetRoles(&orgRolePrefix)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get roles for organization %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	if len(*orgRoles) == 0 {
		http.Error(writer, fmt.Sprintf("No roles found for organization %s", organizationId), http.StatusInternalServerError)
		return
	}
	for _, role := range *orgRoles {
		permissions, err := auth.GetRolePermissions(role.Id)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to get permissions for role %s: %s", role.Id, err.Error()), http.StatusInternalServerError)
			return
		}
		if len(*permissions) != 0 {
			err = auth.DeletePermissions(*permissions)
			if err != nil {
				http.Error(writer, fmt.Sprintf("Failed to delete permissions for role %s: %s", role.Id, err.Error()), http.StatusInternalServerError)
				return
			}
		}
		err = auth.DeleteRole(role.Id)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to delete role %s: %s", role.Id, err.Error()), http.StatusInternalServerError)
			return
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *organizationHandler) GetOrganizationMembers(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg := tokenCustomClaims.HasPermission(readOrgPerm)
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to read organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	users, err := user.GetOrgMembers(organizationId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users in org with id %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(users)
}

func (handler *organizationHandler) AddMemberToOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]

	newMemberId := request.FormValue("user_id")
	if newMemberId == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	addUsersPerm := fmt.Sprintf("org%s:add_members", organizationId)
	canAddUsers := tokenCustomClaims.HasPermission(addUsersPerm)
	if !canAddUsers {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to add users to organization with id: %s", newMemberId, organizationId), http.StatusForbidden)
		return
	}

	err := handler.controller.AddMember(newMemberId, organizationId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to add user with id %s to org with id %s: %s", newMemberId, organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *organizationHandler) RemoveMemberFromOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	memberId := params["memberId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	tokenCustomClaims := token.CustomClaims.(*auth.CustomClaims)
	userId := token.RegisteredClaims.Subject
	removeUsersPerm := fmt.Sprintf("org%s:remove_members", organizationId)
	canRemoveUsers := tokenCustomClaims.HasPermission(removeUsersPerm)
	if !canRemoveUsers && userId != memberId {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to remove users from organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	err := handler.controller.RemoveMember(memberId, organizationId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to remove user with id %s from org with id %s: %s", memberId, organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}
