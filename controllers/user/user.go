package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/controllers/organization"
	"github.com/Sync-Space-49/syncspace-server/db"
	"github.com/jmoiron/sqlx"
)

type Controller struct {
	cfg *config.Config
	db  *db.DB
}

func NewController(cfg *config.Config, db *db.DB) *Controller {
	return &Controller{
		cfg: cfg,
		db:  db,
	}
}

func (c *Controller) GetUserById(userId string) (*User, error) {
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return &User{}, err
	}
	method := "GET"
	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &User{}, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return &User{}, fmt.Errorf("invalid request: %s", string(body))
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (c *Controller) UpdateUserById(userId string, email string, username string, password string, pfpUrl *string) error {
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return fmt.Errorf("failed to get maintenance token: %w", err)
	}

	user, err := c.GetUserById(userId)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	if username == "" {
		username = user.Username
	}
	if pfpUrl == nil {
		pfpUrl = &user.Picture
	}

	method := "PATCH"
	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)

	var payload io.Reader
	if password != "" {
		payload = strings.NewReader(fmt.Sprintf(`{"email":"%s","picture":"%s","password":"%s"}`, email, *pfpUrl, password))
	} else {
		payload = strings.NewReader(fmt.Sprintf(`{"username":"%s","picture":"%s"}`, username, *pfpUrl))
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update user: %s", string(body))
	}

	// seperate request for email because Auth0 won't let you update email and username at the same time
	if email != "" {
		payload = strings.NewReader(fmt.Sprintf(`{"email":"%s"}`, email))
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to update user: %s", string(body))
		}
	}

	return nil
}

func (c *Controller) DeleteUserById(userId string) error {
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return err
	}
	method := "DELETE"
	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete user: %s", string(body))
	}

	return nil
}

func (c *Controller) GetUserOrganizationsById(ctx context.Context, userId string) (*[]organization.Organization, error) {
	usersRoles, err := auth.GetUserRoles(userId)
	if err != nil {
		return nil, err
	}

	var orgIds []string
	findUUIDInRoleRegex := regexp.MustCompile(`org(.*?):`)
	for _, role := range *usersRoles {
		matches := findUUIDInRoleRegex.FindStringSubmatch(role.Name)
		if len(matches) < 2 {
			continue
		}
		organizationId := matches[1]
		alreadyFound := false
		for _, orgId := range orgIds {
			if orgId == organizationId {
				alreadyFound = true
				break
			}
		}
		if !alreadyFound {
			orgIds = append(orgIds, organizationId)
		}
	}

	if len(orgIds) == 0 {
		return &[]organization.Organization{}, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM Organizations WHERE id IN (?)`, orgIds)
	if err != nil {
		return nil, err
	}
	query = c.db.DB.Rebind(query)
	var organizations []organization.Organization
	err = c.db.DB.SelectContext(ctx, &organizations, query, args...)
	if err != nil {
		return nil, err
	}
	return &organizations, nil
}
