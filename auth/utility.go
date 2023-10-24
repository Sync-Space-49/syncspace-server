package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Sync-Space-49/syncspace-server/cache"
	"github.com/Sync-Space-49/syncspace-server/config"
)

func GetManagementToken() (string, error) {
	var managementToken string
	tokenCache := cache.Get()
	token, err := tokenCache.Read(cache.ManagementTokenKey)
	if err != nil {
		return "", err
	}
	if token != nil {
		managementToken = *token
	} else {
		cfg, err := config.Get()
		if err != nil {
			return "", err
		}
		url := fmt.Sprintf("%soauth/token", cfg.Auth0.Domain)
		payload := strings.NewReader(fmt.Sprintf(`{"client_id":"%s","client_secret":"%s","audience":"%s","grant_type":"client_credentials"}`, cfg.Auth0.Server.ClientId, cfg.Auth0.Server.ClientSecret, cfg.Auth0.Management.Audience))
		method := "POST"
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			return "", err
		}
		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		if res.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to get maintenance token: %s", string(body))
		}

		var tokenResponse struct {
			Token string `json:"access_token"`
		}
		err = json.Unmarshal(body, &tokenResponse)
		if err != nil {
			return "", err
		}

		tokenCache.Update(cache.ManagementTokenKey, tokenResponse.Token)
		managementToken = tokenResponse.Token
	}
	return managementToken, nil
}

func GetUserPermissions(userId string) (*[]Permission, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return nil, err
	}
	method := "GET"
	url := fmt.Sprintf("%sapi/v2/users/%s/permissions", cfg.Auth0.Domain, userId)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get permissions: %s", string(body))
	}

	var userPermissions []Permission
	err = json.Unmarshal(body, &userPermissions)
	if err != nil {
		return nil, err
	}
	return &userPermissions, nil
}

func CreatePermission(permissionName string, permissionDescription string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%sapi/v2/resource-servers/%s", cfg.Auth0.Domain, cfg.Auth0.Server.Id)
	payload := strings.NewReader(fmt.Sprintf(`{ "scopes": [ { "value": "%s", "description": "%s" } ] }`, permissionName, permissionDescription))
	method := "PATCH"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf(`failed to create new role: %s`, string(body))
	}

	return nil
}

func CreatePermissions(permissions [][]string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%sapi/v2/resource-servers/%s", cfg.Auth0.Domain, cfg.Auth0.Server.Id)
	formattedPermissions := `{ "scopes": [ `
	for i, permission := range permissions {
		if i == 0 {
			formattedPermissions += fmt.Sprintf(`{ "value": "%s", "description": "%s" }`, permission[0], permission[1])
		} else {
			formattedPermissions += fmt.Sprintf(`, { "value": "%s", "description": "%s" }`, permission[0], permission[1])
		}
	}
	formattedPermissions += ` ] }`
	payload := strings.NewReader(formattedPermissions)
	method := "PATCH"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf(`failed to create new role: %s`, string(body))
	}

	return nil
}

func GetRoles(filter *string) (*[]Role, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return nil, err
	}
	method := "GET"
	url := fmt.Sprintf("%sapi/v2/roles", cfg.Auth0.Domain)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	if filter != nil {
		q := req.URL.Query()
		q.Add("name_filter", *filter)
		req.URL.RawQuery = q.Encode()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get roles: %s", string(body))
	}

	var roles []Role
	err = json.Unmarshal(body, &roles)
	if err != nil {
		return nil, err
	}
	return &roles, nil
}

func CreateRole(roleName string, roleDescription string) (*Role, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%sapi/v2/roles", cfg.Auth0.Domain)
	payload := strings.NewReader(fmt.Sprintf(`{ "name": "%s", "description": "%s" }`, roleName, roleDescription))
	method := "POST"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`failed to create new role: %s`, string(body))
	}

	var newRole Role
	err = json.Unmarshal(body, &newRole)
	if err != nil {
		return nil, err
	}
	return &newRole, nil
}

func AddPermissionToRole(roleId string, permissionName string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%sapi/v2/roles/%s/permissions", cfg.Auth0.Domain, roleId)
	payload := strings.NewReader(fmt.Sprintf(`{ "permissions": [ { "resource_server_identifier": "%s", "permission_name": "%s" } ] }`, cfg.Auth0.Server.Audience, permissionName))
	method := "POST"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf(`failed to assign permission "%s" to role with id "%s": %s`, permissionName, roleId, string(body))
	}

	return nil
}

func AddPermissionsToRole(roleId string, permissionNames []string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%sapi/v2/roles/%s/permissions", cfg.Auth0.Domain, roleId)
	formattedPermissions := `{ "permissions": [ `
	for i, permissionName := range permissionNames {
		if i == 0 {
			formattedPermissions += fmt.Sprintf(`{ "resource_server_identifier": "%s", "permission_name": "%s" }`, cfg.Auth0.Server.Audience, permissionName)
		} else {
			formattedPermissions += fmt.Sprintf(`, { "resource_server_identifier": "%s", "permission_name": "%s" }`, cfg.Auth0.Server.Audience, permissionName)
		}
	}
	formattedPermissions += ` ] }`
	payload := strings.NewReader(formattedPermissions)
	method := "POST"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf(`failed to assign permissions to role with id "%s": %s`, roleId, string(body))
	}

	return nil
}

func GetUserRoles(userId string) (*[]Role, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return nil, err
	}
	method := "GET"
	url := fmt.Sprintf("%sapi/v2/users/%s/roles", cfg.Auth0.Domain, userId)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get roles: %s", string(body))
	}

	var userRoles []Role
	err = json.Unmarshal(body, &userRoles)
	if err != nil {
		return nil, err
	}
	return &userRoles, nil
}

func AddUserToRole(userId string, roleId string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%sapi/v2/users/%s/roles", cfg.Auth0.Domain, userId)
	payload := strings.NewReader(fmt.Sprintf(`{ "roles": [ "%s" ] }`, roleId))
	method := "POST"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf(`failed to add role with id "%s" to user with id "%s": %s`, roleId, userId, string(body))
	}

	return nil
}

func RemoveUserFromRole(userId string, roleId string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%sapi/v2/users/%s/roles", cfg.Auth0.Domain, userId)
	payload := strings.NewReader(fmt.Sprintf(`{ "roles": [ "%s" ] }`, roleId))
	method := "DELETE"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf(`failed to remove role with id "%s" from user with id "%s": %s`, roleId, userId, string(body))
	}

	return nil
}

func RemoveUserFromRoles(userId string, roleIds []string) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}
	managementToken, err := GetManagementToken()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%sapi/v2/users/%s/roles", cfg.Auth0.Domain, userId)
	payload := strings.NewReader(fmt.Sprintf(`{ "roles": [ "%s" ] }`, strings.Join(roleIds, `", "`)))
	method := "DELETE"
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf(`failed to remove roles from user with id "%s": %s`, userId, string(body))
	}

	return nil
}
