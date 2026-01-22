package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoUri string `envconfig:"mongodb_uri"`
}

var cfg Config

var doOnce sync.Once

func GetConfig() *Config {
	doOnce.Do(func() {
		cfg := Config{}
		err := envconfig.Process("", &cfg)
		if err != nil {
			fmt.Printf("Cannot read configuration: %s", err)
			os.Exit(2)
		}
	})
	return &cfg
}
