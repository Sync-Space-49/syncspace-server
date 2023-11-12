package board

import "context"

func (c *Controller) GetStacksByPanelId(ctx context.Context, panelId string) (*[]Stack, error) {
	stacks := make([]Stack, 0)
	err := c.db.DB.SelectContext(ctx, &stacks, `
		SELECT * FROM Stacks WHERE panel_id=$1 ORDER BY position ASC;
	`, panelId)
	if err != nil {
		return nil, err
	}
	return &stacks, nil
}

func (c *Controller) CreateStack(ctx context.Context, title string, panelId string) error {
	var nextPosition int
	err := c.db.DB.GetContext(ctx, &nextPosition, `
		SELECT COALESCE(MAX(position)+1, 0) AS next_position FROM Stacks where panel_id=$1;
	`, panelId)
	if err != nil {
		return err
	}

	_, err = c.db.DB.ExecContext(ctx, `
		INSERT INTO Stacks (title, position, panel_id) VALUES ($1, $2, $3);
	`, title, nextPosition, panelId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) GetStackById(ctx context.Context, stackId string) (*Stack, error) {
	stack := Stack{}
	err := c.db.DB.GetContext(ctx, &stack, `
		SELECT * FROM Stacks WHERE id=$1;
	`, stackId)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

func (c *Controller) UpdateStackById(ctx context.Context, panelId string, stackId string, title string, position *int) error {
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

	return nil
}

func (c *Controller) DeleteStackById(ctx context.Context, panelId string, stackId string) error {
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
	return nil
}

func (c *Controller) GetCompleteStackById(ctx context.Context, boardId string) (*CompleteStack, error) {
	stack, err := c.GetStackById(ctx, boardId)
	if err != nil {
		return nil, err
	}
	completeStack := CopyToCompleteStack(*stack)
	cards, err := c.GetCardsByStackId(ctx, stack.Id.String())
	if err != nil {
		return nil, err
	}
	completeStack.Cards = *cards
	return &completeStack, nil
}
