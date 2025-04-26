package dbo4linkage

import (
	"github.com/sneat-co/sneat-go-core/coretypes"

	"testing"
)

func TestNewContactFullRef(t *testing.T) {
	type args struct {
		spaceID   coretypes.SpaceID
		contactID string
	}
	tests := []struct {
		name string
		args args
		want ItemRef
	}{
		{
			name: "normal",
			args: args{
				spaceID:   "test-space-id",
				contactID: "test-contact-id",
			},
			want: ItemRef{
				Module:     "contactus",
				Collection: "contacts",
				ItemID:     "test-contact-id@test-space-id",
			},
		},
		{
			name: "panic on empty space",
			args: args{
				spaceID:   "",
				contactID: "test-contact-id",
			},
			want: ItemRef{
				Module:     "contactus",
				Collection: "contacts",
				ItemID:     "test-contact-id",
			},
		},
		{
			name: "panic on empty contactID",
			args: args{
				spaceID:   "test-space-id",
				contactID: "",
			},
			want: ItemRef{
				Module:     "contactus",
				Collection: "contacts",
				ItemID:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.contactID == "" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewFullItemRef() did not panic on empty itemID")
					}
				}()
			}
			if tt.args.spaceID == "" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewFullItemRef() did not panic on empty itemID")
					}
				}()
			}
			if got := NewFullItemRef("contactus", "contacts", tt.args.spaceID, tt.args.contactID); got != tt.want {
				t.Errorf("NewFullItemRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
