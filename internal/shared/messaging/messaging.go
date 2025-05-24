package messaging

import "context"

type Broker interface {
	Publish(ctx context.Context, topic string, message []byte) error
}
