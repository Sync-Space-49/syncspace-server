package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Sync-Space-49/syncspace-server/config"
)

func GetManagementToken(cfg *config.Config) (string, error) {
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
	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to get maintenance token: %s", string(body))
	}

	var managementToken struct {
		Token string `json:"access_token"`
	}
	err = json.Unmarshal(body, &managementToken)
	if err != nil {
		return "", err
	}

	return managementToken.Token, nil
}
