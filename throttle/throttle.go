package throttle

import (
	"sync"
	"time"
)

// Throttle is an interface for requesting execution of the throttled resource
// and for stopping the throttler.
type Throttle interface {
	Trigger()
	IsFreeze() bool
}

func NewThrottle(cfg ThrottleConfig) Throttle {
	throttler := newThrottler(cfg)
	if cfg.reached_func != nil {
		go func() {
			for {
				<-throttler.reached
				cfg.reached_func()
			}
		}()
	}
	return throttler
}

type throttler struct {
	cond          *sync.Cond
	count         uint
	period        time.Duration
	last          []int64
	freeze        bool
	reached       chan struct{}
	released_func func()
}

func newThrottler(cfg ThrottleConfig) *throttler {
	return &throttler{
		period:        cfg.period,
		count:         cfg.count,
		last:          make([]int64, cfg.count),
		cond:          sync.NewCond(&sync.Mutex{}),
		reached:       make(chan struct{}),
		released_func: cfg.released_func,
	}
}

func (t *throttler) Trigger() {
	if t.freeze {
		return
	}
	t.cond.L.Lock()
	defer t.cond.L.Unlock()

	t.last = append(t.last[1:], time.Now().UnixNano())
	if t.last[0] > (time.Now().UnixNano() - t.period.Nanoseconds()) {
		t.reached <- struct{}{}
		t.freeze = true
		time.AfterFunc(time.Unix(0, t.last[0]+t.period.Nanoseconds()).Sub(time.Now()), func() {
			t.freeze = false
			if t.released_func != nil {
				t.released_func()
			}
		})
	}
}

func (t *throttler) IsFreeze() bool {
	return t.freeze
}
