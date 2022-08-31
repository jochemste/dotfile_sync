package utils

import (
	"testing"
	"time"
)

func TestGetDateGW(t *testing.T) {
	date := GetDate()
	testtime := time.Now().Format("2006-01-02 15:04:05")

	if date != testtime {
		t.Errorf("Time error: %s is not %s", date, testtime)
	}
}

func TestGetDateBW(t *testing.T) {
	date := GetDate()
	testtime := time.Now().Format("15:04:05 2006-01-02")

	if date == testtime {
		t.Errorf("Time error: %s is %s", date, testtime)
	}
}

func TestGetLatestFormatDateGW(t *testing.T) {
	date1 := GetDate()
	date2 := Layout // Layout is a date in 2006

	res, err := GetLatestFormatDate(date1, date2)
	if err != nil {
		t.Errorf("Could not get latest format: " + err.Error())
	}

	if res != date1 {
		t.Errorf("Returned wrong date: %s", res)
	}
}

func TestGetLatestFormatDateBW(t *testing.T) {
	date1 := GetDate()
	date2 := time.Now().Format("15:04:05 2006-01-02") // Wrong format

	_, err := GetLatestFormatDate(date1, date2)
	if err == nil {
		t.Errorf("Wrong format resulted in undefined behaviour: %s", date2)
	}

}

func TestIsMoreRecentTimeGW(t *testing.T) {
	t1 := time.Now()
	t2 := time.Now().Add(time.Hour * 2)

	if IsMoreRecentTime(t1, t2) {
		t.Errorf("Wrong resulting time. Should be false: %s", t1)
	}

	if !IsMoreRecentTime(t2, t1) {
		t.Errorf("Wrong resulting time. Should be true: %s", t2)
	}

	if IsMoreRecentTime(t1, t1) {
		t.Errorf("Wrong resulting time. Should be false: %s", t2)
	}

	if IsMoreRecentTime(t2, t2) {
		t.Errorf("Wrong resulting time. Should be false: %s", t2)
	}
}

func TestIsSameTimeGW(t *testing.T) {
	t1 := time.Now()
	t2 := time.Now().Add(time.Hour * 2)

	if IsSameTime(t1, t2) {
		t.Errorf("Wrong resulting time. Should be false: %s", t1)
	}

	if IsSameTime(t2, t1) {
		t.Errorf("Wrong resulting time. Should be false: %s", t2)
	}

	if !IsSameTime(t1, t1) {
		t.Errorf("Wrong resulting time. Should be true: %s", t2)
	}

	if !IsSameTime(t2, t2) {
		t.Errorf("Wrong resulting time. Should be true: %s", t2)
	}
}
