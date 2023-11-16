package organization

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

type Organization struct {
	Id          uuid.UUID `db:"id" json:"id"`
	OwnerId     string    `db:"owner_id" json:"owner_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	AiEnabled   bool      `db:"ai_enabled" json:"ai_enabled"`
}
