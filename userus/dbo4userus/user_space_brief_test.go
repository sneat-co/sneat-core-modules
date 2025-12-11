package dbo4userus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"strings"
	"testing"
)

// TestUserSpaceBrief_HasRole tests the HasRole method of UserSpaceBrief
func TestUserSpaceBrief_HasRole(t *testing.T) {
	tests := []struct {
		name     string
		roles    []string
		role     string
		expected bool
	}{
		{
			name:     "has_owner_role",
			roles:    []string{const4contactus.SpaceMemberRoleOwner, const4contactus.SpaceMemberRoleAdmin},
			role:     const4contactus.SpaceMemberRoleOwner,
			expected: true,
		},
		{
			name:     "has_admin_role",
			roles:    []string{const4contactus.SpaceMemberRoleOwner, const4contactus.SpaceMemberRoleAdmin},
			role:     const4contactus.SpaceMemberRoleAdmin,
			expected: true,
		},
		{
			name:     "does_not_have_role",
			roles:    []string{const4contactus.SpaceMemberRoleOwner, const4contactus.SpaceMemberRoleAdmin},
			role:     const4contactus.SpaceMemberRoleMember,
			expected: false,
		},
		{
			name:     "empty_roles",
			roles:    []string{},
			role:     const4contactus.SpaceMemberRoleOwner,
			expected: false,
		},
		{
			name:     "nil_roles",
			roles:    nil,
			role:     const4contactus.SpaceMemberRoleOwner,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a UserSpaceBrief with the test roles
			brief := UserSpaceBrief{
				Roles: tt.roles,
			}

			// Call the HasRole method
			result := brief.HasRole(tt.role)

			// Check the result
			if result != tt.expected {
				t.Errorf("HasRole(%s) = %v, want %v", tt.role, result, tt.expected)
			}
		})
	}
}

// TestUserSpaceBrief_Validate tests the Validate method of UserSpaceBrief
func TestUserSpaceBrief_Validate(t *testing.T) {
	tests := []struct {
		name        string
		brief       UserSpaceBrief
		wantErr     bool
		errContains string
	}{
		{
			name: "valid_brief",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   coretypes.SpaceTypeFamily,
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "contact123",
				Roles:         []string{const4contactus.SpaceMemberRoleOwner},
			},
			wantErr: false,
		},
		{
			name: "missing_user_contact_id",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   coretypes.SpaceTypeFamily,
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "",
				Roles:         []string{const4contactus.SpaceMemberRoleOwner},
			},
			wantErr:     true,
			errContains: "userContactID",
		},
		{
			name: "missing_type",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   "",
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "contact123",
				Roles:         []string{const4contactus.SpaceMemberRoleOwner},
			},
			wantErr:     true,
			errContains: "type",
		},
		{
			name: "invalid_type",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   "invalid_type",
					Title:  "Invalid Space",
					Status: "active",
				},
				UserContactID: "contact123",
				Roles:         []string{const4contactus.SpaceMemberRoleOwner},
			},
			wantErr:     true,
			errContains: "unknown space type",
		},
		{
			name: "empty_roles",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   coretypes.SpaceTypeFamily,
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "contact123",
				Roles:         []string{},
			},
			wantErr:     true,
			errContains: "roles",
		},
		{
			name: "nil_roles",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   coretypes.SpaceTypeFamily,
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "contact123",
				Roles:         nil,
			},
			wantErr:     true,
			errContains: "roles",
		},
		{
			name: "empty_role_in_roles",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   coretypes.SpaceTypeFamily,
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "contact123",
				Roles:         []string{const4contactus.SpaceMemberRoleOwner, ""},
			},
			wantErr:     true,
			errContains: "roles[1]",
		},
		{
			name: "unknown_role",
			brief: UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   coretypes.SpaceTypeFamily,
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "contact123",
				Roles:         []string{const4contactus.SpaceMemberRoleOwner, "unknown_role"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the Validate method
			err := tt.brief.Validate()

			// Check if an error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If an error was expected, check that it contains the expected string
			if err != nil && tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Validate() error = %v, want error containing %v", err, tt.errContains)
			}
		})
	}
}
