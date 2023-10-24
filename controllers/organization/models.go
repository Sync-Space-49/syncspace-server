package organization

import "github.com/google/uuid"

type Organization struct {
	Id          uuid.UUID `db:"id" json:"id"`
	OwnerId     string    `db:"owner_id" json:"owner_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
}
