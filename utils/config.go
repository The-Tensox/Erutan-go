package utils

var (
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
	Config.TickRate = tickRate
}
