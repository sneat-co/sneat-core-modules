package facade4invitus

import (
	"testing"
)

func TestNewMassInviteEntry(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{id: "id1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			massInvite := NewMassInviteEntry(tt.args.id)
			if massInvite.ID != tt.args.id {
				t.Errorf("NewMassInviteEntry().ID = %v, want %v", massInvite.ID, tt.args.id)
			}
			if massInvite.Data == nil {
				t.Errorf("NewMassInviteEntry().Data = nil, want not nil")
			}
			if massInvite.Record == nil {
				t.Errorf("NewMassInviteEntry().Record = nil, want not nil")
			}
			if massInvite.Key == nil {
				t.Errorf("NewMassInviteEntry().Key = nil, want not nil")
			}
		})
	}
}
