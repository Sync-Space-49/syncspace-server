package board

import (
	"context"
	"fmt"

	"github.com/Sync-Space-49/syncspace-server/auth"

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

func (c *Controller) GetBoardById(ctx context.Context, boardId string) (*Board, error) {
	var board Board
	err := c.db.DB.GetContext(ctx, &board, `
		SELECT * FROM Boards WHERE id=$1;
	`, boardId)
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (c *Controller) CreateBoard(ctx context.Context, userId string, name string, is_private bool) (*Board, error) {
	var query string
	query = `INSERT INTO Boards (id, name, is_private) VALUES ($1, $2, $3);`
	orgID := uuid.New().String()
	_, err := c.db.DB.ExecContext(ctx, query, orgID, name, is_private)
	if err != nil {
		return nil, err
	}
	org, err := c.GetBoardById(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (c *Controller) InitializeBoard(ownerId string, boardId string) error {
	ownerRoleName := fmt.Sprintf("org%s:owner", boardId)
	ownerRoleDescription := fmt.Sprintf("Owner of organization with the id: %s", boardId)
	ownerRole, err := auth.CreateRole(ownerRoleName, ownerRoleDescription)
	if err != nil {
		return err
	}
	memberRoleName := fmt.Sprintf("org%s:member", boardId)
	memberRoleDescription := fmt.Sprintf("Member of organization with the id: %s", boardId)
	memberRole, err := auth.CreateRole(memberRoleName, memberRoleDescription)
	if err != nil {
		return err
	}

	readPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:read", boardId),
		Description: fmt.Sprintf("Allows you to read the contents of the organization with id %s", boardId),
	}
	deletePerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:delete", boardId),
		Description: fmt.Sprintf("Allows you to delete the organization with id %s", boardId),
	}
	updatePerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:update", boardId),
		Description: fmt.Sprintf("Allows you to update info about the organization with id %s", boardId),
	}
	addMembersPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:add_members", boardId),
		Description: fmt.Sprintf("Allows you to add members to the organization with id %s", boardId),
	}
	removeMembersPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:remove_members", boardId),
		Description: fmt.Sprintf("Allows you to remove members from the organization with id %s", boardId),
	}

	orgMemberPermissions := []auth.Permission{
		readPerm,
	}
	orgOwnerPermissions := append(orgMemberPermissions, deletePerm, updatePerm, addMembersPerm, removeMembersPerm)

	// Because the owner role has all permissions, we only need to call CreatePermissions once
	err = auth.CreatePermissions(orgOwnerPermissions)
	if err != nil {
		return fmt.Errorf("failed to create permissions for organization: %w", err)
	}

	orgMemberPermissionNames := make([]string, len(orgMemberPermissions))
	for i, permission := range orgMemberPermissions {
		orgMemberPermissionNames[i] = permission.Name
	}
	err = auth.AddPermissionsToRole(memberRole.Id, orgMemberPermissionNames)
	if err != nil {
		return fmt.Errorf("failed to add permissions to member role: %w", err)
	}
	orgOwnerPermissionNames := make([]string, len(orgOwnerPermissions))
	for i, permission := range orgOwnerPermissions {
		orgOwnerPermissionNames[i] = permission.Name
	}
	err = auth.AddPermissionsToRole(ownerRole.Id, orgOwnerPermissionNames)
	if err != nil {
		return fmt.Errorf("failed to add permissions to owner role: %w", err)
	}

	err = auth.AddUserToRole(ownerId, ownerRole.Id)
	if err != nil {
		return fmt.Errorf("failed to add owner role to user: %w", err)
	}
	err = auth.AddUserToRole(ownerId, memberRole.Id)
	if err != nil {
		return fmt.Errorf("failed to add member role to user: %w", err)
	}
	return nil
}
