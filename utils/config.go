package utils

var (
	// Config stores global Config
	Config struct {
		// TimeScale defines the Server's time scale, the higher the faster
		TimeScale float64

		// DebugMode name is self explanatory ...
		DebugMode bool

		// Host name is self explanatory ...
		Host string

		// GroundSize ...
		GroundSize float64
	}
)

// InitializeConfig initialize Config
func InitializeConfig(timeScale float64) {
	Config.TimeScale = timeScale
	Config.GroundSize = 200.0
}
