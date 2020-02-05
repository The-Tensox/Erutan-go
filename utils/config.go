package utils

import "sync"

var (
	// ConfigMtx ...
	ConfigMtx sync.RWMutex

	// Config stores global Config
	Config struct {
		// TimeScale defines the Server's time scale, the higher the faster
		TimeScale float64

		// DebugMode name is self explanatory ...
		DebugMode bool

		// Host name is self explanatory ...
		Host string
	}
)

// InitializeConfig initialize Config
func InitializeConfig(timeScale float64) {
	Config.TimeScale = timeScale
}
