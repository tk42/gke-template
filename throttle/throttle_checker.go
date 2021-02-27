package throttle

import (
	"time"
)

func GetThrottleExit(count uint, period time.Duration) Throttle {
	return NewThrottle(
		ThrottleConfig{
			count, period,
			func() {
				panic("Passed Throttle Exit")
			},
			nil,
		})
}

func GetThrottle(count uint, period time.Duration, reached func(), released func()) Throttle {
	return NewThrottle(
		ThrottleConfig{
			count, period,
			reached,
			released,
		})
}
