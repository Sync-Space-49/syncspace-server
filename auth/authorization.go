package auth

import (
	"strings"
)

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}
	return false
}

func (c CustomClaims) HasPermission(permissionName string) bool {
	for _, permission := range c.Permissions {
		if permission == permissionName {
			return true
		}
	}
	return false
}

func (c CustomClaims) HasAnyPermissions(permissionNames ...string) bool {
	for _, permissionName := range permissionNames {
		if c.HasPermission(permissionName) {
			return true
		}
	}
	return false
}

func (c CustomClaims) HasAllPermissions(permissionNames ...string) bool {
	foundPermissions := make([]string, 0)
	for _, permissionName := range permissionNames {
		if c.HasPermission(permissionName) {
			foundPermissions = append(foundPermissions, permissionName)
		}
	}
	return len(foundPermissions) == len(permissionNames)
}
