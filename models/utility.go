package models

func CopyToCompleteBoard(source Board) CompleteBoard {
	dest := CompleteBoard{}
	dest.Id = source.Id
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
