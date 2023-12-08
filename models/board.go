package models

import (
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
	"github.com/google/uuid"
)

type BoardController struct {
	cfg *config.Config
	db  *db.DB
}

func NewBoardController(cfg *config.Config, db *db.DB) *BoardController {
	return &BoardController{
		cfg: cfg,
		db:  db,
	}
}

type Card struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Points      string    `db:"points" json:"points"`
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
	Description    string    `db:"description" json:"description"`
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
	Id          uuid.UUID       `db:"id" json:"id"`
	OwnerId     string          `db:"owner_id" json:"owner_id"`
	Title       string          `db:"title" json:"title"`
	Description string          `db:"description" json:"description"`
	CreatedAt   string          `db:"created_at" json:"created_at"`
	ModifiedAt  string          `db:"modified_at" json:"modified_at"`
	IsPrivate   bool            `db:"is_private" json:"is_private"`
	Panels      []CompletePanel `json:"panels"`
}

type CompleteCard struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Points      string    `db:"points" json:"points"`
	Position    int       `db:"position" json:"position"`
	StackId     uuid.UUID `db:"stack_id" json:"stack_id"`
	Assignments []User    `json:"assignments"`
}

type AIGeneratedCard struct {
	CardTitle       string      `json:"title"`
	CardDesc        string      `json:"description"`
	CardStoryPoints interface{} `json:"story_points"`
}

type AIGeneratedSprint map[string][]AIGeneratedCard

type SimplifiedCompleteBoard struct {
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Panels      []SimplifiedCompletePanel `json:"panels"`
}

type SimplifiedCompletePanel struct {
	Title  string                    `json:"title"`
	Stacks []SimplifiedCompleteStack `json:"stacks"`
}

type SimplifiedCompleteStack struct {
	Id    uuid.UUID         `db:"id" json:"id"`
	Title string            `json:"title"`
	Cards []AIGeneratedCard `json:"cards"`
}

type DetailedAssignedCard struct {
	UserID      string    `db:"user_id" json:"user_id"`
	ID          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Position    int       `db:"position" json:"position"`
	StackID     uuid.UUID `db:"stack_id" json:"stack_id"`
	Points      string    `db:"points" json:"points"`
	PanelID     uuid.UUID `db:"panel_id" json:"panel_id"`
	BoardID     uuid.UUID `db:"board_id" json:"board_id"`
	OrgID       uuid.UUID `db:"org_id" json:"org_id"`
}
