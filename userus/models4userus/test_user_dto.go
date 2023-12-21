package models4userus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"testing"
	"time"
)

func TestUserDtoValidate(t *testing.T) {
	now := time.Now()
	userDto := UserDto{
		WithCreated: dbmodels.WithCreated{
			CreatedAt: now,
			CreatedBy: "user",
		},
		ContactBase: briefs4contactus.ContactBase{
			ContactBrief: briefs4contactus.ContactBrief{
				Type:   briefs4contactus.ContactTypePerson,
				Gender: "unknown",
				Names: &person.NameFields{
					FirstName: "Firstname",
					LastName:  "Lastname",
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
	}
	userDto.CountryID = dbmodels.UnknownCountryID
	t.Run("empty_record", func(t *testing.T) {
		if err := userDto.Validate(); err != nil {
			t.Fatalf("no error expected, got: %v", err)
		}
	})
}
