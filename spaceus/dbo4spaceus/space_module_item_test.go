package dbo4spaceus

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

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
