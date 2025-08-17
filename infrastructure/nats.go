package bsgostuff_infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	bsgostuff_config "github.com/beavernsticks/go-stuff/config"
	bsgostuff_events "github.com/beavernsticks/go-stuff/events"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type NATSBroker struct {
	conn      *nats.Conn
	js        nats.JetStreamContext
	envPrefix string
	mu        sync.Mutex
	subs      map[string]*nats.Subscription
}

func NewNATSBroker(cfg bsgostuff_config.NATS) (*NATSBroker, error) {
	conn, err := nats.Connect(cfg.URL,
		nats.MaxReconnects(5),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect failed: %w", err)
	}

	js, err := conn.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, fmt.Errorf("jetstream init failed: %w", err)
	}

	var prefix string
	if cfg.TopicPrefix != "" {
		prefix = cfg.TopicPrefix + "."
	}

	return &NATSBroker{
		conn:      conn,
		js:        js,
		envPrefix: prefix,
		subs:      make(map[string]*nats.Subscription),
	}, nil
}

func (b *NATSBroker) Publish(ctx context.Context, topic string, msg proto.Message, opts ...nats.PubOpt) error {
	fullTopic := b.fullTopic(topic)
	payload, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("proto marshal failed: %w", err)
	}

	return b.retry(ctx, 3, 100*time.Millisecond, func() error {
		_, err := b.js.Publish(fullTopic, payload, opts...)
		if errors.Is(err, nats.ErrNoResponders) {
			return fmt.Errorf("publish failed (no responders): %w", err)
		}
		return err
	})
}

func (b *NATSBroker) Subscribe(
	parentCtx context.Context,
	topic string,
	queueGroup string,
	handler func(context.Context, proto.Message) error,
	protoTemplate proto.Message,
) error {
	fullTopic := b.fullTopic(topic)
	durableName := fmt.Sprintf("%s-%s", b.envPrefix, queueGroup)
	dlqTopic := b.fullTopic(fmt.Sprintf("dlq.%s.%s", queueGroup, topic))

	_, err := b.js.AddStream(&nats.StreamConfig{
		Name:      fmt.Sprintf("DLQ_%s", durableName),
		Subjects:  []string{dlqTopic},
		Retention: nats.LimitsPolicy,
	})
	if err != nil && !errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		return fmt.Errorf("create DLQ stream failed: %w", err)
	}

	handlerCtx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	sub, err := b.js.QueueSubscribe(
		fullTopic,
		queueGroup,
		func(msg *nats.Msg) {
			event := proto.Clone(protoTemplate)
			if err := proto.Unmarshal(msg.Data, event); err != nil {
				log.Printf("[ERROR] unmarshal failed: %v", err)
				_ = msg.Term()
				return
			}

			msgCtx, cancel := context.WithTimeout(handlerCtx, 10*time.Second)
			defer cancel()

			processErr := b.retry(msgCtx, 3, 1*time.Second, func() error {
				return handler(msgCtx, event)
			})

			if processErr != nil {
				log.Printf("[WARN] sending to DLQ after retries: %v", processErr)
				_ = b.Publish(context.Background(), dlqTopic, &bsgostuff_events.DeadLetter{
					OriginalTopic: fullTopic,
					Payload:       msg.Data,
					Error:         processErr.Error(),
					Timestamp:     time.Now().Format(time.RFC3339),
				})
				_ = msg.Term()
			} else {
				_ = msg.Ack()
			}
		},
		nats.Durable(durableName),
		nats.ManualAck(),
		nats.AckWait(30*time.Second),
		nats.Context(parentCtx),
	)
	if err != nil {
		return err
	}

	b.mu.Lock()
	b.subs[fmt.Sprintf("%s|%s", fullTopic, queueGroup)] = sub
	b.mu.Unlock()

	return nil
}

func (b *NATSBroker) retry(ctx context.Context, maxAttempts int, initialDelay time.Duration, fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context canceled: %w", err)
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if !isRetriable(lastErr) {
			return fmt.Errorf("non-retriable error: %w", lastErr)
		}

		if attempt < maxAttempts {
			delay := initialDelay * time.Duration(1<<(attempt-1))
			delay += time.Duration(rand.Intn(500)) * time.Millisecond

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}

func isRetriable(err error) bool {
	return errors.Is(err, nats.ErrTimeout) ||
		errors.Is(err, nats.ErrNoResponders) ||
		errors.Is(err, context.DeadlineExceeded)
}

func (b *NATSBroker) fullTopic(topic string) string {
	return fmt.Sprintf("%s.%s", b.envPrefix, topic)
}

func (b *NATSBroker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, sub := range b.subs {
		_ = sub.Drain()
	}
	b.conn.Close()
}
