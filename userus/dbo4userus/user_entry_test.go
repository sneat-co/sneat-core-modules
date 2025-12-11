package dbo4userus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
)

// TestUserEntry_GetID tests the GetID method of UserEntry
func TestUserEntry_GetID(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected string
	}{
		{
			name:     "normal_id",
			userID:   "user123",
			expected: "user123",
		},
		{
			name:     "empty_id",
			userID:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a UserEntry with the test ID
			userEntry := UserEntry{}
			userEntry.ID = tt.userID

			// Call the GetID method
			result := userEntry.GetID()

			// Check the result
			if result != tt.expected {
				t.Errorf("GetID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNewUserEntry tests the NewUserEntry function
func TestNewUserEntry(t *testing.T) {
	tests := []struct {
		name   string
		userID string
	}{
		{
			name:   "normal_id",
			userID: "user123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the NewUserEntry function
			userEntry := NewUserEntry(tt.userID)

			// Check that the ID is set correctly
			if userEntry.ID != tt.userID {
				t.Errorf("NewUserEntry().ID = %v, want %v", userEntry.ID, tt.userID)
			}

			// Check that the FullID is set correctly
			expectedFullID := "users/" + tt.userID
			if userEntry.FullID != expectedFullID {
				t.Errorf("NewUserEntry().FullID = %v, want %v", userEntry.FullID, expectedFullID)
			}

			// Check that the Key is set correctly
			expectedKey := NewUserKey(tt.userID)
			if userEntry.Key.Collection() != expectedKey.Collection() || userEntry.Key.ID != expectedKey.ID {
				t.Errorf("NewUserEntry().Key = %v, want %v", userEntry.Key, expectedKey)
			}

			// Check that the Data is initialized
			if userEntry.Data == nil {
				t.Errorf("NewUserEntry().Data is nil, want initialized UserDbo")
			}

			// Check that the Record is initialized
			if userEntry.Record == nil {
				t.Errorf("NewUserEntry().Record is nil, want initialized Record")
			}
		})
	}
}

// TestNewUserEntryWithDbo tests the NewUserEntryWithDbo function
func TestNewUserEntryWithDbo(t *testing.T) {
	tests := []struct {
		name   string
		userID string
		dto    *UserDbo
	}{
		{
			name:   "normal_id_with_empty_dto",
			userID: "user123",
			dto:    new(UserDbo),
		},
		{
			name:   "normal_id_with_populated_dto",
			userID: "user456",
			dto: &UserDbo{
				ContactBase: briefs4contactus.ContactBase{
					ContactBrief: briefs4contactus.ContactBrief{
						Names: &person.NameFields{
							FirstName: "John",
							LastName:  "Doe",
						},
					},
				},
				CreatedFields: with.CreatedFields{
					CreatedAtField: with.CreatedAtField{
						CreatedAt: time.Now(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the NewUserEntryWithDbo function
			userEntry := NewUserEntryWithDbo(tt.userID, tt.dto)

			// Check that the ID is set correctly
			if userEntry.ID != tt.userID {
				t.Errorf("NewUserEntryWithDbo().ID = %v, want %v", userEntry.ID, tt.userID)
			}

			// Check that the FullID is set correctly
			expectedFullID := "users/" + tt.userID
			if userEntry.FullID != expectedFullID {
				t.Errorf("NewUserEntryWithDbo().FullID = %v, want %v", userEntry.FullID, expectedFullID)
			}

			// Check that the Key is set correctly
			expectedKey := NewUserKey(tt.userID)
			if userEntry.Key.Collection() != expectedKey.Collection() || userEntry.Key.ID != expectedKey.ID {
				t.Errorf("NewUserEntryWithDbo().Key = %v, want %v", userEntry.Key, expectedKey)
			}

			// Check that the Data is the same as the provided dto
			if userEntry.Data != tt.dto {
				t.Errorf("NewUserEntryWithDbo().Data = %v, want %v", userEntry.Data, tt.dto)
			}

			// Check that the Record is initialized
			if userEntry.Record == nil {
				t.Errorf("NewUserEntryWithDbo().Record is nil, want initialized Record")
			}
		})
	}
}
