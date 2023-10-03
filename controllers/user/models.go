package user

type User struct {
	Id             int    `db:"id" json:"id"`
	Username       string `db:"name" json:"name"`
	Email          string `db:"email" json:"email"`
	HashedPassword string `db:"hashed_password" json:"hashed_password"`
	ProfilePicURL  string `db:"pfp_url" json:"pfp_url"`
}
