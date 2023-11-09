package board

import "context"

func (c *Controller) GetCardsByStackId(ctx context.Context, stackId string) (*[]Card, error) {
	var cards []Card
	err := c.db.DB.SelectContext(ctx, &cards, `
		SELECT * FROM Cards WHERE stack_id=$1;
	`, stackId)
	if err != nil {
		return nil, err
	}
	return &cards, nil
}
