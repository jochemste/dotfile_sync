package utils

import (
	"testing"
	"time"
)

func TestGetDate(t *testing.T) {
	date := GetDate()
	testtime := time.Now().Format("2006-01-02 15:04:05")

	if date != testtime {
		t.Errorf("Time error: %s is not %s", date, testtime)
	}
}
