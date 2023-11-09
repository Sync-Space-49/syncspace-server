package board

import "context"

func (c *Controller) GetStacksByPanelId(ctx context.Context, panelId string) (*[]Stack, error) {
	var stacks []Stack
	err := c.db.DB.SelectContext(ctx, &stacks, `
		SELECT * FROM Stacks WHERE panel_id=$1;
	`, panelId)
	if err != nil {
		return nil, err
	}
	return &stacks, nil
}
