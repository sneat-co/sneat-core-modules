package facade4linkage

import (
	"errors"
	"testing"

	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

// TestAssertRelatedItemRefInSpace asserts the B1 authorization boundary: on the
// untrusted user-facing relationship-update write path a related ItemRef must
// stay inside the request's authorized space. A cross-space ("@otherSpace") or
// spaceless (trailing "@") ItemRef must be REJECTED, while bare and
// same-space ("@requestSpace") refs are allowed. See sneat-specs Decision 0002.
func TestAssertRelatedItemRefInSpace(t *testing.T) {
	const requestSpace coretypes.SpaceID = "space1"

	newRef := func(itemID string) dbo4linkage.ItemRef {
		return dbo4linkage.ItemRef{ExtID: "contactus", Collection: "contacts", ItemID: itemID}
	}

	tests := []struct {
		name      string
		ref       dbo4linkage.ItemRef
		wantError bool
	}{
		{name: "bare_same_space", ref: newRef("victim"), wantError: false},
		{name: "explicit_same_space", ref: newRef("victim@space1"), wantError: false},
		{name: "cross_space_rejected", ref: newRef("victim@space2"), wantError: true},
		{name: "spaceless_trailing_at_rejected", ref: newRef("victim@"), wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := assertRelatedItemRefInSpace(requestSpace, tt.ref)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected an authorization error, got nil")
				}
				if !errors.Is(err, facade.ErrUnauthorized) {
					t.Errorf("expected facade.ErrUnauthorized, got %v", err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
