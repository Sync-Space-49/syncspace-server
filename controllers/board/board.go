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

func (c *Controller) CreateBoard(ctx context.Context, userId string, name string, isPrivate bool, orgId string) (*Board, error) {
	var query string
	query = `INSERT INTO Boards (id, title, is_private, organization_id) VALUES ($1, $2, $3, $4);`
	orgID := uuid.New().String()
	_, err := c.db.DB.ExecContext(ctx, query, orgID, name, isPrivate, orgId)
	if err != nil {
		return nil, err
	}
	org, err := c.GetBoardById(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (c *Controller) InitializeBoard(ownerId string, boardId string, orgId string) error {
	ownerRoleName := fmt.Sprintf("org%s:board%s:owner", orgId, boardId)
	ownerRoleDescription := fmt.Sprintf("Owner of board with the id: %s", boardId)
	ownerRole, err := auth.CreateRole(ownerRoleName, ownerRoleDescription)
	if err != nil {
		return err
	}
	memberRoleName := fmt.Sprintf("org%s:board%s:member", orgId, boardId)
	memberRoleDescription := fmt.Sprintf("Member of board with the id: %s", boardId)
	memberRole, err := auth.CreateRole(memberRoleName, memberRoleDescription)
	if err != nil {
		return err
	}

	readPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:read", orgId, boardId),
		Description: fmt.Sprintf("Allows you to read the contents of the board with id %s", boardId),
	}
	deletePerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:delete", orgId, boardId),
		Description: fmt.Sprintf("Allows you to delete the board with id %s", boardId),
	}
	updatePerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:update", orgId, boardId),
		Description: fmt.Sprintf("Allows you to update info about the board with id %s", boardId),
	}
	addMembersPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:add_members", orgId, boardId),
		Description: fmt.Sprintf("Allows you to add members to the board with id %s", boardId),
	}
	removeMembersPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:remove_members", orgId, boardId),
		Description: fmt.Sprintf("Allows you to remove members from the board with id %s", boardId),
	}

	boardMemberPermissions := []auth.Permission{
		readPerm,
	}
	boardOwnerPermissions := append(boardMemberPermissions, deletePerm, updatePerm, addMembersPerm, removeMembersPerm)

	// Because the owner role has all permissions, we only need to call CreatePermissions once
	err = auth.CreatePermissions(boardOwnerPermissions)
	if err != nil {
		return fmt.Errorf("failed to create permissions for board: %w", err)
	}

	boardMemberPermissionNames := make([]string, len(boardMemberPermissions))
	for i, permission := range boardMemberPermissions {
		boardMemberPermissionNames[i] = permission.Name
	}
	err = auth.AddPermissionsToRole(memberRole.Id, boardMemberPermissionNames)
	if err != nil {
		return fmt.Errorf("failed to add permissions to member role: %w", err)
	}
	boardOwnerPermissionNames := make([]string, len(boardOwnerPermissions))
	for i, permission := range boardOwnerPermissions {
		boardOwnerPermissionNames[i] = permission.Name
	}
	err = auth.AddPermissionsToRole(ownerRole.Id, boardOwnerPermissionNames)
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
