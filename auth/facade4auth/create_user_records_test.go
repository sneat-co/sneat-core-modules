package facade4auth

import (
	"context"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mock_dal"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/dto4auth"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/appuser"
	"github.com/strongo/strongoapp/person"
	"go.uber.org/mock/gomock"
)

func Test_InitUserRecord(t *testing.T) {
	ctx := context.Background()
	type args struct {
		user                  facade.UserContext
		userToCreate          dto4auth.DataToCreateUser
		isCreateDefaultSpaces bool
	}
	tests := []struct {
		name     string
		args     args
		wantUser dbo4userus.UserEntry
		wantErr  bool
	}{
		{
			name: "should_create_user_record",
			args: args{
				user:                  facade.NewUserContext("test_user_1"),
				isCreateDefaultSpaces: true,
				userToCreate: dto4auth.DataToCreateUser{
					AuthAccount: appuser.AccountKey{
						Provider: "password",
						ID:       "u1@example.com",
					},
					Names: person.NameFields{
						FirstName: "First",
						LastName:  "UserEntry",
					},
					IanaTimezone: "Europe/Paris",
					Email:        "u1@example.com",
					RemoteClient: dbmodels.RemoteClientInfo{
						HostOrApp:  "unit-test",
						RemoteAddr: "127.0.0.1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// SETUP MOCKS BEGINS

			mockCtrl := gomock.NewController(t)
			db := mock_dal.NewMockDB(mockCtrl)
			facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
				return db, nil
			}

			db.EXPECT().RunReadwriteTransaction(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, f dal.RWTxWorker, options ...dal.TransactionOption) error {
					isContextWithUser := gomock.Cond(func(ctx context.Context) bool {
						return facade.GetUserContext(ctx) != nil
					})
					tx := mock_dal.NewMockReadwriteTransaction(mockCtrl)
					tx.EXPECT().Get(isContextWithUser, gomock.Any()).Return(dal.ErrRecordNotFound).AnyTimes() // TODO: Assert gets
					tx.EXPECT().Insert(isContextWithUser, gomock.Any()).Return(nil)                           // TODO: Assert inserts
					tx.EXPECT().InsertMulti(isContextWithUser, gomock.Any()).Return(nil)                      // TODO: Assert inserts
					return f(ctx, tx)
				})

			sneatauth.GetUserInfo = func(ctx context.Context, uid string) (authUser *sneatauth.AuthUserInfo, err error) {
				authUser = &sneatauth.AuthUserInfo{
					AuthProviderUserInfo: &sneatauth.AuthProviderUserInfo{
						ProviderID: "firebase",
					},
				}
				return
			}
			// SETUP MOCKS ENDS

			// TEST CALL BEGINS
			gotParams, err := CreateUserRecords(facade.NewContextWithUser(ctx, tt.args.user), tt.args.userToCreate, tt.args.isCreateDefaultSpaces)
			// TEST CALL ENDS

			if (err != nil) != tt.wantErr {
				t.Errorf("createOrUpdateUserRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.args.userToCreate.Email, gotParams.User.Data.Email)
		})
	}
}
