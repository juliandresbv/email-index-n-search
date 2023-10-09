package profiling

import (
	"flag"

	"github.com/pkg/profile"
)

const (
	cpuProfileMode       = "cpu"
	memProfileMode       = "mem"
	goRoutineProfileMode = "goroutine"
	threadProfileMode    = "thread"
)

type defaultProfile struct{}

func (defaultProfile) Stop() {}

func SetupProfiling() interface{ Stop() } {
	profilingPath := "./profiling-results"

	profileMode := flag.String("prof.mode", "", "enable profiling mode, one of [cpu, mem, goroutine, thread]")

	flag.Parse()

	switch *profileMode {
	case cpuProfileMode:
		return profile.Start(profile.CPUProfile, profile.ProfilePath(profilingPath))
	case memProfileMode:
		return profile.Start(profile.MemProfile, profile.ProfilePath(profilingPath))
	case goRoutineProfileMode:
		return profile.Start(profile.GoroutineProfile, profile.ProfilePath(profilingPath))
	case threadProfileMode:
		return profile.Start(profile.ThreadcreationProfile, profile.ProfilePath(profilingPath))
	default:
		return defaultProfile{}
	}
}
