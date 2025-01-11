# Real-time Log Consumer and HTTP Server
This repository contains a Go-based application that integrates a Redis-based log subscription system with an HTTP server for efficient log processing and server management. The solution consists of two main components:

## Package: Subscriber
- Redis Integration: Implements a Redis pub/sub mechanism to listen for log updates from multiple channels (e.g., Info, Warning, Error, Debug).
- Concurrent Log Processing: Leverages goroutines and wait groups to handle logs in parallel.
- Error Resilience: Incorporates safe execution methods to recover from panics and aggregate errors.
- Extensibility: Easily expandable for additional log levels or processing logic.

## Package: Main
- HTTP Server: Initializes and runs an HTTP server on port 7889, providing the foundation for potential REST endpoints or health monitoring.
- Dependency Injection: Uses a modular initialization process for the subscriber component to ensure flexibility and maintainability.
- Concurrency: Manages the log subscription service and HTTP server concurrently with synchronization primitives (e.g., wait groups).
- Error Handling: Ensures robust startup and graceful shutdown by capturing and logging critical errors.

Acknowledgments

Special thanks to [Albert Karapetyan](https://github.com/AlbertKarapetyan) for leading the development of this service.