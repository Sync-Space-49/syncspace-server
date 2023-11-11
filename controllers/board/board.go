package board

import (
	"context"
	"fmt"
	"time"

	"github.com/Sync-Space-49/syncspace-server/auth"

	"github.com/google/uuid"
)

func (c *Controller) GetViewableBoardsInOrg(ctx context.Context, orgId string, userId string) (*[]Board, error) {
	var orgBoards []Board
	err := c.db.DB.GetContext(ctx, &orgBoards, `
		SELECT * FROM Boards WHERE organization_id=$1;
	`, orgId)
	if err != nil {
		return nil, err
	}

	var viewableBoards []Board
	for _, board := range orgBoards {
		if board.IsPrivate {
			readBoardPerm := fmt.Sprintf("org%s:board%s:read", orgId, board.Id)
			canReadBoard, err := auth.HasPermission(userId, readBoardPerm)
			if err != nil {
				return nil, err
			}
			if !canReadBoard {
				continue
			}
		}
		viewableBoards = append(viewableBoards, board)
	}
	return &viewableBoards, nil
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

func (c *Controller) GetCompleteBoardById(ctx context.Context, boardId string) (*CompleteBoard, error) {
	var completeBoard CompleteBoard
	board, err := c.GetBoardById(ctx, boardId)
	if err != nil {
		return nil, err
	}
	CopyToCompleteBoard(*board, &completeBoard)
	panels, err := c.GetPanelsByBoardId(ctx, boardId)
	if err != nil {
		return nil, err
	}
	for _, panel := range *panels {
		var completePanel CompletePanel
		CopyToCompletePanel(panel, &completePanel)
		completePanel.Stacks = make([]CompleteStack, 0)
		stacks, err := c.GetStacksByPanelId(ctx, panel.Id.String())
		if err != nil {
			return nil, err
		}
		for _, stack := range *stacks {
			var completeStack CompleteStack
			CopyToCompleteStack(stack, &completeStack)
			cards, err := c.GetCardsByStackId(ctx, stack.Id.String())
			if err != nil {
				return nil, err
			}
			completeStack.Cards = *cards
			completePanel.Stacks = append(completePanel.Stacks, completeStack)
		}
		completeBoard.Panels = append(completeBoard.Panels, completePanel)
	}
	return &completeBoard, nil
}

func (c *Controller) CreateBoard(ctx context.Context, userId string, name string, isPrivate bool, orgId string) (*Board, error) {
	query := `INSERT INTO Boards (id, title, is_private, organization_id, owner_id) VALUES ($1, $2, $3, $4, $5);`
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
	createPanelPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:create_panel", orgId, boardId),
		Description: fmt.Sprintf("Allows you to create a panel on the board with id %s", boardId),
	}
	deletePanelPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:delete_panel", orgId, boardId),
		Description: fmt.Sprintf("Allows you to delete a panel on the board with id %s", boardId),
	}
	updatePanelPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:update_panel", orgId, boardId),
		Description: fmt.Sprintf("Allows you to update a panel on the board with id %s", boardId),
	}
	createStackPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:create_stack", orgId, boardId),
		Description: fmt.Sprintf("Allows you to create a stack on the board with id %s", boardId),
	}
	updateStackPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:update_stack", orgId, boardId),
		Description: fmt.Sprintf("Allows you to update a stack on the board with id %s", boardId),
	}
	deleteStackPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:delete_stack", orgId, boardId),
		Description: fmt.Sprintf("Allows you to delete a stack on the board with id %s", boardId),
	}
	createCardPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:create_card", orgId, boardId),
		Description: fmt.Sprintf("Allows you to create a card on the board with id %s", boardId),
	}
	updateCardPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:update_card", orgId, boardId),
		Description: fmt.Sprintf("Allows you to update a card on the board with id %s", boardId),
	}
	deleteCardPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:board%s:delete_card", orgId, boardId),
		Description: fmt.Sprintf("Allows you to delete a card on the board with id %s", boardId),
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
		readPerm, createCardPerm, updateCardPerm, deleteCardPerm,
	}
	boardOwnerPermissions := append(boardMemberPermissions, deletePerm, updatePerm, createPanelPerm, deletePanelPerm, updatePanelPerm, createStackPerm, updateStackPerm, deleteStackPerm, addMembersPerm, removeMembersPerm)

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
