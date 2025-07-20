package dbo4userus

import (
	"testing"
)

func TestNewUserKey(t *testing.T) {
	// Test with string ID
	stringID := "user123"
	key := NewUserKey(stringID)
	if key == nil {
		t.Error("NewUserKey returned nil for valid string ID")
	}

	// Test with int ID
	intID := 123
	key = NewUserKey(intID)
	if key == nil {
		t.Error("NewUserKey returned nil for valid int ID")
	}

	// Test with empty string ID - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewUserKey did not panic for empty string ID")
		}
	}()
	NewUserKey("")
}

func TestNewUserKeyWithZeroInt(t *testing.T) {
	// Test with zero int ID - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewUserKey did not panic for zero int ID")
		}
	}()
	NewUserKey(0)
}

func TestNewUserKeys(t *testing.T) {
	t.Run("with string IDs", func(t *testing.T) {
		ids := []string{"user1", "user2", "user3"}
		keys := NewUserKeys(ids)
		if len(keys) != len(ids) {
			t.Errorf("NewUserKeys returned %d keys, expected %d", len(keys), len(ids))
		}
		for i, key := range keys {
			if key == nil {
				t.Errorf("NewUserKeys returned nil key for ID %s", ids[i])
			}
			// Verify the key has the correct collection and ID
			if key.Collection() != Kind {
				t.Errorf("Key has incorrect collection: got %s, want %s", key.Collection(), Kind)
			}
			if key.ID != ids[i] {
				t.Errorf("Key has incorrect ID: got %v, want %v", key.ID, ids[i])
			}
		}
	})

	t.Run("with int IDs", func(t *testing.T) {
		ids := []int{101, 102, 103}
		keys := NewUserKeys(ids)
		if len(keys) != len(ids) {
			t.Errorf("NewUserKeys returned %d keys, expected %d", len(keys), len(ids))
		}
		for i, key := range keys {
			if key == nil {
				t.Errorf("NewUserKeys returned nil key for ID %d", ids[i])
			}
			// Verify the key has the correct collection and ID
			if key.Collection() != Kind {
				t.Errorf("Key has incorrect collection: got %s, want %s", key.Collection(), Kind)
			}
			if key.ID != ids[i] {
				t.Errorf("Key has incorrect ID: got %v, want %v", key.ID, ids[i])
			}
		}
	})

	t.Run("with empty slice", func(t *testing.T) {
		var ids []string
		keys := NewUserKeys(ids)
		if len(keys) != 0 {
			t.Errorf("NewUserKeys returned %d keys, expected 0", len(keys))
		}
	})
}

func TestNewUserKeysWithEmptyValue(t *testing.T) {
	// Test with a slice containing an empty string - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewUserKeys did not panic for slice containing empty string")
		}
	}()
	ids := []string{"user1", "", "user3"}
	NewUserKeys(ids)
}

func TestNewUserKeysWithZeroInt(t *testing.T) {
	// Test with a slice containing a zero int - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewUserKeys did not panic for slice containing zero int")
		}
	}()
	ids := []int{101, 0, 103}
	NewUserKeys(ids)
}
