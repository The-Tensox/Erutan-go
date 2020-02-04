package utils

import "sync"

var (
	// ConfigMtx ...
	ConfigMtx sync.RWMutex

	// Config stores global Config
	Config struct {
		// TickRate defines the Server's tick rate, the lower the faster
		TickRate float64

		// DebugMode name is self explanatory ...
		DebugMode bool

		// Host name is self explanatory ...
		Host string
	}
)

// InitializeConfig initialize Config
func InitializeConfig(tickRate float64) {
	ConfigMtx.Lock()
	Config.TickRate = tickRate
	ConfigMtx.Unlock()
}
