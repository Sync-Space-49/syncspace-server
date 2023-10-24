package organization

type Organization struct {
	Id          string `db:"id" json:"id"`
	OwnerId     string `db:"owner_id" json:"owner_id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}
