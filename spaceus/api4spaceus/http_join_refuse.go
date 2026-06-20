package api4spaceus

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// refuseToJoinSpaceRequest request
type refuseToJoinSpaceRequest struct {
	SpaceID string `json:"id"`
	Pin     int32  `json:"pin"`
}

// Validate validates request
func (v *refuseToJoinSpaceRequest) Validate() error {
	if v.SpaceID == "" {
		return validation.NewErrRecordIsMissingRequiredField("space")
	}
	return nil
}

// refuseToJoinSpace refuses to join space
func refuseToJoinSpace(_ facade.ContextWithUser, request refuseToJoinSpaceRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}
	return errors.New("not implemented")
}

// httpPostRefuseToJoinSpace an API endpoint that records user refusal to join a space
func httpPostRefuseToJoinSpace(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithNoAuthRequired)
	if err != nil {
		return
	}
	q := r.URL.Query()
	var pin int
	if pin, err = strconv.Atoi(q.Get("pin")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("pin is expected to be an integer"))
		return
	}
	request := refuseToJoinSpaceRequest{
		SpaceID: q.Get("id"),
		Pin:     int32(pin),
	}
	err = refuseToJoinSpace(ctx, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
