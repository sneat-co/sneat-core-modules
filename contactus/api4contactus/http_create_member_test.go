package api4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/httpmock"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/strongo/strongoapp/with"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpAddMember(t *testing.T) {

	const spaceID = "unit-test"
	request := dal4contactus.CreateMemberRequest{
		SpaceRequest: dto4spaceus.SpaceRequest{
			SpaceID: spaceID,
		},
		WithRelated: dbo4linkage.WithRelated{
			Related: dbo4linkage.RelatedModules{
				string(const4contactus.ModuleID): dbo4linkage.RelatedCollections{
					const4contactus.ContactsCollection: map[string]*dbo4linkage.RelatedItem{
						"c1": {
							RolesOfItem: map[dbo4linkage.RelationshipRoleID]*dbo4linkage.RelationshipRole{
								"spouse": {
									//CreatedField: with.CreatedField{
									//	Created: with.Created{
									//		By: "u1",
									//		At: "2020-01-01",
									//	},
									//},
								},
							},
						},
					},
				},
			},
		},
		CreatePersonRequest: dto4contactus.CreatePersonRequest{
			ContactBase: briefs4contactus.ContactBase{
				ContactBrief: briefs4contactus.ContactBrief{
					Type:     briefs4contactus.ContactTypePerson,
					Gender:   "unknown",
					Title:    "Some new members",
					AgeGroup: "unknown",
					RolesField: with.RolesField{
						Roles: []string{const4contactus.SpaceMemberRoleContributor},
					},
				},
				Status: "active",
				//WithRequiredCountryID: dbmodels.WithRequiredCountryID{
				//	CountryID: dbmodels.UnknownCountryID,
				//},
			},
			EmailsField: with.EmailsField{
				Emails: map[string]with.EmailProps{
					"someone@example.com": {
						Type: "personal",
						CreatedFields: with.CreatedFields{
							CreatedAtField: with.CreatedAtField{
								CreatedAt: time.Now(),
							},
							CreatedByField: with.CreatedByField{
								CreatedBy: "unit-test-user",
							},
						},
					},
				},
			},
		},
	}
	request.CountryID = "IE"

	defer func() {
		apicore.GetAuthTokenFromHttpRequest = nil
	}()
	apicore.GetAuthTokenFromHttpRequest = func(r *http.Request, authRequired bool) (token *sneatauth.Token, err error) {
		return &sneatauth.Token{UID: "TestUserID"}, nil
	}

	//t.Log(buffer.String())

	req := httpmock.NewPostJSONRequest(http.MethodPost, "/v0/space/create_member", request)
	req.Host = "localhost"
	req.Header.Set("Origin", "http://localhost:3000")

	createMember = func(ctx facade.ContextWithUser, request dal4contactus.CreateMemberRequest) (contact dal4contactus.ContactEntry, err error) {
		if request.SpaceID != spaceID {
			t.Fatalf("Expected spaceID=%s, got: %s", spaceID, request.SpaceID)
		}
		contact.ID = "abc1"
		contact.Data = &dbo4contactus.ContactDbo{
			ContactBase: briefs4contactus.ContactBase{
				ContactBrief: briefs4contactus.ContactBrief{
					Type:  briefs4contactus.ContactTypeCompany,
					Title: "Some company",
					OptionalCountryID: with.OptionalCountryID{
						CountryID: "IE",
					},
					RolesField: with.RolesField{
						Roles: []string{const4contactus.SpaceMemberRoleContributor},
					},
				},
				Status: "active",
				//WithRequiredCountryID: dbmodels.WithRequiredCountryID{
				//	CountryID: const4contactus.UnknownCountryID,
				//},
			},
		}
		contact.Data = &dbo4contactus.ContactDbo{
			ContactBase: contact.Data.ContactBase,
		}
		return
	}

	const uid = "unit-test-user"
	apicore.GetAuthTokenFromHttpRequest = func(r *http.Request, authRequired bool) (token *sneatauth.Token, err error) {
		return &sneatauth.Token{UID: uid}, nil
	}
	//sneatfb.NewFirebaseAuthToken = func(ctx context.Context, fbIDToken func() (string, error), authRequired bool) (*auth.Token, error) {
	//	return &auth.Token{UID: uid}, nil
	//}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(httpPostCreateMember)
	handler.ServeHTTP(rr, req)

	responseBody := rr.Body.String()

	if expected := http.StatusCreated; rr.Code != expected {
		t.Fatalf(
			"unexpected status: got (%d) expects (%d): %s",
			rr.Code,
			expected,
			responseBody,
		)
	}
}
