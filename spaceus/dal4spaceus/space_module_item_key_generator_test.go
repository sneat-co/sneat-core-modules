package dal4spaceus

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/stretchr/testify/assert"
)

// fakeExistsTx is a minimal ReadwriteTransaction that only implements Exists
// (the single method GenerateNewSpaceModuleItemKey calls). Every other method
// would panic via the embedded nil interface, which is fine for these tests.
type fakeExistsTx struct {
	dal.ReadwriteTransaction
	// existsResults is consumed one entry per Exists() call, modelling a
	// sequence of collisions (true) followed by a free key (false).
	existsResults []bool
	calls         int
	err           error
}

func (f *fakeExistsTx) Exists(context.Context, *dal.Key) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	exists := false
	if f.calls < len(f.existsResults) {
		exists = f.existsResults[f.calls]
	}
	f.calls++
	return exists, nil
}

// TestGenerateNewSpaceModuleItemKey_HappyPath is the regression for the bug
// where Exists reports a missing doc as (false, nil) — not a not-found error —
// so the old `if err != nil { if dal.IsNotFound(err) … }` form never returned an
// id and always exhausted maxAttempts. A free key must yield a valid id.
func TestGenerateNewSpaceModuleItemKey_HappyPath(t *testing.T) {
	tx := &fakeExistsTx{} // no collisions: every Exists → (false, nil)
	id, key, err := GenerateNewSpaceModuleItemKey(
		context.Background(), tx, "space1", "eventus", "events", 5, 10)
	assert.NoError(t, err)
	assert.Len(t, id, 5)
	assert.NotNil(t, key)
	assert.Equal(t, 1, tx.calls, "should return on the first free key, not loop")
}

// TestGenerateNewSpaceModuleItemKey_RetriesPastCollisions verifies a collision
// on the first attempt is skipped and the next free key is returned.
func TestGenerateNewSpaceModuleItemKey_RetriesPastCollisions(t *testing.T) {
	tx := &fakeExistsTx{existsResults: []bool{true, false}}
	id, _, err := GenerateNewSpaceModuleItemKey(
		context.Background(), tx, "space1", "eventus", "events", 5, 10)
	assert.NoError(t, err)
	assert.Len(t, id, 5)
	assert.Equal(t, 2, tx.calls, "should retry once past the collision")
}

// TestGenerateNewSpaceModuleItemKey_ExhaustsAttempts verifies that when every
// candidate key collides, the generator gives up after maxAttempts.
func TestGenerateNewSpaceModuleItemKey_ExhaustsAttempts(t *testing.T) {
	tx := &fakeExistsTx{existsResults: []bool{true, true, true}}
	_, _, err := GenerateNewSpaceModuleItemKey(
		context.Background(), tx, "space1", "eventus", "events", 5, 3)
	assert.Error(t, err)
	assert.Equal(t, 3, tx.calls)
}

// TestGenerateNewSpaceModuleItemKey_ExistsError verifies a real Exists error is
// propagated rather than swallowed or retried.
func TestGenerateNewSpaceModuleItemKey_ExistsError(t *testing.T) {
	wantErr := errors.New("boom")
	tx := &fakeExistsTx{err: wantErr}
	_, _, err := GenerateNewSpaceModuleItemKey(
		context.Background(), tx, "space1", "eventus", "events", 5, 10)
	assert.ErrorIs(t, err, wantErr)
}

func TestGenerateNewSpaceModuleItemKey(t *testing.T) {
	type args struct {
		ctx         context.Context
		tx          dal.ReadwriteTransaction
		spaceID     coretypes.SpaceID
		moduleID    coretypes.ExtID
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
