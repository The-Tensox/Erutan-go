package cfg

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

var Global Config // Temporary solution xD

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	// DebugMode name is self explanatory ...
	DebugMode bool `yaml:"debug_mode"`

	/******* LOGIC *******/
	Logic struct {
		// TimeScale defines the Server's time scale, the higher the faster
		TimeScale float64 `yaml:"time_scale"`

		// GroundSize ...
		GroundSize float64 `yaml:"ground_size"`

		InitialHerbivorous int `yaml:"initial_herbivorous"`
		InitialHerbs       int `yaml:"initial_herbs"`
	} `yaml:"logic"`
}

func (c Config) String() string {
	return fmt.Sprintf("[TODO]")
}

func Get() Config {
	var cfg Config
	readFile(&cfg)
	readEnv(&cfg)
	return cfg
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