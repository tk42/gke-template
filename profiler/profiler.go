package profiler

import (
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"

	"github.com/tk42/victolinux/env"
)

var once sync.Once
var profiler *Profiler

type Profiler struct{}

func GetProfiler() *Profiler {
	if env.GetBoolean("PROFILER_BLOCKING", true) {
		runtime.SetBlockProfileRate(1)
	}
	once.Do(func() {
		profiler = &Profiler{}
		go http.ListenAndServe(":"+env.GetString("PROFILER_PORT", "6060"), nil)
	})
	return profiler
}
