package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	DBType       string
	DBHost       string
	DBPort       int
	GeneralIndex string
	InfoIndex    string
	WarningIndex string
	ErrorIndex   string
	DebugIndex   string

	RedisHost string
	RedisPort string
	RedisDB   int

	LogsChannel    string
	InfoChannel    string
	WarningChannel string
	ErrorChannel   string
	DebugChannel   string
}

var (
	once   sync.Once
	config *Config
)

func LoadConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		config = &Config{
			DBType:       viper.GetString("db_type"),
			DBHost:       viper.GetString("db_host"),
			DBPort:       viper.GetInt("db_port"),
			GeneralIndex: viper.GetString("general_index"),
			InfoIndex:    viper.GetString("info_index"),
			WarningIndex: viper.GetString("warning_index"),
			ErrorIndex:   viper.GetString("error_index"),
			DebugIndex:   viper.GetString("debug_index"),
			RedisHost:    viper.GetString("redis_host"),
			RedisPort:    viper.GetString("redis_port"),
			RedisDB:      viper.GetInt("redis_db"),

			LogsChannel:    viper.GetString("logs_channel"),
			InfoChannel:    viper.GetString("info_channel"),
			WarningChannel: viper.GetString("warning_channel"),
			ErrorChannel:   viper.GetString("error_channel"),
			DebugChannel:   viper.GetString("debug_channel"),
		}
	})
	return config
}
