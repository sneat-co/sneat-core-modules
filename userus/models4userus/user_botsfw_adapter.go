package models4userus

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

var _ botsfwmodels.AppUserAdapter = (*userBotsFwAdapter)(nil)

type userBotsFwAdapter struct {
	*UserDto
}

func (v *UserDto) BotsFwAdapter() botsfwmodels.AppUserAdapter {
	return &userBotsFwAdapter{UserDto: v}
}

func (u userBotsFwAdapter) SetNames(firstName, lastName, fullName string) error {
	if firstName == "" && lastName == "" && fullName == "" {
		return nil
	}
	if u.Name == nil {
		u.Name = new(dbmodels.Name)
	}
	if firstName != "" && u.Name.First == "" {
		u.Name.First = firstName
	}
	if lastName != "" && u.Name.Last == "" {
		u.Name.Last = lastName
	}
	if fullName != "" && u.Name.Full == "" {
		u.Name.Full = fullName
	}
	return nil
}
