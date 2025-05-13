package dbo4userus

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-core-modules/dbo4all"

	//"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"net/mail"
	"slices"
	"strings"
	"time"
)

type WithUserIDs struct {
	UserIDs map[string]string `json:"userIDs,omitempty" firestore:"userIDs,omitempty"`
}

func (v *WithUserIDs) SetUserID(spaceID coretypes.SpaceID, userID string) {
	if v.UserIDs == nil {
		v.UserIDs = map[string]string{string(spaceID): userID}
	} else {
		v.UserIDs[string(spaceID)] = userID
	}
}

var _ botsfwmodels.AppUserData = (*UserDbo)(nil)
var _ botsfwmodels.AppUserData = (*userBotsFwAdapter)(nil)

// UserDbo is a record that holds information about user
type UserDbo struct {
	briefs4contactus.ContactBase
	with.CreatedFields
	dbmodels.WithPreferredLocale
	dbmodels.WithPrimaryCurrency
	dbmodels.WithLastCurrencies
	botsfwmodels.WithBotUserIDs
	coremodels.SmsStats

	appuser.AccountsOfUser
	appuser.WithLastLogin

	//dbo4linkage.WithRelatedAndIDs

	InvitedByUserID string `firestore:"invitedByUserID,omitempty" ` // TODO: Prevent circular references! see users 6032980589936640 & 5998019824582656

	IsAnonymous bool `json:"isAnonymous" firestore:"isAnonymous"` // Intentionally do not omitempty
	//Title string

	Timezone *dbmodels.Timezone `json:"timezone,omitempty" firestore:"timezone,omitempty"`

	Defaults *UserDefaults `json:"defaults,omitempty" firestore:"defaults,omitempty"`

	Email         string `json:"email,omitempty"  firestore:"email,omitempty"`
	EmailVerified bool   `json:"emailVerified"  firestore:"emailVerified"`

	dbo4all.WithEmails
	dbo4all.WithPhones

	// List of spaces a user belongs to
	Spaces map[string]*UserSpaceBrief `json:"spaces,omitempty"   firestore:"spaces,omitempty"`

	SpaceIDs       []string `json:"spaceIDs,omitempty" firestore:"spaceIDs,omitempty"`
	DefaultSpaceID string   `json:"defaultSpaceRef" firestore:"defaultSpaceRef"`

	Created dbmodels.CreatedInfo `json:"created" firestore:"created"`

	//models4debtus.WithGroups

	// TODO: Should this be moved to company members?
	//models.DatatugUser

	ReferredBy string `firestore:"referredBy,omitempty"`

	LastFeedbackAt   time.Time `firestore:"lastFeedbackAt,omitempty"`
	LastFeedbackRate string    `firestore:"lastFeedbackRate,omitempty"`
}

func (v *UserDbo) GetFullName() string {
	return v.Names.GetFullName()
}

// SetSpaceBrief sets space brief and adds spaceID to the list of space IDs if needed
func (v *UserDbo) SetSpaceBrief(spaceID coretypes.SpaceID, brief *UserSpaceBrief) (updates []update.Update) {
	if spaceID == "" {
		panic("spaceID is empty string")
	}
	if brief == nil {
		panic("brief is nil")
	}
	if v.Spaces == nil {
		v.Spaces = make(map[string]*UserSpaceBrief, 1)
	}
	v.Spaces[string(spaceID)] = brief
	updates = append(updates, update.ByFieldPath([]string{"spaces", string(spaceID)}, brief))
	if !slices.Contains(v.SpaceIDs, string(spaceID)) {
		v.SpaceIDs = append(v.SpaceIDs, string(spaceID))
		updates = append(updates, update.ByFieldName("spaceIDs", v.SpaceIDs))
	}
	return
}

func (v *UserDbo) GetFamilySpaceID() coretypes.SpaceID {
	id, _ := v.GetFirstSpaceBriefBySpaceType(coretypes.SpaceTypeFamily)
	return id
}

// GetSpaceBriefsByType returns the all spaces matching a specific type
func (v *UserDbo) GetSpaceBriefsByType(t coretypes.SpaceType) (spaces map[string]*UserSpaceBrief) {
	for id, brief := range v.Spaces {
		if brief.Type == t {
			if spaces == nil {
				spaces = make(map[string]*UserSpaceBrief)
			}
			spaces[id] = brief
		}
	}
	return
}

func (v *UserDbo) GetFirstSpaceBriefBySpaceType(spaceType coretypes.SpaceType) (spaceID coretypes.SpaceID, spaceBrief *UserSpaceBrief) {
	for id, space := range v.Spaces {
		if space.Type == spaceType {
			return coretypes.SpaceID(id), space
		}
	}
	return "", nil
}

// Validate validates user record
func (v *UserDbo) Validate() error {
	if err := v.ContactBase.Validate(); err != nil {
		return err
	}
	if err := v.SmsStats.Validate(); err != nil {
		return err
	}
	//if v.Avatar != nil {
	//	if err := v.Avatar.Validate(); err != nil {
	//		return validation.NewErrBadRecordFieldValue("avatar", err.Error())
	//	}
	//}
	//if v.Title != "" {
	//	if err := v.Names.Validate(); err != nil {
	//		return err
	//	}
	//}
	if err := v.WithEmails.Validate(); err != nil {
		return err
	}
	if err := v.WithPhones.Validate(); err != nil {
		return err
	}
	if err := v.validateEmails(); err != nil {
		return err
	}
	if err := v.validateSpaces(); err != nil {
		return err
	}
	if err := dbmodels.ValidateGender(v.Gender, true); err != nil {
		return err
	}
	//if v.Datatug != nil {
	//	if err := v.Datatug.Validate(); err != nil {
	//		return err
	//	}
	//}
	if err := v.Created.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("created", err.Error())
	}
	//if err := v.WithRelatedAndIDs.Validate(); err != nil {
	//	return err
	//}
	if err := v.WithBotUserIDs.Validate(); err != nil {
		return err
	}
	return nil
}

func (v *UserDbo) validateEmails() error {
	if strings.TrimSpace(v.Email) != v.Email {
		return validation.NewErrBadRecordFieldValue("email", "contains leading or closing spaces")
	}
	if strings.Contains(v.Email, " ") {
		return validation.NewErrBadRecordFieldValue("email", "contains space")
	}
	if v.Email != "" {
		if _, err := mail.ParseAddress(v.Email); err != nil {
			return validation.NewErrBadRecordFieldValue("email", err.Error())
		}
		if len(v.Emails) == 0 {
			return validation.NewErrBadRecordFieldValue("emails", "user record has 'email' value but 'emails' are empty")
		}
	}
	primaryEmailInEmails := false
	for emailAddress, emailProps := range v.Emails {
		if err := emailProps.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("emails[%s]", emailAddress), err.Error())
		}
		if emailAddress == v.Email {
			primaryEmailInEmails = true
		}
	}
	if v.Email != "" && !primaryEmailInEmails {
		return validation.NewErrBadRecordFieldValue("emails", "user's primary email is not in 'emails' field")
	}
	return nil
}

func (v *UserDbo) validateSpaces() error {
	if len(v.Spaces) != len(v.SpaceIDs) {
		return validation.NewErrBadRecordFieldValue("spaceIDs",
			fmt.Sprintf("len(v.Spaces) != len(v.SpaceIDs): %d != %d", len(v.Spaces), len(v.SpaceIDs)))
	}
	if len(v.Spaces) > 0 {
		spaceIDs := make([]coretypes.SpaceID, 0, len(v.Spaces))
		spaceTitles := make([]string, 0, len(v.Spaces))
		for spaceID, space := range v.Spaces {
			if spaceID == "" {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("spaces['%s']", spaceID), "holds empty id")
			}
			if !slices.Contains(v.SpaceIDs, spaceID) {
				return validation.NewErrBadRecordFieldValue("spaceIDs", "missing space ContactID: "+string(spaceID))
			}
			if err := space.Validate(); err != nil {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("spaces[%s]{title=%s}", spaceID, space.Title), err.Error())
			}
			if space.Title != "" {
				if i := slices.Index(spaceTitles, space.Title); i >= 0 {
					return validation.NewErrBadRecordFieldValue("spaces",
						fmt.Sprintf("at least 2 spaces (%s & %s) with same title: '%s'", spaceID, spaceIDs[i], space.Title))
				}
			}
			spaceIDs = append(spaceIDs, coretypes.SpaceID(spaceID))
			spaceTitles = append(spaceTitles, space.Title)
		}
	}
	if v.DefaultSpaceID != "" {
		if !slices.Contains(v.SpaceIDs, v.DefaultSpaceID) {
			return validation.NewErrBadRecordFieldValue("defaultSpaceID", "not in spaceIDs")
		}
	}
	return nil
}

// GetUserSpaceInfoByID return space info specific to the user by space ContactID
func (v *UserDbo) GetUserSpaceInfoByID(spaceID coretypes.SpaceID) *UserSpaceBrief {
	return v.Spaces[string(spaceID)]
}

func (v *UserDbo) SetBotUserID(platform const4userus.AuthProviderCode, botID, botUserID string) {
	v.AddAccount(appuser.AccountKey{
		Provider: platform,
		App:      botID,
		ID:       botUserID,
	})
}
