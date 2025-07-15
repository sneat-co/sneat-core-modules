package dbo4invitus

type InviteType string

const (
	InviteTypePersonal InviteType = "personal" // To a specific person
	InviteTypePrivate  InviteType = "private"  // To a single person
	InviteTypeMass     InviteType = "mass"     // To a group of people
)
