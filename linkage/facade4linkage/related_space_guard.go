package facade4linkage

import (
	"fmt"
	"strings"

	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/acl/const4acl"
	"github.com/sneat-co/sneat-go-core/acl/dbo4acl"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// specscore: decisions/0002-reserved-extension-space-ids
// classifyRelatedItemRef gates the untrusted, user-facing relationship-update
// WRITE path. A user-supplied related ItemRef resolves either to the request's
// own authorized space or, via the "@{spaceID}" suffix from sneat-specs
// Decision 0002, to another space or the spaceless system namespace (a trailing
// "@" => empty space). It returns the resolved target space and whether that is
// the spaceless /ext/ namespace:
//   - a DIFFERENT, non-empty space ("@otherSpace") is rejected: the caller is
//     authorized only for spaceID and must not write outside it (B1);
//   - the spaceless namespace (target == "") is NOT authorized here — the caller
//     MUST enforce the record's per-record ACL after loading it
//     (see authorizeSpacelessRelatedWrite), because /ext/ is a storage location,
//     not an authorization scope;
//   - a same-space ref (bare itemID, or "itemID@{requestSpace}") is allowed and
//     resolves to the request space exactly as before.
//
// Gating on the RESOLVED target space (not just the presence of "@") also closes
// the empty-request-space case: a bare ref under an empty spaceID resolves to the
// spaceless namespace and is therefore routed through the per-record ACL too.
// https://github.com/sneat-co/sneat-specs/blob/main/spec/decisions/0002-reserved-extension-space-ids.md
func classifyRelatedItemRef(spaceID coretypes.SpaceID, ref dbo4linkage.ItemRef) (spaceless bool, err error) {
	// Use the explicit separator index (NOT ItemRef.SpaceID) so a trailing "@"
	// (empty space => spaceless /ext/) is detected as a space override too.
	target := spaceID
	if i := strings.Index(ref.ItemID, dbo4linkage.SpaceItemIDSeparator); i >= 0 {
		target = coretypes.SpaceID(ref.ItemID[i+1:])
	}
	if target == "" {
		return true, nil // spaceless /ext/ — authorization deferred to per-record ACL
	}
	if target != spaceID {
		return false, fmt.Errorf(
			"%w: related item ref %q resolves to space %q outside the authorized request space %q",
			facade.ErrUnauthorized, ref.ItemID, target, spaceID)
	}
	return false, nil
}

// specscore: decisions/0002-reserved-extension-space-ids
// authorizeSpacelessRelatedWrite enforces per-record access control for a write
// to a spaceless /ext/ record. Authorization comes solely from the record's own
// ACL: the caller must hold an explicit "edit" grant. A record with no ACL (or a
// caller without an edit grant) is denied — there is no blanket "any
// authenticated user may write /ext/".
// https://github.com/sneat-co/sneat-specs/blob/main/spec/decisions/0002-reserved-extension-space-ids.md
func authorizeSpacelessRelatedWrite(userID string, acl *dbo4acl.ACL) error {
	if acl == nil || !acl.UserCan(userID, const4acl.PermittedToEdit) {
		return fmt.Errorf(
			"%w: caller %q is not granted edit on the spaceless /ext/ record",
			facade.ErrUnauthorized, userID)
	}
	return nil
}
