package utils

import (
	"errors"
	"time"
)

var Layout = "2006-01-02 15:04:05"

func GetDate() string {
	return time.Now().Format(Layout)
}

func GetLatestFormatDate(date1 string, date2 string) (string, error) {
	t1, err := time.Parse(Layout, date1)
	if err != nil {
		return "", errors.New("Could not format date1: " + err.Error())
	}
	t2, err := time.Parse(Layout, date2)
	if err != nil {
		return "", errors.New("Could not format date2: " + err.Error())
	}

	delta := t1.Sub(t2)
	if delta > 0 {
		return date1, nil
	} else {
		return date2, nil
	}
}
