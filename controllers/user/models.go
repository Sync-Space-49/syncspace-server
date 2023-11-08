package user

import (
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

type Identity struct {
	Connection string `json:"connection"`
	UserID     string `json:"user_id"`
	Provider   string `json:"provider"`
	IsSocial   bool   `json:"isSocial"`
}

type User struct {
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	Username      string     `json:"username"`
	PhoneNumber   string     `json:"phone_number"`
	PhoneVerified bool       `json:"phone_verified"`
	UserID        string     `json:"user_id"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
	Identities    []Identity `json:"identities"`
	AppMetadata   struct{}   `json:"app_metadata"`
	UserMetadata  struct{}   `json:"user_metadata"`
	Picture       string     `json:"picture"`
	Name          string     `json:"name"`
	Nickname      string     `json:"nickname"`
	Multifactor   []string   `json:"multifactor"`
	LastIP        string     `json:"last_ip"`
	LastLogin     string     `json:"last_login"`
	LoginsCount   int        `json:"logins_count"`
	Blocked       bool       `json:"blocked"`
	GivenName     string     `json:"given_name"`
	FamilyName    string     `json:"family_name"`
}
