package main

import (
	initDI "bbs-logger-consumer/internal/app"
	"log"
	"net/http"
	"sync"
	//"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// Start HTTP server (Gin)
	httpServer := &http.Server{
		Addr: ":7889",
	}

	// Create a wait group for concurrent processing
	var wg sync.WaitGroup
	wg.Add(2)

	s, err := initDI.InitSubscriber()
	if err != nil {
		log.Fatalf("Failed to Init Subscriber: %s", err)
	}

	go func() {
		defer wg.Done()
		if err := s.ListenForLogsUpdates(); err != nil {
			log.Fatal("Failed to subscribe to redis: ", err)
		}
	}()

	// Run HTTP server in a goroutine
	go func() {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Failed to start HTTP server: ", err)
		}
	}()

	// Wait for both servers to finish
	wg.Wait()
}
