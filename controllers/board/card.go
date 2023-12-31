package board

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Sync-Space-49/syncspace-server/controllers/user"
	"github.com/Sync-Space-49/syncspace-server/models"
	"github.com/google/uuid"
)

func (c *Controller) GetCardsByStackId(ctx context.Context, stackId string) (*[]models.Card, error) {
	cards := make([]models.Card, 0)
	err := c.db.DB.SelectContext(ctx, &cards, `
		SELECT * FROM Cards WHERE stack_id=$1 ORDER BY position ASC;
	`, stackId)
	if err != nil {
		return nil, err
	}
	return &cards, nil
}

func (c *Controller) CreateCard(ctx context.Context, title string, description string, points string, boardId string, stackId string) (*models.Card, error) {
	var nextPosition int
	err := c.db.DB.GetContext(ctx, &nextPosition, `
		SELECT COALESCE(MAX(position)+1, 0) AS next_position FROM Cards where stack_id=$1;
	`, stackId)
	if err != nil {
		return nil, err
	}
	cardId := uuid.New().String()
	_, err = c.db.DB.ExecContext(ctx, `
		INSERT INTO Cards (id, title, description, points, position, stack_id) VALUES ($1, $2, $3, $4, $5, $6);
	`, cardId, title, description, points, nextPosition, stackId)
	if err != nil {
		return nil, err
	}
	card, err := c.GetCardById(ctx, cardId)
	if err != nil {
		return nil, err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (c *Controller) GetCardById(ctx context.Context, cardId string) (*models.Card, error) {
	card := models.Card{}
	err := c.db.DB.GetContext(ctx, &card, `
		SELECT * FROM Cards WHERE id=$1;
	`, cardId)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (c *Controller) UpdateCardById(ctx context.Context, boardId string, stackId string, cardId string, newStackId string, title string, description string, points string, position *int) error {
	card, err := c.GetCardById(ctx, cardId)
	if err != nil {
		return err
	}
	if title == "" {
		title = card.Title
	}
	if description == "" {
		description = card.Description
	}
	if position == nil && newStackId != "" {
		var nextPosition int
		err := c.db.DB.GetContext(ctx, &nextPosition, `
			SELECT COALESCE(MAX(position)+1, 0) AS next_position FROM Cards where stack_id=$1;
		`, newStackId)
		if err != nil {
			return err
		}
		position = &nextPosition
	} else {
		if position == nil {
			position = &card.Position
		}
	}
	if newStackId == "" {
		newStackId = stackId
	} else {
		var isStackInBoard bool
		err := c.db.DB.GetContext(ctx, &isStackInBoard, `
			SELECT EXISTS(
				SELECT 1
					FROM Stacks s
						JOIN panels p on s.panel_id = p.id
						JOIN boards b on p.board_id = b.id
					WHERE s.id = $1 AND b.id = $2
			);
		`, newStackId, boardId)
		if err != nil {
			return err
		}
		if !isStackInBoard {
			return errors.New("stack is not in the same board")
		}
	}

	if *position != card.Position {
		var maxPosition int
		err := c.db.DB.GetContext(ctx, &maxPosition, `
			SELECT COALESCE(MAX(position), 0) AS max_position FROM Cards where stack_id=$1;
		`, newStackId)
		if err != nil {
			return err
		}
		if *position > maxPosition || *position < 0 {
			return errors.New("position is out of range")
		}

		if *position > card.Position {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Cards SET position=position-1 WHERE stack_id=$1 AND position>$2 AND position<=$3;
			`, newStackId, card.Position, *position)
		} else {
			_, err = c.db.DB.ExecContext(ctx, `
				UPDATE Cards SET position=position+1 WHERE stack_id=$1 AND position<$2 AND position>=$3;
			`, newStackId, card.Position, *position)
		}
		if err != nil {
			return err
		}
	}
	_, err = c.db.DB.ExecContext(ctx, `
		UPDATE Cards SET title=$1, description=$2, points=$3, position=$4, stack_id=$5 WHERE id=$6;
	`, title, description, points, *position, newStackId, cardId)
	if err != nil {
		return err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeleteCardById(ctx context.Context, boardId string, stackId string, cardId string) error {
	card, err := c.GetCardById(ctx, cardId)
	if err != nil {
		return err
	}
	_, err = c.db.DB.ExecContext(ctx, `
		DELETE FROM assigned_cards WHERE card_id=$1;
	`, cardId)
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

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) CreateCardWithAI(ctx context.Context, boardId string, cardStackId string) (*models.Card, error) {
	requestUrl := fmt.Sprintf("%s/api/generate/card", c.cfg.AI.APIHost)

	detailedBoard, err := c.GetCompleteBoardById(ctx, boardId)
	if err != nil {
		return nil, err
	}
	formattedBoard := models.CopyToSimplifiedCompleteBoard(*detailedBoard)

	type AICardPayload struct {
		StackId string                         `json:"stack_id"`
		Board   models.SimplifiedCompleteBoard `json:"board"`
	}
	aiCardPayload := AICardPayload{
		StackId: cardStackId,
		Board:   formattedBoard,
	}

	var payload bytes.Buffer
	err = json.NewEncoder(&payload).Encode(aiCardPayload)
	if err != nil {
		return nil, err
	}

	method := "POST"
	req, err := http.NewRequest(method, requestUrl, &payload)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error occurred during conversion of HTTP resonse body into bytes. %v", err)
		return nil, err
	}

	var aiCard models.AIGeneratedCard
	err = json.Unmarshal(body, &aiCard)
	if err != nil {
		return nil, err
	}
	CardStoryPointsString := fmt.Sprintf("%v", aiCard.CardStoryPoints)
	card, err := c.CreateCard(ctx, aiCard.CardTitle, aiCard.CardDesc, CardStoryPointsString, boardId, cardStackId)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (c *Controller) AssignCardToUser(ctx context.Context, boardId string, cardId string, userId string) error {
	_, err := c.db.DB.ExecContext(ctx, `
		INSERT INTO assigned_cards (user_id, card_id) VALUES ($1, $2);
	`, userId, cardId)
	if err != nil {
		return err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) UnassignCardFromUser(ctx context.Context, boardId string, cardId string, userId string) error {
	_, err := c.db.DB.ExecContext(ctx, `
		DELETE FROM assigned_cards WHERE user_id=$1 AND card_id=$2;
	`, userId, cardId)
	if err != nil {
		return err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) GetAssignedUsersByCardId(ctx context.Context, cardId string) (*[]string, error) {
	var userIds []string
	err := c.db.DB.SelectContext(ctx, &userIds, `
		SELECT user_id FROM assigned_cards WHERE card_id=$1;
	`, cardId)
	if err != nil {
		return nil, err
	}

	return &userIds, nil
}

func (c *Controller) GetAssignedCardsByUserId(ctx context.Context, boardId string, userId string) (*[]string, error) {
	var cardIds []string
	err := c.db.DB.SelectContext(ctx, &cardIds, `
		SELECT card_id FROM assigned_cards WHERE user_id=$1;
	`, userId)
	if err != nil {
		return nil, err
	}

	err = c.UpdateBoardModifiedAt(ctx, boardId)
	if err != nil {
		return nil, err
	}

	return &cardIds, nil
}

// TODO - change select statement to a join to ensure the card is in the stack

// func (c *Controller) GetAssignedCardsByUserIdOnStack(ctx context.Context, stackId string, userId string) (*[]string, error) {
// 	var cardIds []string
// 	err := c.db.DB.SelectContext(ctx, &cardIds, `
// 		SELECT card_id FROM assigned_cards WHERE user_id=$1;
// 	`, userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &cardIds, nil
// }

func (c *Controller) GetCompleteCardById(ctx context.Context, cardId string) (*models.CompleteCard, error) {
	card, err := c.GetCardById(ctx, cardId)
	if err != nil {
		return nil, err
	}
	completeCard := models.CopyToCompleteCard(*card)
	completeCard.Assignments = make([]models.User, 0)
	assignments, err := c.GetAssignedUsersByCardId(ctx, card.Id.String())
	if err != nil {
		return nil, err
	}
	if len(*assignments) > 0 {
		for _, assignment := range *assignments {
			// completeCard.Assignments = append(completeCard.Assignments, assignment)
			assignedUser, err := user.GetUser(assignment)
			if err != nil {
				return nil, err
			}
			completeCard.Assignments = append(completeCard.Assignments, *assignedUser)
		}
	} else {
		completeCard.Assignments = make([]models.User, 0)
	}

	return &completeCard, nil
}
