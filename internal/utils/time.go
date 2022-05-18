package utils

import "time"

var Now = time.Now
var DefaultMockTime = time.Now()

func MockNow(t time.Time) {
	Now = func() time.Time { return t }
}
