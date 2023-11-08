package board

import (
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
	"github.com/google/uuid"
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

type Board struct {
	Id             uuid.UUID `db:"id" json:"id"`
	OwnerId        string    `db:"owner_id" json:"owner_id"`
	Title          string    `db:"title" json:"title"`
	CreatedAt      string    `db:"created_at" json:"created_at"`
	ModifiedAt     string    `db:"modified_at" json:"modified_at"`
	IsPrivate      bool      `db:"is_private" json:"is_private"`
	OrganizationId uuid.UUID `db:"organization_id" json:"organization_id"`
}

type Panel struct {
	Id        uuid.UUID `db:"id" json:"id"`
	title     string    `db:"title" json:"title"`
	postition int       `db:"position" json:"position"`
	boardId   uuid.UUID `db:"board_id" json:"board_id"`
}

type Stack struct {
	Id        uuid.UUID `db:"id" json:"id"`
	title     string    `db:"title" json:"title"`
	postition int       `db:"position" json:"position"`
	panelId   uuid.UUID `db:"board_id" json:"panel_id"`
}

type Card struct {
	Id          uuid.UUID `db:"id" json:"id"`
	title       string    `db:"title" json:"title"`
	description string    `db:"description" json:"description"`
	postition   int       `db:"position" json:"position"`
	stackId     uuid.UUID `db:"stack_id" json:"stack_id"`
}
