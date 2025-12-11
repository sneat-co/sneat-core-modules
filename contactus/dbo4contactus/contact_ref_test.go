package dbo4contactus

import (
	"reflect"
	"testing"

	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestNewContactFullRef(t *testing.T) {
	type args struct {
		spaceID   coretypes.SpaceID
		contactID string
	}
	tests := []struct {
		name string
		args args
		want dbo4linkage.ItemRef
	}{
		{
			name: "should create a full contact reference",
			args: args{
				spaceID:   "space_1",
				contactID: "contact_1",
			},
			want: dbo4linkage.NewFullItemRef(const4contactus.ExtensionID, const4contactus.ContactsCollection, "space_1", "contact_1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContactFullRef(tt.args.spaceID, tt.args.contactID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContactFullRef() = %v, want %v", got, tt.want)
			}
		})
	}
}
