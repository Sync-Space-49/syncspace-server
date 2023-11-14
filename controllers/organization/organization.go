package organization

import (
	"context"
	"fmt"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/google/uuid"
)

func (c *Controller) CreateOrganization(ctx context.Context, userId string, title string, description *string, aiEnabled bool) (*Organization, error) {
	var query string
	if description == nil {
		query = `INSERT INTO Organizations (id, owner_id, name, ai_enabled) VALUES ($1, $2, $3, $5);`
	} else {
		query = `INSERT INTO Organizations (id, owner_id, name, description, ai_enabled) VALUES ($1, $2, $3, $4, $5);`
	}
	orgID := uuid.New().String()
	_, err := c.db.DB.ExecContext(ctx, query, orgID, userId, title, description, aiEnabled)
	if err != nil {
		return nil, err
	}
	org, err := c.GetOrganizationById(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (c *Controller) InitializeOrganization(ownerId string, organizationId string) error {
	ownerRoleName := fmt.Sprintf("org%s:owner", organizationId)
	ownerRoleDescription := fmt.Sprintf("Owner of organization with the id: %s", organizationId)
	ownerRole, err := auth.CreateRole(ownerRoleName, ownerRoleDescription)
	if err != nil {
		return err
	}
	memberRoleName := fmt.Sprintf("org%s:member", organizationId)
	memberRoleDescription := fmt.Sprintf("Member of organization with the id: %s", organizationId)
	memberRole, err := auth.CreateRole(memberRoleName, memberRoleDescription)
	if err != nil {
		return err
	}

	readPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:read", organizationId),
		Description: fmt.Sprintf("Allows you to read the contents of the organization with id %s", organizationId),
	}
	deletePerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:delete", organizationId),
		Description: fmt.Sprintf("Allows you to delete the organization with id %s", organizationId),
	}
	updatePerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:update", organizationId),
		Description: fmt.Sprintf("Allows you to update info about the organization with id %s", organizationId),
	}
	addMembersPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:add_members", organizationId),
		Description: fmt.Sprintf("Allows you to add members to the organization with id %s", organizationId),
	}
	removeMembersPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:remove_members", organizationId),
		Description: fmt.Sprintf("Allows you to remove members from the organization with id %s", organizationId),
	}
	createRolesPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:create_roles", organizationId),
		Description: fmt.Sprintf("Allows you to use create roles for the organization with id %s", organizationId),
	}
	editRolesPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:edit_roles", organizationId),
		Description: fmt.Sprintf("Allows you to edit role properties, including adding and removing permissions, from roles for the organization with id %s", organizationId),
	}
	deleteRolesPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:delete_roles", organizationId),
		Description: fmt.Sprintf("Allows you to delete roles from the organization with id %s", organizationId),
	}
	addRolesPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:add_roles", organizationId),
		Description: fmt.Sprintf("Allows you to add roles to users in the organization with id %s", organizationId),
	}
	removeRolesPerm := auth.Permission{
		Name:        fmt.Sprintf("org%s:remove_roles", organizationId),
		Description: fmt.Sprintf("Allows you to remove roles from users in the organization with id %s", organizationId),
	}

	orgMemberPermissions := []auth.Permission{
		readPerm,
	}
	orgOwnerPermissions := append(orgMemberPermissions, deletePerm, updatePerm, addMembersPerm, removeMembersPerm, createRolesPerm, editRolesPerm, deleteRolesPerm, addRolesPerm, removeRolesPerm)

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

func (c *Controller) GetOrganizationById(ctx context.Context, organizationId string) (*Organization, error) {
	var organization Organization
	err := c.db.DB.GetContext(ctx, &organization, `
		SELECT * FROM Organizations WHERE id=$1;
	`, organizationId)
	if err != nil {
		return nil, err
	}
	return &organization, nil
}

func (c *Controller) UpdateOrganizationById(ctx context.Context, organizationId string, title string, description string, ai_enabled bool) error {
	org, err := c.GetOrganizationById(ctx, organizationId)
	if err != nil {
		return err
	}
	if title == "" {
		title = org.Name
	}
	if description == "" {
		description = org.Description
	}

	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Organizations SET name=$1, description=$2, ai_enabled=$3 WHERE id=$4;
	`, title, description, ai_enabled, organizationId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) DeleteOrganizationById(ctx context.Context, organizationId string) error {
	_, err := c.db.DB.ExecContext(ctx, `
		DELETE FROM Organizations WHERE id=$1;
	`, organizationId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) AddMember(userId string, organizationId string) error {
	orgMemberRoleName := fmt.Sprintf("org%s:member", organizationId)
	orgMemberRoles, err := auth.GetRoles(&orgMemberRoleName)
	if err != nil {
		return err
	}
	if len(*orgMemberRoles) == 0 {
		return fmt.Errorf("no roles found for organization %s", organizationId)
	}
	err = auth.AddUserToRole(userId, (*orgMemberRoles)[0].Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) RemoveMember(userId string, organizationId string) error {
	orgRolePrefix := fmt.Sprintf("org%s:", organizationId)
	orgRoles, err := auth.GetRoles(&orgRolePrefix)
	if err != nil {
		return err
	}
	if len(*orgRoles) == 0 {
		return fmt.Errorf("no roles found for organization %s", organizationId)
	}
	var orgRoleIds []string
	for _, role := range *orgRoles {
		orgRoleIds = append(orgRoleIds, role.Id)
	}
	err = auth.RemoveUserFromRoles(userId, orgRoleIds)
	if err != nil {
		return err
	}
	return nil
}
