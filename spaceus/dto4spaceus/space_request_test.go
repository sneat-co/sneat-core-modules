package dto4spaceus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestNewSpaceRequest(t *testing.T) {
	type args struct {
		spaceID coretypes.SpaceID
	}
	tests := []struct {
		name string
		args args
		want SpaceRequest
	}{
		{
			name: "normal",
			args: args{
				spaceID: "test-space-id",
			},
			want: SpaceRequest{SpaceID: "test-space-id"},
		},
		{
			name: "panics",
			args: args{
				spaceID: "",
			},
			want: SpaceRequest{SpaceID: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want.SpaceID == "" {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewSpaceRequest() did not panic")
					}
				}()
			}
			if got := NewSpaceRequest(tt.args.spaceID); got.SpaceID != tt.want.SpaceID {
				t.Errorf("NewSpaceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
