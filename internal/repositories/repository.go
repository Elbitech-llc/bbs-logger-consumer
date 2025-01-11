package repositories

import (
	repoI "bbs-logger-consumer/internal/interfaces"
	"bbs-logger-consumer/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type elasticRepository struct {
	esClient *elasticsearch.Client
	indices  map[string]string
}

// NewElasticRepository initializes a new Elasticsearch repository.
func NewElasticRepository(client *elasticsearch.Client, indices map[string]string) repoI.Repository {
	return &elasticRepository{
		esClient: client,
		indices:  indices,
	}
}

func (r *elasticRepository) indexLog(ctx context.Context, log *models.LogMessage, logType string) error {
	log.IndexType = logType
	// Get the index based on logType
	index, exists := r.indices[logType]
	if !exists {
		return fmt.Errorf("index for log type '%s' not found", logType)
	}

	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal log: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: log.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.esClient)
	if err != nil {
		return fmt.Errorf("failed to index %s log: %w", logType, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing %s log: %s", logType, res.String())
	}

	return nil
}

func (r *elasticRepository) AllLogs(ctx context.Context, log *models.LogMessage) error {
	return r.indexLog(ctx, log, "general")
}

func (r *elasticRepository) LogInfo(ctx context.Context, log *models.LogMessage) error {
	return r.indexLog(ctx, log, "info")
}

func (r *elasticRepository) LogWarning(ctx context.Context, log *models.LogMessage) error {
	return r.indexLog(ctx, log, "warning")
}

func (r *elasticRepository) LogError(ctx context.Context, log *models.LogMessage) error {
	return r.indexLog(ctx, log, "error")
}

func (r *elasticRepository) LogDebug(ctx context.Context, log *models.LogMessage) error {
	return r.indexLog(ctx, log, "debug")
}
