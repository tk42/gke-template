package throttle

import (
	"fmt"
	"testing"
	"time"
)

var (
	period = 500 * time.Millisecond
)

func ThrottleTemplate(thres_count uint, interval int64, try_count int) {
	throttle := NewThrottle(ThrottleParameter(thres_count, period,
		Reached(func() {
			fmt.Printf("reached\t%v\n", time.Now())
		}),
		Released(func() {
			fmt.Printf("released\t%v\n", time.Now())
		}),
	))

	go func() {
		for i := 0; i < try_count; i++ {
			fmt.Printf("loop\t%v\n", time.Now())
			throttle.Trigger()
			time.Sleep(period / time.Duration(interval))
		}
	}()

	time.Sleep(2 * period)
	fmt.Printf("end\t%v\n", time.Now())
}

func TestThrottleNotFulfill(t *testing.T) {
	ThrottleTemplate(5, 5, 3) // 5/5 = 1
}

func TestThrottleNotReached(t *testing.T) {
	ThrottleTemplate(5, 4, 10) // 5/4 = 1.25 > 1
}

func TestThrottleReached(t *testing.T) {
	ThrottleTemplate(5, 6, 10) // 5/10 = 0.5 < 1
}
