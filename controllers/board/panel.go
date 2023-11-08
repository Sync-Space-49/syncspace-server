package board

import "context"

func (c *Controller) GetPanelsByBoardId(ctx context.Context, boardId string) (*[]Panel, error) {
	var panels []Panel
	err := c.db.DB.SelectContext(ctx, &panels, `
		SELECT * FROM Panels WHERE board_id=$1;
	`, boardId)
	if err != nil {
		return nil, err
	}
	return &panels, nil
}

func (c *Controller) GetPanelById(ctx context.Context, panelId string) (*Panel, error) {
	var panel Panel
	err := c.db.DB.GetContext(ctx, &panel, `
		SELECT * FROM Panels WHERE id=$1;
	`, panelId)
	if err != nil {
		return nil, err
	}
	return &panel, nil
}
