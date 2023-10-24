package organization

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

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

func (c *Controller) CreateOrganization(ctx context.Context, userId string, title string, description *string) (*Organization, error) {
	var query string
	if description == nil {
		query = `INSERT INTO Organizations (id, owner_id, name) VALUES ($1, $2, $3);`
	} else {
		query = `INSERT INTO Organizations (id, owner_id, name, description) VALUES ($1, $2, $3, $4);`
	}
	orgID := uuid.New().String()
	_, err := c.db.DB.ExecContext(ctx, query, orgID, userId, title, description)
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

	readPerm := []string{
		fmt.Sprintf("org%s:read", organizationId),
		fmt.Sprintf("Allows you to read the contents of the organization with id %s", organizationId),
	}
	deletePerm := []string{
		fmt.Sprintf("org%s:delete", organizationId),
		fmt.Sprintf("Allows you to delete the organization with id %s", organizationId),
	}
	updatePerm := []string{
		fmt.Sprintf("org%s:update", organizationId),
		fmt.Sprintf("Allows you to update info about the organization with id %s", organizationId),
	}
	addMembersPerm := []string{
		fmt.Sprintf("org%s:add_members", organizationId),
		fmt.Sprintf("Allows you to add members to the organization with id %s", organizationId),
	}
	removeMembersPerm := []string{
		fmt.Sprintf("org%s:remove_members", organizationId),
		fmt.Sprintf("Allows you to remove members from the organization with id %s", organizationId),
	}

	orgMemberPermissions := [][]string{
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
		orgMemberPermissionNames[i] = permission[0]
	}
	err = auth.AddPermissionsToRole(memberRole.Id, orgMemberPermissionNames)
	if err != nil {
		return fmt.Errorf("failed to add permissions to member role: %w", err)
	}
	orgOwnerPermissionNames := make([]string, len(orgOwnerPermissions))
	for i, permission := range orgOwnerPermissions {
		orgOwnerPermissionNames[i] = permission[0]
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

func (c *Controller) UpdateOrganizationById(ctx context.Context, organizationId string, title string, description string) error {
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
		UPDATE Organizations SET name=$1, description=$2 WHERE id=$3;
	`, title, description, organizationId)
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

func (c *Controller) GetUserOrganizations(ctx context.Context, userId string) (*[]Organization, error) {
	usersRoles, err := auth.GetUserRoles(userId)
	if err != nil {
		return nil, err
	}

	var orgIds []int
	re := regexp.MustCompile("[0-9]+")
	for _, role := range *usersRoles {
		organizationId, err := strconv.Atoi(re.FindString(role.Name))
		if err != nil {
			return nil, err
		}
		orgIds = append(orgIds, organizationId)
	}

	var organizations []Organization
	err = c.db.DB.SelectContext(ctx, &organizations, `
		SELECT * FROM Organizations WHERE id IN (?);
	`, orgIds)
	if err != nil {
		return nil, err
	}
	return &organizations, nil
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
