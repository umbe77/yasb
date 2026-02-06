package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoUri string `envconfig:"mongo_db_uri"`
}

var (
	cfg    Config
	doOnce sync.Once
)

func GetConfig() *Config {
	doOnce.Do(func() {
		err := envconfig.Process("", &cfg)
		if err != nil {
			fmt.Printf("Cannot read configuration: %s", err)
			os.Exit(2)
		}
	})
	return &cfg
}
