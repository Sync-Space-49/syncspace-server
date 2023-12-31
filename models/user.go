package models

import "github.com/google/uuid"

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

type FavoriteBoards struct {
	ID      uuid.UUID `db:"id" json:"id"`
	UserID  string    `db:"user_id" json:"user_id"`
	BoardID uuid.UUID `db:"board_id" json:"board_id"`
}
