package models4auth

import (
	"time"

	"github.com/dal-go/dalgo/record"
)

const GaClientKind = "UserGaClient"

type GaClientEntity struct {
	Created   time.Time
	UserAgent string `firestore:",omitempty"`
	IpAddress string `firestore:",omitempty"`
}

type GaClient struct {
	record.WithID[string]
	*GaClientEntity
}

func (GaClient) Kind() string {
	return GaClientKind
}

func (gaClient GaClient) Entity() any {
	return gaClient.GaClientEntity
}

func (GaClient) NewEntity() any {
	return new(GaClientEntity)
}

func (gaClient *GaClient) SetEntity(entity any) {
	if entity == nil {
		gaClient.GaClientEntity = nil

	} else {
		gaClient.GaClientEntity = entity.(*GaClientEntity)

	}
}

//var _ db.EntityHolder = (*GaClient)(nil)
