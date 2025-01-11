//go:build wireinject
// +build wireinject

package app

import (
	"bbs-logger-consumer/internal/config"
	sub "bbs-logger-consumer/internal/subscriber"
	"fmt"
	"net/http"

	iRepo "bbs-logger-consumer/internal/interfaces"
	"bbs-logger-consumer/internal/repositories"

	"bbs-logger-consumer/internal/services"
	serviceI "bbs-logger-consumer/internal/services/interfaces"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/wire"
)

// Provide Elasticsearch client
func provideElasticClient() (*elasticsearch.Client, error) {
	// Access configuration statically
	configs := config.LoadConfig()

	cfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("%s:%d", configs.DBHost, configs.DBPort)},
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
	}
	return elasticsearch.NewClient(cfg)
}

// Provide the indexes for logging
func provideIndexes() map[string]string {
	// Access configuration statically
	configs := config.LoadConfig()

	return map[string]string{
		"general": configs.GeneralIndex,
		"info":    configs.InfoIndex,
		"warning": configs.WarningIndex,
		"error":   configs.ErrorIndex,
		"debug":   configs.DebugIndex,
	}
}

// Provide the configuration
func provideConfigs() *config.Config {
	// Access configuration statically
	return config.LoadConfig()
}

// Provide the repository implementation
func provideRepository(esClient *elasticsearch.Client, indexes map[string]string) iRepo.Repository {
	return repositories.NewElasticRepository(esClient, indexes)
}

// InitializeService sets up the dependency injection for Service.
func InitService() (serviceI.Service, error) {
	wire.Build(
		provideElasticClient,
		provideIndexes,
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
		InitService,
		provideConfigs,
		provideSubscriber,
	)
	return nil, nil
}
