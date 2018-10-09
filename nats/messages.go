package nats

import "time"

type Message interface {
	Key() string
}

type LogsCreatedMessage struct {
	ID         string    `json:"id"`
	LogContent string    `json:"logContent"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (m *LogsCreatedMessage) Key() string {
	return "log.created"
}
