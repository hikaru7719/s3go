package time

import "time"

var Default = &UTCTime{}

type UTCTime struct {
}

func (u *UTCTime) Now() string {
	return time.Now().In(time.UTC).Format("20060102T150405Z")
}

func (u *UTCTime) Date() string {
	return time.Now().In(time.UTC).Format("20060102")
}
