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
	userPermissions := c.Permissions
	for _, permission := range userPermissions {
		if permission == permissionName {
			return true
		}
	}
	return false
}
