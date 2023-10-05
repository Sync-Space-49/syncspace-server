package user

import "database/sql"

type User struct {
	Id             int            `db:"id" json:"id"`
	Username       string         `db:"username" json:"username"`
	Email          string         `db:"email" json:"email"`
	HashedPassword string         `db:"password" json:"password"`
	ProfilePicURL  sql.NullString `db:"pfp_url" json:"pfp_url"`
}
