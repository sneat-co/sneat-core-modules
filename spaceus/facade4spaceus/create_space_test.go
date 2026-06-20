package facade4spaceus

import (
	"context"
	"testing"

	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/mocks/mock_dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactusmodels/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"github.com/strongo/strongoapp/person"
	"go.uber.org/mock/gomock"
)

func TestCreateSpace(t *testing.T) { // TODO: Implement unit tests
	ctx := context.Background()
	user := facade.NewUserContext("TestUser")
	ctxWithUser := facade.NewContextWithUser(ctx, user)

	setupMockDb := func(insertMultiTimes int) {
		mockCtrl := gomock.NewController(t)
		db := mock_dal.NewMockDB(mockCtrl)
		facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
			return db, nil
		}

		tx := mock_dal.NewMockReadwriteTransaction(mockCtrl)
		assertContextWithDeadLine := gomock.Cond(func(x context.Context) bool {
			_, ok := x.Deadline()
			return ok
		})
		tx.EXPECT().Get(assertContextWithDeadLine, gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record) error {
			switch record.Key().Collection() {
			case dbo4userus.UsersCollection:
				record.SetError(nil)
				userDto := record.Data().(*dbo4userus.UserDbo)
				userDto.CountryID = "--"
				userDto.Status = "active"
				userDto.Gender = dbmodels.GenderMale
				userDto.AgeGroup = dbmodels.AgeGroupAdult
				userDto.Type = briefs4contactus.ContactTypePerson
				userDto.Names = &person.NameFields{
					FirstName: "1st",
					LastName:  "Lastnameoff",
				}
				userDto.Created = dbmodels.CreatedInfo{
					Client: dbmodels.RemoteClientInfo{
						HostOrApp:  "sneat.app",
						RemoteAddr: "127.0.0.1",
					},
				}
				return nil
			default:
				err := dal.ErrRecordNotFound
				record.SetError(err)
				return err
			}
		}).AnyTimes()
		tx.EXPECT().Insert(assertContextWithDeadLine, gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record, opts ...dal.InsertOption) error {
			return nil
		}).AnyTimes()
		if insertMultiTimes > 0 {
			tx.EXPECT().InsertMulti(assertContextWithDeadLine, gomock.Any()).DoAndReturn(func(ctx context.Context, records []dal.Record, opts ...dal.InsertOption) error {
				return nil
			}).Times(insertMultiTimes)
		}
		tx.EXPECT().Update(assertContextWithDeadLine, gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, key *dal.Key, updates []update.Update, preconditions ...dal.Precondition) error {
			return nil
		}).AnyTimes()
		db.EXPECT().RunReadwriteTransaction(assertContextWithDeadLine, gomock.Any()).DoAndReturn(func(ctx context.Context, worker func(ctx context.Context, tx dal.ReadwriteTransaction) error, o ...dal.TransactionOption) error {
			return worker(ctx, tx)
		}).AnyTimes()

		facade.GetSneatDB = func(ctx context.Context) (dal.DB, error) {
			return db, nil
		}
	}

	t.Run("error on bad request", func(t *testing.T) {
		setupMockDb(0)
		fakeContributor := &fakeContactusSpaceContributor{}
		RegisterContactusSpaceContributor(fakeContributor)
		result, err := CreateSpace(ctxWithUser, dto4spaceus.CreateSpaceRequest{})
		assert.Error(t, err)
		assert.Equal(t, coretypes.SpaceID(""), result.Space.ID)
		assert.False(t, fakeContributor.called)
	})

	t.Run("user's 1st space", func(t *testing.T) {
		setupMockDb(1)
		fakeContributor := &fakeContactusSpaceContributor{}
		RegisterContactusSpaceContributor(fakeContributor)

		result, err := CreateSpace(ctxWithUser, dto4spaceus.CreateSpaceRequest{Type: coretypes.SpaceTypeFamily})
		assert.Nil(t, err)

		assert.NotEqual(t, coretypes.SpaceID(""), result.Space.ID)
		assert.Nil(t, result.Space.Data.Validate())
		assert.Equal(t, 1, len(result.Space.Data.UserIDs))
		assert.Equal(t, 1, result.Space.Data.Version)
		assert.True(t, fakeContributor.called)
		assert.Equal(t, result.Space.ID, fakeContributor.gotSpaceID)
	})

}

// fakeContactusSpaceContributor is a test double for ContactusSpaceContributor so that
// spaceus tests stay decoupled from the contactus module.
type fakeContactusSpaceContributor struct {
	called     bool
	gotSpaceID coretypes.SpaceID
}

func (f *fakeContactusSpaceContributor) BuildSpaceCreationRecords(
	spaceID coretypes.SpaceID,
	_ string,
	_ briefs4contactus.ContactBrief,
	_ time.Time,
	_ string,
) ([]dal.Record, error) {
	f.called = true
	f.gotSpaceID = spaceID
	rec := dal.NewRecordWithData(dal.NewKeyWithID("contactus", "contactus"), &struct{}{})
	rec.MarkAsChanged()
	rec.SetError(nil)
	return []dal.Record{rec}, nil
}

func Test_getUniqueSpaceID(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	readSession := mock_dal.NewMockReadSession(mockCtrl)
	readSession.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dal.ErrRecordNotFound)
	spaceID, err := getUniqueSpaceID(ctx, readSession, "TestCompany LTD")
	assert.NoError(t, err)
	assert.Equal(t, coretypes.SpaceID("testcompanyltd"), spaceID)
}
