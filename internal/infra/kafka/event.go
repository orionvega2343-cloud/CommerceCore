package kafka

import (
	"encoding/json"
	"time"
)

type Event struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewEvent(eventType string, payload any) (Event, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return Event{}, err
	}
	event := Event{EventType: eventType, Payload: b, Timestamp: time.Now()}
	return event, nil
}
