package auth

import (
	"database/sql"

	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Id            int            `json:"id"`
	Username      string         `json:"username"`
	ProfilePicURL sql.NullString `json:"pfp_url"`
	jwt.StandardClaims
}

type OranizationMemberClaims struct {
	Organizations []struct {
		Id           string `json:"id"`
		MemberId     string `json:"member_id"`
		PrivilegeIds []int  `json:"privileges"`
	} `json:"organizations"`
	jwt.StandardClaims
}

type BoardMemberClaims struct {
	BoardId       string `json:"board_id"`
	BoardMemberId string `json:"board_member_id"`
	PrivilegeIds  []int  `json:"privileges"`
	jwt.StandardClaims
}
