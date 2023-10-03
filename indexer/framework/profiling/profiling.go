package profiling

import (
	"flag"

	"github.com/pkg/profile"
)

const (
	cpuProfMode       = "cpu"
	memProfMode       = "mem"
	goRoutineProfMode = "goroutine"
	threadProfMode    = "thread"
)

type DefaultProfile struct{}

func (DefaultProfile) Stop() {}

func SetupProfiling() interface{ Stop() } {
	profilingPath := "./profiling-results"

	profileMode := flag.String("prof.mode", "", "enable profiling mode, one of [cpu, mem, goroutine, thread]")

	flag.Parse()

	switch *profileMode {
	case cpuProfMode:
		return profile.Start(profile.CPUProfile, profile.ProfilePath(profilingPath))
	case memProfMode:
		return profile.Start(profile.MemProfile, profile.ProfilePath(profilingPath))
	case goRoutineProfMode:
		return profile.Start(profile.GoroutineProfile, profile.ProfilePath(profilingPath))
	case threadProfMode:
		return profile.Start(profile.ThreadcreationProfile, profile.ProfilePath(profilingPath))
	default:
		return DefaultProfile{}
	}
}
