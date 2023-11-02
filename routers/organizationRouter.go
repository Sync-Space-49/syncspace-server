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
	"github.com/Sync-Space-49/syncspace-server/controllers/organization"
	"github.com/Sync-Space-49/syncspace-server/controllers/user"
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
	handler.router.Handle(organizationsPrefix, auth.EnsureValidToken()(http.HandlerFunc(handler.CreateOrganization))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetOrganization))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateOrganization))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteOrganization))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/members", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetOrganizationMembers))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/members", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.AddMemberToOrganization))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/members/{memberId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.RemoveMemberFromOrganization))).Methods("DELETE")
	// TODO: controller methods below
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles", organizationsPrefix), handler.GetOrganizationRoles).Methods("GET")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles", organizationsPrefix), handler.AddOrganizationRole).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), handler.GetOrganizationRole).Methods("GET")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), handler.UpdateOrganizationRole).Methods("PUT")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), handler.DeleteOrganizationRole).Methods("DELETE")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/privileges", organizationsPrefix), handler.GetOrganizationRolePrivileges).Methods("GET")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/privileges/{privilegeId}", organizationsPrefix), handler.AddOrganizationRolePrivilege).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/privileges/{privilegeId}", organizationsPrefix), handler.RemoveOrganizationRolePrivilege).Methods("DELETE")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/{memberId}", organizationsPrefix), handler.AddMemberToRole).Methods("POST")
	handler.router.HandleFunc(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/{memberId}", organizationsPrefix), handler.RemoveMemberFromRole).Methods("DELETE")
	fmt.Println("Returning orgs router...")
	return handler.router
}

func (handler *organizationHandler) CreateOrganization(writer http.ResponseWriter, request *http.Request) {
	title := request.FormValue("title")
	description := request.FormValue("description")
	if title == "" {
		http.Error(writer, "No Title Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	ctx := request.Context()
	org, err := handler.controller.CreateOrganization(ctx, userId, title, &description)
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
}

func (handler *organizationHandler) GetOrganization(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Testing Hit GetOrganization")
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to read organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	org, err := handler.controller.GetOrganizationById(ctx, organizationId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get organization: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(org)
}

func (handler *organizationHandler) UpdateOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	title := request.FormValue("title")
	description := request.FormValue("description")

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	updateOrgPerm := fmt.Sprintf("org%s:update", organizationId)
	canReadOrg, err := auth.HasPermission(userId, updateOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to update organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err = handler.controller.UpdateOrganizationById(ctx, organizationId, title, description)
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
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	deleteOrgPerm := fmt.Sprintf("org%s:delete", organizationId)
	canDeleteOrg, err := auth.HasPermission(userId, deleteOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canDeleteOrg {
		http.Error(writer, fmt.Sprintf("User does not have permission to delete organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	ctx := request.Context()
	err = handler.controller.DeleteOrganizationById(ctx, organizationId)
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
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
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
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	userId := request.FormValue("user_id")
	if userId == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	signedInUserId := token.RegisteredClaims.Subject
	addUsersPerm := fmt.Sprintf("org%s:add_members", organizationId)
	canAddUsers, err := auth.HasPermission(signedInUserId, addUsersPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canAddUsers {
		http.Error(writer, fmt.Sprintf("User does not have permission to add users to organization with id: %s", organizationId), http.StatusForbidden)
		return
	}
	err = handler.controller.AddMember(userId, organizationId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users in org with id %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *organizationHandler) RemoveMemberFromOrganization(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	memberId := params["memberId"]
	if memberId == "" {
		http.Error(writer, "No Member ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	removeUsersPerm := fmt.Sprintf("org%s:remove_members", organizationId)
	canRemoveUsers, err := auth.HasPermission(userId, removeUsersPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canRemoveUsers {
		http.Error(writer, fmt.Sprintf("User does not have permission to remove users from organization with id: %s", organizationId), http.StatusForbidden)
		return
	}

	err = handler.controller.RemoveMember(memberId, organizationId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users in org with id %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

// TODO: below
func (handler *organizationHandler) GetOrganizationRoles(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get roles by organization ID
	// TODO: Return roles with status code 200
}

func (handler *organizationHandler) AddOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request can add roles to organization (posisbly with middleware)
	// TODO: Add role to organization
	// TODO: Return role with status code 201
}

func (handler *organizationHandler) GetOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get role by organization ID and role ID
	// TODO: Return role with status code 200
}

func (handler *organizationHandler) UpdateOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request can update roles in organization (posisbly with middleware)
	// TODO: Find what fields are being updated
	// TODO: Get role by organization ID and role ID
	// TODO: update role in database
	// TODO: send back 204
}

func (handler *organizationHandler) DeleteOrganizationRole(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request can delete roles in organization (posisbly with middleware)
	// TODO: Delete role from database
	// TODO: send back 204
}

func (handler *organizationHandler) GetOrganizationRolePrivileges(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// roleId, err := strconv.Atoi(params["roleId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request is apart of organization (posisbly with middleware)
	// TODO: Get privileges by organization ID and role ID
	// TODO: Return role privileges with status code 200
}

func (handler *organizationHandler) AddOrganizationRolePrivilege(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// roleId, err := strconv.Atoi(params["roleId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request can add privileges to role in organization (posisbly with middleware)
	// TODO: Add privilege to role
	// TODO: Return privilege with status code 201
	// TODO: send back 204
}

func (handler *organizationHandler) RemoveOrganizationRolePrivilege(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// roleId, err := strconv.Atoi(params["roleId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request can remove privileges from role in organization (posisbly with middleware)
	// TODO: Remove privilege from role
	// TODO: send back 204
}

func (handler *organizationHandler) AddMemberToRole(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// memberId, err := strconv.Atoi(params["memberId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Member ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// roleId, err := strconv.Atoi(params["roleId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request can add members to role in organization (posisbly with middleware)
	// TODO: Add member to role
	// TODO: Return with status code 204
}

func (handler *organizationHandler) RemoveMemberFromRole(writer http.ResponseWriter, request *http.Request) {
	// params := mux.Vars(request)
	// organizationId, err := strconv.Atoi(params["organizationId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Organization ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// memberId, err := strconv.Atoi(params["memberId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Member ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// roleId, err := strconv.Atoi(params["roleId"])
	// if err != nil {
	// 	http.Error(writer, fmt.Sprintf("Invalid Role ID: %s", err.Error()), http.StatusBadRequest)
	// 	return
	// }
	// TODO: verify user sending request can remove members from role in organization (posisbly with middleware)
	// TODO: Remove member from role
	// TODO: Return with status code 204
}
