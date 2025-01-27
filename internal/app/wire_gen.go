// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"bbs-logger-consumer/internal/config"
	"bbs-logger-consumer/internal/interfaces"
	"bbs-logger-consumer/internal/repositories"
	services2 "bbs-logger-consumer/internal/services"
	"bbs-logger-consumer/internal/services/interfaces"
	"bbs-logger-consumer/internal/subscriber"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"os"
)

// Injectors from wire.go:

// InitializeService sets up the dependency injection for Service.
func InitService() (services.Service, error) {
	client, err := provideElasticClient()
	if err != nil {
		return nil, err
	}
	v := provideIndices()
	repository := provideRepository(client, v)
	service := services2.NewService(repository)
	return service, nil
}

// Initialize Subscriber sets up the dependency injection for Service.
func InitSubscriber() (interfaces.Subscriber, error) {
	service, err := InitService()
	if err != nil {
		return nil, err
	}
	config := provideConfigs()
	subscriber := provideSubscriber(service, config)
	return subscriber, nil
}

// wire.go:

var (
	ENV string
)

func getENV() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	ENV = env
}

// Provide the configuration
func provideConfigs() *config.Config {
	getENV()

	return config.LoadConfig(ENV)
}

// Provide Elasticsearch client
func provideElasticClient() (*elasticsearch.Client, error) {

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
func provideRepository(esClient *elasticsearch.Client, indices map[string]string) interfaces.Repository {
	return repositories.NewElasticRepository(esClient, indices)
}

// Provide the subscriber implementation
func provideSubscriber(s services.Service, c *config.Config) interfaces.Subscriber {
	return subscriber.NewRedisSubscriber(s, c)
}
