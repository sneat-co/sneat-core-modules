package checks4contactus

import (
	"slices"

	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
)

func IsSpaceMember(roles []string) bool {
	return slices.Contains(roles, const4contactus.SpaceMemberRoleMember)
}
