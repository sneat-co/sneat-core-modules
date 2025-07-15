package dbo4invitus

type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusSending  InviteStatus = "sending"
	InviteStatusSent     InviteStatus = "sent"
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusDeclined InviteStatus = "declined"
	InviteStatusRevoked  InviteStatus = "revoked"
	InviteStatusExpired  InviteStatus = "expired"
)

func IsKnownInviteStatus(s InviteStatus) bool {
	return s == InviteStatusPending ||
		s == InviteStatusSending ||
		s == InviteStatusSent ||
		s == InviteStatusAccepted ||
		s == InviteStatusDeclined ||
		s == InviteStatusRevoked ||
		s == InviteStatusExpired
}
