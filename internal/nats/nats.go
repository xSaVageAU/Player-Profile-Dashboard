package nats

import (
	"fmt"
	"log"
	"playerprofile/internal/config"

	"github.com/nats-io/nats.go"
)

type Client struct {
	Conn *nats.Conn
	JS   nats.JetStreamContext
}

func Connect(cfg config.NATSConfig) (*Client, error) {
	opts := []nats.Option{
		nats.Token(cfg.Token),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	log.Printf("Connected to NATS JetStream at %s", cfg.URL)

	return &Client{
		Conn: nc,
		JS:   js,
	}, nil
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
