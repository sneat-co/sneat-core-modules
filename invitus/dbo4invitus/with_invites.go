package dbo4invitus

import (
	"fmt"
	"github.com/strongo/validation"
)

// WithInvites holds a list of active invites the member has created.
// Obsolete: use InviteChannel instead
type WithInvites struct {
	// Holds a list of active invites the member has created.
	Invites []*MemberInviteBrief `json:"invites" firestore:"invites,omitempty"`
}

func (v WithInvites) Validate() error {
	for i, mi := range v.Invites {
		if mi == nil {
			return validation.NewErrBadRecordFieldValue("invites", fmt.Sprintf("nil invite at index %d", i))
		}
		if err := mi.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("invites[%d]", i), err.Error())
		}
	}
	return nil
}

func (v WithInvites) GetInviteBriefByChannelAndToContactID(channel InviteChannel, toContactID string) *MemberInviteBrief {
	for _, mi := range v.Invites {
		if mi.To.Channel == channel && mi.To.ContactID == toContactID {
			return mi
		}
	}
	return nil
}
