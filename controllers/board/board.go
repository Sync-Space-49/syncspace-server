package board

import (
	"context"
	"fmt"
	"time"

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
	query = `INSERT INTO Boards (id, title, is_private, organization_id, owner_id) VALUES ($1, $2, $3, $4, $5);`
	boardId := uuid.New().String()
	_, err := c.db.DB.ExecContext(ctx, query, boardId, name, isPrivate, orgId, userId)
	if err != nil {
		return nil, err
	}
	org, err := c.GetBoardById(ctx, boardId)
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

func (c *Controller) UpdateBoardById(ctx context.Context, orgId string, boardId string, title string, isPrivate bool, ownerId string, previousOwnerId string) error {
	board, err := c.GetBoardById(ctx, boardId)
	if err != nil {
		return err
	}
	if title == "" {
		title = board.Title
	}
	if ownerId == "" {
		ownerId = board.OwnerId
		// fmt.Printf("2 ownerId: %s", ownerId)
	}
	modified_at := time.Now().UTC()
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Boards SET title=$1, is_private=$2, owner_id=$3, modified_at=$4 WHERE id=$5;
	`, title, isPrivate, ownerId, modified_at, boardId)
	if err != nil {
		return err
	}

	boardOwnerRoleName := fmt.Sprintf("org%s:board%s:owner", orgId, boardId)
	// fmt.Printf("org%s:board%s:owner", orgId, boardId)
	boardOwnerRoles, err := auth.GetRoles(&boardOwnerRoleName)
	if err != nil {
		return err
	}
	if len(*boardOwnerRoles) == 0 {
		return fmt.Errorf("no roles found for board %s", boardOwnerRoles)
	}
	err = auth.RemoveUserFromRole(previousOwnerId, (*boardOwnerRoles)[0].Id)
	// fmt.Printf("org%s:board%s:owner removed from user %s", orgId, boardId, previousOwnerId)
	if err != nil {
		return err
	}
	err = auth.AddUserToRole(ownerId, (*boardOwnerRoles)[0].Id)
	// fmt.Printf("org%s:board%s:owner added to user %s", orgId, boardId, ownerId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeleteBoardById(ctx context.Context, boardId string) error {
	_, err := c.db.DB.ExecContext(ctx, `
		DELETE FROM Boards WHERE id=$1;
	`, boardId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) AddMemberToBoard(userId string, orgId string, boardId string) error {
	boardMemberRoleName := fmt.Sprintf("org%s:board%s:member", orgId, boardId)
	// fmt.Printf("org%s:board%s:owner", orgId, boardId)
	boardMemberRoles, err := auth.GetRoles(&boardMemberRoleName)
	if err != nil {
		return err
	}
	err = auth.AddUserToRole(userId, (*boardMemberRoles)[0].Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) RemoveMemberFromBoard(userId string, orgId string, boardId string) error {
	boardRolePrefix := fmt.Sprintf("org%s:board%s:", orgId, boardId)
	boardRoles, err := auth.GetRoles(&boardRolePrefix)
	// fmt.Printf("boardRoles %s", boardRoles)
	if err != nil {
		return err
	}
	if len(*boardRoles) == 0 {
		return fmt.Errorf("no roles found for board %s", boardId)
	}
	var boardRoleIds []string
	for _, role := range *boardRoles {
		boardRoleIds = append(boardRoleIds, role.Id)
	}
	err = auth.RemoveUserFromRoles(userId, boardRoleIds)
	if err != nil {
		return err
	}
	return nil
}
