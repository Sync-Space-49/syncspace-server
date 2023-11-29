package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/models"
)

func GetUser(userId string) (*models.User, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance token: %w", err)
	}

	method := "GET"
	url := fmt.Sprintf("%sapi/v2/users/%s", cfg.Auth0.Domain, userId)
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get users: %s", string(body))
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUsersWithRole(roleId string) (*[]models.User, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance token: %w", err)
	}

	method := "GET"
	url := fmt.Sprintf("%sapi/v2/roles/%s/users", cfg.Auth0.Domain, roleId)
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get users: %s", string(body))
	}
	var usersWithRole []struct {
		UserId string `json:"user_id"`
	}
	err = json.Unmarshal(body, &usersWithRole)
	if err != nil {
		return nil, err
	}

	if len(usersWithRole) == 0 {
		return &[]models.User{}, nil
	}

	var users []models.User
	for _, user := range usersWithRole {
		user, err := GetUser(user.UserId)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return &users, nil
}

func GetOrgMembers(organizationId string) (*[]models.User, error) {
	orgMemberRoleName := fmt.Sprintf("org%s:member", organizationId)
	roles, err := auth.GetRoles(&orgMemberRoleName)
	if err != nil {
		return nil, err
	}
	orgMemberRole := (*roles)[0]
	users, err := GetUsersWithRole(orgMemberRole.Id)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetOrgOwners(organizationId string) (*[]models.User, error) {
	orgOwnerRoleName := fmt.Sprintf("org%s:owner", organizationId)
	roles, err := auth.GetRoles(&orgOwnerRoleName)
	if err != nil {
		return nil, err
	}
	orgOwnerRole := (*roles)[0]
	users, err := GetUsersWithRole(orgOwnerRole.Id)
	if err != nil {
		return nil, err
	}
	return users, nil
}
