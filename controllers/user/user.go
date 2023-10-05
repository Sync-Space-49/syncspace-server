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

func (c *Controller) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (c *Controller) ComparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("failed to compare password: %w", err)
	}
	return nil
}

func (c *Controller) GetUserById(ctx context.Context, userId int) (*User, error) {
	var user User
	row := c.db.DB.QueryRowxContext(ctx, `
		SELECT * FROM users
		WHERE id = $1
	`, userId)

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

	return &user, nil
}

func (c *Controller) CreateUser(ctx context.Context, username string, email string, password string, pfpUrl *string) error {
	hashedPassword, err := c.HashPassword(password)
	if err != nil {
		return err
	}

	if pfpUrl == nil {
		_, err = c.db.DB.ExecContext(ctx, `
			INSERT INTO users (username, email, password)
			VALUES ($1, $2, $3)
		`, username, email, hashedPassword)
	} else {
		_, err = c.db.DB.ExecContext(ctx, `
			INSERT INTO users (username, email, password, pfp_url)
			VALUES ($1, $2, $3, $4)
		`, username, email, hashedPassword, pfpUrl)
	}

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

	err = c.ComparePassword(user.HashedPassword, password)
	if err != nil {
		return nil, fmt.Errorf("incorrect password: %w", err)
	}
	return &user, nil
}

func (c *Controller) UpdateUser(ctx context.Context, userId int, email string, username string, password string, pfpUrl *string) error {
	query := `
		UPDATE users
		SET email = $1, username = $2, password = $3, pfp_url = $4
		WHERE id = $5
	`
	_, err := c.db.DB.ExecContext(ctx, query, email, username, password, pfpUrl, userId)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
