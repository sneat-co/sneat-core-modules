package api4spaceus

import (
	"net/http"
	"strings"

	"github.com/sneat-co/sneat-core-modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

// httpPostAddMetric is an API endpoint that adds a metric
func httpPostAddMetric(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request facade4spaceus.AddSpaceMetricRequest
	q := r.URL.Query()
	if request.SpaceID = coretypes.SpaceID(q.Get("id")); strings.TrimSpace(string(request.SpaceID)) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("space 'id' should be passed as query parameter"))
		return
	}
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = addMetric(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}

var addMetric = facade4spaceus.AddMetric
