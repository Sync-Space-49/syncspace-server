package models

func CopyToCompleteBoard(source Board) CompleteBoard {
	dest := CompleteBoard{}
	dest.Id = source.Id
	dest.Description = source.Description
	dest.Title = source.Title
	dest.OwnerId = source.OwnerId
	dest.CreatedAt = source.CreatedAt
	dest.ModifiedAt = source.ModifiedAt
	dest.IsPrivate = source.IsPrivate
	return dest
}

func CopyToCompletePanel(source Panel) CompletePanel {
	dest := CompletePanel{}
	dest.Id = source.Id
	dest.Title = source.Title
	dest.Position = source.Position
	dest.BoardId = source.BoardId
	return dest
}

func CopyToCompleteStack(source Stack) CompleteStack {
	dest := CompleteStack{}
	dest.Id = source.Id
	dest.Title = source.Title
	dest.Position = source.Position
	dest.PanelId = source.PanelId
	return dest
}

func CopyToCompleteCard(source Card) CompleteCard {
	dest := CompleteCard{}
	dest.Id = source.Id
	dest.Title = source.Title
	dest.Description = source.Description
	dest.Points = source.Points
	dest.Position = source.Position
	dest.StackId = source.StackId
	return dest
}

func CopyToSimplifiedCompleteBoard(source CompleteBoard) SimplifiedCompleteBoard {
	dest := SimplifiedCompleteBoard{}
	dest.Title = source.Title
	dest.Description = source.Description
	dest.Panels = make([]SimplifiedCompletePanel, len(source.Panels))
	for i, panel := range source.Panels {
		dest.Panels[i] = CopyToSimplifiedCompletePanel(panel)
	}
	return dest
}

func CopyToSimplifiedCompletePanel(source CompletePanel) SimplifiedCompletePanel {
	dest := SimplifiedCompletePanel{}
	dest.Title = source.Title
	dest.Stacks = make([]SimplifiedCompleteStack, len(source.Stacks))
	for i, stack := range source.Stacks {
		dest.Stacks[i] = CopyToSimplifiedCompleteStack(stack)
	}
	return dest
}

func CopyToSimplifiedCompleteStack(source CompleteStack) SimplifiedCompleteStack {
	dest := SimplifiedCompleteStack{}
	dest.Cards = make([]AIGeneratedCard, len(source.Cards))
	for i, card := range source.Cards {
		dest.Cards[i] = CopyToAIGeneratedCard(card)
	}
	return dest
}

func CopyToAIGeneratedCard(source CompleteCard) AIGeneratedCard {
	dest := AIGeneratedCard{}
	dest.CardTitle = source.Title
	dest.CardDesc = source.Description
	dest.CardStoryPoints = source.Points
	return dest
}
