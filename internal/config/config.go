package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Type       string
		Host       string
		Port       int
		LogIndices struct {
			General string
			Info    string
			Warning string
			Error   string
			Debug   string
		}
	}

	Redis struct {
		Host   string
		Port   string
		LogsDB int
	}

	Logging struct {
		LogsFilePath string
		Channels     struct {
			Logs    string
			Info    string
			Warning string
			Error   string
			Debug   string
		}
	}
}

var (
	once   sync.Once
	config *Config
)

func LoadConfig(env string) *Config {
	once.Do(func() {
		fileName := fmt.Sprintf("config.%s.json", env)

		viper.SetConfigName(fileName)
		viper.SetConfigType("json")
		viper.AddConfigPath("./internal/config/")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		config = &Config{}

		err = viper.Unmarshal(&config)
		if err != nil {
			log.Fatalf("Unable to parse configuration into struct: %v", err)
		}
	})
	return config
}
