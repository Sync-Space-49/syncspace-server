package auth

import (
	"fmt"
	"net/http"
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

func AuthenticateLoginToken(writer http.ResponseWriter, request *http.Request, accessToken string) *UserClaims {
	cfg, err := config.Get()
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get config: %s", err.Error()), http.StatusBadRequest)
		return nil
	}
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	return parsedAccessToken.Claims.(*UserClaims)
}
