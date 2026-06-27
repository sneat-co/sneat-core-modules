package dto4spaceus

import (
	"testing"

	"github.com/sneat-co/sneat-go-core/coretypes"
)

func TestCreateSpaceRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request CreateSpaceRequest
		wantErr bool
	}{
		{name: "family_ok", request: CreateSpaceRequest{Type: coretypes.SpaceTypeFamily}, wantErr: false},
		{name: "company_with_title_ok", request: CreateSpaceRequest{Type: coretypes.SpaceTypeCompany, Title: "Acme"}, wantErr: false},
		{name: "empty_type_rejected", request: CreateSpaceRequest{Type: ""}, wantErr: true},
		{name: "system_rejected", request: CreateSpaceRequest{Type: coretypes.SpaceTypeSystem, Title: "Games"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
