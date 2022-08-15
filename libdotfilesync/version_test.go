package libdotfilesync

import "testing"

func TestIsInitialised(t *testing.T) {
	if len(NAME) <= 0 {
		t.Errorf("NAME was not initialised")
	}
	if len(VERSION) <= 0 {
		t.Errorf("VERSION was not initialised")
	}
	if len(AUTHOR) <= 0 {
		t.Errorf("AUTHOR was not initialised")
	}
}
