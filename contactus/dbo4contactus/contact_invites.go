package dbo4contactus

import (
	"github.com/dal-go/dalgo/dal"
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

func (v WithInvitesToContactBriefs) GetInviteBriefByChannelAndInviterUserID(channel dbo4invitus.InviteChannel, creatorUserID string) (id string, brief *InviteToContactBrief) {
	var b InviteToContactBrief
	for id, b = range v.Invites {
		if b.CreatedByUserID == creatorUserID && b.Channel == channel {
			return id, &b
		}
	}
	return "", nil
}

func (v WithInvitesToContactBriefs) DeleteInviteBrief(id string) dal.Update {
	delete(v.Invites, id)
	return dal.Update{
		Field: "invites." + id,
		Value: dal.DeleteField,
	}
}

func (v WithInvitesToContactBriefs) AddInviteBrief(inviteID, createdByUserID string, channel dbo4invitus.InviteChannel, createdTime time.Time) dal.Update {
	brief := InviteToContactBrief{
		Channel:         channel,
		CreatedTime:     createdTime,
		CreatedByUserID: createdByUserID,
	}
	v.Invites[inviteID] = brief
	return dal.Update{
		Field: "invites." + inviteID,
		Value: brief,
	}
}
