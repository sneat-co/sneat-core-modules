package sneat_core_modules

import (
	"testing"
)

func TestCoreModules(t *testing.T) {
	if coreModules := CoreModules(); len(coreModules) == 0 {
		t.Errorf("CoreModules() returned empty list")
	}
}
