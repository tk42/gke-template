package profiler

import (
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/jimako1989/gke-template/env"
)

var once sync.Once
var profiler *Profiler

type Profiler struct{}

func GetProfiler() *Profiler {
	once.Do(func() {
		profiler = &Profiler{}
		go http.ListenAndServe(":"+env.GetString("PROFILER_PORT", "6060"), nil)
	})
	return profiler
}
