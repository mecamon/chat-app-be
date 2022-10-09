package presenters

import "github.com/mecamon/chat-app-be/models"

func FormatGroups(groups []models.GroupChat) []models.GroupChatDTO {
	var groupsFormatted []models.GroupChatDTO

	for _, g := range groups {
		gf := models.GroupChatDTO{
			ID:           g.ID.Hex(),
			Name:         g.Name,
			Description:  g.Description,
			ImageURL:     g.ImageURL,
			GroupOwner:   g.GroupOwner.Hex(),
			Participants: g.Participants,
		}
		groupsFormatted = append(groupsFormatted, gf)
	}

	return groupsFormatted
}
