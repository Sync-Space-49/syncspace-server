package board

import "context"

func (c *Controller) GetCardsByStackId(ctx context.Context, stackId string) (*[]Card, error) {
	cards := make([]Card, 0)
	err := c.db.DB.SelectContext(ctx, &cards, `
		SELECT * FROM Cards WHERE stack_id=$1;
	`, stackId)
	if err != nil {
		return nil, err
	}
	return &cards, nil
}

func (c *Controller) CreateCard(ctx context.Context, title string, description string, stackId string) error {
	var nextPosition int
	err := c.db.DB.GetContext(ctx, &nextPosition, `
		SELECT COALESCE(MAX(position)+1, 0) AS next_position FROM Cards where stack_id=$1;
	`, stackId)
	if err != nil {
		return err
	}

	_, err = c.db.DB.ExecContext(ctx, `
		INSERT INTO Cards (title, description, position, stack_id) VALUES ($1, $2, $3. $4);
	`, title, description, nextPosition, stackId)
	if err != nil {
		return err
	}
	return nil
}
