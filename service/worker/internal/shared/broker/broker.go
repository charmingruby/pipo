package broker

import "context"

type Broker interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error
}
