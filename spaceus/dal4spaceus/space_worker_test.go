package dal4spaceus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/mocks/mock_dal"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSpaceWorkerParams_GetRecords(t *testing.T) {
	type args struct {
		records []dal.Record
	}
	const userID = "user1"
	const spaceID = "space1"

	tests := []struct {
		name    string
		params  SpaceWorkerParams
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no user context - space record not added",
			params: SpaceWorkerParams{
				Space:   dbo4spaceus.NewSpaceEntry(spaceID),
				UserCtx: nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.NoError(t, err, i...)
			},
		},
		{
			name: "with user context - space record added",
			params: SpaceWorkerParams{
				Space:   dbo4spaceus.NewSpaceEntry(spaceID),
				UserCtx: facade.NewUserContext(userID),
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				return assert.NoError(t, err, i...)
			},
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			tx := mock_dal.NewMockReadTransaction(mockController)

			// Assert that GetMulti is called with the records
			tx.EXPECT().
				GetMulti(ctx, gomock.Any()).
				Times(1).
				Do(func(ctx context.Context, records []dal.Record) error {
					if tt.params.UserCtx != nil {
						assert.Equal(t, 1, len(records))
					} else {
						assert.Equal(t, 0, len(records))
					}
					for i := range records {
						records[i].SetError(nil)
					}
					if tt.params.Space.Data != nil {
						tt.params.Space.Data.Type = coretypes.SpaceTypePrivate
						tt.params.Space.Data.Status = dbmodels.StatusActive
						tt.params.Space.Data.UserIDs = []string{userID}
						tt.params.Space.Data.CreatedBy = userID
						tt.params.Space.Data.CreatedAt = time.Now()
						tt.params.Space.Data.UpdatedBy = tt.params.Space.Data.CreatedBy
						tt.params.Space.Data.UpdatedAt = tt.params.Space.Data.CreatedAt
						tt.params.Space.Data.Version = 1
					}
					return nil
				})

			tt.wantErr(t,
				tt.params.GetRecords(ctx, tx, tt.args.records...),
				fmt.Sprintf("GetRecords(ctx, tx, %+v)", tt.args.records))
		})
	}
}
