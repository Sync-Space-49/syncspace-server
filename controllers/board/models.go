package board

import "github.com/google/uuid"

type Board struct {
	Id             uuid.UUID `db:"id" json:"id"`
	OwnerId        string    `db:"owner_id" json:"owner_id"`
	Title          string    `db:"title" json:"title"`
	CreatedAt      string    `db:"created_at" json:"created_at"`
	ModifiedAt     string    `db:"modified_at" json:"modified_at"`
	IsPrivate      bool      `db:"is_private" json:"is_private"`
	OrganizationId int       `db:"organization_id" json:"organization_id"`
}
