package models4schedulus

import "github.com/dal-go/dalgo/record"

type HappeningContext struct {
	record.WithID[string]
	Dto *HappeningDto
}
