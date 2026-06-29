package dbo4spaceus

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// TestNewSpaceModuleItemKeyFromItemRef verifies both storage shapes defined by
// sneat-specs Decision 0002 (spaceless system namespace):
//   - space-bound:       /spaces/{space-id}/ext/{ext-id}/{collection}/{item-id}
//   - system namespace:  /ext/{ext-id}/{collection}/{item-id}
//
// The presence/absence of an "@{space-id}" suffix on the itemID is the sole
// discriminator.
func TestNewSpaceModuleItemKeyFromItemRef(t *testing.T) {
	tests := []struct {
		name     string
		spaceID  coretypes.SpaceID
		itemRef  dbo4linkage.ItemRef
		wantPath string
	}{
		{
			name:     "space_bound",
			spaceID:  "space1",
			itemRef:  dbo4linkage.ItemRef{ExtID: "contactus", Collection: "contacts", ItemID: "item1"},
			wantPath: "spaces/space1/ext/contactus/contacts/item1",
		},
		{
			name:     "system_namespace_empty_space",
			spaceID:  "",
			itemRef:  dbo4linkage.ItemRef{ExtID: "contactus", Collection: "contacts", ItemID: "item1"},
			wantPath: "ext/contactus/contacts/item1",
		},
		{
			name:     "space_bound_via_at_suffix",
			spaceID:  "space1",
			itemRef:  dbo4linkage.ItemRef{ExtID: "contactus", Collection: "contacts", ItemID: "item1@space2"},
			wantPath: "spaces/space2/ext/contactus/contacts/item1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewSpaceModuleItemKeyFromItemRef(tt.spaceID, tt.itemRef)
			if got := key.String(); got != tt.wantPath {
				t.Errorf("NewSpaceModuleItemKeyFromItemRef() path = %q, want %q", got, tt.wantPath)
			}
		})
	}
}

// TestNewSpaceModuleKey verifies the space-bound and spaceless module keys.
func TestNewSpaceModuleKey(t *testing.T) {
	if got := NewSpaceModuleKey("space1", "contactus").String(); got != "spaces/space1/ext/contactus" {
		t.Errorf("NewSpaceModuleKey(space1) = %q", got)
	}
	if got := NewSpaceModuleKey("", "contactus").String(); got != "ext/contactus" {
		t.Errorf("NewSpaceModuleKey(empty) = %q", got)
	}
}

func TestNewSpaceModuleItemIncompleteKey(t *testing.T) {
	type args struct {
		spaceID    coretypes.SpaceID
		moduleID   coretypes.ExtID
		collection string
	}
	tests := []struct {
		name string
		args args
		//want *dal.Key
	}{
		{
			name: "test1",
			args: args{
				spaceID:    "space1",
				moduleID:   "module1",
				collection: "collection1",
			},
			//want: dal.NewIncompleteKey("collection1", reflect.TypeOf(""), dal.NewKeyWithParent("space1", "module1")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSpaceModuleItemIncompleteKey[string](tt.args.spaceID, tt.args.moduleID, tt.args.collection)
			if got == nil {
				t.Errorf("NewSpaceModuleItemIncompleteKey() = nil, want not nil")
				return
			}
		})
	}
}
