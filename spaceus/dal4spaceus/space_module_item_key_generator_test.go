package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateNewSpaceModuleItemKey(t *testing.T) {
	type args struct {
		ctx         context.Context
		tx          dal.ReadwriteTransaction
		spaceID     coretypes.SpaceID
		moduleID    coretypes.ModuleID
		collection  string
		length      int
		maxAttempts int
	}
	tests := []struct {
		name      string
		args      args
		wantId    string
		wantKey   *dal.Key
		wantErr   assert.ErrorAssertionFunc
		wantPanic string
	}{
		{
			name:      "tx_nil",
			wantPanic: "tx nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic != "" {
				assert.Panics(t, func() {
					_, _, _ = GenerateNewSpaceModuleItemKey(tt.args.ctx, tt.args.tx, tt.args.spaceID, tt.args.moduleID, tt.args.collection, tt.args.length, tt.args.maxAttempts)
				})
			} else {
				gotId, gotKey, err := GenerateNewSpaceModuleItemKey(tt.args.ctx, tt.args.tx, tt.args.spaceID, tt.args.moduleID, tt.args.collection, tt.args.length, tt.args.maxAttempts)
				if !tt.wantErr(t, err, fmt.Sprintf("GenerateNewSpaceModuleItemKey(%v, %v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.tx, tt.args.spaceID, tt.args.moduleID, tt.args.collection, tt.args.length, tt.args.maxAttempts)) {
					return
				}
				assert.Equalf(t, tt.wantId, gotId, "GenerateNewSpaceModuleItemKey(%v, %v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.tx, tt.args.spaceID, tt.args.moduleID, tt.args.collection, tt.args.length, tt.args.maxAttempts)
				assert.Equalf(t, tt.wantKey, gotKey, "GenerateNewSpaceModuleItemKey(%v, %v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.tx, tt.args.spaceID, tt.args.moduleID, tt.args.collection, tt.args.length, tt.args.maxAttempts)
			}
		})
	}
}
