package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Sync-Space-49/syncspace-server/auth"
	"github.com/Sync-Space-49/syncspace-server/models"
	"github.com/jmoiron/sqlx"
)

func (c *Controller) GetUsers() (*[]models.User, error) {
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return &[]models.User{}, err
	}
	method := "GET"
	url := fmt.Sprintf("%sapi/v2/users", c.cfg.Auth0.Domain)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &[]models.User{}, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return &[]models.User{}, fmt.Errorf("invalid request: %s", string(body))
	}

	var users []models.User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return &[]models.User{}, err
	}

	return &users, nil
}

func (c *Controller) GetUserById(userId string) (*models.User, error) {
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return &models.User{}, err
	}
	method := "GET"
	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &models.User{}, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return &models.User{}, fmt.Errorf("invalid request: %s", string(body))
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}

func (c *Controller) UpdateUserById(userId string, email string, username string, password string, pfpUrl *string) error {
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return fmt.Errorf("failed to get maintenance token: %w", err)
	}

	user, err := c.GetUserById(userId)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	if username == "" {
		username = user.Username
	}
	if pfpUrl == nil {
		pfpUrl = &user.Picture
	}

	method := "PATCH"
	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)

	var payload io.Reader
	if password != "" {
		payload = strings.NewReader(fmt.Sprintf(`{"email":"%s","picture":"%s","password":"%s"}`, email, *pfpUrl, password))
	} else {
		payload = strings.NewReader(fmt.Sprintf(`{"username":"%s","picture":"%s"}`, username, *pfpUrl))
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update user: %s", string(body))
	}

	// seperate request for email because Auth0 won't let you update email and username at the same time
	if email != "" {
		payload = strings.NewReader(fmt.Sprintf(`{"email":"%s"}`, email))
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to update user: %s", string(body))
		}
	}

	return nil
}

func (c *Controller) DeleteUserById(ctx context.Context, userId string) error {
	managementToken, err := auth.GetManagementToken()
	if err != nil {
		return err
	}
	method := "DELETE"
	url := fmt.Sprintf("%sapi/v2/users/%s", c.cfg.Auth0.Domain, userId)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("BEARER %s", managementToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete user: %s", string(body))
	}

	_, err = c.db.DB.ExecContext(ctx, `DELETE FROM Organizations WHERE owner_id=$1`, userId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `DELETE FROM Boards WHERE owner_id=$1`, userId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `DELETE FROM Assigned_Cards WHERE user_id=$1`, userId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) GetUserOrganizationsById(ctx context.Context, userId string) (*[]models.Organization, error) {
	usersRoles, err := auth.GetUserRoles(userId)
	if err != nil {
		return nil, err
	}

	var orgIds []string
	findUUIDInRoleRegex := regexp.MustCompile(`org(.*?):`)
	for _, role := range *usersRoles {
		matches := findUUIDInRoleRegex.FindStringSubmatch(role.Name)
		if len(matches) < 2 {
			continue
		}
		organizationId := matches[1]
		alreadyFound := false
		for _, orgId := range orgIds {
			if orgId == organizationId {
				alreadyFound = true
				break
			}
		}
		if !alreadyFound {
			orgIds = append(orgIds, organizationId)
		}
	}

	if len(orgIds) == 0 {
		return &[]models.Organization{}, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM Organizations WHERE id IN (?)`, orgIds)
	if err != nil {
		return nil, err
	}
	query = c.db.DB.Rebind(query)
	var organizations []models.Organization
	err = c.db.DB.SelectContext(ctx, &organizations, query, args...)
	if err != nil {
		return nil, err
	}
	return &organizations, nil
}

func (c *Controller) GetUserOwnedOrganizationsById(ctx context.Context, userId string) (*[]models.Organization, error) {
	var organizations []models.Organization
	err := c.db.DB.SelectContext(ctx, &organizations, `
		SELECT * FROM Organizations WHERE owner_id=$1;
	`, userId)
	if err != nil {
		return nil, err
	}
	return &organizations, nil
}

func (c *Controller) GetUserBoardsById(ctx context.Context, userId string) (*[]models.Board, error) {
	usersRoles, err := auth.GetUserRoles(userId)
	if err != nil {
		return nil, err
	}

	var boardIds []string
	findUUIDInRoleRegex := regexp.MustCompile(`org.*?:board(.*?):`)
	for _, role := range *usersRoles {
		matches := findUUIDInRoleRegex.FindStringSubmatch(role.Name)
		if len(matches) < 2 {
			continue
		}
		boardId := matches[1]
		alreadyFound := false
		for _, bId := range boardIds {
			if bId == boardId {
				alreadyFound = true
				break
			}
		}
		if !alreadyFound {
			boardIds = append(boardIds, boardId)
		}
	}

	if len(boardIds) == 0 {
		return &[]models.Board{}, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM Boards WHERE id IN (?)`, boardIds)
	if err != nil {
		return nil, err
	}
	query = c.db.DB.Rebind(query)
	var boards []models.Board
	err = c.db.DB.SelectContext(ctx, &boards, query, args...)
	if err != nil {
		return nil, err
	}
	return &boards, nil
}

func (c *Controller) GetUserOwnedBoardsById(ctx context.Context, userId string) (*[]models.Board, error) {
	var boards []models.Board
	err := c.db.DB.SelectContext(ctx, &boards, `
		SELECT * FROM Boards WHERE owner_id=$1;
	`, userId)
	if err != nil {
		return nil, err
	}
	return &boards, nil
}

func (c *Controller) GetUserAssignedCardsById(ctx context.Context, userId string) (*[]models.Card, error) {
	var cards []models.Card
	err := c.db.DB.SelectContext(ctx, &cards, `
	SELECT ac.user_id, c.*, s.id as stack_id, p.id as panel_id, b.id as board_id, o.id as org_id
	FROM Assigned_Cards ac
	JOIN Cards c on ac.card_id = c.id
	JOIN Stacks s ON c.stack_id = s.id
	JOIN Panels p ON s.panel_id = p.id
	JOIN Boards b ON p.board_id = b.id
	JOIN organizations o ON b.organization_id = o.id
	WHERE ac.user_id = $1;
	`, userId)
	if err != nil {
		return nil, err
	}
	return &cards, nil
}

func (c *Controller) GetFavouriteBoards(ctx context.Context, userId string) (*[]models.Board, error) {
	var favouriteBoards []models.Board
	//  This Join will return board objects that are favorited by the user with the given userId
	err := c.db.DB.SelectContext(ctx, &favouriteBoards, `
		SELECT b.*
		FROM favorite_boards AS fb
		JOIN boards AS b ON fb.board_id = b.id
		WHERE fb.user_id = $1;
	`, userId)
	if err != nil {
		return nil, err
	}
	return &favouriteBoards, nil
}

func (c *Controller) AddFavouriteBoard(ctx context.Context, userId string, boardId string) error {
	_, err := c.db.DB.ExecContext(ctx, `
		INSERT INTO favorite_boards (user_id, board_id) VALUES ($1, $2);
	`, userId, boardId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) RemoveFavouriteBoard(ctx context.Context, userId string, boardId string) error {
	_, err := c.db.DB.ExecContext(ctx, `
		DELETE FROM favorite_boards WHERE user_id=$1 AND board_id=$2;
	`, userId, boardId)
	if err != nil {
		return err
	}
	return nil
}
