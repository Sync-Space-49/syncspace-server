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
