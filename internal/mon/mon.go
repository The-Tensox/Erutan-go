package mon

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(NetworkAddCounter)
	prometheus.MustRegister(NetworkRemoveCounter)
	prometheus.MustRegister(LifeGauge)
	prometheus.MustRegister(SpeedGauge)
	prometheus.MustRegister(CollisionCounter)
	prometheus.MustRegister(PhysicalObjectsGauge)
	prometheus.MustRegister(EatCounter)
	prometheus.MustRegister(ReproductionCounter)
	prometheus.MustRegister(VolumeGauge)
	prometheus.MustRegister(ObserverEventCounter)
}

var (
	// Physics
	CollisionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "physics_collision",
			Help: "collision",
		},
	)
	PhysicalObjectsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "physics_object",
			Help: "physics_object",
		},
	)

	// Living being
	LifeGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "herbivorous_life",
			Help: "herbivorous_life",
		},
	)
	SpeedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "herbivorous_speed",
			Help: "herbivorous_speed",
		},
	)
	EatCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "herbivorous_eat",
			Help: "herbivorous_eat",
		},
	)
	ReproductionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "herbivorous_reproduction",
			Help: "herbivorous_reproduction",
		},
	)
	VolumeGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "herbivorous_volume",
			Help: "herbivorous_volume",
		},
	)

	// Network high-level
	NetworkAddCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "network_add",
			Help: "network_add",
		},
	)
	NetworkRemoveCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "network_remove",
			Help: "network_remove",
		},
	)

	// Observer
	ObserverEventCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "observer_event",
			Help: "observer_event",
		},
	)
)
