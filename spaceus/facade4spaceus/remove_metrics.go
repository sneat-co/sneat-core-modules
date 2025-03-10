package facade4spaceus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// RemoveMetrics removes a metric
func RemoveMetrics(ctx facade.ContextWithUser, request dto4spaceus.SpaceMetricsRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4spaceus.RunSpaceWorkerWithUserContext(ctx, ctx.User(), request.SpaceID,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceWorkerParams) (err error) {
			changed := false
			space := params.Space

			metrics := make([]*dbo4spaceus.SpaceMetric, 0, len(space.Data.Metrics))
		Metrics:
			for _, metric := range space.Data.Metrics {
				for i, metricID := range request.Metrics {
					if metric.ID == metricID {
						changed = true
						request.Metrics = append(request.Metrics[:i], request.Metrics[i+1:]...)
						continue Metrics
					}
				}
				metrics = append(metrics, metric)
			}
			if changed {
				var updates []update.Update
				if len(metrics) == 0 {
					updates = []update.Update{update.ByFieldName("metrics", update.DeleteField)}
				} else {
					updates = []update.Update{update.ByFieldName("metrics", metrics)}
				}
				if err = dal4spaceus.TxUpdateSpace(ctx, tx, params.Started, params.Space, updates); err != nil {
					return err
				}
			}
			return nil
		})
	return
}
