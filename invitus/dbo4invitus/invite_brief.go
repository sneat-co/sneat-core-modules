package dbo4invitus

import "github.com/strongo/validation"

// InviteBrief summary about invite
type InviteBrief struct {
	ID   string      `json:"id" firestore:"id"`
	Pin  string      `json:"pin,omitempty" firestore:"pin,omitempty"`
	From *InviteFrom `json:"from,omitempty" firestore:"from,omitempty"`
	To   *InviteTo   `json:"to,omitempty" firestore:"to,omitempty"`
}

// Validate returns error if not valid
func (v InviteBrief) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("from", err.Error())
	}
	if v.To != nil {
		if err := v.To.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue("to", err.Error())
		}
	}
	return nil
}

// NewInviteBriefFromDbo creates brief from DTO
func NewInviteBriefFromDbo(id string, dto InviteDbo) InviteBrief {
	from := dto.From
	to := *dto.To
	return InviteBrief{ID: id, From: &from, To: &to, Pin: dto.Pin}
}
