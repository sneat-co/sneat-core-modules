package facade4auth

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
)

type CreateUserWorkerParams struct {
	*dal4userus.UserWorkerParams
	record.WithRecordChanges
}
