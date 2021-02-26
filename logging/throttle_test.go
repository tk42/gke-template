package logging

import (
	"fmt"
	"testing"
	"time"
)

var (
	period = 200 * time.Millisecond
)

func TestThrottle(t *testing.T) {
	throttle := GetThrottle(period, func() { fmt.Println("Throttle triggered") })
	defer throttle.Stop()

	go func() {
		for {
			fmt.Println("looped")
			throttle.Trigger()
			time.Sleep(period / 5)
		}
	}()

	time.Sleep(2 * period)
}
