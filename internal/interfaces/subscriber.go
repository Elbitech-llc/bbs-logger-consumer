package interfaces

type Subscriber interface {
	ListenForLogsUpdates() error
}
