package dbo4userus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"reflect"
	"strings"
	"testing"
	"time"
)

// createTestUserDbo creates a UserDbo instance for testing
func createTestUserDbo() *UserDbo {
	now := time.Now()
	userDbo := &UserDbo{
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: now,
			},
			CreatedByField: with.CreatedByField{
				CreatedBy: "user",
			},
		},
		ContactBase: briefs4contactus.ContactBase{
			ContactBrief: briefs4contactus.ContactBrief{
				Type:   briefs4contactus.ContactTypePerson,
				Gender: "unknown",
				Names: &person.NameFields{
					FirstName: "Firstname",
					LastName:  "Lastname",
					FullName:  "Firstname Lastname",
				},
				AgeGroup: "unknown",
			},
			Status: "active",
		},
		Created: dbmodels.CreatedInfo{
			Client: dbmodels.RemoteClientInfo{
				HostOrApp:  "unit-test",
				RemoteAddr: "127.0.0.1",
			},
		},
		WithPreferredLocale: dbmodels.WithPreferredLocale{
			PreferredLocale: "en-US",
		},
	}
	return userDbo
}

// TestUserDbo_GetUserSpaceInfoByID tests the GetUserSpaceInfoByID method
func TestUserDbo_GetUserSpaceInfoByID(t *testing.T) {
	tests := []struct {
		name     string
		spaces   map[string]*UserSpaceBrief
		spaceID  coretypes.SpaceID
		expected *UserSpaceBrief
	}{
		{
			name:     "nil_spaces",
			spaces:   nil,
			spaceID:  "space1",
			expected: nil,
		},
		{
			name: "space_exists",
			spaces: map[string]*UserSpaceBrief{
				"space1": {
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:   "family",
						Title:  "Family Space",
						Status: "active",
					},
					UserContactID: "contact1",
					Roles:         []string{const4contactus.SpaceMemberRoleOwner},
				},
			},
			spaceID: "space1",
			expected: &UserSpaceBrief{
				SpaceBrief: dbo4spaceus.SpaceBrief{
					Type:   "family",
					Title:  "Family Space",
					Status: "active",
				},
				UserContactID: "contact1",
				Roles:         []string{const4contactus.SpaceMemberRoleOwner},
			},
		},
		{
			name: "space_does_not_exist",
			spaces: map[string]*UserSpaceBrief{
				"space1": {
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:   "family",
						Title:  "Family Space",
						Status: "active",
					},
					UserContactID: "contact1",
					Roles:         []string{const4contactus.SpaceMemberRoleOwner},
				},
			},
			spaceID:  "space2",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createTestUserDbo()
			u.Spaces = tt.spaces

			result := u.GetUserSpaceInfoByID(tt.spaceID)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetUserSpaceInfoByID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestUserDbo_GetFullName tests the GetFullName method
func TestUserDbo_GetFullName(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*UserDbo)
		expected string
	}{
		{
			name: "with_full_name",
			setup: func(u *UserDbo) {
				u.Names.FullName = "John Doe"
			},
			expected: "John Doe",
		},
		{
			name: "with_first_and_last_name",
			setup: func(u *UserDbo) {
				u.Names.FirstName = "John"
				u.Names.LastName = "Doe"
				u.Names.FullName = ""
			},
			expected: "John Doe",
		},
		{
			name: "with_first_name_only",
			setup: func(u *UserDbo) {
				u.Names.FirstName = "John"
				u.Names.LastName = ""
				u.Names.FullName = ""
			},
			expected: "John",
		},
		{
			name: "with_last_name_only",
			setup: func(u *UserDbo) {
				u.Names.FirstName = ""
				u.Names.LastName = "Doe"
				u.Names.FullName = ""
			},
			expected: "Doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createTestUserDbo()
			tt.setup(u)

			result := u.GetFullName()

			if result != tt.expected {
				t.Errorf("GetFullName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestUserDbo_validateSpaces tests the validateSpaces method
func TestUserDbo_validateSpaces(t *testing.T) {
	tests := []struct {
		name           string
		spaces         map[string]*UserSpaceBrief
		spaceIDs       []string
		defaultSpaceID string
		wantErr        bool
		errMsg         string
	}{
		{
			name:     "empty_spaces",
			spaces:   nil,
			spaceIDs: nil,
			wantErr:  false,
		},
		{
			name: "valid_spaces",
			spaces: map[string]*UserSpaceBrief{
				"space1": {
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:   "family",
						Title:  "Family Space",
						Status: "active",
					},
					UserContactID: "contact1",
					Roles:         []string{const4contactus.SpaceMemberRoleOwner},
				},
				"space2": {
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:   "private",
						Title:  "Private Space",
						Status: "active",
					},
					UserContactID: "contact2",
					Roles:         []string{const4contactus.SpaceMemberRoleAdmin},
				},
			},
			spaceIDs:       []string{"space1", "space2"},
			defaultSpaceID: "space1",
			wantErr:        false,
		},
		{
			name: "spaces_length_mismatch",
			spaces: map[string]*UserSpaceBrief{
				"space1": {
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:   "family",
						Title:  "Family Space",
						Status: "active",
					},
					UserContactID: "contact1",
					Roles:         []string{const4contactus.SpaceMemberRoleOwner},
				},
			},
			spaceIDs: []string{"space1", "space2"},
			wantErr:  true,
			errMsg:   "len(v.Spaces) != len(v.SpaceIDs)",
		},
		{
			name: "default_space_id_not_in_spaceIDs",
			spaces: map[string]*UserSpaceBrief{
				"space1": {
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:   "family",
						Title:  "Family Space",
						Status: "active",
					},
					UserContactID: "contact1",
					Roles:         []string{const4contactus.SpaceMemberRoleOwner},
				},
			},
			spaceIDs:       []string{"space1"},
			defaultSpaceID: "space2",
			wantErr:        true,
			errMsg:         "not in spaceIDs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := createTestUserDbo()
			u.Spaces = tt.spaces
			u.SpaceIDs = tt.spaceIDs
			u.DefaultSpaceID = tt.defaultSpaceID

			err := u.validateSpaces()

			if (err != nil) != tt.wantErr {
				t.Errorf("validateSpaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateSpaces() error message = %v, want to contain %v", err.Error(), tt.errMsg)
			}
		})
	}
}
