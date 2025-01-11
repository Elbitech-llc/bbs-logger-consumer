package services

import (
	"bbs-logger-consumer/internal/models"
	"context"
)

type Service interface {
	AllLogs(ctx context.Context, log *models.LogMessage) error
	LogInfo(ctx context.Context, log *models.LogMessage) error
	LogError(ctx context.Context, log *models.LogMessage) error
	LogWarning(ctx context.Context, log *models.LogMessage) error
	LogDebug(ctx context.Context, log *models.LogMessage) error
}
