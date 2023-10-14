package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
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

func (c *Controller) GetUserById(userId string, tokenString string) (*User, error) {
	managementToken, err := auth.GetManagementToken(c.cfg, tokenString)
	if err != nil {
		return &User{}, err
	}
	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &User{}, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return &User{}, err
	}

	return &user, nil
}

func (c *Controller) UpdateUserById(tokenString string, userId string, email string, username string, password string, pfpUrl string) error {
	managementToken, err := auth.GetManagementToken(c.cfg, tokenString)
	if err != nil {
		return fmt.Errorf("failed to get maintenance token: %w", err)
	}
	userInfo, err := c.GetUserById(userId, tokenString)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	if username == "" {
		username = userInfo.Username
	}
	if pfpUrl == "" {
		pfpUrl = userInfo.Picture
	}

	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)
	method := "PATCH"

	var payload io.Reader
	if password != "" {
		payload = strings.NewReader(fmt.Sprintf(`{"email":"%s","picture":"%s","password":"%s"}`, email, pfpUrl, password))
	} else {
		payload = strings.NewReader(fmt.Sprintf(`{"username":"%s","picture":"%s"}`, username, pfpUrl))
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", managementToken))
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to update user: %s", string(body))
	}

	if email != "" {
		payload = strings.NewReader(fmt.Sprintf(`{"email":"%s"}`, email))
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", managementToken))
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("failed to update user: %s", string(body))
		}
	}

	return nil
}
