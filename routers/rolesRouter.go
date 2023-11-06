package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/controllers/user"
	"github.com/Sync-Space-49/syncspace-server/db"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gorilla/mux"
)

type roleHandler struct {
	router *mux.Router
}

func registerRoleRoutes(parentRouter *mux.Router, cfg *config.Config, db *db.DB) *mux.Router {
	handler := &roleHandler{
		router: parentRouter.NewRoute().Subrouter(),
	}

	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetOrganizationRoles))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/permissions", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetOrganizationPermissions))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.CreateRole))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetRole))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.UpdateRole))).Methods("PUT")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/{roleId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.DeleteRole))).Methods("DELETE")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/permissions", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetRolePermissions))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/members", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.GetMembersWithRole))).Methods("GET")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/members/{memberId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.AddMemberToRole))).Methods("POST")
	handler.router.Handle(fmt.Sprintf("%s/{organizationId}/roles/{roleId}/members/{memberId}", organizationsPrefix), auth.EnsureValidToken()(http.HandlerFunc(handler.RemoveMemberFromRole))).Methods("DELETE")

	return handler.router
}

func (handler *roleHandler) GetOrganizationRoles(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	orgRolePrefix := fmt.Sprintf("org%s", organizationId)
	readOrgPerm := fmt.Sprintf("%s:read", orgRolePrefix)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	roles, err := auth.GetRoles(&orgRolePrefix)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get roles for organization %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(roles)
}

func (handler *roleHandler) GetOrganizationPermissions(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	orgRolePrefix := fmt.Sprintf("org%s", organizationId)
	readOrgPerm := fmt.Sprintf("%s:read", orgRolePrefix)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	roles, err := auth.GetRoles(&orgRolePrefix)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get roles for organization %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	permissions := make([]auth.Permission, 0)
	for _, role := range *roles {
		rolePermissions, err := auth.GetRolePermissions(role.Id)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to get permissions for role %s: %s", role.Id, err.Error()), http.StatusInternalServerError)
			return
		}
		permissions = append(permissions, *rolePermissions...)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(permissions)
}

func (handler *roleHandler) CreateRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	orgPrefix := fmt.Sprintf("org%s", organizationId)
	creatRolesPerm := fmt.Sprintf("%s:create_roles", orgPrefix)
	canCreateRoles, err := auth.HasPermission(userId, creatRolesPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user permissions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if !canCreateRoles {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to add roles to organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	roleName := request.FormValue("name")
	if roleName == "" {
		http.Error(writer, "No name Found", http.StatusBadRequest)
		return
	}
	roleDescription := request.FormValue("description")
	if roleDescription == "" {
		http.Error(writer, "No description Found", http.StatusBadRequest)
		return
	}
	permissionNames := request.Form["permission_names"]
	if len(permissionNames) == 0 {
		http.Error(writer, "No Permission IDs Found", http.StatusBadRequest)
		return
	}

	roleName = fmt.Sprintf("%s:%s", orgPrefix, roleName)
	auth.CreateRole(roleName, roleDescription)
	role, err := auth.GetRoles(&roleName)
	roleId := (*role)[0].Id
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get role with query string %s: %s", roleName, err.Error()), http.StatusInternalServerError)
		return
	}

	addMemberToSpecificRolePermission := auth.Permission{
		Name:        fmt.Sprintf("org%s:role%s:add_member", organizationId, roleId),
		Description: fmt.Sprintf("Permission to add users to role with id %s to organization with id %s", roleId, organizationId),
	}
	removeMemberFromSpecificRolePermission := auth.Permission{
		Name:        fmt.Sprintf("org%s:role%s:remove_member", organizationId, roleId),
		Description: fmt.Sprintf("Permission to remove users from role with id %s to organization with id %s", roleId, organizationId),
	}
	permissions := []auth.Permission{addMemberToSpecificRolePermission, removeMemberFromSpecificRolePermission}
	err = auth.CreatePermissions(permissions)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to create permissions for role %s: %s", roleId, err.Error()), http.StatusInternalServerError)
		return
	}

	permissionNames = append(permissionNames, addMemberToSpecificRolePermission.Name, removeMemberFromSpecificRolePermission.Name)
	for _, permissionName := range permissionNames {
		if !strings.Contains(permissionName, orgPrefix) {
			http.Error(writer, fmt.Sprintf("Permission %s does not have org prefix %s (maybe meant for another org?)", permissionName, orgPrefix), http.StatusBadRequest)
			return
		}
	}
	err = auth.AddPermissionsToRole(roleId, permissionNames)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to add permissions to role %s: %s", roleName, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(role)
}

func (handler *roleHandler) GetRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	roleId := params["roleId"]
	if roleId == "" {
		http.Error(writer, "No Role ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	role, err := auth.GetRoleById(roleId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get role for organization %s: %s", organizationId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(role)
}

func (handler *roleHandler) UpdateRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	roleId := params["roleId"]
	if roleId == "" {
		http.Error(writer, "No Role ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	orgPrefix := fmt.Sprintf("org%s", organizationId)
	editRolesPerm := fmt.Sprintf("%s:edit_roles", orgPrefix)
	canEditRoles, err := auth.HasPermission(userId, editRolesPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canEditRoles {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to edit roles to organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	roleName := request.FormValue("name")
	if roleName == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}
	roleDescription := request.FormValue("description")
	if roleDescription == "" {
		http.Error(writer, "No User ID Found", http.StatusBadRequest)
		return
	}

	roleName = fmt.Sprintf("%s:%s", orgPrefix, roleName)
	auth.UpdateRole(roleId, roleName, roleDescription)
	currentRolePermissions, err := auth.GetRolePermissions(roleId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get role permissions for role %s: %s", roleId, err.Error()), http.StatusInternalServerError)
		return
	}

	permissionNames := request.Form["permission_names"]
	for _, permissionName := range permissionNames {
		if !strings.Contains(permissionName, orgPrefix) {
			http.Error(writer, fmt.Sprintf("Permission %s does not have org prefix %s (maybe meant for another org?)", permissionName, orgPrefix), http.StatusBadRequest)
			return
		}
	}
	addPermissionNames := make([]string, 0)
	for _, newPermissionName := range permissionNames {
		isNewPerm := false
		for _, currentPermissionName := range *currentRolePermissions {
			if newPermissionName != currentPermissionName.Name {
				isNewPerm = true
				break
			}
		}
		if isNewPerm {
			addPermissionNames = append(addPermissionNames, newPermissionName)
		}
	}
	err = auth.AddPermissionsToRole(roleId, addPermissionNames)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to add permissions %v to role %s: %s", addPermissionNames, roleId, err.Error()), http.StatusInternalServerError)
		return
	}

	deletePermissionNames := make([]string, 0)
	addMemberToRoleName := fmt.Sprintf("org%s:role%s:add_member", organizationId, roleId)
	removeMemberFromRoleName := fmt.Sprintf("org%s:role%s:remove_member", organizationId, roleId)
	for _, currentPermissionName := range *currentRolePermissions {
		isDeletedPerm := true
		for _, newPermissionName := range permissionNames {
			if currentPermissionName.Name == newPermissionName && currentPermissionName.Name != addMemberToRoleName && currentPermissionName.Name != removeMemberFromRoleName {
				isDeletedPerm = false
				break
			}
		}
		if isDeletedPerm {
			deletePermissionNames = append(deletePermissionNames, currentPermissionName.Name)
		}
	}
	err = auth.RemovePermissionsFromRole(roleId, deletePermissionNames)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to remove permissions %v from role %s: %s", deletePermissionNames, roleId, err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *roleHandler) DeleteRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	roleId := params["roleId"]
	if roleId == "" {
		http.Error(writer, "No Role ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	deleteRolesPerm := fmt.Sprintf("org%s:delete_roles", organizationId)
	canDeleteRoles, err := auth.HasPermission(userId, deleteRolesPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canDeleteRoles {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to edit roles to organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}
	err = auth.DeleteRole(roleId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to delete role %s: %s", roleId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *roleHandler) GetRolePermissions(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	roleId := params["roleId"]
	if roleId == "" {
		http.Error(writer, "No Role ID Found", http.StatusBadRequest)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read roles to organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	permissions, err := auth.GetRolePermissions(roleId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get permissions for role %s: %s", roleId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(permissions)
}

func (handler *roleHandler) GetMembersWithRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	roleId := params["roleId"]

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	readOrgPerm := fmt.Sprintf("org%s:read", organizationId)
	canReadOrg, err := auth.HasPermission(userId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User with id %s does not have permission to read roles to organization with id: %s", userId, organizationId), http.StatusForbidden)
		return
	}

	members, err := user.GetUsersWithRole(roleId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get users with role %s: %s", roleId, err.Error()), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(members)
}

func (handler *roleHandler) AddMemberToRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	roleId := params["roleId"]
	if roleId == "" {
		http.Error(writer, "No Role ID Found", http.StatusBadRequest)
		return
	}
	memberId := params["memberId"]
	if memberId == "" {
		http.Error(writer, "No Member ID Found", http.StatusBadRequest)
		return
	}

	orgPrefix := fmt.Sprintf("org%s", organizationId)
	readOrgPerm := fmt.Sprintf("%s:read", orgPrefix)
	canReadOrg, err := auth.HasPermission(memberId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", memberId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User you're trying to give role to (%s) does not have permission to read organization with id: %s", memberId, organizationId), http.StatusForbidden)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	addMemberToSpecificRole := fmt.Sprintf("org%s:role%s:add_member", organizationId, roleId)
	canAddMemberToSpecificRole, err := auth.HasPermission(userId, addMemberToSpecificRole)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canAddMemberToSpecificRole {
		addToRolesPerm := fmt.Sprintf("%s:add_roles", orgPrefix)
		canAddToRoles, err := auth.HasPermission(userId, addToRolesPerm)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
			return
		}
		if !canAddToRoles {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to add roles to organization with id: %s", userId, organizationId), http.StatusForbidden)
			return
		}
	}

	err = auth.AddUserToRole(memberId, roleId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to add member %s to role %s: %s", memberId, roleId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}

func (handler *roleHandler) RemoveMemberFromRole(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	organizationId := params["organizationId"]
	if organizationId == "" {
		http.Error(writer, "No Organization ID Found", http.StatusBadRequest)
		return
	}
	roleId := params["roleId"]
	if roleId == "" {
		http.Error(writer, "No Role ID Found", http.StatusBadRequest)
		return
	}
	memberId := params["memberId"]
	if memberId == "" {
		http.Error(writer, "No Member ID Found", http.StatusBadRequest)
		return
	}

	orgPrefix := fmt.Sprintf("org%s", organizationId)
	readOrgPerm := fmt.Sprintf("%s:read", orgPrefix)
	canReadOrg, err := auth.HasPermission(memberId, readOrgPerm)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", memberId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canReadOrg {
		http.Error(writer, fmt.Sprintf("User you're trying to remove role from (%s) does not have permission to read organization with id: %s", memberId, organizationId), http.StatusForbidden)
		return
	}

	token := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	userId := token.RegisteredClaims.Subject
	removeMemberFromSpecificRole := fmt.Sprintf("org%s:role%s:remove_member", organizationId, roleId)
	canRemoveMemberFromSpecificRole, err := auth.HasPermission(userId, removeMemberFromSpecificRole)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
		return
	}
	if !canRemoveMemberFromSpecificRole {
		removeFromRolesPerm := fmt.Sprintf("%s:add_roles", orgPrefix)
		canRemoveFromRoles, err := auth.HasPermission(userId, removeFromRolesPerm)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to get user with id %s permissions: %s", userId, err.Error()), http.StatusInternalServerError)
			return
		}
		if !canRemoveFromRoles {
			http.Error(writer, fmt.Sprintf("User with id %s does not have permission to add roles to organization with id: %s", userId, organizationId), http.StatusForbidden)
			return
		}
	}

	err = auth.RemoveUserFromRole(memberId, roleId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to remove member %s from role %s: %s", memberId, roleId, err.Error()), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusNoContent)
}
