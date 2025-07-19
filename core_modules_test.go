package sneat_core_modules

import (
	"testing"
)

func TestCoreModules(t *testing.T) {
	if coreModules := CoreExtensions(); len(coreModules) == 0 {
		t.Errorf("CoreExtensions() returned empty list")
	}
}
