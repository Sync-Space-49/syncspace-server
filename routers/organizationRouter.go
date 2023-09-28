package routers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/controllers/organization"
	"github.com/Sync-Space-49/syncspace-server/db"
)

type organizationHandler struct {
	router     *mux.Router
	controller *organization.Controller
}

func registerOrganizationRoutes(cfg *config.Config, db *db.DB) *mux.Router {
	handler := &organizationHandler{
		router:     mux.NewRouter(),
		controller: organization.NewController(cfg, db),
	}

	// Middleware to be added for authorization
	handler.router.HandleFunc(fmt.Sprintf("%s", organizationsPrefix), handler.CreateOrganization).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), handler.GetOrganization).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), handler.UpdateOrganization).Methods("PUT")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), handler.DeleteOrganization).Methods("DELETE")
	handler.router.HandleFunc(fmt.Sprintf("%s/{userId}", organizationsPrefix), handler.GetUserOrganizations).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/members", organizationsPrefix), handler.GetOrganizationMembers).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/members/new", organizationsPrefix), handler.AddMemberToOrganization).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/members/{memberId}", organizationsPrefix), handler.GetOrganizationMember).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/members/{memberId}", organizationsPrefix), handler.RemoveMemberFromOrganization).Methods("DELETE")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles", organizationsPrefix), handler.GetOrganizationRoles).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/new", organizationsPrefix), handler.AddOrganizationRole).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), handler.GetOrganizationRole).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), handler.UpdateOrganizationRole).Methods("PUT")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), handler.DeleteOrganizationRole).Methods("DELETE")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/privileges", organizationsPrefix), handler.GetOrganizationRolePrivileges).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/privileges/{privilegeId}", organizationsPrefix), handler.AddOrganizationRolePrivilege).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/privileges/{privilegeId}", organizationsPrefix), handler.RemoveOrganizationRolePrivilege).Methods("DELETE")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/:organizationId/roles/:roleId/:memberId", organizationsPrefix), handler.AddMemberToRole).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/:organizationId/roles/:roleId/:memberId", organizationsPrefix), handler.RemoveMemberFromRole).Methods("DELETE")

	return handler.router
}

func (handler *organizationHandler) CreateOrganization(writer http.ResponseWriter, request *http.Request) {
}

func (handler *organizationHandler) GetOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: Get organization by ID
	// TODO: Return organization with status code 201
}

func (handler *organizationHandler) UpdateOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: Find what fields are being updated
	// TODO: Get organization by ID
	// TODO: update organization in database
	// TODO: send back 204
}

func (handler *organizationHandler) DeleteOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user is owner of organization
	// TODO: Delete organization from database
	// TODO: send back 204
}

func (handler *organizationHandler) GetUserOrganizations(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	userId, err := strconv.Atoi(params["userId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid User ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user is signed in user (posisbly with middleware)
	// TODO: Get organizations by user ID
	// TODO: Return organizations with status code 200
}

func (handler *organizationHandler) GetOrganizationMembers(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get members by organization ID
	// TODO: Return members with status code 200
}

func (handler *organizationHandler) AddMemberToOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can add members to organization (posisbly with middleware)
	// TODO: Add member to organization
	// TODO: Return member with status code 201
}

func (handler *organizationHandler) GetOrganizationMember(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	memberId, err := strconv.Atoi(params["memberId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Member ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get member by organization ID and member ID
	// TODO: Return member with status code 200
}

func (handler *organizationHandler) RemoveMemberFromOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	memberId, err := strconv.Atoi(params["memberId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Member ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can remove members from organization (posisbly with middleware)
	// TODO: Remove member from organization
	// TODO: send back 204
}

func (handler *organizationHandler) GetOrganizationRoles(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get roles by organization ID
	// TODO: Return roles with status code 200
}

func (handler *organizationHandler) AddOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can add roles to organization (posisbly with middleware)
	// TODO: Add role to organization
	// TODO: Return role with status code 201
}

func (handler *organizationHandler) GetOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get role by organization ID and role ID
	// TODO: Return role with status code 200
}

func (handler *organizationHandler) UpdateOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can update roles in organization (posisbly with middleware)
	// TODO: Find what fields are being updated
	// TODO: Get role by organization ID and role ID
	// TODO: update role in database
	// TODO: send back 204
}

func (handler *organizationHandler) DeleteOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can delete roles in organization (posisbly with middleware)
	// TODO: Delete role from database
	// TODO: send back 204
}

func (handler *organizationHandler) GetOrganizationRolePrivileges(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	roleId, err := strconv.Atoi(params["roleId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get privileges by organization ID and role ID
	// TODO: Return role privileges with status code 200
}

func (handler *organizationHandler) AddOrganizationRolePrivilege(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	roleId, err := strconv.Atoi(params["roleId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can add privileges to role in organization (posisbly with middleware)
	// TODO: Add privilege to role
	// TODO: Return privilege with status code 201
	// TODO: send back 204
}

func (handler *organizationHandler) RemoveOrganizationRolePrivilege(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	roleId, err := strconv.Atoi(params["roleId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can remove privileges from role in organization (posisbly with middleware)
	// TODO: Remove privilege from role
	// TODO: send back 204
}

func (handler *organizationHandler) AddMemberToRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	memberId, err := strconv.Atoi(params["memberId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Member ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	roleId, err := strconv.Atoi(params["roleId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can add members to role in organization (posisbly with middleware)
	// TODO: Add member to role
	// TODO: Return with status code 204
}

func (handler *organizationHandler) RemoveMemberFromRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId, err := strconv.Atoi(params["organizationId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	memberId, err := strconv.Atoi(params["memberId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Member ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	roleId, err := strconv.Atoi(params["roleId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
		return
	}
	// TODO: verify user sending request can remove members from role in organization (posisbly with middleware)
	// TODO: Remove member from role
	// TODO: Return with status code 204
}
