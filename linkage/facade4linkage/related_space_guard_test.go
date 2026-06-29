package facade4linkage

import (
	"errors"
	"testing"
	"time"

	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/acl/const4acl"
	"github.com/sneat-co/sneat-go-core/acl/dbo4acl"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/strongoapp/with"
)

// TestClassifyRelatedItemRef asserts the authorization boundary on the untrusted
// user-facing relationship-update write path. A cross-space ("@otherSpace") ref
// is REJECTED (B1). A spaceless (/ext/) ref — a trailing "@", or a bare ref under
// an empty request space — is classified spaceless (authorization is then
// deferred to the per-record ACL). Bare and "@requestSpace" refs are same-space.
// See sneat-specs Decision 0002.
func TestClassifyRelatedItemRef(t *testing.T) {
	newRef := func(itemID string) dbo4linkage.ItemRef {
		return dbo4linkage.ItemRef{ExtID: "contactus", Collection: "contacts", ItemID: itemID}
	}

	tests := []struct {
		name          string
		requestSpace  coretypes.SpaceID
		ref           dbo4linkage.ItemRef
		wantSpaceless bool
		wantError     bool
	}{
		{name: "bare_same_space", requestSpace: "space1", ref: newRef("victim"), wantSpaceless: false, wantError: false},
		{name: "explicit_same_space", requestSpace: "space1", ref: newRef("victim@space1"), wantSpaceless: false, wantError: false},
		{name: "cross_space_rejected", requestSpace: "space1", ref: newRef("victim@space2"), wantSpaceless: false, wantError: true},
		{name: "spaceless_trailing_at", requestSpace: "space1", ref: newRef("victim@"), wantSpaceless: true, wantError: false},
		{name: "empty_request_bare_is_spaceless", requestSpace: "", ref: newRef("victim"), wantSpaceless: true, wantError: false},
		{name: "empty_request_cross_space_rejected", requestSpace: "", ref: newRef("victim@space2"), wantSpaceless: false, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spaceless, err := classifyRelatedItemRef(tt.requestSpace, tt.ref)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected an authorization error, got nil")
				}
				if !errors.Is(err, facade.ErrUnauthorized) {
					t.Errorf("expected facade.ErrUnauthorized, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if spaceless != tt.wantSpaceless {
				t.Errorf("spaceless = %v, want %v", spaceless, tt.wantSpaceless)
			}
		})
	}
}

// TestAuthorizeSpacelessRelatedWrite asserts per-record ACL enforcement for a
// spaceless /ext/ write: only a caller holding an explicit edit grant in the
// record's own ACL is permitted; a record with no ACL, or a caller without an
// edit grant, is rejected — there is no blanket "/ext/ is writable" rule.
// See sneat-specs feature reserved-extension-space-ids (REQ:record-level-acl).
func TestAuthorizeSpacelessRelatedWrite(t *testing.T) {
	grant := func(by string) with.CreatedFields {
		return with.CreatedFields{
			CreatedAtField: with.CreatedAtField{CreatedAt: time.Now()},
			CreatedByField: with.CreatedByField{CreatedBy: by},
		}
	}
	aclWith := func(users map[string]dbo4acl.Permissions) *dbo4acl.ACL {
		return &dbo4acl.ACL{Users: users}
	}

	tests := []struct {
		name      string
		userID    string
		acl       *dbo4acl.ACL
		wantError bool
	}{
		{
			name:   "editor_permitted",
			userID: "editor",
			acl: aclWith(map[string]dbo4acl.Permissions{
				"editor": {const4acl.PermittedToEdit: grant("owner")},
			}),
			wantError: false,
		},
		{
			name:   "viewer_rejected",
			userID: "viewer",
			acl: aclWith(map[string]dbo4acl.Permissions{
				"viewer": {const4acl.PermittedToView: grant("owner")},
			}),
			wantError: true,
		},
		{
			name:   "stranger_rejected",
			userID: "stranger",
			acl: aclWith(map[string]dbo4acl.Permissions{
				"editor": {const4acl.PermittedToEdit: grant("owner")},
			}),
			wantError: true,
		},
		{name: "nil_acl_rejected", userID: "editor", acl: nil, wantError: true},
		{name: "empty_user_rejected", userID: "", acl: aclWith(map[string]dbo4acl.Permissions{"editor": {const4acl.PermittedToEdit: grant("owner")}}), wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authorizeSpacelessRelatedWrite(tt.userID, tt.acl)
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
