package time

import "time"

// Default is package variable
var Default = &UTCTime{}

// UTCTime represents utc time .
type UTCTime struct {
}

// Now function returns now formatted UTC time
func (u *UTCTime) Now() string {
	return time.Now().In(time.UTC).Format("20060102T150405Z")
}

// Date function returns date formatted UTC time
func (u *UTCTime) Date() string {
	return time.Now().In(time.UTC).Format("20060102")
}
