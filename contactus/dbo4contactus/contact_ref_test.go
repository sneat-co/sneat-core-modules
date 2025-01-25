package dbo4contactus

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"testing"
)

func TestNewContactFullRef(t *testing.T) {
	type args struct {
		spaceID   string
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
				Space:      "test-space-id",
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
				Space:      "",
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
				Space:      "test-space-id",
				Module:     "contactus",
				Collection: "contacts",
				ItemID:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want.Space == "" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewContactFullRef() did not panic")
					}
				}()
			}
			if tt.want.ItemID == "" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewContactFullRef() did not panic")
					}
				}()
			}
			if got := NewContactFullRef(tt.args.spaceID, tt.args.contactID); got != tt.want {
				t.Errorf("NewContactFullRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
