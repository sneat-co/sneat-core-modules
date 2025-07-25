package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mock_dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/person"
	"github.com/strongo/strongoapp/with"
	"go.uber.org/mock/gomock"
	"slices"
	"testing"
	"time"
)

func TestAcceptPersonalInvite(t *testing.T) {
	type args struct {
		ctx     facade.ContextWithUser
		request ClaimPersonalInviteRequest
	}
	ctx := facade.NewContextWithUserID(context.Background(), "123")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "nil_params",
			args:    args{ctx: ctx},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ClaimPersonalInvite(tt.args.ctx, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("ClaimPersonalInvite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAcceptPersonalInviteRequest_Validate(t *testing.T) {
	type fields struct {
		InviteRequest InviteRequest
		Member        dbmodels.DtoWithID[*briefs4contactus.ContactBase]
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "should_return_error_for_empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ClaimPersonalInviteRequest{
				ClaimInviteRequest: ClaimInviteRequest{
					InviteRequest: tt.fields.InviteRequest,
				},
			}
			if err := v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createOrUpdateUserRecord(t *testing.T) {
	ctx := context.Background()
	type args struct {
		user              dbo4userus.UserEntry
		userRecordError   error
		spaceRecordError  error
		inviteRecordError error
		request           ClaimPersonalInviteRequest
		space             dbo4spaceus.SpaceEntry
		spaceMember       dbmodels.DtoWithID[*briefs4contactus.ContactBase]
		invite            InviteEntry
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil_params",
			args: args{
				user:            dbo4userus.NewUserEntry("test_user_id"),
				userRecordError: dal.ErrRecordNotFound,
				space: dbo4spaceus.NewSpaceEntryWithDbo("testspaceid", &dbo4spaceus.SpaceDbo{
					SpaceBrief: dbo4spaceus.SpaceBrief{
						OptionalCountryID: with.OptionalCountryID{
							CountryID: with.UnknownCountryID,
						},
						Type:  "family",
						Title: "Family",
					},
				}),
				spaceMember: dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
					ID: "test_member_id2",
					Data: &briefs4contactus.ContactBase{
						ContactBrief: briefs4contactus.ContactBrief{
							Type:   briefs4contactus.ContactTypePerson,
							Gender: "unknown",
							Names: &person.NameFields{
								FirstName: "First",
							},
							//Status:   "active",
							AgeGroup: "unknown",
						},
						//WithRequiredCountryID: dbmodels.WithRequiredCountryID{
					},
				},
				invite: NewInviteEntryWithDbo("test_personal_invite_id", &dbo4invitus.InviteDbo{
					Roles: []string{"contributor"},
				}),
				request: ClaimPersonalInviteRequest{
					ClaimInviteRequest: ClaimInviteRequest{
						RemoteClient: dbmodels.RemoteClientInfo{
							HostOrApp:  "unit-test",
							RemoteAddr: "localhost",
						},
						InviteRequest: InviteRequest{
							InviteID: "test_personal_invite_id",
							Pin:      "1234",
						},
					},
					SpaceRequest: dto4spaceus.SpaceRequest{
						SpaceID: "testspaceid",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			//
			tt.args.user.Record.SetError(tt.args.userRecordError)
			tt.args.space.Record.SetError(tt.args.spaceRecordError)
			tt.args.invite.Record.SetError(tt.args.inviteRecordError)
			//
			tx := mock_dal.NewMockReadwriteTransaction(mockCtrl)
			if tt.args.userRecordError == nil && tt.args.spaceRecordError == nil && tt.args.inviteRecordError == nil {
				tx.EXPECT().Insert(gomock.Any(), tt.args.user.Record).Return(nil)
			}
			now := time.Now()
			params := dal4contactus.NewContactusSpaceWorkerParams(facade.NewUserContext(tt.args.user.ID), tt.args.space.ID)
			var member dal4contactus.ContactEntry
			if err := createOrUpdateUserRecord(ctx, tx, now, tt.args.user, tt.args.request, member, params, tt.args.spaceMember.Data, tt.args.invite); err != nil {
				if !tt.wantErr {
					t.Errorf("createOrUpdateUserRecord() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			userDto := tt.args.user.Data
			assert.Equal(t, now, userDto.CreatedAt, "CreatedAt")
			assert.Equal(t, member.Data.Gender, userDto.Gender, "Gender")
			assert.Equal(t, 1, len(userDto.Spaces), "len(Spaces)")
			assert.Equal(t, 1, len(userDto.SpaceIDs), "len(SpaceIDs)")
			assert.True(t, slices.Contains(userDto.SpaceIDs, string(tt.args.request.SpaceID)), "SpaceIDs contains tt.args.request.SpaceID")
			spaceBrief := userDto.Spaces[string(tt.args.request.SpaceID)]
			assert.NotNil(t, spaceBrief, "Spaces[tt.args.request.SpaceID]")
		})
	}
}

func Test_updateInviteRecord(t *testing.T) {
	ctx := context.Background()
	type args struct {
		uid    string
		invite InviteEntry
		status dbo4invitus.InviteStatus
	}
	now := time.Now()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should_pass",
			args: args{
				status: dbo4invitus.InviteStatusAccepted,
				invite: NewInviteEntryWithDbo("test_invite_id1", &dbo4invitus.InviteDbo{
					ToSpaceContactID: "to_member_id2",
					Address:          "to.test.user@example.com",
					Pin:              "1234",
					SpaceID:          "testspaceid1",
					Space: &dbo4invitus.InviteSpace{
						Type:  "family",
						Title: "Family",
					},
					CreatedAt: time.Now(),
					Created: dbmodels.CreatedInfo{
						Client: dbmodels.RemoteClientInfo{
							HostOrApp:  "unit-test",
							RemoteAddr: "127.0.0.1",
						},
					},
					FromUserID: "from_user_id1",
					InviteBase: dbo4invitus.InviteBase{
						Type:    "personal",
						Channel: "email",
						From: dbo4invitus.InviteFrom{
							InviteContact: dbo4invitus.InviteContact{
								UserID:    "from_user_id1",
								ContactID: "from_contact_id1",
								Title:     "From ContactID 1",
							},
						},
						To: &dbo4invitus.InviteTo{
							InviteContact: dbo4invitus.InviteContact{
								Title:     "To ContactID 2",
								ContactID: "to_contact_id2",
								Channel:   "email",
								Address:   "to.test.user@example.com",
							},
						},
					},
					Roles: []string{"contributor"}}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			tx := mock_dal.NewMockReadwriteTransaction(mockCtrl)
			tx.EXPECT().Update(ctx, tt.args.invite.Key, gomock.Any()).Return(nil)
			assert.Equal(t, "", tt.args.invite.Data.To.UserID)
			if err := updateInviteStatus(ctx, tx, tt.args.uid, now, tt.args.invite, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("updateInviteStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.args.status, tt.args.invite.Data.Status)
			assert.Equal(t, tt.args.uid, tt.args.invite.Data.To.UserID)
		})
	}
}

func Test_updateSpaceRecord(t *testing.T) {
	type args struct {
		uid            string
		memberID       string
		space          dbo4spaceus.SpaceEntry
		contactusSpace dal4contactus.ContactusSpaceEntry
		requestMember  dbmodels.DtoWithID[*briefs4contactus.ContactBase]
	}
	testMember := dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
		ID:   "test_member_id1",
		Data: &briefs4contactus.ContactBase{},
	}
	tests := []struct {
		name            string
		spaceRecordErr  error
		args            args
		wantSpaceMember dbmodels.DtoWithID[*briefs4contactus.ContactBase]
		wantErr         bool
	}{
		{
			name:           "should_pass",
			spaceRecordErr: nil,
			args: args{
				uid:      "test_user_id",
				memberID: "test_member_id1",
				space: dbo4spaceus.NewSpaceEntryWithDbo("testspaceid", &dbo4spaceus.SpaceDbo{
					SpaceBrief: dbo4spaceus.SpaceBrief{
						Type:  "family",
						Title: "Family",
					},
				}),
				contactusSpace: dal4contactus.NewContactusSpaceEntryWithData("testspaceid", &dbo4contactus.ContactusSpaceDbo{
					WithSingleSpaceContactsWithoutContactIDs: briefs4contactus.WithSingleSpaceContactsWithoutContactIDs[*briefs4contactus.ContactBrief]{
						WithContactsBase: briefs4contactus.WithContactsBase[*briefs4contactus.ContactBrief]{
							WithContactBriefs: briefs4contactus.WithContactBriefs[*briefs4contactus.ContactBrief]{
								Contacts: map[string]*briefs4contactus.ContactBrief{
									testMember.ID: &testMember.Data.ContactBrief,
								},
							},
						},
					},
				}),
				requestMember: dbmodels.DtoWithID[*briefs4contactus.ContactBase]{
					ID: testMember.ID,
					Data: &briefs4contactus.ContactBase{
						ContactBrief: briefs4contactus.ContactBrief{
							Names: &person.NameFields{
								FirstName: "First name",
							},
						},
					},
				},
			},
			wantErr:         false,
			wantSpaceMember: testMember,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//mockCtrl := gomock.NewController(t)
			//tx := mock_dal.NewMockReadwriteTransaction(mockCtrl)
			//tx.EXPECT().Update(gomock.Any(), tt.args.space.Key, gomock.Any()).Return(nil)
			//tx.EXPECT().Update(gomock.Any(), tt.args.contactusSpace.Key, gomock.Any()).Return(nil)
			tt.args.contactusSpace.Record.SetError(tt.spaceRecordErr)
			params := dal4contactus.NewContactusSpaceWorkerParams(facade.NewUserContext(tt.args.uid), tt.args.space.ID)
			params.SpaceModuleEntry.Data.AddContact(tt.args.memberID, &tt.args.requestMember.Data.ContactBrief)
			params.SpaceModuleEntry.Data.AddUserID(tt.args.uid)
			params.Space.Data.AddUserID(tt.args.uid)
			member := dal4contactus.NewContactEntry(tt.args.space.ID, "member1")
			gotSpaceMember, err := updateContactusSpaceRecord(tt.args.uid, tt.args.memberID, params, member)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateSpaceRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, gotSpaceMember, "gotSpaceMember is nil")
			//if !reflect.DeepEqual(gotSpaceMember, tt.wantSpaceMember) {
			//	t.Errorf("updateSpaceRecord() gotSpaceMember = %v, want %v", gotSpaceMember, tt.wantSpaceMember)
			//}
		})
	}
}
