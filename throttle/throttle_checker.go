package throttle

import (
	"fmt"
	"time"
)

func GetThrottleExit(count uint, period time.Duration) Throttle {
	return NewThrottle(
		ThrottleConfig{
			count, period,
			func() {
				panic("DETECTED THROTTLE CHCK!")
			},
			nil,
		})
}

func GetThrottleSuppress(count uint, period time.Duration) Throttle {
	return NewThrottle(
		ThrottleConfig{
			count, period,
			func() {
				fmt.Print("DETECTED THROTTLE CHECK")
			},
			func() {
				fmt.Print("RECOVERED THROTTLE CHECK")
			},
		})
}
