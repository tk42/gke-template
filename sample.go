package main

import (
	_ "net/http/pprof"
	"time"

	"github.com/jimako1989/gke-template/profiler"
)

func NewCounter() <-chan int {
	c := make(chan int, 1)
	go func() {
		for i := 1; ; i++ {
			c <- i
		}
	}()
	return c
}

func worker() {
	counter := NewCounter()
	for c := range counter {
		if c == 47 {
			return
		}
	}
}

func main() {
	profiler.GetProfiler()

	for {
		go worker()
		time.Sleep(time.Millisecond * 20)
	}
}
