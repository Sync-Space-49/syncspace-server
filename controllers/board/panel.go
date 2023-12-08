package board

import (
	"context"
	"errors"

	"github.com/Sync-Space-49/syncspace-server/controllers/user"
	"github.com/Sync-Space-49/syncspace-server/models"
	"github.com/google/uuid"
)

func (c *Controller) GetPanelsByBoardId(ctx context.Context, boardId string) (*[]models.Panel, error) {
	panels := make([]models.Panel, 0)
	err := c.db.DB.SelectContext(ctx, &panels, `
		SELECT * FROM Panels WHERE board_id=$1 ORDER BY position ASC;
	`, boardId)
	if err != nil {
		return nil, err
	}
	return &panels, nil
}

func (c *Controller) CreatePanel(ctx context.Context, title string, boardId string) (*models.Panel, error) {
	var nextPosition int
	err := c.db.DB.GetContext(ctx, &nextPosition, `
		SELECT COALESCE(MAX(position)+1, 0) AS next_position FROM Panels where board_id=$1;
	`, boardId)
	if err != nil {
		return nil, err
	}
	panelId := uuid.New().String()
	_, err = c.db.DB.ExecContext(ctx, `
		INSERT INTO Panels (id, title, position, board_id) VALUES ($1, $2, $3, $4);
	`, panelId, title, nextPosition, boardId)
	if err != nil {
		return nil, err
	}
	panel, err := c.GetPanelById(ctx, panelId)
	if err != nil {
		return nil, err
	}
	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return nil, err
	}
	return panel, nil
}

func (c *Controller) GetPanelById(ctx context.Context, panelId string) (*models.Panel, error) {
	var panel models.Panel
	err := c.db.DB.GetContext(ctx, &panel, `
		SELECT * FROM Panels WHERE id=$1;
	`, panelId)
	if err != nil {
		return nil, err
	}
	return &panel, nil
}

func (c *Controller) UpdatePanelById(ctx context.Context, boardId string, panelId string, title string, position *int) error {
	panel, err := c.GetPanelById(ctx, panelId)
	if err != nil {
		return err
	}
	if title == "" {
		title = panel.Title
	}
	if position == nil {
		position = &panel.Position
	}

	if *position != panel.Position {
		var maxPosition int
		err := c.db.DB.GetContext(ctx, &maxPosition, `
			SELECT COALESCE(MAX(position), 0) AS max_position FROM Panels where board_id=$1;
		`, boardId)
		if err != nil {
			return err
		}
		if *position > maxPosition || *position < 0 {
			return errors.New("position is out of range")
		}

		if *position > panel.Position {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Panels SET position=position-1 WHERE board_id=$1 AND position>$2 AND position<=$3;
			`, boardId, panel.Position, *position)
		} else {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Panels SET position=position+1 WHERE board_id=$1 AND position<$2 AND position>=$3;
			`, boardId, panel.Position, *position)
		}
		if err != nil {
			return err
		}
	}
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Panels SET title=$1, position=$2 WHERE id=$3;
	`, title, *position, panelId)
	if err != nil {
		return err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeletePanelById(ctx context.Context, boardId string, panelId string) error {
	panel, err := c.GetPanelById(ctx, panelId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `
		DELETE FROM Panels WHERE id=$1;
	`, panelId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Panels SET position=position-1 WHERE board_id=$1 AND position>$2;
	`, boardId, panel.Position)
	if err != nil {
		return err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) GetCompletePanelById(ctx context.Context, panelId string) (*models.CompletePanel, error) {
	panel, err := c.GetPanelById(ctx, panelId)
	if err != nil {
		return nil, err
	}
	completePanel := models.CopyToCompletePanel(*panel)
	completePanel.Stacks = make([]models.CompleteStack, 0)
	stacks, err := c.GetStacksByPanelId(ctx, panel.Id.String())
	if err != nil {
		return nil, err
	}
	if len(*stacks) > 0 {
		for _, stack := range *stacks {
			completeStack := models.CopyToCompleteStack(stack)
			completeStack.Cards = make([]models.CompleteCard, 0)
			cards, err := c.GetCardsByStackId(ctx, stack.Id.String())
			if err != nil {
				return nil, err
			}
			if len(*cards) > 0 {
				for _, card := range *cards {
					completeCard := models.CopyToCompleteCard(card)
					completeCard.Assignments = make([]models.User, 0)
					assignments, err := c.GetAssignedUsersByCardId(ctx, card.Id.String())
					if err != nil {
						return nil, err
					}
					if len(*assignments) > 0 {
						for _, assignment := range *assignments {
							assignedUser, err := user.GetUser(assignment)
							if err != nil {
								return nil, err
							}
							completeCard.Assignments = append(completeCard.Assignments, *assignedUser)
						}
					} else {
						completeCard.Assignments = make([]models.User, 0)
					}
					completeStack.Cards = append(completeStack.Cards, completeCard)
				}
			} else {
				completeStack.Cards = make([]models.CompleteCard, 0)
			}
			completePanel.Stacks = append(completePanel.Stacks, completeStack)
		}
	} else {
		completePanel.Stacks = make([]models.CompleteStack, 0)
	}
	return &completePanel, nil
}
