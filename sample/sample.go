package main

import (
	"log"
	_ "net/http/pprof"
	"time"

	"github.com/tk42/gke-template/profiler"
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
	log.Println("Hello world!")

	profiler.GetProfiler()

	log.Println("Finished set up profiler")

	for {
		go worker()
		time.Sleep(time.Millisecond * 20)
	}
}
