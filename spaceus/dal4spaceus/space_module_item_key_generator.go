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
		// tx.Exists reports a missing document as (false, nil) — it does NOT
		// return a not-found error (dalgo2firestore clears it). So branch on the
		// returned bool: a non-existent key is the free id we want. The earlier
		// `if err != nil { if dal.IsNotFound(err) … }` form never fired (err was
		// always nil) and so always exhausted maxAttempts — breaking every
		// space-module item create (e.g. eventus event → calendarius happening).
		var exists bool
		if exists, err = tx.Exists(ctx, key); err != nil {
			return "", nil, err
		}
		if !exists {
			return id, key, nil
		}
	}
	return "", nil, errors.New("too many attempts to generate a random happening ContactID")
}
