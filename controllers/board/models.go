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

type Card struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Position    int       `db:"position" json:"position"`
	StackId     uuid.UUID `db:"stack_id" json:"stack_id"`
}

type Stack struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Title    string    `db:"title" json:"title"`
	Position int       `db:"position" json:"position"`
	PanelId  uuid.UUID `db:"panel_id" json:"panel_id"`
}
type Panel struct {
	Id       uuid.UUID `db:"id" json:"id"`
	Title    string    `db:"title" json:"title"`
	Position int       `db:"position" json:"position"`
	BoardId  uuid.UUID `db:"board_id" json:"board_id"`
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

type CompleteStack struct {
	Id       uuid.UUID      `db:"id" json:"id"`
	Title    string         `db:"title" json:"title"`
	Position int            `db:"position" json:"position"`
	PanelId  uuid.UUID      `db:"panel_id" json:"panel_id"`
	Cards    []CompleteCard `json:"cards"`
}

type CompletePanel struct {
	Id       uuid.UUID       `db:"id" json:"id"`
	Title    string          `db:"title" json:"title"`
	Position int             `db:"position" json:"position"`
	BoardId  uuid.UUID       `db:"board_id" json:"board_id"`
	Stacks   []CompleteStack `json:"stacks"`
}

type CompleteBoard struct {
	Id         uuid.UUID       `db:"id" json:"id"`
	OwnerId    string          `db:"owner_id" json:"owner_id"`
	Title      string          `db:"title" json:"title"`
	CreatedAt  string          `db:"created_at" json:"created_at"`
	ModifiedAt string          `db:"modified_at" json:"modified_at"`
	IsPrivate  bool            `db:"is_private" json:"is_private"`
	Panels     []CompletePanel `json:"panels"`
}

type CompleteCard struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Position    int       `db:"position" json:"position"`
	StackId     uuid.UUID `db:"stack_id" json:"stack_id"`
	Assignments []string  `json:"assignments"`
}
