package board

import (
	"context"
	"errors"
)

func (c *Controller) GetCardsByStackId(ctx context.Context, stackId string) (*[]Card, error) {
	cards := make([]Card, 0)
	err := c.db.DB.SelectContext(ctx, &cards, `
		SELECT * FROM Cards WHERE stack_id=$1 ORDER BY position ASC;
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

func (c *Controller) GetCardById(ctx context.Context, cardId string) (*Card, error) {
	card := Card{}
	err := c.db.DB.GetContext(ctx, &card, `
		SELECT * FROM Cards WHERE id=$1;
	`, cardId)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (c *Controller) UpdateCardById(ctx context.Context, stackId string, cardId string, title string, description string, position *int) error {
	card, err := c.GetCardById(ctx, stackId)
	if err != nil {
		return err
	}
	if title == "" {
		title = card.Title
	}
	if description == "" {
		description = card.Description
	}
	if position == nil {
		position = &card.Position
	}

	if *position != card.Position {
		var maxPosition int
		err := c.db.DB.GetContext(ctx, &maxPosition, `
			SELECT COALESCE(MAX(position), 0) AS max_position FROM Cards where stack_id=$1;
		`, stackId)
		if err != nil {
			return err
		}
		if *position > maxPosition || *position < 0 {
			return errors.New("position is out of range")
		}

		if *position > card.Position {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Cards SET position=position-1, stack_id=$1 WHERE stack_id=$1 AND position>$2 AND position<=$3;
			`, stackId, card.Position, *position)
		} else {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Cards SET position=position+1, stack_id=$1 WHERE stack_id=$1 AND position<$2 AND position>=$3;
			`, stackId, card.Position, *position)
		}
		if err != nil {
			return err
		}
	}
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Cards SET title=$1, description=$2, position=$3 WHERE id=$4;
	`, title, description, *position, cardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeleteCardById(ctx context.Context, stackId string, cardId string) error {
	card, err := c.GetCardById(ctx, cardId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `
		DELETE FROM Cards WHERE id=$1;
	`, cardId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Cards SET position=position-1 WHERE stack_id=$1 AND position>$2;
	`, stackId, card.Position)
	if err != nil {
		return err
	}
	return nil
}
