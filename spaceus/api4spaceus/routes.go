package api4spaceus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/extension"
)

// RegisterHttpRoutes registers space routes
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/spaces/create_space", httpPostCreateSpace)
	//
	handle(http.MethodPost, "/v0/space/join_info", httpPostGetSpaceJoinInfo)
	handle(http.MethodPost, "/v0/space/join_space", httpPostJoinSpace)
	handle(http.MethodPost, "/v0/space/refuse_to_join_space", httpPostRefuseToJoinSpace)
	handle(http.MethodPost, "/v0/space/add_metric", httpPostAddMetric)
	handle(http.MethodPost, "/v0/space/remove_metrics", httpPostRemoveMetrics)
	handle(http.MethodGet, "/v0/space", httpGetSpace)

	handle(http.MethodPost, "/v0/space/update_related", httpPostUpdateRelated)
}
