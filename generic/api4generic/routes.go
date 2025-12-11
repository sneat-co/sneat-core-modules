package api4generic

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/extension"
)

// RegisterHttpRoutes registers HTTP handlers
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/api4invitus/$generic/create", create)
	handle(http.MethodPut, "/api4invitus/$generic/update", update)
	handle(http.MethodDelete, "/api4invitus/$generic/delete", delete)
}
