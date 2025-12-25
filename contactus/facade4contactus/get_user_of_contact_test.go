package facade4contactus

import (
	"context"
	"reflect"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/mocks/mock_dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"go.uber.org/mock/gomock"
)

func TestGetUserOfContact(t *testing.T) {
	type args struct {
		ctx      context.Context
		contact  UserAccountsProvider
		platform string
		app      string
	}
	tests := []struct {
		name     string
		args     args
		wantUser dbo4userus.UserEntry
		wantErr  bool
	}{
		{
			name: "by_userID",
			args: args{
				ctx: context.Background(),
				contact: &briefs4contactus.ContactBrief{
					WithUserID: dbmodels.WithUserID{
						UserID: "user1",
					},
				},
				platform: "telegram",
				app:      "",
			},
			wantUser: dbo4userus.NewUserEntry("user1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			db := mock_dal.NewMockDB(ctrl)
			db.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil)
			facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
				return db, nil
			}
			gotUser, err := GetUserOfContact(tt.args.ctx, tt.args.contact, tt.args.platform, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserOfContact() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("GetUserOfContact() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
