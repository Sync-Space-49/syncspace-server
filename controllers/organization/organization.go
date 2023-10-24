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

func (c *Controller) CreateOrganization(ctx context.Context, userId string, title string, description string) (*Organization, error) {
	orgID := uuid.New().String()
	_, err := c.db.DB.ExecContext(ctx, `
		INSERT INTO Organizations (id, owner_id, name, description) VALUES ($1, $2, $3, $4);
	`, orgID, userId, title, description)
	if err != nil {
		return nil, err
	}
	org, err := c.GetOrganizationById(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return org, nil
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
