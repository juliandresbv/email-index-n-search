package profiling

import (
	"flag"

	"github.com/pkg/profile"
)

const (
	cpu       = "cpu"
	mem       = "mem"
	goroutine = "goroutine"
	thread    = "thread"
)

type DefaultProfile struct{}

func (DefaultProfile) Stop() {}

func SetupProfiling() interface{ Stop() } {
	profilingPath := "./profiling-results"

	profileMode := flag.String("profile.mode", "", "enable profiling mode, one of [cpu, mem, goroutine, thread]")

	flag.Parse()

	switch *profileMode {
	case cpu:
		return profile.Start(profile.CPUProfile, profile.ProfilePath(profilingPath))
	case mem:
		return profile.Start(profile.MemProfile, profile.ProfilePath(profilingPath))
	case goroutine:
		return profile.Start(profile.GoroutineProfile, profile.ProfilePath(profilingPath))
	case thread:
		return profile.Start(profile.ThreadcreationProfile, profile.ProfilePath(profilingPath))
	default:
		return DefaultProfile{}
	}
}
