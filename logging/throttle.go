package logging

import (
	"time"

	throttle "github.com/boz/go-throttle"
)

func GetThrottleExit(period time.Duration) throttle.ThrottleDriver {
	return throttle.ThrottleFunc(period, false, func() {
		panic("Passed Throttle Exit")
	})
}

func GetThrottle(period time.Duration, f func()) throttle.ThrottleDriver {
	return throttle.ThrottleFunc(period, false, f)
}
