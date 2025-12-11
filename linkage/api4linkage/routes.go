package api4linkage

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/extension"
)

func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/linkage/update_item_relationships", httpUpdateItemRelationships)
}
