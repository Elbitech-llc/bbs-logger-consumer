package main

import (
	initDI "bbs-logger-consumer/internal/app"
	"log"
	"os"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev" // Default to development if not set
	}

	log.Printf("Logger consumer[%s] is starting...", env)

	s, err := initDI.InitSubscriber()
	if err != nil {
		log.Fatalf("Failed to Init Subscriber: %s", err)
	}

	if err := s.ListenForLogsUpdates(); err != nil {
		log.Fatal("Failed to subscribe to redis: ", err)
	}
}
