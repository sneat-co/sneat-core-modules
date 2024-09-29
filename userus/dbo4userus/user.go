package dbo4userus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
)

// UsersCollection a name of the user's db table
const UsersCollection = "users"

// UserDefaults keeps user's defaults
type UserDefaults struct {
	ShortNames []briefs4contactus.ShortName `json:"shortNames,omitempty" firestore:"shortNames,omitempty"`
}
