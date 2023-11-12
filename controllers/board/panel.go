package board

import (
	"context"
)

func (c *Controller) GetPanelsByBoardId(ctx context.Context, boardId string) (*[]Panel, error) {
	panels := make([]Panel, 0)
	err := c.db.DB.SelectContext(ctx, &panels, `
		SELECT * FROM Panels WHERE board_id=$1 ORDER BY position ASC;
	`, boardId)
	if err != nil {
		return nil, err
	}
	return &panels, nil
}

func (c *Controller) CreatePanel(ctx context.Context, title string, boardId string) error {
	var nextPosition int
	err := c.db.DB.GetContext(ctx, &nextPosition, `
		SELECT COALESCE(MAX(position)+1, 0) AS next_position FROM Panels where board_id=$1;
	`, boardId)
	if err != nil {
		return err
	}

	_, err = c.db.DB.ExecContext(ctx, `
		INSERT INTO Panels (title, position, board_id) VALUES ($1, $2, $3);
	`, title, nextPosition, boardId)
	if err != nil {
		return err
	}

	return nil
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
	return nil
}
