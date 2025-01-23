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

	"bbs-logger-consumer/pkg"

	"github.com/go-redis/redis/v8"
)

type redisSubscriber struct {
	s   services.Service
	ctx context.Context
	rdb *redis.Client
	cfg *config.Config
}

func NewRedisSubscriber(s services.Service, c *config.Config) sub.Subscriber {
	if c == nil || c.Redis.Host == "" || c.Redis.Port == "" {
		log.Println("Invalid configuration for RedisSubscriber")
		return nil
	}

	// Initialize Redis
	rdb, ctx, err := pkg.InitializeRedis(c.Redis.Host, c.Redis.Port, c.Redis.LogsDB)
	if err != nil {
		fmt.Printf("failed to initialize Redis: %v", err)
		return nil
	}

	return &redisSubscriber{s: s, ctx: ctx, rdb: rdb, cfg: c}
}

func (rs *redisSubscriber) listenForLogsUpdates(ctx context.Context) {
	if rs.rdb == nil || rs.cfg == nil || rs.cfg.Logging.Channels.Logs == "" {
		log.Println("Redis client or channel not configured for Logs")
		return
	}

	ls := rs.rdb.Subscribe(ctx, rs.cfg.Logging.Channels.Logs)
	defer ls.Close()

	log.Println("Starting to listen for the channel LOGS")

	for msg := range ls.Channel() {
		fmt.Printf("Received: %s", msg.Payload)
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
	ls := rs.rdb.Subscribe(ctx, rs.cfg.Logging.Channels.Info)
	defer ls.Close()

	log.Println("Starting to listen for the channel INFO logs")

	for msg := range ls.Channel() {
		fmt.Printf("Received: %s", msg.Payload)
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
	ls := rs.rdb.Subscribe(ctx, rs.cfg.Logging.Channels.Warning)
	defer ls.Close()

	log.Println("Starting to listen for the channel WARNING logs")

	for msg := range ls.Channel() {
		fmt.Printf("Received: %s", msg.Payload)
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
	ls := rs.rdb.Subscribe(ctx, rs.cfg.Logging.Channels.Error)
	defer ls.Close()

	log.Println("Starting to listen for the channel ERROR logs")

	for msg := range ls.Channel() {
		fmt.Printf("Received: %s", msg.Payload)
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
	ls := rs.rdb.Subscribe(ctx, rs.cfg.Logging.Channels.Error)
	defer ls.Close()

	log.Println("Starting to listen for the channel CUSTOM logs")

	for msg := range ls.Channel() {
		fmt.Printf("Received: %s", msg.Payload)
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
	if rs.cfg.Logging.Channels.Logs == "" ||
		rs.cfg.Logging.Channels.Info == "" ||
		rs.cfg.Logging.Channels.Warning == "" ||
		rs.cfg.Logging.Channels.Error == "" {
		log.Println("One or more Redis channel names are not configured")
		return fmt.Errorf("invalid channel configuration")
	}

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

	if rs.ctx == nil {
		return fmt.Errorf("context is nil in safeListen")
	}

	fn(rs.ctx)
	return nil
}
