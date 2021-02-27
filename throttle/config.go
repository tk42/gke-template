package throttle

import (
	"time"
)

type ThrottleConfig struct {
	count         uint
	period        time.Duration
	reached_func  func()
	released_func func()
}

type ThrottleOption func(*ThrottleConfig)

func Reached(f func()) ThrottleOption {
	return func(op *ThrottleConfig) {
		op.reached_func = f
	}
}

func Released(f func()) ThrottleOption {
	return func(op *ThrottleConfig) {
		op.released_func = f
	}
}

func ThrottleParameter(count uint, period time.Duration, ops ...ThrottleOption) ThrottleConfig {
	params := ThrottleConfig{
		count:         count,
		period:        period,
		reached_func:  nil,
		released_func: nil,
	}
	for _, option := range ops {
		option(&params)
	}
	return params
}
