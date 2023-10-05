package user

import (
	"context"
	"fmt"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"

	"golang.org/x/crypto/bcrypt"
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

func (c *Controller) GetUserById(userId int) (User, error) {
	var user User
	return user, nil
}

func (c *Controller) CreateUser(username string, email string, password string, pfpUrl *string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	var query string
	if pfpUrl == nil {
		query = `
			INSERT INTO users (username, email, password)
			VALUES (:username, :email, :hashed_password)
		`
	} else {
		query = `
			INSERT INTO users (username, email, password, pfp_url)
			VALUES (:username, :email, :hashed_password, :pfp_url)
		`
	}
	_, err = c.db.DB.NamedExec(query, map[string]interface{}{
		"username":        username,
		"email":           email,
		"hashed_password": hashedPassword,
		"pfp_url":         pfpUrl,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (c *Controller) GetUserByCredentials(ctx context.Context, credential string, password string) (*User, error) {
	// find user with email or username
	row := c.db.DB.QueryRowxContext(ctx, `
		SELECT * FROM users
		WHERE email = $1 OR username = $1
	`, credential)

	var user User
	// check if pfp is null
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.ProfilePicURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("incorrect password: %w", err)
	}
	return &user, nil
}
