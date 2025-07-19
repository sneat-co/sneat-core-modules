package api4linkage

import (
	"github.com/sneat-co/sneat-go-core/extension"
	"net/http"
)

func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/linkage/update_item_relationships", httpUpdateItemRelationships)
}
