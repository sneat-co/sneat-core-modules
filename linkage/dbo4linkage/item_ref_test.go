package dbo4linkage

import "testing"

// TestItemRef_Validate covers the "@" reserved-separator rules from sneat-specs
// Decision 0002: a bare itemID and a well-formed "itemID@{spaceID}" composite
// are valid; a trailing "@" (empty space), a leading "@" (empty doc id), and
// more than one "@" (an embedded separator in the document id) are rejected.
func TestItemRef_Validate(t *testing.T) {
	newRef := func(itemID string) ItemRef {
		return ItemRef{ExtID: "contactus", Collection: "contacts", ItemID: itemID}
	}
	tests := []struct {
		name    string
		ref     ItemRef
		wantErr bool
	}{
		{name: "bare", ref: newRef("item1"), wantErr: false},
		{name: "space_bound_composite", ref: newRef("item1@space1"), wantErr: false},
		{name: "trailing_at_empty_space", ref: newRef("item1@"), wantErr: true},
		{name: "leading_at_empty_doc", ref: newRef("@space1"), wantErr: true},
		{name: "double_at", ref: newRef("item1@space1@x"), wantErr: true},
		{name: "missing_item_id", ref: newRef(""), wantErr: true},
		{name: "missing_ext_id", ref: ItemRef{Collection: "contacts", ItemID: "item1"}, wantErr: true},
		{name: "missing_collection", ref: ItemRef{ExtID: "contactus", ItemID: "item1"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ref.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ItemRef.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestItemRef_DocID verifies stripping of the optional "@{spaceID}" suffix.
func TestItemRef_DocID(t *testing.T) {
	tests := []struct {
		itemID string
		want   string
	}{
		{itemID: "item1", want: "item1"},
		{itemID: "item1@space1", want: "item1"},
		{itemID: "item1@", want: "item1"},
	}
	for _, tt := range tests {
		ref := ItemRef{ExtID: "contactus", Collection: "contacts", ItemID: tt.itemID}
		if got := ref.DocID(); got != tt.want {
			t.Errorf("ItemRef.DocID() for %q = %q, want %q", tt.itemID, got, tt.want)
		}
	}
}
