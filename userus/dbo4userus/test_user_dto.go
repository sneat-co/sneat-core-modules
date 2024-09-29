package dbo4userus

import (
	briefs4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"testing"
	"time"
)

func TestUserDtoValidate(t *testing.T) {
	now := time.Now()
	userDto := UserDbo{
		CreatedFields: with.CreatedFields{
			CreatedAtField: with.CreatedAtField{
				CreatedAt: now,
			},
			CreatedByField: with.CreatedByField{
				CreatedBy: "user",
			},
		},
		ContactBase: briefs4contactus2.ContactBase{
			ContactBrief: briefs4contactus2.ContactBrief{
				Type:   briefs4contactus2.ContactTypePerson,
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
	userDto.CountryID = with.UnknownCountryID
	t.Run("empty_record", func(t *testing.T) {
		if err := userDto.Validate(); err != nil {
			t.Fatalf("no error expected, got: %v", err)
		}
	})
}
