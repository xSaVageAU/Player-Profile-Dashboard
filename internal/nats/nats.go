package nats

import (
	"encoding/json"
	"fmt"
	"log"
	"playerprofile/internal/config"
	"playerprofile/internal/models"

	"github.com/fxamacker/cbor/v2"
	"github.com/klauspost/compress/zstd"
	"github.com/nats-io/nats.go"
	"strconv"
	"strings"
)

type Client struct {
	Conn      *nats.Conn
	JS        nats.JetStreamContext
	KV        nats.KeyValue
	EconomyKV nats.KeyValue
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

	kv, err := js.KeyValue(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get KV bucket %s: %w", cfg.Bucket, err)
	}

	econKV, err := js.KeyValue(cfg.EconomyBucket)
	if err != nil {
		log.Printf("Warning: Failed to get economy KV bucket %s: %v", cfg.EconomyBucket, err)
	}

	log.Printf("Connected to NATS. Buckets: %s, %s", cfg.Bucket, cfg.EconomyBucket)

	return &Client{
		Conn:      nc,
		JS:        js,
		KV:        kv,
		EconomyKV: econKV,
	}, nil
}

// FetchBalance retrieves player economy data from the economy KV bucket
func (c *Client) FetchBalance(uuid string) (float64, error) {
	if c.EconomyKV == nil {
		return 0, fmt.Errorf("economy KV not initialized")
	}

	log.Printf("Fetching economy for: %s", uuid)
	entry, err := c.EconomyKV.Get(uuid)
	if err != nil {
		// Try undashed
		undashed := strings.ReplaceAll(uuid, "-", "")
		log.Printf("Dashed economy key failed, trying undashed: %s", undashed)
		entry, err = c.EconomyKV.Get(undashed)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get economy for %s: %w", uuid, err)
	}

	log.Printf("Successfully fetched economy data: %s", string(entry.Value()))

	var data map[string]interface{}
	if err := json.Unmarshal(entry.Value(), &data); err != nil {
		return 0, fmt.Errorf("failed to decode economy JSON: %w", err)
	}

	// Try common keys: "balance", "money", "amount"
	for _, key := range []string{"balance", "money", "amount"} {
		if val, ok := data[key]; ok {
			switch v := val.(type) {
			case float64:
				return v, nil
			case int64:
				return float64(v), nil
			case string:
				f, err := strconv.ParseFloat(v, 64)
				if err == nil {
					return f, nil
				}
			}
		}
	}

	return 0, nil
}

// FetchBundle retrieves a player data bundle from NATS, decompresses it, and decodes CBOR
func (c *Client) FetchBundle(uuid string) (*models.PlayerDataBundle, error) {
	key := fmt.Sprintf("bundle.%s", uuid)
	log.Printf("Fetching NATS key: %s", key)
	entry, err := c.KV.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get bundle for %s (key: %s): %w", uuid, key, err)
	}

	log.Printf("Successfully fetched key %s, size: %d bytes", key, len(entry.Value()))

	// 1. Decompress Zstd
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	decompressed, err := decoder.DecodeAll(entry.Value(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress bundle: %w", err)
	}

	log.Printf("Decompressed size: %d bytes", len(decompressed))

	// 2. Decode CBOR
	var bundle models.PlayerDataBundle
	if err := cbor.Unmarshal(decompressed, &bundle); err != nil {
		return nil, fmt.Errorf("failed to decode CBOR bundle: %w", err)
	}

	log.Printf("Decoded bundle for UUID: %s, NBT size: %d", bundle.UUID, len(bundle.NBT))
	return &bundle, nil
}

// FetchSession retrieves a player session state from NATS
func (c *Client) FetchSession(uuid string) (*models.SessionState, error) {
	key := fmt.Sprintf("session.%s", uuid)
	entry, err := c.KV.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get session for %s: %w", uuid, err)
	}

	// Note: Check if sessions are also Zstd compressed.
	// If the helper said "JSON-serialized", it might be raw.
	// But the helper also said "The data is compressed with Zstd... so you'll need to decompress"
	// which might apply to everything in that bucket.
	
	data := entry.Value()
	
	// Check for Zstd magic header (0xFD2FB528)
	if len(data) > 4 && data[0] == 0x28 && data[1] == 0xB5 && data[2] == 0x2F && data[3] == 0xFD {
		decoder, _ := zstd.NewReader(nil)
		data, err = decoder.DecodeAll(data, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress session: %w", err)
		}
	}

	var session models.SessionState
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to decode JSON session: %w", err)
	}

	return &session, nil
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
