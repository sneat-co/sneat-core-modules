package dbo4contactus

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
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
		want dbo4linkage.SpaceModuleItemRef
	}{
		{
			name: "normal",
			args: args{
				spaceID:   "test-space-id",
				contactID: "test-contact-id",
			},
			want: dbo4linkage.SpaceModuleItemRef{
				Module:     "contactus",
				Collection: "contacts",
				ItemID:     "test-contact-id",
			},
		},
		{
			name: "panic on empty space",
			args: args{
				spaceID:   "",
				contactID: "test-contact-id",
			},
			want: dbo4linkage.SpaceModuleItemRef{
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
			want: dbo4linkage.SpaceModuleItemRef{
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
						t.Errorf("NewContactFullRef() did not panic on empty contactID")
					}
				}()
			}
			if tt.args.spaceID == "" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewContactFullRef() did not panic on empty spaceID")
					}
				}()
			}
			if got := NewContactFullRef(tt.args.spaceID, tt.args.contactID); got != tt.want {
				t.Errorf("NewContactFullRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
