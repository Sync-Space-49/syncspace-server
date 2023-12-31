package models

import (
	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/Sync-Space-49/syncspace-server/db"
	"github.com/google/uuid"
)

type OrganizationController struct {
	cfg *config.Config
	db  *db.DB
}

func NewOrganizationController(cfg *config.Config, db *db.DB) *OrganizationController {
	return &OrganizationController{
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
