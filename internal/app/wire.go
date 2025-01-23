//go:build wireinject
// +build wireinject

package app

import (
	"bbs-logger-consumer/internal/config"
	sub "bbs-logger-consumer/internal/subscriber"
	"fmt"
	"net/http"
	"os"

	iRepo "bbs-logger-consumer/internal/interfaces"
	"bbs-logger-consumer/internal/repositories"

	"bbs-logger-consumer/internal/services"
	serviceI "bbs-logger-consumer/internal/services/interfaces"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
)

var (
	ENV string
)

func getENV() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // Default to development if not set
	}
	ENV = env
}

// Provide the configuration
func provideConfigs() *config.Config {
	getENV()
	// Access configuration statically
	return config.LoadConfig(ENV)
}

// Provide Elasticsearch client
func provideElasticClient() (*elasticsearch.Client, error) {
	// Access configuration statically
	configs := provideConfigs()

	cfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("%s:%d", configs.DB.Host, configs.DB.Port)},
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
	}
	return elasticsearch.NewClient(cfg)
}

// Provide the indices for logging
func provideIndices() map[string]string {
	// Access configuration statically
	configs := provideConfigs()

	return map[string]string{
		"general": configs.DB.LogIndices.General,
		"info":    configs.DB.LogIndices.Info,
		"warning": configs.DB.LogIndices.Warning,
		"error":   configs.DB.LogIndices.Error,
		"debug":   configs.DB.LogIndices.Debug,
	}
}

// Provide the repository implementation
func provideRepository(esClient *elasticsearch.Client, indices map[string]string) iRepo.Repository {
	return repositories.NewElasticRepository(esClient, indices)
}

// InitializeService sets up the dependency injection for Service.
func InitService() (serviceI.Service, error) {
	wire.Build(
		provideElasticClient,
		provideIndices,
		provideRepository,
		services.NewService,
	)
	return nil, nil
}

// Provide the subscriber implementation
func provideSubscriber(s serviceI.Service, c *config.Config) iRepo.Subscriber {
	return sub.NewRedisSubscriber(s, c)
}

// Initialize Subscriber sets up the dependency injection for Service.
func InitSubscriber() (iRepo.Subscriber, error) {
	wire.Build(
		provideConfigs,
		InitService,
		provideSubscriber,
	)
	return nil, nil
}
