package facade4linkage

import (
	"fmt"
	"strings"

	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// specscore: decisions/0002-reserved-extension-space-ids
// assertRelatedItemRefInSpace fails closed on the untrusted, user-facing
// relationship-update WRITE path. A user-supplied related ItemRef must resolve
// to the request's own authorized space. The "@{spaceID}" suffix introduced by
// sneat-specs Decision 0002 lets a ref name a DIFFERENT space ("@otherSpace")
// or the spaceless system namespace (a trailing "@" => empty space). Honouring
// either here would let a caller authorized only for spaceID write outside it,
// so both are rejected. Legitimate same-space refs (a bare itemID, or
// "itemID@{requestSpace}") are allowed and resolve to the request space exactly
// as before. Spaceless / cross-space writes are expected only from trusted
// server/worker paths under the deferred system-namespace ACL work, never from
// this endpoint.
// https://github.com/sneat-co/sneat-specs/blob/main/spec/decisions/0002-reserved-extension-space-ids.md
func assertRelatedItemRefInSpace(spaceID coretypes.SpaceID, ref dbo4linkage.ItemRef) error {
	// Use the explicit separator index (NOT ItemRef.SpaceID) so a trailing "@"
	// (empty space => spaceless /ext/) is detected as a space override too.
	if i := strings.Index(ref.ItemID, dbo4linkage.SpaceItemIDSeparator); i >= 0 {
		targetSpace := coretypes.SpaceID(ref.ItemID[i+1:])
		if targetSpace != spaceID {
			return fmt.Errorf(
				"%w: related item ref %q resolves to space %q outside the authorized request space %q",
				facade.ErrUnauthorized, ref.ItemID, targetSpace, spaceID)
		}
	}
	return nil
}
