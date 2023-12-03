package board

import (
	"context"
	"errors"

	"github.com/Sync-Space-49/syncspace-server/models"
	"github.com/google/uuid"
)

func (c *Controller) GetStacksByPanelId(ctx context.Context, panelId string) (*[]models.Stack, error) {
	stacks := make([]models.Stack, 0)
	err := c.db.DB.SelectContext(ctx, &stacks, `
		SELECT * FROM Stacks WHERE panel_id=$1 ORDER BY position ASC;
	`, panelId)
	if err != nil {
		return nil, err
	}
	return &stacks, nil
}

func (c *Controller) CreateStack(ctx context.Context, title string, boardId string, panelId string) (*models.Stack, error) {
	var nextPosition int
	err := c.db.DB.GetContext(ctx, &nextPosition, `
		SELECT COALESCE(MAX(position)+1, 0) AS next_position FROM Stacks where panel_id=$1;
	`, panelId)
	if err != nil {
		return nil, err
	}
	stackId := uuid.New().String()
	_, err = c.db.DB.ExecContext(ctx, `
		INSERT INTO Stacks (id, title, position, panel_id) VALUES ($1, $2, $3, $4);
	`, stackId, title, nextPosition, panelId)
	if err != nil {
		return nil, err
	}
	stack, err := c.GetStackById(ctx, stackId)
	if err != nil {
		return nil, err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return nil, err
	}

	return stack, nil
}

func (c *Controller) GetStackById(ctx context.Context, stackId string) (*models.Stack, error) {
	stack := models.Stack{}
	err := c.db.DB.GetContext(ctx, &stack, `
		SELECT * FROM Stacks WHERE id=$1;
	`, stackId)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

func (c *Controller) UpdateStackById(ctx context.Context, boardId string, panelId string, stackId string, title string, position *int) error {
	stack, err := c.GetStackById(ctx, stackId)
	if err != nil {
		return err
	}
	if title == "" {
		title = stack.Title
	}
	if position == nil {
		position = &stack.Position
	}

	if *position != stack.Position {
		var maxPosition int
		err := c.db.DB.GetContext(ctx, &maxPosition, `
			SELECT COALESCE(MAX(position), 0) AS max_position FROM Stacks where panel_id=$1;
		`, panelId)
		if err != nil {
			return err
		}
		if *position > maxPosition || *position < 0 {
			return errors.New("position is out of range")
		}

		if *position > stack.Position {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Stacks SET position=position-1 WHERE panel_id=$1 AND position>$2 AND position<=$3;
			`, panelId, stack.Position, *position)
		} else {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Stacks SET position=position+1 WHERE panel_id=$1 AND position<$2 AND position>=$3;
			`, panelId, stack.Position, *position)
		}
		if err != nil {
			return err
		}
	}
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Stacks SET title=$1, position=$2 WHERE id=$3;
	`, title, *position, stackId)
	if err != nil {
		return err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeleteStackById(ctx context.Context, boardId string, panelId string, stackId string) error {
	stack, err := c.GetStackById(ctx, stackId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `
		DELETE FROM Stacks WHERE id=$1;
	`, stackId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Stacks SET position=position-1 WHERE panel_id=$1 AND position>$2;
	`, panelId, stack.Position)
	if err != nil {
		return err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) GetCompleteStackById(ctx context.Context, stackId string) (*models.CompleteStack, error) {
	stack, err := c.GetStackById(ctx, stackId)
	if err != nil {
		return nil, err
	}
	completeStack := models.CopyToCompleteStack(*stack)
	completeStack.Cards = make([]models.CompleteCard, 0)
	cards, err := c.GetCardsByStackId(ctx, stack.Id.String())
	if err != nil {
		return nil, err
	}
	if len(*cards) > 0 {
		for _, card := range *cards {
			completeCard := models.CopyToCompleteCard(card)
			completeCard.Assignments = make([]string, 0)
			assignments, err := c.GetAssignedUsersByCardId(ctx, card.Id.String())
			if err != nil {
				return nil, err
			}
			completeCard.Assignments = *assignments
			completeStack.Cards = append(completeStack.Cards, completeCard)
		}
	} else {
		completeStack.Cards = make([]models.CompleteCard, 0)
	}
	return &completeStack, nil
}
