package util

import (
	"time"
)

const location = "Asia/Tokyo"

var (
	// TimeNowFunc :
	TimeNowFunc = func() time.Time { return time.Now().Local() }
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// InitLocale :
func InitLocale() {
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}
