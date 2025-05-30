package facade4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/gosimple/slug"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
)

// AddSpaceMetricRequest request
type AddSpaceMetricRequest struct {
	dto4spaceus.SpaceRequest
	Metric dbo4spaceus.SpaceMetric `json:"metric"`
}

// Validate validates request
func (v *AddSpaceMetricRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if err := v.Metric.Validate(); err != nil {
		return err
	}
	return nil
}

// AddMetric adds metric
func AddMetric(ctx facade.ContextWithUser, request AddSpaceMetricRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	err = dal4spaceus.RunSpaceWorkerWithUserContext(ctx, request.SpaceID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4spaceus.SpaceWorkerParams) (err error) {
			request.Metric.ID = strings.ReplaceAll(slug.Make(request.Metric.Title), "-", "_")
			for _, m := range params.Space.Data.Metrics {
				if m.ID == request.Metric.ID {
					err = validation.NewErrBadRequestFieldValue("title", "duplicate slug(title)")
					return
				}
			}
			params.Space.Data.Metrics = append(params.Space.Data.Metrics, &request.Metric)
			if err = dal4spaceus.TxUpdateSpace(ctx, tx, params.Started, params.Space, []update.Update{
				update.ByFieldName("metrics", params.Space.Data.Metrics),
			}); err != nil {
				return err
			}
			return nil
		})
	return
}
