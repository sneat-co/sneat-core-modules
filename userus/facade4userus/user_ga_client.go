package facade4userus

import (
	"context"
	"fmt"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-go-core/facade"
)

func SaveGaClient(ctx context.Context, gaClientId, userAgent, ipAddress string) (gaClient models4auth.GaClient, err error) {
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		data := new(models4auth.GaClientEntity)
		key := dal.NewKeyWithID(models4auth.GaClientKind, gaClientId)
		r := dal.NewRecordWithData(key, data)

		if err = tx.Get(ctx, r); err != nil && !dal.IsNotFound(err) {
			return fmt.Errorf("failed to get UserGaClient by ContactID: %w", err)
		}
		if data.UserAgent != userAgent || data.IpAddress != ipAddress {
			data.UserAgent = userAgent
			data.IpAddress = ipAddress
			data.Created = time.Now()
			if err = tx.Set(ctx, r); err != nil {
				err = fmt.Errorf("failed to save UserGaClient: %w", err)
				return err
			}
		}
		return nil
	}, nil)
	return
}
