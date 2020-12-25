package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	eventstore "github.com/go-event-store/eventstore"
)

func EventLogger(action eventstore.DomainEventAction) eventstore.DomainEventMiddleware {
	return func(ctx context.Context, event eventstore.DomainEvent) (eventstore.DomainEvent, error) {
		payload, err := json.Marshal(event.Payload())
		if err != nil {
			fmt.Printf("[%s] %s: Failed to convert Payload: %s\n", event.CreatedAt().Format(time.RFC3339), event.Name(), err.Error())

			return event, err
		}

		fmt.Printf("[EventStore][%s][%s] %s: Payload: %s\n", action, event.CreatedAt().Format(time.RFC3339), event.Name(), payload)

		return event, nil
	}
}
