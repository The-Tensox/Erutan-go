package mon

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(NetworkActionIgnoreCounter)
	prometheus.MustRegister(NetworkActionUpdateCounter)
	prometheus.MustRegister(NetworkActionDestroyCounter)
	prometheus.MustRegister(LifeGauge)
	prometheus.MustRegister(SpeedGauge)
	prometheus.MustRegister(CollisionCounter)
	prometheus.MustRegister(PhysicalObjectsGauge)
	prometheus.MustRegister(EatCounter)
	prometheus.MustRegister(ReproductionCounter)
}

var (
	// Physics
	CollisionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "collision",
			Help: "collision",
		},
	)
	PhysicalObjectsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "physical_object",
			Help: "physical_object",
		},
	)

	// Living being
	LifeGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "life",
			Help: "life",
		},
	)
	SpeedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "speed",
			Help: "speed",
		},
	)
	EatCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "eat",
			Help: "eat",
		},
	)
	ReproductionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "reproduction",
			Help: "reproduction",
		},
	)

	// Network high-level
	NetworkActionIgnoreCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "network_action_ignore",
			Help: "Network action ignore",
		},
	)

	NetworkActionUpdateCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "network_action_update",
			Help: "Network action update",
		},
	)

	NetworkActionDestroyCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "network_action_destroy",
			Help: "Network action destroy",
		},
	)
)
