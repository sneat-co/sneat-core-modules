package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-core-modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSpaceWorkerParams_GetRecords(t *testing.T) {
	type args struct {
		records []dal.Record
	}
	tests := []struct {
		name    string
		params  SpaceWorkerParams
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no user context - space record not added",
			params: SpaceWorkerParams{
				Space:   dbo4spaceus.NewSpaceEntry("space1"),
				UserCtx: nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i...)
			},
		},
		{
			name: "with user context - space record added",
			params: SpaceWorkerParams{
				Space:   dbo4spaceus.NewSpaceEntry("space1"),
				UserCtx: facade.NewUserContext("user1"),
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err, i...)
			},
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			tx := mocks4dal.NewMockReadTransaction(mockController)

			// Using gomock assert that GetMulti is called with the records

			tx.EXPECT().
				GetMulti(ctx, gomock.Any()).
				Do(func(ctx context.Context, records []dal.Record) {
					if tt.params.UserCtx != nil {
						assert.Equal(t, 1, len(records))
					} else {
						assert.Equal(t, 0, len(records))
					}
					for i := range records {
						records[i].SetError(nil)
					}
					if tt.params.Space.Data != nil {
						tt.params.Space.Data.Type = core4spaceus.SpaceTypePrivate
						tt.params.Space.Data.Status = dbmodels.StatusActive
						tt.params.Space.Data.UserIDs = []string{"user1"}
						tt.params.Space.Data.CreatedBy = "user1"
						tt.params.Space.Data.CreatedAt = time.Now()
						tt.params.Space.Data.UpdatedBy = tt.params.Space.Data.CreatedBy
						tt.params.Space.Data.UpdatedAt = tt.params.Space.Data.CreatedAt
						tt.params.Space.Data.Version = 1
					}
				}).
				Return(nil)

			tt.wantErr(t,
				tt.params.GetRecords(ctx, tx, tt.args.records...),
				fmt.Sprintf("GetRecords(ctx, tx, %+v)", tt.args.records))
		})
	}
}
