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
	spaceID coretypes.SpaceID, moduleID coretypes.ExtID, collection string, length, maxAttempts int,
) (
	id string, key *dal.Key, err error,
) {
	if tx == nil {
		panic("tx nil transaction")
	}
	for i := 0; i < maxAttempts; i++ {
		id = random.ID(length)
		key = dbo4spaceus.NewSpaceModuleItemKey(spaceID, moduleID, collection, id)
		if _, err = tx.Exists(ctx, key); err != nil { // TODO: use tx.Exists()
			if dal.IsNotFound(err) {
				return id, key, nil
			}
			return "", nil, err
		}
	}
	return "", nil, errors.New("too many attempts  to generate a random happening ContactID")
}
