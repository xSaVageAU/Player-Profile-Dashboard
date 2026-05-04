package models

// InventorySlot represents a single slot in an inventory (Minecraft style)
type InventorySlot struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
	Name  string `json:"name"`
	Icon  string `json:"icon"` // URL or CSS class
}

// Stats represents various player statistics
type Stats struct {
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	BlocksBroken int     `json:"blocks_broken"`
	TimePlayed   string  `json:"time_played"`
	Balance      float64 `json:"balance"`
}

// Advancement represents a player's achievement
type Advancement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Completed   bool   `json:"completed"`
}

// PlayerProfile represents the full data for a single player
type PlayerProfile struct {
	UUID         string          `json:"uuid"`
	Username     string          `json:"username"`
	SkinURL      string          `json:"skin_url"`
	Health       float64         `json:"health"`
	MaxHealth    float64         `json:"max_health"`
	Hunger       int             `json:"hunger"`
	Level        int             `json:"level"`
	Experience   float64         `json:"experience"`
	Stats        Stats           `json:"stats"`
	Inventory    []InventorySlot `json:"inventory"`
	Armor        []InventorySlot `json:"armor"`
	EnderChest   []InventorySlot `json:"ender_chest"`
	Advancements []Advancement   `json:"advancements"`
}
