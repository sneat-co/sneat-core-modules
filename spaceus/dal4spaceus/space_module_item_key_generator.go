package dal4spaceus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/strongo/random"
)

func GenerateNewSpaceModuleItemKey(ctx context.Context, tx dal.ReadwriteTransaction,
	spaceID coretypes.SpaceID, moduleID coretypes.ModuleID, collection string, length, maxAttempts int,
) (
	id string, key *dal.Key, err error,
) {
	for i := 0; i < maxAttempts; i++ {
		id = random.ID(length)
		key = dbo4spaceus.NewSpaceModuleItemKey(spaceID, moduleID, collection, id)
		record := dal.NewRecordWithData(key, make(map[string]any))
		if err := tx.Get(ctx, record); err != nil { // TODO: use tx.Exists()
			if dal.IsNotFound(err) {
				return id, key, nil
			}
			return "", nil, err
		}
	}
	return "", nil, errors.New("too many attempts  to generate a random happening ContactID")
}
