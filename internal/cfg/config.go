package cfg

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

var current *Config // Good old global variable :D

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`


	/******* LOGIC *******/
	Logic struct {
		// TimeScale defines the Server's time scale, the higher the faster
		TimeScale float64 `yaml:"time_scale"`

		OctreeSize float64 `yaml:"octree_size"`

		// GroundSize ...
		GroundSize float64 `yaml:"ground_size"`

		Herbivorous struct {
			ReproductionThreshold float64 `yaml:"reproduction_threshold"`
			ReproductionLifeLoss float64 `yaml:"reproduction_life_loss"`
			EatLifeGain float64 `yaml:"eat_life_gain"`
			LifeLossRate float64 `yaml:"life_loss_rate"`
			Quantity int `yaml:"quantity"`
			Life float64 `yaml:"life"`
		} `yaml:"herbivorous"`

		InitialHerbs       int `yaml:"initial_herbs"`

		Player struct {
			// Can't use protometry.Vector3 because can't extend type outside package
			Spawn []float64 `yaml:"spawn"`
			Culling float64 `yaml:"culling"`
		} `yaml:"player"`
	} `yaml:"logic"`

	// DebugMode name is self explanatory ...
	DebugMode   bool    `yaml:"debug_mode"`
	Ssl bool `yaml:"ssl"`
	SslCert     string  `yaml:"ssl_cert"`
	SslKey      string  `yaml:"ssl_key"`
	MetricsPort string  `yaml:"metrics_port"`
	NetworkRate float64 `yaml:"network_rate"`
	UpdatesRate float64 `yaml:"updates_rate"`
}

func (c Config) String() string {
	return fmt.Sprintf("{ Server: %v, Logic: %v,DebugMode: %v, SslCert: %v, SslKey: " +
		"%v, MetricsPort: %v, NetworkRate: %v, FPS: %v }",
		c.Server, c.Logic, c.DebugMode, c.SslCert, c.SslKey, c.MetricsPort, c.NetworkRate, c.UpdatesRate)
}

func initialize() {
	var cfg Config
	readFile(&cfg)
	readEnv(&cfg)
	current = &cfg
}

func Get() *Config {
	if current == nil {
		initialize()
	}
	return current
}



func readFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		os.Exit(2)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		os.Exit(2)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		os.Exit(2)
	}
}