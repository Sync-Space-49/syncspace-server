package board

func CopyToCompleteBoard(source Board, dest *CompleteBoard) {
	dest.Id = source.Id
	dest.Title = source.Title
	dest.OwnerId = source.OwnerId
	dest.CreatedAt = source.CreatedAt
	dest.ModifiedAt = source.ModifiedAt
	dest.IsPrivate = source.IsPrivate
}

func CopyToCompletePanel(source Panel, dest *CompletePanel) {
	dest.Id = source.Id
	dest.Title = source.Title
	dest.Postition = source.Position
	dest.BoardId = source.BoardId
}

func CopyToCompleteStack(source Stack, dest *CompleteStack) {
	dest.Id = source.Id
	dest.Title = source.Title
	dest.Postition = source.Position
	dest.PanelId = source.PanelId
}
