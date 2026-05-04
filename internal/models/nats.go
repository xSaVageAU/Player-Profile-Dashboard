package models

import "time"

// PlayerDataBundle is the CBOR-serialized structure stored in NATS
// Field tags match the Java addon's short keys (u, n, s, a, t)
type PlayerDataBundle struct {
	UUID         interface{}            `cbor:"u" json:"uuid"`
	NBT          []byte                 `cbor:"n" json:"nbt"`
	Stats        map[string]interface{} `cbor:"s" json:"stats"`
	Advancements map[string]interface{} `cbor:"a" json:"advancements"`
	Timestamp    time.Time              `cbor:"t" json:"timestamp"`
}

// SessionState is the JSON-serialized structure for locking
type SessionState struct {
	UUID       string `json:"uuid"`
	State      string `json:"state"` // CLEAN, DIRTY, RESTORING
	LastServer string `json:"lastServer"`
}
