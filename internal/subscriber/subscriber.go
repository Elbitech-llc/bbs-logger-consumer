package subscriber

import (
	"bbs-logger-consumer/internal/config"
	sub "bbs-logger-consumer/internal/interfaces"
	"bbs-logger-consumer/internal/models"
	services "bbs-logger-consumer/internal/services/interfaces"
	"context"

	"encoding/json"
	"fmt"
	"log"
	"sync"

	caching "bbs-logger-consumer/pkg"

	"github.com/go-redis/redis/v8"
)

// type Subscriber interface {
// 	ListenForLogsUpdates() error
// }

type redisSubscriber struct {
	s   services.Service
	ctx context.Context
	rdb *redis.Client
	cfg *config.Config
}

func NewRedisSubscriber(s services.Service, c *config.Config) sub.Subscriber {
	// Initialize Redis
	rdb, ctx, err := caching.InitializeRedis(c.RedisHost, c.RedisPort, c.RedisDB)
	if err != nil {
		fmt.Printf("failed to initialize Redis: %v", err)
		return nil
	}

	return &redisSubscriber{s: s, ctx: ctx, rdb: rdb, cfg: c}
}

func (rs *redisSubscriber) listenForLogsUpdates(ctx context.Context) {
	ls := rs.rdb.Subscribe(ctx, rs.cfg.LogsChannel)
	defer ls.Close()

	for msg := range ls.Channel() {
		fmt.Printf("Received log message: %s", msg.Payload)
		var logMsg models.LogMessage

		err := json.Unmarshal([]byte(msg.Payload), &logMsg)
		if err != nil {
			log.Printf("Error deserializing log: %v", err)
			continue
		}

		err = rs.s.AllLogs(ctx, &models.LogMessage{
			ID:        logMsg.ID,
			Timestamp: logMsg.Timestamp,
			Level:     logMsg.Level,
			IndexType: logMsg.IndexType,
			Message:   logMsg.Message,
		})

		if err != nil {
			log.Printf("Error creating log message model: %v", err)
			continue
		}
	}
}

func (rs *redisSubscriber) listenForInfoLogsUpdates(ctx context.Context) {
	ls := rs.rdb.Subscribe(ctx, rs.cfg.InfoChannel)
	defer ls.Close()

	for msg := range ls.Channel() {
		fmt.Printf("Received log message: %s", msg.Payload)
		var logMsg models.LogMessage

		err := json.Unmarshal([]byte(msg.Payload), &logMsg)
		if err != nil {
			log.Printf("Error deserializing log: %v", err)
			continue
		}

		err = rs.s.LogInfo(ctx, &models.LogMessage{
			ID:        logMsg.ID,
			Timestamp: logMsg.Timestamp,
			Level:     logMsg.Level,
			IndexType: logMsg.IndexType,
			Message:   logMsg.Message,
		})

		if err != nil {
			log.Printf("Error creating log message model: %v", err)
			continue
		}
	}
}

func (rs *redisSubscriber) listenForWarningLogsUpdates(ctx context.Context) {
	ls := rs.rdb.Subscribe(ctx, rs.cfg.WarningChannel)
	defer ls.Close()

	for msg := range ls.Channel() {
		fmt.Printf("Received log message: %s", msg.Payload)
		var logMsg models.LogMessage

		err := json.Unmarshal([]byte(msg.Payload), &logMsg)
		if err != nil {
			log.Printf("Error deserializing log: %v", err)
			continue
		}

		err = rs.s.LogWarning(ctx, &models.LogMessage{
			ID:        logMsg.ID,
			Timestamp: logMsg.Timestamp,
			Level:     logMsg.Level,
			IndexType: logMsg.IndexType,
			Message:   logMsg.Message,
		})

		if err != nil {
			log.Printf("Error creating warning log message model: %v", err)
			continue
		}
	}
}

func (rs *redisSubscriber) listenForErrorLogsUpdates(ctx context.Context) {
	ls := rs.rdb.Subscribe(ctx, rs.cfg.ErrorChannel)
	defer ls.Close()

	for msg := range ls.Channel() {
		fmt.Printf("Received error log message: %s", msg.Payload)
		var logMsg models.LogMessage

		err := json.Unmarshal([]byte(msg.Payload), &logMsg)
		if err != nil {
			log.Printf("Error deserializing log: %v", err)
			continue
		}

		err = rs.s.LogError(ctx, &models.LogMessage{
			ID:        logMsg.ID,
			Timestamp: logMsg.Timestamp,
			Level:     logMsg.Level,
			IndexType: logMsg.IndexType,
			Message:   logMsg.Message,
		})

		if err != nil {
			log.Printf("Error creating error log message model: %v", err)
			continue
		}
	}
}

func (rs *redisSubscriber) listenForCustomLogsUpdates(ctx context.Context) {
	ls := rs.rdb.Subscribe(ctx, rs.cfg.DebugChannel)
	defer ls.Close()

	for msg := range ls.Channel() {
		fmt.Printf("Received custom log message: %s", msg.Payload)
		var logMsg models.LogMessage

		err := json.Unmarshal([]byte(msg.Payload), &logMsg)
		if err != nil {
			log.Printf("Error deserializing log: %v", err)
			continue
		}

		err = rs.s.LogDebug(ctx, &models.LogMessage{
			ID:        logMsg.ID,
			Timestamp: logMsg.Timestamp,
			Level:     logMsg.Level,
			IndexType: logMsg.IndexType,
			Message:   logMsg.Message,
		})

		if err != nil {
			log.Printf("Error creating custom log message model: %v", err)
			continue
		}
	}
}

// ListenForLogsUpdates subscribes to Redis channels and processes logs
func (rs *redisSubscriber) ListenForLogsUpdates() error {
	// Create a wait group for concurrent processing
	var wg sync.WaitGroup
	errCh := make(chan error, 5) // Buffer size matches the number of goroutines

	wg.Add(5)

	// Process logs updates
	go func() {
		defer wg.Done()
		if err := rs.safeListen(rs.listenForLogsUpdates); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := rs.safeListen(rs.listenForInfoLogsUpdates); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := rs.safeListen(rs.listenForErrorLogsUpdates); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := rs.safeListen(rs.listenForWarningLogsUpdates); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		if err := rs.safeListen(rs.listenForCustomLogsUpdates); err != nil {
			errCh <- err
		}
	}()

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Collect errors from the channel
	var combinedErr error
	for err := range errCh {
		if combinedErr == nil {
			combinedErr = err
		} else {
			combinedErr = fmt.Errorf("%v; %w", combinedErr, err)
		}
	}

	return combinedErr
}

// A helper to safely execute a function and recover from panics
func (rs *redisSubscriber) safeListen(fn func(ctx context.Context)) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in listener: %v", r)
		}
	}()

	fn(rs.ctx)
	return nil
}
