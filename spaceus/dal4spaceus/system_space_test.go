package dal4spaceus

import (
	"context"
	"errors"
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

// fillSpace sets a loaded space record to the given type with the given members.
func fillSpace(data *dbo4spaceus.SpaceDbo, spaceType coretypes.SpaceType, userIDs []string) {
	data.Type = spaceType
	data.Title = "Test Space"
	data.Status = dbmodels.StatusActive
	data.UserIDs = userIDs
	data.CreatedBy = "creator"
	data.CreatedAt = time.Now()
	data.UpdatedBy = data.CreatedBy
	data.UpdatedAt = data.CreatedAt
	data.Version = 1
}

// Task 3 (non-member write) + Task 4 (authenticated read): an authenticated user
// who is NOT a member may access records in a System space, whereas a non-member
// in a non-System space is rejected.
func TestSpaceWorkerParams_GetRecords_SystemSpaceSkipsMembership(t *testing.T) {
	const spaceID = "games"
	const nonMemberUserID = "outsider"

	tests := []struct {
		name      string
		spaceType coretypes.SpaceType
		wantErr   bool
	}{
		{name: "system_non_member_allowed", spaceType: coretypes.SpaceTypeSystem, wantErr: false},
		{name: "private_non_member_rejected", spaceType: coretypes.SpaceTypePrivate, wantErr: true},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := SpaceWorkerParams{
				Space:   dbo4spaceus.NewSpaceEntry(spaceID),
				UserCtx: facade.NewUserContext(nonMemberUserID),
			}
			ctrl := gomock.NewController(t)
			tx := mock_dal.NewMockReadTransaction(ctrl)
			tx.EXPECT().GetMulti(ctx, gomock.Any()).Times(1).Do(
				func(ctx context.Context, records []dal.Record) error {
					for i := range records {
						records[i].SetError(nil)
					}
					// the space record carries members that exclude nonMemberUserID
					fillSpace(params.Space.Data, tt.spaceType, []string{"member1"})
					return nil
				})
			err := params.GetRecords(ctx, tx, params.Space.Record)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, facade.ErrUnauthorized))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Task 3 (unauthenticated write rejected): a write to a System space with no
// authenticated user is rejected by the pre-transaction check before any
// transaction runs.
func TestRunModuleSpaceWorkerWithUserCtx_SystemSpaceRejectsUnauthenticated(t *testing.T) {
	const spaceID coretypes.SpaceID = "games"
	const moduleID = "test_module"
	ctxNoUser := facade.NewContextWithUser(context.Background(), nil)

	facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
		ctrl := gomock.NewController(t)
		db := mock_dal.NewMockDB(ctrl)
		// The space is loaded outside the transaction; mark it as a System space.
		db.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
			func(ctx context.Context, record dal.Record) error {
				record.SetError(nil)
				fillSpace(record.Data().(*dbo4spaceus.SpaceDbo), coretypes.SpaceTypeSystem, []string{"member1"})
				return nil
			})
		// RunReadwriteTransaction must NOT be called: the unauthenticated write is
		// rejected before any transaction starts (no EXPECT() registered for it).
		return db, nil
	}

	worker := func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *ModuleSpaceWorkerParams[*fooModuleSpaceData]) error {
		return nil
	}
	err := RunModuleSpaceWorkerWithUserCtx(ctxNoUser, spaceID, moduleID, new(fooModuleSpaceData), worker)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, facade.ErrUnauthorized))
}
