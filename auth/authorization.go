package auth

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func NewOrganizationAuthToken(claims OranizationMemberClaims) (string, error) {
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return authToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseOrganizationAuthToken(authToken string) *OranizationMemberClaims {
	parsedAuthToken, _ := jwt.ParseWithClaims(authToken, &OranizationMemberClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAuthToken.Claims.(*OranizationMemberClaims)
}

func NewBoardAuthToken(claims BoardMemberClaims) (string, error) {
	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return authToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseBoardAuthToken(authToken string) *BoardMemberClaims {
	parsedAuthToken, _ := jwt.ParseWithClaims(authToken, &BoardMemberClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAuthToken.Claims.(*BoardMemberClaims)
}
