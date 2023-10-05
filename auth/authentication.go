package auth

import (
	"fmt"
	"time"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/controllers/user"
	"github.com/golang-jwt/jwt"
)

func CreateLoginToken(user user.User) (*string, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	claims := UserClaims{
		Id:            user.Id,
		Username:      user.Username,
		ProfilePicURL: user.ProfilePicURL,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	return &signedString, err
}

func AuthenticateLoginToken(accessToken string) (*UserClaims, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse access token: %w", err)
	}

	return parsedAccessToken.Claims.(*UserClaims), nil
}
