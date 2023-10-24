package auth

import (
	"fmt"
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

func HasPermission(userId string, permissionName string) (bool, error) {
	userPermissions, err := GetUserPermissions(userId)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %v", err)
	}
	for _, permission := range *userPermissions {
		if permission.Name == permissionName {
			return true, nil
		}
	}
	return false, nil
}
