package models

type LogMessage struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	IndexType string `json:"-"`
}
