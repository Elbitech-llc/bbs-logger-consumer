// Package services contains the business logic implementations.
// It interacts with repository layers to perform operations on data.
package services

import (
	repositories "bbs-logger-consumer/internal/interfaces"      // Importing repository interfaces
	"bbs-logger-consumer/internal/models"                       // Importing model
	serviceI "bbs-logger-consumer/internal/services/interfaces" // Importing service interface definitions
	"context"
)

// Service is a struct that holds a reference to the RepositoryI interface,
// which is used to interact with the data in the database.

type service struct {
	repo repositories.Repository
}

// NewService initializes and returns an instance of service.
func NewService(repo repositories.Repository) serviceI.Service {
	return &service{repo: repo}
}

func (s *service) AllLogs(ctx context.Context, log *models.LogMessage) error {
	return s.repo.AllLogs(ctx, log)
}

func (s *service) LogInfo(ctx context.Context, log *models.LogMessage) error {
	return s.repo.LogInfo(ctx, log)
}

func (s *service) LogWarning(ctx context.Context, log *models.LogMessage) error {
	return s.repo.LogWarning(ctx, log)
}

func (s *service) LogError(ctx context.Context, log *models.LogMessage) error {
	return s.repo.LogError(ctx, log)
}

func (s *service) LogDebug(ctx context.Context, log *models.LogMessage) error {
	return s.repo.LogDebug(ctx, log)
}
