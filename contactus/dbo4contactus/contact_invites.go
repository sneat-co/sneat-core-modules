package dbo4contactus

import (
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/strongo/validation"
	"time"
)

type InviteToContactBrief struct {
	Channel dbo4invitus.InviteChannel `json:"channel" firestore:"channel"`
	//
	CreatedByUserID string    `json:"createdByUserID" firestore:"createdByUserID"`
	CreatedTime     time.Time `json:"createdTime" firestore:"createdTime"`
}

func (v InviteToContactBrief) Validate() error {
	if v.CreatedByUserID == "" {
		return validation.NewErrRecordIsMissingRequiredField("createdByUserID")
	}
	if v.CreatedTime.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("createdTime")
	}
	return nil
}

type WithInvitesToContactBriefs struct {
	Invites map[string]InviteToContactBrief `json:"invites,omitempty" firestore:"invites,omitempty"`
}

func (v WithInvitesToContactBriefs) Validate() error {
	for id, brief := range v.Invites {
		if err := brief.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("invites["+id+"]", err.Error())
		}
	}
	return nil
}

func (v WithInvitesToContactBriefs) GetInviteBriefByChannelAndToContactID(channel dbo4invitus.InviteChannel) (id string, brief *InviteToContactBrief) {
	var b InviteToContactBrief
	for id, b = range v.Invites {
		if b.Channel == channel {
			return id, &b
		}
	}
	return "", nil
}
